package filter

type Option func(Filter) Filter

func WithLimit(r Range) Option {
	return func(f Filter) Filter {
		f.Limit = r
		return f
	}
}

type Filter struct {
	Limit Range
}

func New(opts ...Option) Filter {
	f := Filter{}
	for _, o := range opts {
		f = o(f)
	}
	return f
}
