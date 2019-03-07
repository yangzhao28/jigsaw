package jigsaw

import (
	"errors"
	"fmt"

	"github.com/yangzhao28/jigsaw/src/common"
)

type ModType int

// Mod types
const (
	ModTypePlaceHolder ModType = iota
	ModTypeBuiltin
	ModTypePlugin
)

type mNode struct {
	// deps other
	Deps []common.Mod
	// deps by other
	DepsBy []common.Mod
	Type   ModType
	Mod    common.Mod
}

// ModManager manage mod loading / deps parsing
type ModManager struct {
	// preload map[string]common.Mod
	refs map[string]*mNode
}

func NewModManager() *ModManager {
	return &ModManager{
		// preload: make(map[string]common.Mod),
		refs: make(map[string]*mNode),
	}
}

func (mm *ModManager) testVersion(m common.Mod, dep common.ModDepDescriptor) bool {
	v := m.Meta().Version
	return v.GreaterEqual(dep.MinVersion) && dep.MaxVersion.GreaterEqual(v)
}

func (mm *ModManager) linkDeps(m common.Mod, deps []common.ModDepDescriptor) {
	for _, dep := range deps {
		r, ok := mm.refs[dep.PackageName]
		if !ok {
			r = &mNode{
				Type: ModTypePlaceHolder,
			}
			mm.refs[dep.PackageName] = r
		}
		r.DepsBy = append(r.DepsBy, m)
	}
}

func (mm *ModManager) Register(m common.Mod) error {
	// validation
	if m == nil {
		return errors.New("invalid mod")
	}
	for _, dep := range m.Deps() {
		if dep.PackageName == m.Meta().Package {
			return errors.New("cycle deps")
		}
	}
	// loop deps
	pName := m.Meta().Package
	r, ok := mm.refs[pName]
	if !ok {
		r = &mNode{Mod: m}
		mm.refs[pName] = r
	}
	if r.Type != ModTypePlaceHolder {
		return errors.New("duplicated mod")
	}
	// update mod
	r.Mod = m
	r.Type = ModTypePlugin
	// update deps, bidirectional link
	mm.linkDeps(m, m.Deps())
	return nil
}

func (mm *ModManager) Load(m common.Mod) error {
	err := mm.parseDeps()
	if err != nil {
		return err
	}
	return nil
}

func (mm *ModManager) getDepDiscriptor(m common.Mod, pName string) common.ModDepDescriptor {
	for _, dep := range m.Deps() {
		if dep.PackageName == pName {
			return dep
		}
	}
	panic("should not reach here")
}

func (mm *ModManager) parseDeps() error {
	for pName, ref := range mm.refs {
		if ref.Type == ModTypePlaceHolder {
			return fmt.Errorf("missing mod %s", pName)
		}
		// check deps version requirements
		for _, depBy := range ref.DepsBy {
			if !mm.testVersion(ref.Mod, mm.getDepDiscriptor(depBy, pName)) {
				return fmt.Errorf("mode %s version not meets requirement by %s", pName, depBy.Meta().Package)
			}
		}
	}
	return nil
}
