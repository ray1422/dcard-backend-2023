package controllers

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/ray1422/dcard-backend-2023/models"
	"github.com/ray1422/dcard-backend-2023/utils/db"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// HeadHandler returns the head of list /list/<key>/<user_id>
func HeadHandler(c *gin.Context) {
	head := models.List{}
	listKey, exists := c.Params.Get("key")
	userIDStr, _ := c.Params.Get("user_id")
	userID, err := strconv.Atoi(userIDStr)
	if err != nil || !exists {
		c.Status(400)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	db := db.MongoDB()
	err = db.Collection("list").FindOne(ctx, bson.M{
		"key":     listKey,
		"user_id": userID,
	}).Decode(&head)
	if err == mongo.ErrNoDocuments {
		c.Status(404)
		return
	} else if err != nil {
		c.Status(500)
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return
	}
	c.JSON(200, head)

}
