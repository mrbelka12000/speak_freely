package ai

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient_GenerateTopics(t *testing.T) {
	t.Skip("dev purpose only")

	resp, err := testClient.GenerateTheme(context.Background(), GenerateThemeRequest{
		Level: "B2",
	})
	assert.NoError(t, err)
	assert.NotNil(t, resp)
}
