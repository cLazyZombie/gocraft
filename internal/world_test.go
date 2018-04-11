package internal_test

import (
	"testing"

	"github.com/cLazyZombie/gocraft/gocrafttest"
	. "github.com/cLazyZombie/gocraft/internal"
	"github.com/stretchr/testify/assert"
)

func TestWorld_Init(t *testing.T) {
	assert.True(t, true, "true")
	store := &gocrafttest.StoreMock{}
	world := NewWorld(store)

	assert.NotNil(t, world, "not nil")
}
