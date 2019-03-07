package common

// ModEventDescriptor use to describe events accepted by a mod
type ModEventDescriptor struct {
	EventPath string
	Callback  Entry
}

type ModDepDescriptor struct {
	PackageName string
	MinVersion  Version
	MaxVersion  Version
}

type Publisher interface {
	Publish(subject string, payload interface{})
}

// Mod means a module with additional features
type Mod interface {
	UUID() string
	Name() string
	Meta() Meta
	Initialize() error
	Enable()
	Disable()
	Deps() []ModDepDescriptor
	EventSpecs() []ModEventDescriptor
}
