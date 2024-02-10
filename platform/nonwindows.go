//go:build !windows

package platform

func IsDoubleClickRun() bool {
	return false
}
