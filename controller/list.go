package controller

import (
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pilagod/gorm-cursor-paginator/v2/cursor"
	"github.com/pilagod/gorm-cursor-paginator/v2/paginator"
	"github.com/ray1422/dcard-backend-2023/model"
	"github.com/ray1422/dcard-backend-2023/utils/db"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// listHandler handles /list/<id>/<version>?next=[cursor]
func listHandler(c *gin.Context) {
	listIDStr := c.Param("id")
	listID, err := strconv.Atoi(listIDStr)
	if err != nil {
		c.Status(400)
		return
	}
	versionStr := c.Param("version")
	version, err := strconv.Atoi(versionStr)
	if err != nil {
		c.Status(400)
		return
	}
	cursorStr := c.Query("next")
	listNodes, cursor, err := list(uint(listID), uint(version), cursorStr)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.Status(404)
			return
		}
		c.Status(500)
		logger.Default.Error(c, "%s", err)

	}
	articles := model.SerializeCursorPagedItems[model.ArticleSerializer](listNodes, cursor)
	c.JSON(200, articles)
}

// List nodes of the list
func list(listID uint, version uint, nextCursor string) ([]model.ListNode, cursor.Cursor, error) {
	nodes := []model.ListNode{}
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
		Select([]string{"list_nodes.node_order"}).
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
