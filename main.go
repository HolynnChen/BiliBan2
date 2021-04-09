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
	db.AutoMigrate(&SaveData{})

	banWindowFilter := Filter.NewBanWindowFilter(100, 3600, 0.75) //创建容量为100，窗口有效时间为3600秒，相似度要求为0.75的封禁窗口
	banProcess := &CustomBanProcess{
		db:       db,
		filter:   banWindowFilter,
		reporter: Helper.Reporter(env["cookie"].(string), env["csrf"].(string)),
	} //创建自定义封禁处理结构体
	banProcess.Restore(100) //从数据库恢复最多100条因频繁发言封禁的记录导入到窗口

	center := DanmuCenter.NewDanmuCenter(&DanmuCenter.DanmuCenterConfig{
		TimeRange:      16,
		MonitorNumber:  100,
		SpecialFocusOn: []int{1370218}, //1237390
		Silent:         true,
	},
		DanmuCenter.SetSaveFilter( //是否入库检测
			Filter.NewLenFilter(8),        // 简易长度过滤
			banWindowFilter,               // 移除高等级的窗口
			Filter.NewUserLevelFilter(5),  // 过滤掉用户等级>=5的
			Filter.NewFansMedalFilter(3),  // 过滤掉粉丝勋章等级>=3的
			Filter.NewHaveBeenBanFilter(), // 过滤掉已被Ban的弹幕
			Filter.NewKeyWordFilter([]string{"谢谢", "感谢", "多谢"}),               // 关键词匹配过滤
			Filter.NewLenFilter(9, Filter.SetLenFilterCompressRepeatGroup(3)), // 过滤掉重复词压缩后长度小于9的弹幕
		),
		DanmuCenter.SetSafeFilter( //是否正常弹幕
			Filter.NewHighReatWordFilter(0.75), //单字符重复率>0.75视作正常弹幕
		),
		DanmuCenter.SetBanFilter( //是否异常弹幕
			banWindowFilter, //与封禁窗口比较
			Filter.NewHighSimilarityAndSpeedFilter(0.75, 3), //时间范围内达到startCheck后检测最新几组的相似率
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
			banWindowFilter.Ban(&DanmuCenter.BanData{Content: bodyStr})
			io.WriteString(w, "add string ["+bodyStr+"] success")
		})
		log.Println(http.ListenAndServe("0.0.0.0:6060", nil))
	}()

	center.Start()
}

type CustomBanProcess struct {
	filter   DanmuCenter.BanProcess
	db       *gorm.DB
	reporter func(*DanmuCenter.BanData)
}

type SaveData struct {
	ID        uint                `gorm:"primaryKey"`
	Data      DanmuCenter.BanData `gorm:"embedded"`
	CreatedAt time.Time
}

func (process *CustomBanProcess) Ban(banData *DanmuCenter.BanData) {
	log.Printf("%+v\n", banData)
	process.filter.Ban(banData)
	//异步掉耗时操作
	go process.db.Create(&SaveData{
		Data: *banData,
	})
	go syncBan(banData)
	go process.reporter(banData)
}

func (process *CustomBanProcess) Restore(limit int) {
	saveDatas := make([]SaveData, 0)
	err := process.db.Model(&SaveData{}).Where("reason = ?", "时间范围内近似发言过多").Order("created_at desc").Limit(limit).Find(&saveDatas).Error
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("尝试导入%d条频繁发言封禁记录到封禁窗口", len(saveDatas))
	for _, data := range saveDatas {
		process.filter.Ban(&data.Data)
	}
}

type SyncData struct {
	UserID    int64  `json:"UserId"`
	UserName  string `json:"UserName"`
	RoomID    int    `json:"RoomId"`
	Content   string `json:"Content"`
	TimeStamp int64  `json:"TimeStamp"`
	Reason    string `json:"Reason"`
}

func syncBan(banData *DanmuCenter.BanData) {
	jsonData, _ := json.Marshal(SyncData{
		UserID:    banData.UserID,
		UserName:  banData.UserName,
		RoomID:    banData.RoomID,
		Content:   banData.Content,
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
