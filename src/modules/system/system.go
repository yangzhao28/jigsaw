package system

import (
	"os"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/satori/go.uuid"
	"github.com/yangzhao28/jigsaw/src/common"
)

var (
	meta = common.Meta{
		Package: "builtin.system",
		Author:  "yangz",
		Version: common.NewVersion("0.1.0"),
	}
	deps = []common.ModDepDescriptor{}
)

type System struct {
	id  string
	pub common.Publisher
}

func New() common.Mod {
	return &System{
		id: uuid.Must(uuid.NewV4()).String(),
	}
}

func (s *System) Meta() common.Meta {
	return meta
}

func (s *System) Name() string {
	n := strings.Split(s.Meta().Package, ".")
	return n[len(n)-1]
}

func (s *System) UUID() string {
	return s.id
}

func (s *System) Deps() []common.ModDepDescriptor {
	return deps
}

func (s *System) Initialize() error {
	return nil
}

func (s *System) Enable() {
}

func (s *System) Disable() {
}

func (s *System) EventSpecs() []common.ModEventDescriptor {
	return []common.ModEventDescriptor{
		common.ModEventDescriptor{EventPath: "system.ready", Callback: s.eventReady},
		common.ModEventDescriptor{EventPath: "system.quit", Callback: s.eventQuit},
	}
}

func (s *System) eventReady(arg interface{}) {
	s.pub = arg.(common.Publisher)
}

func (s *System) eventQuit(arg interface{}) {
	code, _ := arg.(int)
	logrus.Infof("recv quit event, quit (%v)", code)
	os.Exit(code)
}
