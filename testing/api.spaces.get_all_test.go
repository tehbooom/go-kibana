package testing

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tehbooom/go-kibana/kbapi"
)

func TestSpacesGetAll(t *testing.T) {
	client := NewTestClient(t)

	t.Run("GetAllSpaces", func(t *testing.T) {

		req := &kbapi.SpacesGetAllRequest{}

		resp, err := client.KibanaClient.Spaces.GetAll(client.Context, req)

		require.NoError(t, err)
		assert.NotEmpty(t, resp.Body)
		assert.IsType(t, &kbapi.SpacesGetAllResponseBody{}, resp.Body)
	})
}
