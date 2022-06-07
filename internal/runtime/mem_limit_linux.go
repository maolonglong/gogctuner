//go:build linux
// +build linux

package runtime

import (
	"errors"

	cg "go.chensl.me/gogctuner/internal/cgroups"
)

func MemoryLimit() (int64, bool, error) {
	cgroups, err := newQueryer()
	if err != nil {
		return -1, false, err
	}

	limit, defined, err := cgroups.MemoryLimit()
	if !defined || err != nil {
		return -1, false, err
	}

	return limit, true, nil
}

type queryer interface {
	MemoryLimit() (int64, bool, error)
}

var (
	_newCgroups2 = cg.NewCGroups2ForCurrentProcess
	_newCgroups  = cg.NewCGroupsForCurrentProcess
)

func newQueryer() (queryer, error) {
	cgroups, err := _newCgroups2()
	if err == nil {
		return cgroups, nil
	}
	if errors.Is(err, cg.ErrNotV2) {
		return _newCgroups()
	}
	return nil, err
}
