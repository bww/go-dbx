package query

type ReadConfig struct {
	Limit     Range
	Order     Order
	Timeframe Timeframe
	Params    Params
}

func New(opts []ReadOption) ReadConfig {
	return ReadConfig{}.WithOptions(opts)
}

func (f ReadConfig) WithOptions(opts []ReadOption) ReadConfig {
	for _, o := range opts {
		f = o(f)
	}
	return f
}

// Param is deprecated: use `Params.Get` instead
func (f ReadConfig) Param(n string) (interface{}, bool) {
	if f.Params != nil {
		v, ok := f.Params[n]
		return v, ok
	}
	return nil, false
}

type ReadOption func(ReadConfig) ReadConfig

func UseReadConfig(c ReadConfig) ReadOption {
	return func(_ ReadConfig) ReadConfig {
		return c // just replace with the provided config
	}
}

func WithLimit(r Range) ReadOption {
	return func(f ReadConfig) ReadConfig {
		f.Limit = r
		return f
	}
}

func WithOrder(o Order) ReadOption {
	return func(f ReadConfig) ReadConfig {
		f.Order = o
		return f
	}
}

func WithTimeframe(t Timeframe) ReadOption {
	return func(f ReadConfig) ReadConfig {
		f.Timeframe = t
		return f
	}
}

func WithReadParams(p Params) ReadOption {
	return func(c ReadConfig) ReadConfig {
		if c.Params == nil {
			c.Params = make(Params)
		}
		for k, v := range p {
			c.Params[k] = v
		}
		return c
	}
}
