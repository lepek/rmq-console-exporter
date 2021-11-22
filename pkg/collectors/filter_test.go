package collectors

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFilterOk(t *testing.T) {
	rules := []string{
		`^.*\.dev$`,
	}
	filter, err := NewFilter(rules)
	assert.Nil(t, err)
	assert.Equal(t, true, filter.Filter("object_test.dev"))
	assert.Equal(t, false, filter.Filter("object_test.dev.something-else"))
}
