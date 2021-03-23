package Filter

import (
	"time"

	"github.com/Holynnchen/BiliBan2/DanmuCenter"
	"github.com/Holynnchen/BiliBan2/DanmuCenter/Utils"
)

type banWindowFilter struct {
	banWindow     []*banWindowData
	banWindowSize int
	banWindowTime int64
	writeMark     int
	nowSize       int
	similarity    float32
}

type banWindowData struct {
	banString  string
	enableTime int64
	disable    bool
}

func (filter *banWindowFilter) SaveCheck(center *DanmuCenter.DanmuCenter, danmu *DanmuCenter.Danmu) (bool, string) {
	if danmu.MedalLevel == 0 && danmu.UserLevel < 3 {
		return false, ""
	}
	content := Utils.ReplaceSimilarAndNumberRune(danmu.Content)
	for i := 1; i < filter.nowSize+1; i++ {
		banWindowData := filter.banWindow[(filter.writeMark-i+filter.banWindowSize)%filter.banWindowSize]
		if time.Now().Unix()-banWindowData.enableTime > filter.banWindowTime {
			break
		}
		if banWindowData.disable {
			continue
		}
		if Utils.GetSimilarity(banWindowData.banString, content) > 0.9 { //固定在0.9
			banWindowData.disable = true
			break
		}
	}
	return false, ""
}

func (filter *banWindowFilter) BanCheck(center *DanmuCenter.DanmuCenter, danmu *DanmuCenter.Danmu) (bool, string) {
	content := Utils.ReplaceSimilarAndNumberRune(danmu.Content)
	for i := 1; i < filter.nowSize+1; i++ {
		banWindowData := filter.banWindow[(filter.writeMark-i+filter.banWindowSize)%filter.banWindowSize]
		if time.Now().Unix()-banWindowData.enableTime > filter.banWindowTime {
			break
		}
		if banWindowData.disable || banWindowData.enableTime > time.Now().Unix() {
			continue
		}
		if Utils.GetSimilarity(banWindowData.banString, content) > filter.similarity {
			banWindowData.enableTime = time.Now().Unix() //时间续期
			return true, "匹配封禁窗口"
		}
	}
	return false, ""
}

func (filter *banWindowFilter) Ban(banData *DanmuCenter.BanData) {
	content := Utils.ReplaceSimilarAndNumberRune(banData.Content)
	for i := 1; i < filter.nowSize+1; i++ {
		banWindowData := filter.banWindow[(filter.writeMark-i+filter.banWindowSize)%filter.banWindowSize]
		if time.Now().Unix()-banWindowData.enableTime > filter.banWindowTime {
			break
		}
		if Utils.GetSimilarity(banWindowData.banString, content) > filter.similarity {
			return
		}
	}
	filter.banWindow[filter.writeMark] = &banWindowData{banString: content, enableTime: time.Now().Unix() + 10, disable: false}
	filter.writeMark = (filter.writeMark + 1) % filter.banWindowSize
	if filter.nowSize < filter.banWindowSize {
		filter.nowSize++
	}
}

/*
生成一个窗口大小为banWindowSize，最短有效时间为banWindowTime，要求相似率为similarity的封禁窗口
该窗口用于快速匹配相似语句，达到快速封杀的目的
为了防止手动节奏的情况，增加了后悔机制
    加入窗口后，要过10s后才会生效
    加入窗口后，若有粉丝勋章等级>0 或 用户等级>=3的相似率大于0.9的发言，那么将禁用此条记录
*/
func NewBanWindowFilter(banWindowSize int, banWindowTime int64, similarity float32) *banWindowFilter {
	return &banWindowFilter{
		banWindow:     make([]*banWindowData, banWindowSize, banWindowSize),
		banWindowSize: banWindowSize,
		banWindowTime: banWindowTime,
		similarity:    similarity,
		writeMark:     0,
		nowSize:       0,
	}
}
