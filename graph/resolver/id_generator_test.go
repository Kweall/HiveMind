package resolver_test

import (
	"hivemind/graph/resolver"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateID(t *testing.T) {
	t.Run("ID length", func(t *testing.T) {
		id := resolver.GenerateID()

		assert.Len(t, id, 16, "Generated ID should be 16 characters long")
	})

	t.Run("ID uniqueness", func(t *testing.T) {
		id1 := resolver.GenerateID()
		id2 := resolver.GenerateID()

		assert.NotEqual(t, id1, id2, "Generated IDs should be unique")
	})

	t.Run("No errors during generation", func(t *testing.T) {
		id := resolver.GenerateID()

		assert.NotEmpty(t, id, "Generated ID should not be empty")
	})
}
