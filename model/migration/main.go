package main

import (
	"fmt"
	"os"
	"time"

	"github.com/ray1422/dcard-backend-2023/utils/db"
	"gorm.io/gorm"
)

// Migrate Migrate
func Migrate() (err error) {
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
		ArticleID uint
	}

	// Article is a sample model for the list which contains title and content.
	type Article struct {
		gorm.Model
		Title   string `json:"title"`
		Content string `json:"content"`
	}

	err = db.GormDB().AutoMigrate(&Article{})
	if err != nil {
		return
	}
	err = db.GormDB().AutoMigrate(&ListNode{})
	if err != nil {
		return
	}

	err = db.GormDB().AutoMigrate(&List{})
	if err != nil {
		return
	}
	return
}

// Rollback Rollback
func Rollback() (err error) {
	err = db.GormDB().Migrator().DropTable("articles")
	if err != nil {
		return
	}
	err = db.GormDB().Migrator().DropTable("list_nodes")
	if err != nil {
		return
	}
	err = db.GormDB().Migrator().DropTable("lists")
	if err != nil {
		return
	}

	return
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "usage: %s <m|r>\n", os.Args[0])
		os.Exit(2)
		return
	}
	var err error = nil
	if os.Args[1] != "r" {
		err = Migrate()
	} else {
		err = Rollback()
	}
	if err != nil {
		fmt.Println("error occurred:\n", err)
	} else {
		fmt.Println("operation successfully completed.")
	}
}
