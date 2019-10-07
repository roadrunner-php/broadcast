package broadcast

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConnContext_ManageTopics(t *testing.T) {
	ctx := &ConnContext{Topics: make([]string, 0)}

	assert.Equal(t, []string{}, ctx.Topics)

	ctx.addTopic("a", "b")
	assert.Equal(t, []string{"a", "b"}, ctx.Topics)

	ctx.addTopic("a", "c")
	assert.Equal(t, []string{"a", "b", "c"}, ctx.Topics)

	ctx.dropTopic("b", "c")
	assert.Equal(t, []string{"a"}, ctx.Topics)

	ctx.dropTopic("b", "c")
	assert.Equal(t, []string{"a"}, ctx.Topics)

	ctx.dropTopic("a")
	assert.Equal(t, []string{}, ctx.Topics)
}
