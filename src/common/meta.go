package common

type Meta struct {
	Package   string
	Version   Version
	Author    string
	Extension map[string]string
}

func NewMeta() *Meta {
	return &Meta{
		Version:   NewVersion(),
		Extension: make(map[string]string),
	}
}
