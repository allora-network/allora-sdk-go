package pool

import (
	"context"
	"errors"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// retryDecision is what ExecuteWithRetry does with a failed attempt.
type retryDecision int

const (
	// failOverAndPenalise means the endpoint is at fault. ExecuteWithRetry
	// reports it unhealthy, applies a backoff penalty, and tries the next
	// endpoint.
	failOverAndPenalise retryDecision = iota

	// failOverOnly means the endpoint is healthy but may not be able to answer
	// this particular request. ExecuteWithRetry tries the next endpoint without
	// a health penalty or a backoff.
	failOverOnly

	// returnToCaller means the endpoint gave a definitive answer about the
	// request. No other endpoint can answer differently, so ExecuteWithRetry
	// returns at once without a health penalty or a backoff.
	returnToCaller
)

// classifyRetry decides how ExecuteWithRetry treats a failed attempt.
//
// The question is not "did the call fail" but "would asking a different
// endpoint help, and is this endpoint to blame". Two healthy nodes serving the
// same chain can still differ in ways that matter:
//
//   - State pruning. A pruned node cannot serve a height it has discarded; an
//     archive node can. cosmos-sdk v0.50.14 reports this from
//     BaseApp.CreateQueryContext as errorsmod.Wrapf(sdkerrors.ErrInvalidRequest,
//     "failed to load state at height ..."), and errorsmod registers
//     ErrInvalidRequest with codes.Unknown, so it arrives as Unknown.
//   - Sync lag. A node behind the head answers about an older state.
//   - Node configuration. Gas limits, maximum query results, rate limits, and
//     auth on a provider's gateway are per-endpoint.
//   - Software version. During a rolling upgrade an older binary does not
//     register a newer query method and answers Unimplemented.
//
// A proxy in front of the node produces its own codes: a load balancer answers
// 429/502/503/504 and an auth gateway answers 401/403. The REST client maps the
// unambiguous ones back to gRPC codes.
//
// callerCtx is the context the caller passed to ExecuteWithRetry, not the
// per-attempt context. It separates "the caller gave up" from "this attempt ran
// out of time", which share a gRPC code but take opposite decisions.
//
// The table, with the evidence for each row:
//
//	code               decision             reason
//	OK                 returnToCaller       not an error
//	Canceled           failOverOnly         the caller's own cancellation is
//	                                        caught before this switch and
//	                                        returns at once. A Canceled that
//	                                        reaches here came from the node or
//	                                        its proxy, which says nothing about
//	                                        the endpoint's health.
//	Unknown            failOverOnly         cosmos-sdk application errors land
//	                                        here, including a pruned height and
//	                                        a height above the node's head:
//	                                        errorsmod registers
//	                                        ErrInvalidRequest and
//	                                        ErrInvalidHeight with
//	                                        codes.Unknown. Another node may
//	                                        hold the state and this one is
//	                                        healthy. A non-status error also
//	                                        reports Unknown, but
//	                                        status.FromError returns ok=false
//	                                        for it, and that is a local or
//	                                        transport fault: penalise.
//	                                        This row is a judgement call: a
//	                                        node bug also reports Unknown. Not
//	                                        penalising costs a broken endpoint
//	                                        its cooling, which the health loop
//	                                        catches; penalising costs a healthy
//	                                        endpoint its health for a
//	                                        legitimate request, which is the
//	                                        defect this policy fixes.
//	InvalidArgument    returnToCaller       the request is malformed. Every node
//	                                        rejects it identically.
//	DeadlineExceeded   failOverAndPenalise  callerCtx is alive, so the
//	                                        per-attempt deadline fired: this
//	                                        node did not answer in time. With
//	                                        callerCtx done: returnToCaller.
//	NotFound           returnToCaller       the resource does not exist and every
//	                                        node agrees. A node behind the head
//	                                        answers NotFound for a resource that
//	                                        exists; pin a height with
//	                                        config.Height to make the answer
//	                                        authoritative.
//	AlreadyExists      returnToCaller       definitive.
//	PermissionDenied   failOverAndPenalise  a gateway rejected our credentials.
//	                                        It does so on every call until
//	                                        someone fixes the config, so the
//	                                        endpoint is persistently unusable
//	                                        to us. Without a penalty the pool
//	                                        keeps choosing it first and pays a
//	                                        wasted round trip on every request.
//	                                        Cosmos query handlers do not return
//	                                        this code for a legitimate request,
//	                                        so the penalty only lands on a
//	                                        misconfigured gateway.
//	Unauthenticated    failOverAndPenalise  same.
//	ResourceExhausted  failOverAndPenalise  the node or its proxy is rate
//	                                        limiting us, so back off from this
//	                                        endpoint and try another. When the
//	                                        message reports a message larger
//	                                        than the send or receive limit
//	                                        ("larger than max", grpc-go's
//	                                        wording for both), no node will
//	                                        accept it: returnToCaller.
//	FailedPrecondition failOverAndPenalise  the node is not in a state to serve:
//	                                        not ready, still syncing, app hash
//	                                        mismatch.
//	Aborted            failOverAndPenalise  a concurrency conflict on the node;
//	                                        another endpoint is the fix.
//	OutOfRange         returnToCaller       the request is outside the valid
//	                                        range. Definitive.
//	Unimplemented      returnToCaller       this node's binary does not register
//	                                        the method. During a rolling upgrade
//	                                        another endpoint may serve it, but a
//	                                        pool normally holds one chain
//	                                        version and the caller needs the
//	                                        diagnosis more than the extra round
//	                                        trips.
//	Internal           failOverAndPenalise  the node's own bug, panic, or storage
//	                                        error.
//	Unavailable        failOverAndPenalise  the node is down or the transport
//	                                        failed.
//	DataLoss           failOverAndPenalise  the node's storage is corrupt.
func classifyRetry(err error, callerCtx context.Context) retryDecision {
	if err == nil {
		return returnToCaller
	}

	// The caller's cancellation outranks every other consideration: retrying
	// would ignore an explicit request to stop.
	if callerCtx != nil && callerCtx.Err() != nil {
		return returnToCaller
	}

	// With a live callerCtx these can only come from the attempt itself.
	if errors.Is(err, context.Canceled) {
		return failOverOnly
	}
	if errors.Is(err, context.DeadlineExceeded) {
		return failOverAndPenalise
	}

	st, ok := status.FromError(err)
	if !ok {
		// Not a gRPC status: a local or transport fault.
		return failOverAndPenalise
	}

	switch st.Code() {
	case codes.OK:
		return returnToCaller
	case codes.Unknown:
		return failOverOnly
	case codes.InvalidArgument, codes.NotFound, codes.AlreadyExists,
		codes.OutOfRange, codes.Unimplemented:
		return returnToCaller
	case codes.PermissionDenied, codes.Unauthenticated:
		return failOverAndPenalise
	case codes.Canceled:
		return failOverOnly
	case codes.ResourceExhausted:
		if isMessageTooLarge(st.Message()) {
			return returnToCaller
		}
		return failOverAndPenalise
	default:
		return failOverAndPenalise
	}
}

// isMessageTooLarge reports whether a ResourceExhausted message describes a
// payload that exceeded the send or receive limit rather than a rate limit.
// grpc-go words both limits with "larger than max"; the limit belongs to the
// request or response size, so no other endpoint will accept the payload.
func isMessageTooLarge(msg string) bool {
	return strings.Contains(msg, "larger than max")
}
