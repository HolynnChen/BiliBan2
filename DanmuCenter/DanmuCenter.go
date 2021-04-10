package DanmuCenter

import (
	"context"
	"errors"
	"log"
	"sync"
	"time"

	"github.com/Holynnchen/BiliBan2/DanmuCenter/Utils"
	"github.com/holynnchen/bililive"
)

type Action int

const (
	Continue Action = iota
	Break
	Ban
)

func NewDanmuCenter(config *DanmuCenterConfig, options ...DanmuCenterOption) *DanmuCenter {
	danmuCenter := &DanmuCenter{
		DanmuDB:     new(sync.Map),
		BanDB:       new(sync.Map),
		config:      config,
		preFilter:   make([]Filter, 0),
		afterFilter: make([]Filter, 0),
		banProcess:  nil,
		banIndex:    make([]int64, 0),
		roomIDs:     make(map[int]Utils.Empty),
	}
	live := &bililive.Live{
		Debug:              0,
		AnalysisRoutineNum: 1, //顺序性处理有利于小规模下充分利用封禁窗口
		StormFilter:        true,
		LotteryDanmuFilter: true,
		ReceiveMsg:         danmuCenter.liveReceiveMsg,
		End:                danmuCenter.liveEnd,
	}
	danmuCenter.Live = live
	live.Start(context.Background())
	for _, option := range options {
		option(danmuCenter)
	}
	return danmuCenter
}

func (c *DanmuCenter) liveReceiveMsg(roomID int, msg *bililive.MsgModel) {
	danmu := &Danmu{
		RoomID:     roomID,
		UserID:     msg.UserID,
		UserName:   msg.UserName,
		UserLevel:  msg.UserLevel,
		MedalLevel: msg.MedalLevel,
		Content:    msg.Content,
		Timestamp:  msg.Timestamp,
	}
	//是否入库前检测
	if ok := c.runFilters(&c.preFilter, danmu); ok {
		return
	}
	//过滤过期
	c.DanmuDB.Store(danmu.UserID, append(c.GetRecentDanmu(danmu.UserID), danmu))

	//判断是否异常弹幕
	c.runFilters(&c.afterFilter, danmu)
}

func (c *DanmuCenter) liveEnd(roomID int) {
	delete(c.roomIDs, roomID)
	c.Live.Remove(roomID)
}

func (c *DanmuCenter) runFilters(filters *[]Filter, danmu *Danmu) bool {
	for _, filter := range c.preFilter {
		action, reason := filter(c, danmu)
		switch action {
		case Break:
			return true
		case Ban:
			c.ban(danmu, danmu.RoomID, reason)
			return true
		case Continue:
			continue
		default:
			log.Println("error: fail to identify action")
		}
	}
	return false
}

func (c *DanmuCenter) ban(danmu *Danmu, roomID int, reason string) {
	banData := &BanData{
		UserID:    danmu.UserID,
		UserName:  danmu.UserName,
		RoomID:    roomID,
		Content:   danmu.Content,
		Timestamp: danmu.Timestamp,
		Reason:    reason,
	}
	if _, hasBan := c.BanDB.LoadOrStore(danmu.UserID, banData); hasBan {
		return
	}
	c.banIndex = append(c.banIndex, danmu.UserID)
	c.banProcess.Ban(banData)
	return
}

func (c *DanmuCenter) filterValidDanmu(danmuList []*Danmu, timeRange int64) []*Danmu {
	index := len(danmuList)
	for ; index > 0; index-- {
		if time.Now().Unix()-danmuList[index-1].Timestamp > timeRange {
			break
		}
	}
	return danmuList[index:]
}

func (c *DanmuCenter) GetRecentDanmu(UserID int64) []*Danmu {
	rawData, ok := c.DanmuDB.Load(UserID)
	if !ok {
		return []*Danmu{}
	}
	danmuList := rawData.([]*Danmu)
	return c.filterValidDanmu(danmuList, c.config.TimeRange)
}

func (c *DanmuCenter) cleanDanmuDB(timeRange int64) {
	count := 0
	all := 0
	c.DanmuDB.Range(func(key, value interface{}) bool {
		all++
		danmuList := value.([]*Danmu)
		if danmuLen := len(danmuList); danmuLen != 0 {
			if time.Now().Unix()-danmuList[danmuLen-1].Timestamp > timeRange {
				c.DanmuDB.Delete(key)
				count++
				all--
			}
		}
		return true
	})
	if !c.config.Silent {
		log.Printf("定时清理弹幕DB：移除%d个过期key，当前%d个key\n", count, all)
	}
}

func (c *DanmuCenter) updateRoom(monitorNumber int) error {
	var (
		newRoomIDs []int
		err        error
	)
	switch c.config.RankType {
	case 0:
		newRoomIDs, err = Utils.GetTopRoom(monitorNumber)
	case 1:
		newRoomIDs, err = Utils.GetTop50HotRoom()
	default:
		err = errors.New("no such rank type ")
	}
	if err != nil {
		return err
	}
	if c.config.SpecialFocusOn != nil {
		newRoomIDs = append(newRoomIDs, c.config.SpecialFocusOn...)
	}
	addList := make([]int, 0)
	removeList := make([]int, 0)
	for _, id := range newRoomIDs {
		if _, ok := c.roomIDs[id]; !ok {
			c.roomIDs[id] = Utils.Empty{}
			addList = append(addList, id)
		}
	}
	newRoomIDsMap := Utils.IntArrayToMap(newRoomIDs)
	for id := range c.roomIDs {
		if _, ok := newRoomIDsMap[id]; !ok {
			removeList = append(removeList, id)
		}
	}
	if !c.config.Silent && (len(removeList) > 0 || len(addList) > 0) {
		log.Printf("同步热门榜：移除%d个房间(%v)，新增%d个房间(%v)\n", len(removeList), removeList, len(addList), addList)
	}
	c.Live.Remove(removeList...)
	c.Live.Join(addList...)
	c.roomIDs = newRoomIDsMap
	return nil
}

func (c *DanmuCenter) tickerTask() {
	tickerCleanDanmu := time.NewTicker(time.Minute)
	tickerUpdateRoom := time.NewTicker(time.Minute)
	for {
		select {
		case <-context.Background().Done():
			tickerCleanDanmu.Stop()
			tickerUpdateRoom.Stop()
			return
		case <-tickerCleanDanmu.C:
			go c.cleanDanmuDB(c.config.TimeRange)
		case <-tickerUpdateRoom.C:
			go c.updateRoom(c.config.MonitorNumber)
		}
	}

}

type defaultBan struct{}

func (defaultBan) Ban(banData *BanData) {
	log.Printf("%+v\n", banData)
}

func (c *DanmuCenter) setDefaultConfig() {
	if c.banProcess == nil {
		c.banProcess = defaultBan{}
	}
	if c.config.TimeRange == 0 {
		c.config.TimeRange = 10
	}
	if c.config.MonitorNumber == 0 {
		c.config.MonitorNumber = 10
	}
}

func (c *DanmuCenter) Start() {
	c.setDefaultConfig()
	err := c.updateRoom(c.config.MonitorNumber)
	if err != nil {
		log.Fatal(err)
	}
	go c.tickerTask()
	c.Live.Wait()
}
