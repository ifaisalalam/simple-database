package store

type Result struct {
	value interface{}
}

func (r *Result) GetValue() interface{} {
	return r.value
}
