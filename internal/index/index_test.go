package index

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidVersion(t *testing.T) {
	assert.True(t, isValidVersion("1.5"))
	assert.True(t, isValidVersion("1.6.4"))
	assert.True(t, isValidVersion("1.7.10"))
	assert.True(t, isValidVersion("1.12.2"))
	assert.True(t, isValidVersion("1.19.2"))
	assert.True(t, isValidVersion("1.20"))
	assert.True(t, isValidVersion("1.20.4"))
}
