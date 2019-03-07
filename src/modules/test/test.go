package main

import (
	"strings"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/satori/go.uuid"
	"github.com/yangzhao28/jigsaw/src/common"
)

var (
	meta = common.Meta{
		Package: "test.quit",
		Author:  "yangz",
		Version: common.NewVersion("0.1.0"),
	}
	deps = []common.ModDepDescriptor{}
)

type Test struct {
	id  string
	pub common.Publisher
}

func New() common.Mod {
	return &Test{
		id: uuid.Must(uuid.NewV4()).String(),
	}
}

func (s *Test) Meta() common.Meta {
	return meta
}

func (s *Test) Name() string {
	n := strings.Split(s.Meta().Package, ".")
	return n[len(n)-1]
}

func (s *Test) UUID() string {
	return s.id
}

func (s *Test) Deps() []common.ModDepDescriptor {
	return deps
}

func (s *Test) Initialize() error {
	return nil
}

func (s *Test) Enable() {
}

func (s *Test) Disable() {
}

func (s *Test) EventSpecs() []common.ModEventDescriptor {
	return []common.ModEventDescriptor{
		common.ModEventDescriptor{EventPath: "system.ready", Callback: s.eventRun},
	}
}

func (s *Test) eventRun(arg interface{}) {
	s.pub = arg.(common.Publisher)
	go func() {
		for i := 0; i < 10; i++ {
			time.Sleep(1 * time.Second)
			logrus.Infof("timer triggered  %v", i)
		}
		logrus.Infof("quit")
		s.pub.Publish("system.quit", 0)
	}()
}
