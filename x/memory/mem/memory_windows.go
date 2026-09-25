//go:build windows

// WARNING: THIS PACKAGE IS EXPERIMENTAL AND MAY BE MODIFIED OR REMOVED WITHOUT
// NOTICE! USE WITH EXTREME CAUTION!
package mem

import (
	"syscall"
	"unsafe"
)

// omitting a few fields for brevity...
// https://msdn.microsoft.com/en-us/library/windows/desktop/aa366589(v=vs.85).aspx
type memStatusEx struct {
	dwLength     uint32
	dwMemoryLoad uint32
	ullTotalPhys uint64
	ullAvailPhys uint64
	unused       [5]uint64
}

func loadKernel32() *memStatusEx {
	kernel32, err := syscall.LoadDLL("kernel32.dll")
	if err != nil {
		return nil
	}
	// GetPhysicallyInstalledSystemMemory is simpler, but broken on
	// older versions of windows (and uses this under the hood anyway).
	globalMemoryStatusEx, err := kernel32.FindProc("GlobalMemoryStatusEx")
	if err != nil {
		return nil
	}
	msx := &memStatusEx{
		dwLength: 64,
	}
	r, _, _ := globalMemoryStatusEx.Call(uintptr(unsafe.Pointer(msx)))
	if r == 0 {
		return nil
	}
	return msx
}

func sysTotalMemory() uint64 {
	msx := loadKernel32()
	if msx != nil {
		return msx.ullTotalPhys
	}
	return 0
}

func sysFreeMemory() uint64 {
	msx := loadKernel32()
	if msx != nil {
		return msx.ullAvailPhys
	}
	return 0
}
