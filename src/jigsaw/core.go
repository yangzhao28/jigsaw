package jigsaw

import (
	"fmt"
	"os"
	"path/filepath"
	"plugin"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/yangzhao28/jigsaw/src/common"
	"github.com/yangzhao28/jigsaw/src/executor"
	"github.com/yangzhao28/jigsaw/src/listener"
)

type Config struct {
	WorkerNum       int
	InputQueueSize  int
	WorkerQueueSize int
	Listener        executor.EventListener
}

func DefaultConfig() Config {
	return Config{
		WorkerNum:       10,
		InputQueueSize:  256,
		WorkerQueueSize: 32,
		Listener:        listener.New(),
	}
}

type Jigsaw struct {
	executor *executor.Executor
	quit     chan struct{}
}

// New create new jigsaw instance
func New(config Config) *Jigsaw {
	return &Jigsaw{
		executor: executor.New(
			executor.Config{
				WorkerNum:       config.WorkerNum,
				InputQueueSize:  config.InputQueueSize,
				WorkerQueueSize: config.WorkerQueueSize,
			},
			config.Listener,
		),
		quit: make(chan struct{}),
	}
}

func (j *Jigsaw) InstallEvents(m common.Mod) error {
	descs := m.EventSpecs()
	for _, desc := range descs {
		j.executor.Register(desc.EventPath, desc.Callback)
	}
	return nil
}

func (j *Jigsaw) Serve() {
	j.executor.Serve()
	j.Publish("system.ready", j)
}

func (j *Jigsaw) Quit() {
	j.executor.Quit()
	close(j.quit)
}

func (j *Jigsaw) Wait() {
	<-j.quit
}

func (j *Jigsaw) Publish(subject string, arg interface{}) {
	j.executor.Publish(subject, arg)
}

func (j *Jigsaw) LoadMod(path string) (common.Mod, error) {
	p, err := plugin.Open(path)
	if err != nil {
		return nil, err
	}
	s, err := p.Lookup("New")
	if err != nil {
		return nil, err
	}
	createFunc, ok := s.(func() common.Mod)
	if !ok {
		return nil, fmt.Errorf("invalid function signature %v, expect [func() Mod]", s)
	}
	mod := createFunc()
	err = j.InstallEvents(mod)
	if err != nil {
		return nil, err
	}
	return mod, err
}

func (j *Jigsaw) LoadPlugins(pluginDir string) []string {
	var pluginPaths []string
	filepath.Walk(pluginDir, func(path string, info os.FileInfo, err error) error {
		if strings.HasSuffix(info.Name(), ".so") {
			pluginPaths = append(pluginPaths, path)
		}
		return nil
	})
	for _, path := range pluginPaths {
		mod, err := j.LoadMod(path)
		if err != nil {
			logrus.Warnf("fail to load mod %s: %s", path, err)
			continue
		}
		logrus.Infof("plugin %s loaded", mod.Meta().Package)
	}
	return nil
}
