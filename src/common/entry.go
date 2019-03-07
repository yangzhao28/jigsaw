package common

// Entry place holder for values, may be module or calls
type Entry func(arg interface{})

func (e Entry) Call(arg interface{}) {
	e(arg)
}
