package query

type Params map[string]interface{}

func (p *Params) Set(k string, v interface{}) {
	if *p == nil {
		*p = make(Params)
	}
	(*p)[k] = v
}

func (p Params) Get(n string) (interface{}, bool) {
	if p == nil {
		return nil, false
	}
	v, ok := p[n]
	return v, ok
}
