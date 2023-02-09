package model

import (
	"time"

	"github.com/pilagod/gorm-cursor-paginator/v2/cursor"
	"github.com/ray1422/dcard-backend-2023/utils"
	"gorm.io/gorm"
)

// Serializable Serializable
type Serializable[t any] interface {
	Serialize() t
}

// CursorPagedSerializer is a generic serializer for cursor paged lis
type CursorPagedSerializer struct {
	Items []any  `json:"items"`
	Next  string `json:"next,omitempty"`
	Prev  string `json:"prev,omitempty"`
}

// SerializeCursorPagedItems serializes items.
func SerializeCursorPagedItems[t any, u Serializable[t]](items []u, cursor cursor.Cursor) (ret CursorPagedSerializer) {
	ret.Items = utils.Map(func(v u) any {
		return v.Serialize()
	}, items)

	if cursor.After != nil {
		ret.Next = *cursor.After
	}

	if cursor.Before != nil {
		ret.Prev = *cursor.Before
	}
	return
}

// List points to the start node to the ListNode
type List struct {
	ID      uint   `gorm:"primarykey" json:"id"`
	Key     string `json:"key" gorm:"index"`
	Version uint32
}

// ListNode is the base class expected to be extended by other model
type ListNode struct {
	ID        uint      `gorm:"primarykey"`
	Version   uint32    `gorm:"index"`
	ListID    uint32    `gorm:"index"`
	NodeOrder uint32    `gorm:"index"`
	CreatedAt time.Time `gorm:"index"`
	ArticleID int
	Article   Article `gorm:"foreignKey:ArticleID"`
}

// Serialize ListNode to ArticleSerializer
func (n ListNode) Serialize() ArticleSerializer {
	return n.Article.Serialize()
}

// check ListNode impl the interface
var _ Serializable[ArticleSerializer] = ListNode{}

// ListNodesSerializer is the serializer of ListNode
type ListNodesSerializer struct {
	Article []ArticleSerializer `json:"articles,omitempty"`
	Next    string              `json:"next"`
}

// Article is a sample model for the list which contains title and content.
type Article struct {
	gorm.Model
	Title   string `json:"title"`
	Content string `json:"content"`
}

// ArticleSerializer is the serializer of article
type ArticleSerializer struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

// Serialize article
func (article Article) Serialize() ArticleSerializer {
	return ArticleSerializer{
		ID:      int(article.ID),
		Title:   article.Title,
		Content: article.Content,
	}
}
