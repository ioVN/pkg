// WARNING: THIS PACKAGE IS EXPERIMENTAL AND MAY BE MODIFIED OR REMOVED WITHOUT
// NOTICE! USE WITH EXTREME CAUTION!
package mem

import "testing"

func TestNonZero(t *testing.T) {
	if TotalMemory() == 0 {
		t.Fatal("TotalMemory returned 0")
	}
	if FreeMemory() == 0 {
		t.Fatal("FreeMemory returned 0")
	}
}
