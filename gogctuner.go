package gogctuner

import (
	"log"
	"os"
	"runtime"
	"runtime/debug"

	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/process"
	iruntime "go.chensl.me/gogctuner/internal/runtime"
)

const (
	_hardTarget = 0.7
	_minGOGC    = 25
	_maxGOGC    = 500
)

var prevgogc = -1

func init() {
	f := &finalizer{}

	f.ref = &finalizerRef{parent: f}
	runtime.SetFinalizer(f.ref, finalizerHandler)
	f.ref = nil
}

type finalizer struct {
	ref *finalizerRef
}

type finalizerRef struct {
	parent *finalizer
}

func finalizerHandler(f *finalizerRef) {
	tuner()
	runtime.SetFinalizer(f, finalizerHandler)
}

func tuner() {
	usedPercent, err := getUsedPercent()
	if err != nil {
		log.Printf("gogctuner: Failed to get used percent: %v", err)
		return
	}

	newgogc := int((_hardTarget - usedPercent) / usedPercent * 100)
	if usedPercent > _hardTarget || newgogc < _minGOGC {
		newgogc = _minGOGC
	}
	if newgogc > _maxGOGC {
		newgogc = _maxGOGC
	}

	if newgogc != prevgogc {
		log.Printf("gogctuner: Updating GOGC=%v", newgogc)
		debug.SetGCPercent(newgogc)
		prevgogc = newgogc
	}
}

func getUsedPercent() (float64, error) {
	v, err := mem.VirtualMemory()
	if err != nil {
		return -1, err
	}
	total := int64(v.Total)

	p, err := process.NewProcess(int32(os.Getpid()))
	if err != nil {
		return -1, err
	}

	mem, err := p.MemoryInfo()
	if err != nil {
		return -1, err
	}
	used := mem.RSS

	limit, defined, _ := iruntime.MemoryLimit()
	if !defined || limit > total {
		limit = total
	}

	return float64(used) / float64(limit), nil
}
