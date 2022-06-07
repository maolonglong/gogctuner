//go:build !linux
// +build !linux

package runtime

func MemoryLimit() (int64, bool, error) {
	return -1, false, nil
}
