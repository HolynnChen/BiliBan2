package main

import (
	"log"
	"time"

	"github.com/Holynnchen/BiliBan2/DanmuCenter"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	//性能调优
	"net/http"
	_ "net/http/pprof"
)

func main() {
	//性能调优
	go func() {
		log.Println(http.ListenAndServe("0.0.0.0:6060", nil))
	}()

	db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	db.AutoMigrate(&SaveData{})

	banWindowFilter := DanmuCenter.NewBanWindowFilter(100, 3600, 0.75) //创建容量为100，窗口有效时间为3600秒，相似度要求为0.75的封禁窗口
	banProcess := &CustomBanProcess{db: db, filter: banWindowFilter}   //创建自定义封禁处理结构体
	banProcess.Restore(100)                                            //从数据库恢复最多100条因频繁发言封禁的记录导入到窗口

	center := DanmuCenter.NewDanmuCenter(&DanmuCenter.DanmuCenterConfig{
		TimeRange:      15,
		MonitorNumber:  50,
		SpecialFocusOn: []int{1370218}, //1237390
		Silent:         true,
	},
		DanmuCenter.SetSaveFilter( //是否入库检测
			DanmuCenter.NewUserLevelFilter(5),                                            // 过滤掉用户等级>=15的
			DanmuCenter.NewFansMedalFilter(3),                                            // 过滤掉粉丝勋章等级>=10的
			DanmuCenter.NewKeyWordFilter([]string{"谢谢", "感谢", "多谢"}),                     // 关键词匹配过滤
			DanmuCenter.NewHaveBeenBanFilter(),                                           //过滤掉已被Ban的弹幕
			DanmuCenter.NewLenFilter(10, DanmuCenter.SetLenFilterCompressRepeatGroup(3)), //过滤掉重复词压缩后长度小于9的弹幕
		),
		DanmuCenter.SetSafeFilter( //是否正常弹幕
			DanmuCenter.NewHighReatWordFilter(0.75), //单字符重复率>0.75视作正常弹幕
		),
		DanmuCenter.SetBanFilter( //是否异常弹幕
			banWindowFilter, //与封禁窗口比较
			DanmuCenter.NewHighSimilarityAndSpeedFilter(0.75, 3), //时间范围内达到startCheck后检测最新几组的相似率
		),
		DanmuCenter.SetBanProcess(banProcess), //处理封禁情况
	)
	center.Start()
}

type CustomBanProcess struct {
	filter DanmuCenter.BanProcess
	db     *gorm.DB
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
