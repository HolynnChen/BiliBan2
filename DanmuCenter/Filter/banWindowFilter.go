package Filter

import (
	"log"
	"time"

	"github.com/Holynnchen/BiliBan2/DanmuCenter"
	"github.com/Holynnchen/BiliBan2/DanmuCenter/Utils"
)

type BanWindowFilter struct {
	banWindow     []*banWindowData
	banWindowSize int
	banWindowTime int64
	writeMark     int
	nowSize       int
	similarity    float32
	fuzzy         bool
}

type banWindowData struct {
	banString  string
	enableTime int64
	disable    bool
}

const hightSaveCheck float32 = 0.85 // 固定安全期望

func (filter *BanWindowFilter) UnlockCheck(center *DanmuCenter.DanmuCenter, danmu *DanmuCenter.Danmu) (bool, string) {
	if danmu.MedalLevel == 0 && danmu.UserLevel < 3 {
		return false, ""
	}
	content := Utils.SimpleReplaceSimilar(danmu.Content)
	if filter.fuzzy {
		content = Utils.ReplaceSimilarAndNumberRune(content)
	}
	for i := 1; i < filter.nowSize+1; i++ {
		banWindowData := filter.banWindow[(filter.writeMark-i+filter.banWindowSize)%filter.banWindowSize]
		if filter.banWindowTime > 0 && time.Now().Unix()-banWindowData.enableTime > filter.banWindowTime {
			break
		}
		if banWindowData.disable {
			continue
		}
		if checkValue := Utils.GetSimilarity(banWindowData.banString, content); checkValue > hightSaveCheck {
			log.Printf("解封窗口 %.4f %+v %+v\n", checkValue, banWindowData, danmu)
			banWindowData.disable = true
			break
		}
	}
	return false, ""
}

func (filter *BanWindowFilter) MatchCheck(center *DanmuCenter.DanmuCenter, danmu *DanmuCenter.Danmu) (bool, string) {
	content := Utils.SimpleReplaceSimilar(danmu.Content)
	if filter.fuzzy {
		content = Utils.ReplaceSimilarAndNumberRune(content)
	}
	for i := 1; i < filter.nowSize+1; i++ {
		banWindowData := filter.banWindow[(filter.writeMark-i+filter.banWindowSize)%filter.banWindowSize]
		if filter.banWindowTime > 0 && time.Now().Unix()-banWindowData.enableTime > filter.banWindowTime {
			break
		}
		if banWindowData.enableTime > time.Now().Unix() {
			continue
		}
		similarity := Utils.GetSimilarity(banWindowData.banString, content)
		if !banWindowData.disable && similarity > filter.similarity {
			banWindowData.enableTime = time.Now().Unix() //时间续期
			return true, "匹配封禁窗口【" + banWindowData.banString + "】"
		}
		if banWindowData.disable && similarity > hightSaveCheck {
			return false, ""
		}
	}
	return false, ""
}

func (filter *BanWindowFilter) Add(content string) bool {
	content = Utils.SimpleReplaceSimilar(content)
	if filter.fuzzy {
		content = Utils.ReplaceSimilarAndNumberRune(content)
	}
	for i := 1; i < filter.nowSize+1; i++ {
		banWindowData := filter.banWindow[(filter.writeMark-i+filter.banWindowSize)%filter.banWindowSize]
		if filter.banWindowTime > 0 && time.Now().Unix()-banWindowData.enableTime > filter.banWindowTime {
			break
		}
		if Utils.GetSimilarity(banWindowData.banString, content) > filter.similarity {
			return false
		}
	}
	filter.banWindow[filter.writeMark] = &banWindowData{banString: content, enableTime: time.Now().Unix() + 10, disable: false}
	log.Printf("加入窗口 %+v\n", filter.banWindow[filter.writeMark])
	filter.writeMark = (filter.writeMark + 1) % filter.banWindowSize
	if filter.nowSize < filter.banWindowSize {
		filter.nowSize++
	}
	return true
}

/*
生成一个窗口大小为banWindowSize，最短有效时间为banWindowTime，要求相似率为similarity的封禁窗口
该窗口用于快速匹配相似语句，达到快速封杀的目的
为了防止手动节奏的情况，增加了后悔机制
    加入窗口后，要过10s后才会生效
    加入窗口后，若有粉丝勋章等级>0 或 用户等级>=3的相似率大于0.9的发言，那么将禁用此条记录
*/
func NewBanWindowFilter(banWindowSize int, banWindowTime int64, similarity float32, fuzzy bool) *BanWindowFilter {
	return &BanWindowFilter{
		banWindow:     make([]*banWindowData, banWindowSize, banWindowSize),
		banWindowSize: banWindowSize,
		banWindowTime: banWindowTime,
		similarity:    similarity,
		writeMark:     0,
		nowSize:       0,
		fuzzy:         fuzzy,
	}
}
