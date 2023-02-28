package model

import (
	"errors"
	"log"
	"strings"
	"time"

	"github.com/pilagod/gorm-cursor-paginator/v2/cursor"
	"github.com/pilagod/gorm-cursor-paginator/v2/paginator"
	"github.com/ray1422/dcard-backend-2023/utils"
	"github.com/ray1422/dcard-backend-2023/utils/db"
	"gorm.io/gorm"
)

// ErrDuplicateKeyError DuplicateKeyError
var ErrDuplicateKeyError = errors.New("duplicate key")

// Serializable Serializable
type Serializable[t any] interface {
	Serialize() t
}

// CursorPagedSerializer is a generic serializer for cursor paged lis
type CursorPagedSerializer[t any] struct {
	Items []t    `json:"items"`
	Next  string `json:"next,omitempty"`
	Prev  string `json:"prev,omitempty"`
}

// SerializeCursorPagedItems serializes items.
func SerializeCursorPagedItems[t any, u Serializable[t]](items []u, cursor cursor.Cursor) (ret CursorPagedSerializer[t]) {
	ret.Items = utils.Map(func(v u) t {
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
	Version uint32 `json:"version"`
}

// ListSerializer ListSerializer
type ListSerializer List

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

// NewList creates list
func NewList(key string) (uint, error) {
	list := List{
		Key: key,
	}
	err := db.GormDB().Create(&list).Error
	// pgErr := &pgconn.PgError{}
	if err != nil {
		log.Println("failed to create list:", err)
		// not sure what's wrong with pgconn.PgError{}.
		// It seems GORM uses an incompatible version of pgconn that I can't cast the error and catch the error code.
		// temporary catch by error message string lol
		if strings.HasPrefix(err.Error(), "ERROR: duplicate key") {
			return 0, ErrDuplicateKeyError
		}
		return 0, errors.New("internal error")
	}
	return list.ID, nil
}

// Serialize list
// List is self-serializable, which needn't additional serializer.
// This function is for the compatibility of `Serializable` and consistency for the whole program.
func (list List) Serialize() ListSerializer {
	return ListSerializer(list)
}

// DeleteList deletes list by ID
func DeleteList(id uint) error {
	db := db.GormDB().Delete(&List{
		ID: id,
	})
	if db.Error != nil {
		return db.Error
	}
	if db.RowsAffected < 1 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// SetListVersion sets list's version
func SetListVersion(listID uint, listVersion uint32) error {
	db := db.GormDB().Model(&List{}).Where("id", listID).Update("version", listVersion)
	if db.Error != nil {
		return db.Error
	}
	if db.RowsAffected < 1 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// FindListByKey Find List By Key
func FindListByKey(key string) (*List, error) {
	list := List{}
	err := db.GormDB().Take(&list, "key = ?", key).Error
	if err != nil {
		return nil, err
	}
	return &list, nil
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

// GetListNodes nodes of the list
func GetListNodes(listID uint, version uint, nextCursor string) ([]ListNode, cursor.Cursor, error) {
	nodes := []ListNode{}
	cursor := cursor.Cursor{}
	if nextCursor != "" {
		cursor.After = &nextCursor
	}
	pg := createListNodePaginator(cursor, paginator.ASC, nil)
	// query := `
	// 	SELECT articles.id, list_nodes.node_order
	// 	FROM list_nodes
	// 	INNER JOIN articles ON (articles.id = list_nodes.article_id)
	// 	WHERE list_nodes.list_id = ? AND list_nodes.version = ?
	// `
	stmt := db.GormDB().
		Select([]string{"list_nodes.node_order", "list_nodes.id"}).
		InnerJoins("Article", db.GormDB().Select([]string{"id", "title", "content"})).
		Where("list_id = ?", listID).Where("version = ?", version)
	result, cursor, err := pg.Paginate(stmt, &nodes)
	if err != nil {
		return nil, cursor, err
	}
	if result.Error != nil {
		return nil, cursor, result.Error
	}

	return nodes, cursor, nil

}

func createListNodePaginator(
	cursor paginator.Cursor,
	order paginator.Order,
	limit *int,
) *paginator.Paginator {
	opts := []paginator.Option{
		&paginator.Config{
			// Keys: []string{"ID"},
			Keys:  []string{"NodeOrder", "ID"},
			Limit: 10,
			Order: order,
		},
	}
	if limit != nil {
		opts = append(opts, paginator.WithLimit(*limit))
	}
	if cursor.After != nil {
		opts = append(opts, paginator.WithAfter(*cursor.After))
	}
	if cursor.Before != nil {
		opts = append(opts, paginator.WithBefore(*cursor.Before))
	}
	return paginator.New(opts...)
}

// InsertNodes inserts nodes.
func InsertNodes(nodes []ListNode) error {
	return db.GormDB().Create(nodes).Error
}
