package testing

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tehbooom/go-kibana/kbapi"
)

func TestSpacesGet(t *testing.T) {
	client := NewTestClient(t)

	t.Run("GetSpaces", func(t *testing.T) {

		req := &kbapi.SpacesGetRequest{ID: "default"}

		resp, err := client.KibanaClient.Spaces.Get(client.Context, req)

		require.NoError(t, err)
		assert.NotEmpty(t, resp.Body)
		assert.IsType(t, &kbapi.Space{}, resp.Body)
	})
}
