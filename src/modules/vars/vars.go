package vars

import (
	"errors"

	"github.com/yangzhao28/sketch/common"
	"github.com/yangzhao28/sketch/evbus"
)

type Vars struct {
	meta    common.Meta
	ident   common.Ident
	bus     *evbus.EvBus
	storage Storage
}

func New(bus *evbus.EvBus) common.Module {
	return &Vars{
		bus:     bus,
		storage: NewMemoryStorage(),
	}
}

func (v *Vars) Meta() common.Meta {
	return v.meta
}

func (v *Vars) Identity() *common.Ident {
	return &v.ident
}

func (v *Vars) Initialize() error {
	v.bus.Register(v, evbus.ReceiverDescription{
		Path:     "vars.set",
		Priority: 1,
		Callback: v.OnSetVar,
	})
	v.bus.Register(v, evbus.ReceiverDescription{
		Path:     "vars.get",
		Priority: 1,
		Callback: v.OnGetVar,
	})
	return nil
}

var (
	ErrInvalidKey   = errors.New("invalid key")
	ErrInvalidValue = errors.New("invalid value")
)

func (v *Vars) OnSetVar(ev *evbus.Event) error {
	key, ok := ev.Payload["key"].(string)
	if !ok {
		return ErrInvalidKey
	}
	val := ev.Payload["value"]
	v.storage.Put(key, val)
	return nil
}

func (v *Vars) OnGetVar(ev *evbus.Event) error {
	key, ok := ev.Payload["key"].(string)
	if !ok {
		return ErrInvalidKey
	}
	val := v.storage.Get(key)
	v.bus.Send(&evbus.Event{
		
	})

	return nil
}
