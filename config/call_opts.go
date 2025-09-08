package config

// CallOpt is a functional option for configuring client calls
type CallOpt func(*CallOpts)

// CallOpts holds configuration options for client calls
type CallOpts struct {
	Height int64
	// Add other options here as needed
}

// Apply applies the provided options to CallOpts
func (c *CallOpts) Apply(opts ...CallOpt) {
	for _, opt := range opts {
		opt(c)
	}
}

// DefaultCallOpts returns default call options
func DefaultCallOpts() *CallOpts {
	return &CallOpts{
		Height: 0, // 0 means latest height
	}
}

// Height sets the block height for the call
func Height(height int64) CallOpt {
	return func(opts *CallOpts) {
		opts.Height = height
	}
}
