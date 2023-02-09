package utils_test

import (
	"testing"

	"github.com/ray1422/dcard-backend-2023/utils"
	"github.com/stretchr/testify/assert"
)

func TestMap(t *testing.T) {
	u := []uint{1, 2, 3}
	v := utils.Map(func(val uint) int {
		return -int(val)
	}, u)
	assert.Equal(t, -int(u[0]), v[0])
	assert.Equal(t, -int(u[1]), v[1])
	assert.Equal(t, -int(u[2]), v[2])

	uBak := make([]uint, len(v))
	assert.Equal(t, copy(uBak, u), len(u))

	utils.MapInPlace(func(t *uint) {
		*t = *t << 1
	}, u)
	assert.Equal(t, uBak[0]*2, u[0])
	assert.Equal(t, uBak[1]*2, u[1])
	assert.Equal(t, uBak[2]*2, u[2])
}
