package option

import (
	"github.com/bww/go-dbx/v1/filter"
)

type ReadOption func(ReadConfig) (ReadConfig, error)

type ReadConfig struct {
	Filter filter.Filter
	Params Params
}

func NewReadConfig(opts []ReadOption) (ReadConfig, error) {
	return ReadConfig{}.WithOptions(opts...)
}

func (c ReadConfig) WithOptions(opts ...ReadOption) (ReadConfig, error) {
	var err error
	for _, f := range opts {
		c, err = f(c)
		if err != nil {
			return c, err
		}
	}
	return c, nil
}

func UseReadConfig(c ReadConfig) ReadOption {
	return func(_ ReadConfig) (ReadConfig, error) {
		return c, nil
	}
}

func WithFilter(f filter.Filter) ReadOption {
	return func(c ReadConfig) (ReadConfig, error) {
		c.Filter = f
		return c, nil
	}
}

func ReadParams(p Params) ReadOption {
	return func(c ReadConfig) (ReadConfig, error) {
		if c.Params == nil {
			c.Params = make(Params)
		}
		for k, v := range p {
			c.Params[k] = v
		}
		return c, nil
	}
}
