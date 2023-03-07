package filter

type Option func(Filter) Filter

func UseFilter(f Filter) Option {
	return func(_ Filter) Filter {
		return f // just replace with the provided filter
	}
}

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

func New(opts []Option) Filter {
	return Filter{}.WithOptions(opts...)
}

func (f Filter) WithOptions(opts ...Option) Filter {
	for _, o := range opts {
		f = o(f)
	}
	return f
}

func (f Filter) Param(n string) (interface{}, bool) {
	if f.Params != nil {
		v, ok := f.Params[n]
		return v, ok
	}
	return nil, false
}
