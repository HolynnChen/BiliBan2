package DanmuCenter

import (
	"net/http"
	"net/url"
	"sync"

	"github.com/Holynnchen/BiliBan2/DanmuCenter/Utils"
	"github.com/holynnchen/bililive"
)

type Process func(content string) string

type BanProcess interface {
	Ban(banData *BanData)
}

type Filter func(center *DanmuCenter, danmu *Danmu) (Action, string)

type DanmuCenter struct {
	DanmuDB *sync.Map      // 弹幕储存
	BanDB   *sync.Map      // 封禁记录储存
	Live    *bililive.Live // 直播间实例
	//private
	config      *DanmuCenterConfig
	preFilter   []Filter
	afterFilter []Filter
	banProcess  BanProcess
	banIndex    []int64
	roomIDs     map[int]Utils.Empty
}

type DanmuCenterConfig struct {
	TimeRange      int64 // 弹幕储存时间范围，秒级
	MonitorNumber  int   // 热门榜前几
	SpecialFocusOn []int // 特别关注的直播间
	RankType       int   // 获取直播间的途径
	Silent         bool  // 安静模式
}

type DanmuCenterOption func(center *DanmuCenter)

type Danmu struct {
	UserID     int64  // 用户uid
	RoomID     int    // 房间id
	UserName   string // 用户名
	UserLevel  int    // 用户等级
	IsAdmin    bool   // 是否房管
	MedalLevel int    // 勋章等级
	Content    string // 弹幕内容
	CT         string // 弹幕token
	Timestamp  int64  // 时间戳
}

func (danmu *Danmu) DeepCopy() Danmu {
	return *danmu
}

type BanData struct {
	UserID    int64  `json:"user_id"`
	UserName  string `json:"user_name"`
	RoomID    int    `json:"room_id"`
	Content   string `json:"content"`
	CT        string `json:"ct"`
	Timestamp int64  `json:"timestamp"`
	Reason    string `json:"reason"`
}

func SetPreFilter(filters ...Filter) DanmuCenterOption {
	return func(center *DanmuCenter) {
		center.preFilter = append(center.preFilter, filters...)
	}
}

func SetAfterFilter(filters ...Filter) DanmuCenterOption {
	return func(center *DanmuCenter) {
		center.afterFilter = append(center.afterFilter, filters...)
	}
}

func SetFilter(filters ...Filter) DanmuCenterOption {
	return SetAfterFilter(filters...)
}

// Process 处理封禁情况
func SetBanProcess(process BanProcess) DanmuCenterOption {
	return func(center *DanmuCenter) {
		center.banProcess = process
	}
}

func SetProxy(proxyFunc func() func(*http.Request) (*url.URL, error)) DanmuCenterOption {
	return func(center *DanmuCenter) {
		center.Live.Proxy = proxyFunc
	}
}
