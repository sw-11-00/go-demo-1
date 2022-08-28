package database

import (
	"flag"
	"fmt"
	"go-demo-1/config"
	"testing"
	"time"
)

type Chains struct {
	ID          int64
	BlockHeight int       `gorm:"column:block_height"`
	BlockHash   string    `gorm:"column:block_hash"`
	CreatedAt   time.Time `gorm:"column:created_at"`
	UpdatedAt   time.Time `gorm:"column:updated_at"`
}

var (
	configFlag = flag.String("config", "/Users/jiansui/Downloads/work/png/go/go-demo-1/config.yml", "Config file")
)

func TestGorm(t *testing.T) {
	flag.Parse()
	cfg := config.LoadConfig(*configFlag)

	db := Open(cfg.Database.Username, cfg.Database.Password, cfg.Database.Host, cfg.Database.Port, cfg.Database.Database)
	u := Chains{
		BlockHeight: 1,
		BlockHash:   "000000",
		CreatedAt:   time.Now().Local(),
		UpdatedAt:   time.Now().UTC(),
	}

	if err := db.Create(&u).Error; err != nil {
		fmt.Println("插入失败", err)
		return
	}
}
