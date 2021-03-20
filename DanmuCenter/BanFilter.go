package DanmuCenter

import (
	"log"
	"time"

	"github.com/TheFutureIsOurs/ahocorasick"
)

//封禁高速发言且整体连续(三条)相似度大于similarity的账户
type highSimilarityAndSpeedFilter struct {
	similarity float32
	startCheck int
}

func (filter *highSimilarityAndSpeedFilter) Check(center *DanmuCenter, danmu *Danmu) (bool, string) {
	danmuList := center.GetRecentDanmu(danmu.UserID)
	dataLen := len(danmuList)
	if dataLen < filter.startCheck {
		return false, ""
	}
	var allCompare float32 = 0
	for i := 1; i < dataLen; i++ {
		allCompare = (allCompare*float32(i-1) + GetSimilarity(danmuList[dataLen-i].Content, danmuList[dataLen-i-1].Content)) / float32(i)
		if i > filter.startCheck-1 && allCompare > filter.similarity {
			return true, "时间范围内近似发言过多"
		}
	}
	return false, ""
}

func NewHighSimilarityAndSpeedFilter(similarity float32, startCheck int) *highSimilarityAndSpeedFilter {
	return &highSimilarityAndSpeedFilter{
		similarity: similarity,
		startCheck: startCheck,
	}
}

type banWindowFilter struct {
	banWindow     []*banWindowData
	banWindowSize int
	banWindowTime int64
	writeMark     int
	nowSize       int
	similarity    float32
}

type banWindowData struct {
	banString string
	banTime   int64
}

func (filter *banWindowFilter) Check(center *DanmuCenter, danmu *Danmu) (bool, string) {
	content := ReplaceSimilarAndNumberRune(danmu.Content)
	for i := 1; i < filter.nowSize+1; i++ {
		banWindowData := filter.banWindow[(filter.writeMark-i+filter.banWindowSize)%filter.banWindowSize]
		if time.Now().Unix()-banWindowData.banTime > filter.banWindowTime {
			break
		}
		if GetSimilarity(banWindowData.banString, content) > filter.similarity {
			banWindowData.banTime = time.Now().Unix() //时间续期
			return true, "匹配封禁窗口"
		}
	}
	return false, ""
}

func (filter *banWindowFilter) Ban(banData *BanData) {
	content := ReplaceSimilarAndNumberRune(banData.Content)
	for i := 1; i < filter.nowSize+1; i++ {
		banWindowData := filter.banWindow[(filter.writeMark-i+filter.banWindowSize)%filter.banWindowSize]
		if time.Now().Unix()-banWindowData.banTime > filter.banWindowTime {
			break
		}
		if GetSimilarity(banWindowData.banString, content) > filter.similarity {
			return
		}
	}
	filter.banWindow[filter.writeMark] = &banWindowData{banString: content, banTime: time.Now().Unix()}
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

type keyWordFilter struct {
	*ahocorasick.Ac
}

func (filter *keyWordFilter) Check(center *DanmuCenter, danmu *Danmu) (bool, string) {
	return filter.MultiPatternHit([]rune(danmu.Content)), "关键词匹配"
}

func NewKeyWordFilter(keywords []string) *keyWordFilter {
	ac, err := ahocorasick.Build(keywords)
	if err != nil {
		log.Fatal(err)
	}
	return &keyWordFilter{ac}
}
