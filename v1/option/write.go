package option

type WriteOption func(WriteConfig) (WriteConfig, error)

type WriteConfig struct {
	Cascade bool
	Params  Params
}

func NewWriteConfig(opts []WriteOption) (WriteConfig, error) {
	return WriteConfig{}.WithOptions(opts...)
}

func (c WriteConfig) WithOptions(opts ...WriteOption) (WriteConfig, error) {
	var err error
	for _, f := range opts {
		c, err = f(c)
		if err != nil {
			return c, err
		}
	}
	return c, nil
}

func UseWriteConfig(c WriteConfig) WriteOption {
	return func(_ WriteConfig) (WriteConfig, error) {
		return c, nil
	}
}

func WithCascade() WriteOption {
	return func(c WriteConfig) (WriteConfig, error) {
		c.Cascade = true
		return c, nil
	}
}

func WriteParams(p Params) WriteOption {
	return func(c WriteConfig) (WriteConfig, error) {
		if c.Params == nil {
			c.Params = make(Params)
		}
		for k, v := range p {
			c.Params[k] = v
		}
		return c, nil
	}
}
