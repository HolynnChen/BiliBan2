package main

import (
	"log"
	"time"

	"github.com/Holynnchen/BiliBan2/DanmuCenter"
	"github.com/Holynnchen/BiliBan2/DanmuCenter/Filter"
	"gorm.io/gorm"
)

type CustomBanProcess struct {
	localFilter  *Filter.BanWindowFilter
	systemFilter *Filter.BanWindowFilter
	db           *gorm.DB
	reporter     func(*DanmuCenter.BanData)

	nowCursorID int64
}
type SaveData struct {
	ID        uint                `gorm:"primaryKey"`
	Data      DanmuCenter.BanData `gorm:"embedded"`
	CreatedAt time.Time
}

type SystemBanData struct {
	ID        uint   `gorm:"primaryKey"`
	Content   string `json:"content"`
	CreatedAt time.Time
}

func (process *CustomBanProcess) Ban(banData *DanmuCenter.BanData) {
	log.Printf("%+v\n", banData)
	process.localFilter.Add(banData.Content)
	//异步掉耗时操作
	go process.db.Create(&SaveData{
		Data: *banData,
	})
	go syncBan(banData)
	//go process.reporter(banData)
}

func (process *CustomBanProcess) RestoreLocalFilter(limit int) {
	saveDatas := make([]SaveData, 0)
	err := process.db.Model(&SaveData{}).Where("reason = ?", "时间范围内近似发言过多").Order("created_at desc").Limit(limit).Find(&saveDatas).Error
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("尝试导入%d条频繁发言封禁记录到本地封禁窗口\n", len(saveDatas))
	for _, data := range saveDatas {
		process.localFilter.Add(data.Data.Content)
	}
}

func (process *CustomBanProcess) RestoreSystemFilter(limit int) {
	systemBanDatas := make([]SystemBanData, 0)
	err := process.db.Model(&SystemBanData{}).Order("created_at desc").Limit(limit).Find(&systemBanDatas).Error
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("尝试导入%d条系统确认封禁记录到系统封禁窗口\n", len(systemBanDatas))
	for _, data := range systemBanDatas {
		process.systemFilter.Add(data.Content)
	}
}

func (process *CustomBanProcess) UpdateSystemFilter() error {
	queryData, err := queryBan(process.nowCursorID)
	if err != nil {
		log.Println(err)
		return err
	}
	if len(queryData) == 0 {
		return nil
	}
	count := 0
	for i := 0; i < len(queryData); i++ {
		if queryData[i].Reason == "垃圾广告" && queryData[i].Danmaku.Comment != "" {
			if success := process.systemFilter.Add(queryData[i].Danmaku.Comment); success {
				count++
			}
		}
	}
	if count > 0 {
		log.Printf("同步系统封禁新增规则%d条\n", count)
	}
	process.nowCursorID = queryData[0].CursorID
	return nil
}

func (process *CustomBanProcess) TimingUpdataSystemFilter(d time.Duration) {
	ticker := time.NewTicker(d)
	for {
		select {
		case <-ticker.C:
			go process.UpdateSystemFilter()
		}
	}
}
