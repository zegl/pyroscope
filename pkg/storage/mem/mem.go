//go:build !linux
// +build !linux

package mem

import (
	"github.com/shirou/gopsutil/mem"
)

// on linux we also look at cgroup mem limit
func Total() (uint64, error) {
	vm, err := mem.VirtualMemory()
	if err != nil {
		return 0, err
	}

	return vm.Total, nil
}
