package auth

var defaultOptions = &options{
	shouldAuth: DefaultDeciderMethod,
}

type options struct {
	shouldAuth Decider
}

type Decider func(fullMethodName string) bool

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

func WithDecider(f Decider) Option {
	return func(o *options) {
		o.shouldAuth = f
	}
}
