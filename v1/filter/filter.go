package filter

type Option func(Filter) Filter

func WithLimit(r Range) Option {
	return func(f Filter) Filter {
		f.Limit = r
		return f
	}
}

func WithOrder(o Order) Option {
	return func(f Filter) Filter {
		f.Order = o
		return f
	}
}

func WithTimeframe(t Timeframe) Option {
	return func(f Filter) Filter {
		f.Timeframe = t
		return f
	}
}

type Filter struct {
	Limit     Range
	Order     Order
	Timeframe Timeframe
	Params    map[string]interface{}
}

func New(opts ...Option) Filter {
	f := Filter{}
	for _, o := range opts {
		f = o(f)
	}
	return f
}
