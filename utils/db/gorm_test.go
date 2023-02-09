package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDBConnection(t *testing.T) {
	db := GormDB()
	assert.NotNil(t, db)
}
