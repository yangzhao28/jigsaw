package main

import (
	"github.com/yangzhao28/jigsaw/src/jigsaw"
	"github.com/yangzhao28/jigsaw/src/modules/system"
)

func main() {
	j := jigsaw.New(jigsaw.DefaultConfig())
	j.InstallEvents(system.New())
	j.LoadPlugins("./plugin")
	j.Serve()
	j.Wait()
}
