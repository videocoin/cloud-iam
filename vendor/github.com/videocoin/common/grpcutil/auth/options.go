package auth

var (
	defaultOptions = &options{
		shouldAuth: DefaultDeciderMethod,
	}
)

type options struct {
	shouldAuth Decider
}

// Decider function defines rules for suppressing any interceptor logs
type Decider func(fullMethodName string) bool

// DefaultDeciderMethod is the default implementation of decider to see if you should log the call
// by default this if always true so all calls are logged
func DefaultDeciderMethod(fullMethodName string) bool {
	return true
}

func evaluateServerOpt(opts []Option) *options {
	optCopy := &options{}
	*optCopy = *defaultOptions
	for _, o := range opts {
		o(optCopy)
	}
	return optCopy
}

type Option func(*options)

// WithDecider customizes the function for deciding if the gRPC interceptor logs should log.
func WithDecider(f Decider) Option {
	return func(o *options) {
		o.shouldAuth = f
	}
}
