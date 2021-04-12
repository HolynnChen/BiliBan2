package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/Holynnchen/BiliBan2/DanmuCenter"
	"github.com/Holynnchen/BiliBan2/DanmuCenter/Filter"
	"github.com/Holynnchen/BiliBan2/DanmuCenter/Helper"
	"github.com/goccy/go-json"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	//性能调优
	"net/http"
)

var env = make(map[string]interface{})

func init() {
	if _, err := os.Stat("env.toml"); err != nil {
		return
	}
	if _, err := toml.DecodeFile("env.toml", &env); err != nil {
		log.Panic(err)
	}
	fmt.Printf("变量值: %+v\n", env)
}

func main() {
	db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	db.AutoMigrate(&SaveData{}, &SystemBanData{})

	localBanWindowFilter := Filter.NewBanWindowFilter(100, 3600, 0.75) //创建容量为100，窗口有效时间为3600秒，相似度要求为0.75的封禁窗口
	systemBanWindowFilter := Filter.NewBanWindowFilter(100, 3600, 0.75)
	banProcess := &CustomBanProcess{
		db:          db,
		localFilter: localBanWindowFilter,
		reporter:    Helper.Reporter(env["cookie"].(string), env["csrf"].(string)),
	} //创建自定义封禁处理结构体
	banProcess.RestoreLocalFilter(100) //从数据库恢复最多100条因频繁发言封禁的记录导入到窗口

	center := DanmuCenter.NewDanmuCenter(&DanmuCenter.DanmuCenterConfig{
		TimeRange:      16,
		MonitorNumber:  1000,
		SpecialFocusOn: []int{1370218}, //1237390
		Silent:         true,
	},
		DanmuCenter.SetPreFilter( //入库前检测
			Helper.Safe(Filter.NewLenFilter(8).Check),         // 简易长度过滤
			Helper.Ban(systemBanWindowFilter.MatchCheck),      // 匹配系统确认封禁记录
			Helper.Continue(localBanWindowFilter.UnlockCheck), // 移除高等级的窗口
			Helper.Safe(Filter.NewUserLevelFilter(5).Check),   // 过滤掉用户等级>=5的
			Helper.Safe(Filter.NewFansMedalFilter(3).Check),   // 过滤掉粉丝勋章等级>=3的
			Helper.Break(Filter.NewHaveBeenBanFilter().Check), // 过滤掉已被Ban的弹幕
			// Filter.NewKeyWordFilter([]string{"谢谢", "感谢", "多谢"}),               // 关键词匹配过滤
			Helper.Safe(Filter.NewLenFilter(9, Filter.SetLenFilterCompressRepeatGroup(3)).Check), // 过滤掉重复词压缩后长度小于9的弹幕
		),
		DanmuCenter.SetAfterFilter( //入库后检测
			Helper.Safe(Filter.NewHighReatWordFilter(0.75).Check),             //单字符重复率>0.75视作正常弹幕
			Helper.Ban(localBanWindowFilter.MatchCheck),                       //与封禁窗口比较
			Helper.Ban(Filter.NewHighSimilarityAndSpeedFilter(0.75, 3).Check), //时间范围内达到startCheck后检测最新几组的相似率
		),
		DanmuCenter.SetBanProcess(banProcess), //处理封禁情况
	)

	go func() {
		//提供导入窗口
		http.HandleFunc("/addWindow", func(w http.ResponseWriter, r *http.Request) {
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				http.Error(w, "can't get body", http.StatusBadRequest)
				return
			}
			bodyStr := string(body)
			log.Println("add Window", bodyStr)
			localBanWindowFilter.Add(bodyStr)
			io.WriteString(w, "add string ["+bodyStr+"] success")
		})
		log.Println(http.ListenAndServe("0.0.0.0:6060", nil))
	}()

	center.Start()
}

type CustomBanProcess struct {
	localFilter  *Filter.BanWindowFilter
	systemFilter *Filter.BanWindowFilter
	db           *gorm.DB
	reporter     func(*DanmuCenter.BanData)
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
		process.localFilter.Add(data.Content)
	}
}

type SyncData struct {
	UserID    int64  `json:"UserId"`
	UserName  string `json:"UserName"`
	RoomID    int    `json:"RoomId"`
	Content   string `json:"Content"`
	CT        string `json:"ct"`
	TimeStamp int64  `json:"TimeStamp"`
	Reason    string `json:"Reason"`
}

func syncBan(banData *DanmuCenter.BanData) {
	jsonData, _ := json.Marshal(SyncData{
		UserID:    banData.UserID,
		UserName:  banData.UserName,
		RoomID:    banData.RoomID,
		Content:   banData.Content,
		CT:        banData.CT,
		TimeStamp: banData.Timestamp,
		Reason:    banData.Reason,
	})
	resp, err := http.DefaultClient.Post(env["sync_url"].(string), "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
}
