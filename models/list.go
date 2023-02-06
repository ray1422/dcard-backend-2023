package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// List points to the start node to the ListNode
type List struct {
	Key    string             `json:"key" bson:"key"`
	UserID uint               `json:"user_id" bson:"user_id"`
	Head   primitive.ObjectID `bson:"head" json:"head"`
}

// ListNode is the base class expected to be extended by other model
type ListNode struct {
	ID       primitive.ObjectID  `bson:"_id,omitempty" json:"_id,omitempty"`
	NextID   *primitive.ObjectID `bson:"next_id"`
	RefID    primitive.ObjectID  `bson:"ref_id"`
	ExpireOn *time.Time          `json:"expire_on" bson:"expire_on"`
}

// Article is a sample model for the list which contains title and content.
type Article struct {
	ID      primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Title   string             `json:"title" bson:"title"`
	Content string             `json:"content" bson:"content"`
}
