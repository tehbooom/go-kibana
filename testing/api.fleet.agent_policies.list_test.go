package testing

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tehbooom/go-kibana/kbapi"
)

func TestAgentPoliciesList(t *testing.T) {
	client := NewTestClient(t)

	t.Run("ListAgentPoliciesSuccess", func(t *testing.T) {

		req := &kbapi.FleetAgentPoliciesRequest{}

		resp, err := client.KibanaClient.Fleet.AgentPolicies.List(client.Context, req)

		require.NoError(t, err)
		assert.NotEmpty(t, resp.Body)
		assert.IsType(t, &kbapi.GetFleetAgentPoliciesResponseBody{}, resp.Body)

		foundPolicies := make(map[string]bool)
		for _, policy := range resp.Body.Items {
			foundPolicies[policy.Id] = true
		}

		t.Log(foundPolicies)
	})
	// Write more tests with pagination, kquery, and getting Full back
}
