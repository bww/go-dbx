package option

type Config struct {
	FetchRelated  bool
	StoreRelated  bool
	DeleteRelated bool
	Params        map[string]interface{}
}

func (c Config) Param(n string) (interface{}, bool) {
	if c.Params != nil {
		v, ok := c.Params[n]
		return v, ok
	}
	return nil, false
}

type Option func(Config) Config

func NewConfig(base Config, opts []Option) Config {
	var err error
	c := base
	for _, f := range opts {
		c, err = f(c)
		if err != nil {
			return c, err
		}
	}
	return c, nil
}

func UseConfig(c Config) Option {
	return func(_ Config) Config {
		return c
	}
}

func FetchRelated(on bool) Option {
	return func(c Config) Config {
		c.FetchRelated = on
		return c
	}
}

func StoreRelated(on bool) Option {
	return func(c Config) Config {
		c.StoreRelated = on
		return c
	}
}

func DeleteRelated(on bool) Option {
	return func(c Config) Config {
		c.DeleteRelated = on
		return c
	}
}

func Params(p map[string]interface{}) Option {
	return func(c Config) Config {
		if c.Params == nil {
			c.Params = make(map[string]interface{})
		}
		for k, v := range p {
			c.Params[k] = v
		}
		return c
	}
}
