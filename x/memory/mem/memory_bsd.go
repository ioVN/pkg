//go:build freebsd || openbsd || dragonfly || netbsd

// WARNING: THIS PACKAGE IS EXPERIMENTAL AND MAY BE MODIFIED OR REMOVED WITHOUT
// NOTICE! USE WITH EXTREME CAUTION!
package mem

func sysTotalMemory() uint64 {
	s, err := sysctlUint64("hw.physmem")
	if err != nil {
		return 0
	}
	return s
}

func sysFreeMemory() uint64 {
	s, err := sysctlUint64("hw.usermem")
	if err != nil {
		return 0
	}
	return s
}
