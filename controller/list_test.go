package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ray1422/dcard-backend-2023/model"
	"github.com/ray1422/dcard-backend-2023/utils/db"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/test/bufconn"
)

// orderTable is an random table. It's used for make sure records' order is not as same as the order of insertion
var orderTable [20]int
var orderTableReverse [20]int

// this test tests `get_head` and `get_page` APIs.
// write with chatGPT
func init() {
	orderTable = [20]int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19}
	rand.Shuffle(len(orderTable), func(i, j int) { orderTable[i], orderTable[j] = orderTable[j], orderTable[i] })
	for i, v := range orderTable {
		orderTableReverse[v] = i
	}
}

func TestListRESTful(t *testing.T) {
	fmt.Println("order table")
	fmt.Println(orderTable)
	fmt.Println("--")
	myVersion := uint32(time.Now().Unix() & 0xffffffff)
	r := gin.Default()
	lst := bufconn.Listen(0xfffff + 1)
	bufDialer := func(context.Context, string) (net.Conn, error) {
		return lst.Dial()
	}
	_ = bufDialer
	var err error = nil
	go func() {
		err := listenGRPC(lst)
		assert.Nil(t, err)
	}()

	var listID uint = 0

	// creates some data
	{
		listID, err = model.NewList("yet-another-list-test")
		if !assert.Nil(t, err) {
			t.Fatal("failed to create list")
		}
		if !assert.NotZero(t, listID) {
			t.Fatal("listID returned from NewList shouldn't be zero")
		}

		model.SetListVersion(listID, myVersion)
		// creates some article
		articleIDs := []uint{}
		for i := 0; i < len(orderTable); i++ {
			arc := model.Article{
				Title: "Title_" + fmt.Sprint(i),
			}
			assert.Nil(t, db.GormDB().Create(&arc).Error)
			articleIDs = append(articleIDs, arc.ID)
		}
		notHereArticle := model.Article{
			Title: "Not Here!!!!",
		}
		err := db.GormDB().Create(&notHereArticle).Error
		assert.Nil(t, err)
		nodes := []model.ListNode{
			{
				ArticleID: int(notHereArticle.ID),
				NodeOrder: 0,
				Version:   myVersion - 1,
			},
		}

		for i, v := range articleIDs {
			nodes = append(nodes, model.ListNode{
				ListID:    uint32(listID),
				ArticleID: int(v),
				NodeOrder: uint32(orderTable[i]),
				Version:   myVersion,
			})
		}
		assert.Nil(t, model.InsertNodes(nodes))
	}

	defer func() {
		err := model.DeleteList(listID)
		assert.Nil(t, err)
		assert.Nil(t, db.GormDB().Where("created_at <= ?", time.Now()).Delete(&model.ListNode{}).Error)
	}()
	Register(r.Group("/list"), lst)

	// test /list/<key>
	{
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/list/yet-another-list-test", nil)
		r.ServeHTTP(w, req)
		assert.Equal(t, http.StatusOK, w.Code)
		ret := model.List{}
		assert.Nil(t, json.Unmarshal(w.Body.Bytes(), &ret))
		assert.NotZero(t, ret.ID)
		assert.Equal(t, listID, ret.ID)
	}
	// test /list/<id>/<version>
	{
		cursor := ""
		// do-while pattern
		// cnt is the counter for article ID
		cnt := 0
		for next := true; next; next = cursor != "" {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", fmt.Sprintf("/list/%d/%d?cursor=%s", listID, myVersion, cursor), nil)
			r.ServeHTTP(w, req)
			assert.Equal(t, http.StatusOK, w.Code)
			ret := model.CursorPagedSerializer{}
			json.Unmarshal(w.Body.Bytes(), &ret)
			{
				cursor = ret.Next
				for _, item := range ret.Items {
					jsonBody, _ := json.Marshal(item)
					article := &model.ArticleSerializer{}
					json.Unmarshal(jsonBody, &article)
					assert.NotContains(t, article.Title, "not here")
					// fmt.Printf("Title: %s\n", article.Title)
					assert.Equalf(t, fmt.Sprintf("Title_%d", orderTableReverse[cnt]), article.Title, "cnt=%d", cnt)
					cnt++
				}

			}
		}
	}

}
