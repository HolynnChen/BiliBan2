package DanmuCenter

import (
	"sync"

	"github.com/holynnchen/bililive"
)

type Process func(content string) string

type Filter interface {
	Check(center *DanmuCenter, danmu *Danmu) (bool, string)
}

type BanProcess interface {
	Ban(banData *BanData)
}

type Empty struct{}

type DanmuCenter struct {
	DanmuDB *sync.Map      // 弹幕储存
	BanDB   *sync.Map      // 封禁记录储存
	Live    *bililive.Live // 直播间实例
	//private
	config     *DanmuCenterConfig
	saveFilter []Filter
	safeFilter []Filter
	banFilter  []Filter
	banProcess BanProcess
	banIndex   []int64
	roomIDs    map[int]Empty
}

type DanmuCenterConfig struct {
	TimeRange      int64 // 弹幕储存时间范围，秒级
	MonitorNumber  int   // 热门榜前几
	SpecialFocusOn []int // 特别关注的直播间
	Silent         bool  // 安静模式
}

type DanmuCenterOption func(center *DanmuCenter)

type Danmu struct {
	UserID     int64  // 用户uid
	RoomID     int    // 房间id
	UserName   string // 用户名
	UserLevel  int    //用户等级
	MedalLevel int    // 勋章等级
	Content    string // 弹幕内容
	Timestamp  int64  // 时间戳
}

type BanData struct {
	UserID    int64  `json:"user_id"`
	UserName  string `json:"user_name"`
	RoomID    int    `json:"room_id"`
	Content   string `json:"content"`
	Timestamp int64  `json:"timestamp"`
	Reason    string `json:"reason"`
}

// Filter->true 正常弹幕
func SetSafeFilter(filters ...Filter) DanmuCenterOption {
	return func(center *DanmuCenter) {
		center.safeFilter = append(center.safeFilter, filters...)
	}
}

// Filter->true 封禁弹幕
func SetBanFilter(filters ...Filter) DanmuCenterOption {
	return func(center *DanmuCenter) {
		center.banFilter = append(center.banFilter, filters...)
	}
}

// Filter->true 过滤不入库不检测
func SetSaveFilter(filter ...Filter) DanmuCenterOption {
	return func(center *DanmuCenter) {
		center.saveFilter = append(center.saveFilter, filter...)
	}
}

// Process 处理封禁情况
func SetBanProcess(process BanProcess) DanmuCenterOption {
	return func(center *DanmuCenter) {
		center.banProcess = process
	}
}
