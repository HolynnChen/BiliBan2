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
