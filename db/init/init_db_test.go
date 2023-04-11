package init_db

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	assert.Equal(t, len(tokens), 5)
	assert.Equal(t, len(amms), 3)
	assert.Greater(t, len(pools), 0)
}
