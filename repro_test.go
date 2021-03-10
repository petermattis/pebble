package pebble

import (
	"fmt"
	"runtime"
	"testing"

	"github.com/cockroachdb/pebble/vfs"
	"github.com/stretchr/testify/require"
)

func TestWithFinalizer(t *testing.T) {
	testFoo(t, true)
}

func TestWithoutFinalizer(t *testing.T) {
	testFoo(t, false)
}

func testFoo(t *testing.T, finalizer bool) {
	for i := 0; i < 1000; i++ {
		db, err := Open("", &Options{FS: vfs.NewMem()})
		require.NoError(t, err)

		if finalizer {
			runtime.SetFinalizer(db, func(obj interface{}) {
				fmt.Printf("%p run finalizer\n", obj)
			})
		}

		require.NoError(t, db.Close())
		db = nil

		runtime.GC()

		var memStats runtime.MemStats
		runtime.ReadMemStats(&memStats)
		fmt.Printf("heap objects: %d\n", memStats.HeapObjects)
	}
}
