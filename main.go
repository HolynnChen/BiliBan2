package main

import (
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
	"gorm.io/driver/mysql"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	//性能调优
	"net/http"
)

var env = make(map[string]string)

func init() {
	if _, err := os.Stat("env.toml"); err != nil {
		return
	}
	if _, err := toml.DecodeFile("env.toml", &env); err != nil {
		log.Panic(err)
	}
	fmt.Printf("变量值: %+v\n", env)
}

var dbType = map[string]func(string) gorm.Dialector{
	"sqlite": sqlite.Open,
	"mysql":  mysql.Open,
}

func main() {
	db, err := gorm.Open(dbType[env["db_type"]](env["db_ddns"]), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	db.AutoMigrate(&SaveData{}, &SystemBanData{})

	localBanWindowFilter := Filter.NewBanWindowFilter(100, 180, 0.75, -1, true)  //创建容量为100，窗口有效时间为180秒，相似度要求为0.75的封禁窗口
	systemBanWindowFilter := Filter.NewBanWindowFilter(200, -1, 0.75, -1, false) //-1为永久生效
	banProcess := &CustomBanProcess{
		db:           db,
		localFilter:  localBanWindowFilter,
		systemFilter: systemBanWindowFilter,
		reporter:     Helper.Reporter(env["cookie"], env["csrf"]),
	} //创建自定义封禁处理结构体
	banProcess.RestoreLocalFilter(100)  //从数据库恢复最多100条因频繁发言封禁的记录导入到窗口
	banProcess.RestoreSystemFilter(100) //从数据库恢复最多100条系统封禁的记录导入到窗口
	if err := banProcess.UpdateSystemFilter(); err != nil {
		log.Fatal(err)
	}
	go banProcess.TimingUpdataSystemFilter(30 * time.Second)

	center := DanmuCenter.NewDanmuCenter(&DanmuCenter.DanmuCenterConfig{
		TimeRange:      16,
		MonitorNumber:  1000,
		SpecialFocusOn: []int{1370218}, //1237390
		Silent:         true,
	},
		DanmuCenter.SetProxy(getProxy), // 设置代理
		DanmuCenter.SetPreFilter( //入库前检测
			Helper.Safe(Filter.NewLenFilter(8).Check),                                            // 简易长度过滤
			Helper.Safe(Filter.NewAdminFilter().Check),                                           // 过滤掉房管
			Helper.Break(Filter.NewHaveBeenBanFilter().Check),                                    // 过滤掉已被Ban的弹幕
			Helper.Safe(Filter.NewHighReatWordFilter(0.5).Check),                                 // 单字符重复率>0.75视作正常弹幕
			Helper.Safe(Filter.NewLenFilter(9, Filter.SetLenFilterCompressRepeatGroup(2)).Check), // 过滤掉重复词压缩后长度小于9的弹幕
			Helper.Ban(systemBanWindowFilter.MatchCheck),                                         // 匹配系统确认封禁记录
			Helper.Continue(localBanWindowFilter.UnlockCheck),                                    // 移除高等级的窗口
			Helper.Safe(Filter.NewUserLevelFilter(5).Check),                                      // 过滤掉用户等级>=5的
			Helper.Safe(Filter.NewFansMedalFilter(3).Check),                                      // 过滤掉粉丝勋章等级>=3的
			Helper.Safe(Filter.NewKeyWordFilter([]string{"谢谢", "感谢", "多谢", "欢迎", "点歌"}).Check),   // 关键词匹配过滤
			Helper.Ban(localBanWindowFilter.MatchCheck),                                          // 与封禁窗口比较
		),
		DanmuCenter.SetAfterFilter( //入库后检测
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
