package tests_test

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestServer(t *testing.T) {
	t.Run("it responds to requests", func(t *testing.T) {
		resp, err := http.Get("http://server:12345")
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, resp.StatusCode)
	})
}
