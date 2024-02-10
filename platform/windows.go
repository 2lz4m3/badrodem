//go:build windows

package platform

import (
	"syscall"
	"unsafe"
)

// from https://gist.github.com/yougg/213250cc04a52e2b853590b06f49d865
func IsDoubleClickRun() bool {
	kernel32 := syscall.NewLazyDLL("kernel32.dll")
	lp := kernel32.NewProc("GetConsoleProcessList")
	if lp != nil {
		var pids [2]uint32
		var maxCount uint32 = 2
		ret, _, _ := lp.Call(uintptr(unsafe.Pointer(&pids)), uintptr(maxCount))
		if ret > 1 {
			return false
		}
	}
	return true
}
