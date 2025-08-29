//go:build !linux

package filesystememulation

import "testing"

func PickDifferentFSDir(t *testing.T, _ string) string {
	t.Skip("skipping filesystem test on non-linux systems")
	panic("unreachable")
}
