package errorsutil_test

import (
	"fmt"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"

	errorsutil "github.com/lolopinto/go-errors-util"
)

func TestSyncError(t *testing.T) {
	var serr errorsutil.SyncError

	var wg sync.WaitGroup
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go func(i int) {
			defer wg.Done()
			if i%2 == 0 {
				serr.Append(nil)
			} else {
				serr.Append(fmt.Errorf("err %d", i))
			}
		}(i)
	}
	wg.Wait()

	err := serr.Err()
	require.NotNil(t, err)

	_, ok := err.(*errorsutil.ErrorList)
	require.True(t, ok)

	parts := strings.Split(err.Error(), "\n")
	require.Len(t, parts, 5)
}
