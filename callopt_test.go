package allora

import (
	"testing"
)

func TestCallOpts(t *testing.T) {
	// Test default options
	opts := DefaultCallOpts()
	if opts.Height != 0 {
		t.Errorf("Expected default height to be 0, got %d", opts.Height)
	}

	// Test Height option
	heightOpt := Height(123)
	callOpts := DefaultCallOpts()
	heightOpt(callOpts)

	if callOpts.Height != 123 {
		t.Errorf("Expected height to be 123, got %d", callOpts.Height)
	}

	// Test Apply method
	callOpts2 := DefaultCallOpts()
	callOpts2.Apply(Height(456))

	if callOpts2.Height != 456 {
		t.Errorf("Expected height to be 456, got %d", callOpts2.Height)
	}

	// Test multiple options (if we add more in the future)
	callOpts3 := DefaultCallOpts()
	callOpts3.Apply(Height(789))

	if callOpts3.Height != 789 {
		t.Errorf("Expected height to be 789, got %d", callOpts3.Height)
	}
}

func TestCallOptsUsage(t *testing.T) {
	// Example of how the new API would be used

	// Old way (would be generated with height int64):
	// client.Emissions().GetParams(ctx, req, 123)

	// New way (will be generated with opts ...CallOpt):
	// client.Emissions().GetParams(ctx, req, Height(123))
	// client.Emissions().GetParams(ctx, req) // Uses default (height 0)

	// Test that we can create the options
	opts := []CallOpt{
		Height(100),
	}

	callOpts := DefaultCallOpts()
	callOpts.Apply(opts...)

	if callOpts.Height != 100 {
		t.Errorf("Expected height to be 100, got %d", callOpts.Height)
	}
}

