package main_test

import (
	"context"
	"testing"
	"time"

	"github.com/ray1422/dcard-backend-2023/models"
	"github.com/ray1422/dcard-backend-2023/utils/db"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TestTTL is just a simple test that creates a document in list_node
func TestTTL(t *testing.T) {
	refID, err := primitive.ObjectIDFromHex("5c7452c7aeb4c97e0cdb75bf")
	assert.Nil(t, err)
	now := time.Now().Add(time.Second * 10)
	_, err = db.MongoDB().Collection("list_node").InsertOne(context.TODO(), models.ListNode{
		ID:       primitive.NewObjectID(),
		RefID:    refID,
		ExpireOn: &now,
	})
	assert.Nil(t, err)

}
