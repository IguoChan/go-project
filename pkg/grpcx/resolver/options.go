package resolver

type options struct {
	scheme string
}

type Option func(opts *options)

func defaultOptions() *options {
	return &options{}
}

func SetScheme(scheme string) Option {
	return func(opts *options) {
		if scheme == "" {
			return
		}
		opts.scheme = scheme
	}
}
