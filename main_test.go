package main

import (
	"log"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestDB(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	saveDatas := make([]SaveData, 0)
	err = db.Model(&SaveData{}).Where("reason = ?", "时间范围内近似发言过多").Order("created_at desc").Limit(100).Find(&saveDatas).Error
	if err != nil {
		log.Fatal(err)
	}
	log.Println(len(saveDatas))
}
