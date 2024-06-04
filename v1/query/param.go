package query

type Params map[string]interface{}

func (p *Params) Set(k string, v interface{}) Params {
	if *p == nil {
		*p = make(Params)
	}
	(*p)[k] = v
	return *p
}

func (p Params) Get(n string) (interface{}, bool) {
	if p == nil {
		return nil, false
	}
	v, ok := p[n]
	return v, ok
}

func (p Params) Bool(n string) (bool, bool) {
	if p == nil {
		return false, false
	}
	v, ok := p[n].(bool)
	if !ok {
		return false, false
	}
	return v, true
}

func (p Params) Int(n string) (int, bool) {
	if p == nil {
		return 0, false
	}
	v, ok := p[n].(int)
	if !ok {
		return 0, false
	}
	return v, true
}

func (p Params) Float64(n string) (float64, bool) {
	if p == nil {
		return 0, false
	}
	v, ok := p[n].(float64)
	if !ok {
		return 0, false
	}
	return v, true
}

func (p Params) String(n string) (string, bool) {
	if p == nil {
		return "", false
	}
	v, ok := p[n].(string)
	if !ok {
		return "", false
	}
	return v, true
}
