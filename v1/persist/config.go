package persist

import (
	"fmt"
	"sync"
)

var cascadeOptionDeprecatedWarningOnce sync.Once

func warnCascadeOptionDeprecated() {
	cascadeOptionDeprecatedWarningOnce.Do(func() {
		fmt.Println("dbx: CASCADE OPERATIONS ARE DEPRECATED; DO SUBFETCHING/CASCADING IN CLIENT CODE INSTEAD")
	})
}

func (c Config) Param(n string) (interface{}, bool) {
	if c.Params != nil {
		v, ok := c.Params[n]
		return v, ok
	}
	return nil, false
}

type Config struct {
	FetchRelated  bool
	StoreRelated  bool
	DeleteRelated bool
	Params        map[string]interface{}
}

func (c Config) WithOptions(opts []Option) Config {
	for _, f := range opts {
		c = f(c)
	}
	return c
}

type Option func(Config) Config

func UseConfig(c Config) Option {
	return func(_ Config) Config {
		return c
	}
}

func Cascade(on bool) Option {
	return func(c Config) Config {
		warnCascadeOptionDeprecated()
		c.FetchRelated = on
		c.StoreRelated = on
		c.DeleteRelated = on
		return c
	}
}

func FetchRelated(on bool) Option {
	return func(c Config) Config {
		warnCascadeOptionDeprecated()
		c.FetchRelated = on
		return c
	}
}

func StoreRelated(on bool) Option {
	return func(c Config) Config {
		warnCascadeOptionDeprecated()
		c.StoreRelated = on
		return c
	}
}

func DeleteRelated(on bool) Option {
	return func(c Config) Config {
		warnCascadeOptionDeprecated()
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
