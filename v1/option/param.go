package option

type Params map[string]interface{}

func (p Params) Get(n string) (interface{}, bool) {
	if p == nil {
		return nil, false
	}
	v, ok := p[n]
	return v, ok
}
