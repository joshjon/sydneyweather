package api

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestValueCache(t *testing.T) {
	expiry := 100 * time.Millisecond
	c := newValueCache[string](expiry)
	wantVal := "some-val"
	c.put(&wantVal)

	// Present
	require.False(t, c.expired())
	gotVal, ok := c.get()
	require.True(t, ok)
	require.Equal(t, wantVal, *gotVal)

	// Expired
	time.Sleep(expiry + time.Millisecond)
	require.True(t, c.expired())

	// Stale value should still be present
	gotVal, ok = c.get()
	require.True(t, ok)
	require.Equal(t, wantVal, *gotVal)
}
