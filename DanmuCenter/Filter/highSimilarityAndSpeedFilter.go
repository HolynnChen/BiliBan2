package Filter

import (
	"log"

	"github.com/Holynnchen/BiliBan2/DanmuCenter"
	"github.com/Holynnchen/BiliBan2/DanmuCenter/Utils"
)

//封禁高速发言且整体连续(startCheck条)相似度大于similarity的账户
type HighSimilarityAndSpeedFilter struct {
	similarity float32
	startCheck int
}

func (filter *HighSimilarityAndSpeedFilter) Check(center *DanmuCenter.DanmuCenter, danmu *DanmuCenter.Danmu) (bool, string) {
	danmuList := center.GetRecentDanmu(danmu.UserID)
	dataLen := len(danmuList)
	if dataLen < filter.startCheck {
		return false, ""
	}
	var allCompare float32 = 0
	for i := 1; i < dataLen; i++ {
		allCompare = (allCompare*float32(i-1) + Utils.GetSimilarity(danmuList[dataLen-i].Content, danmuList[dataLen-i-1].Content)) / float32(i)
		if i > filter.startCheck-1 && allCompare > filter.similarity {
			for _,danmu := range danmuList{
				log.Printf("%+v ",*danmu)
			}
			log.Println()
			return true, "时间范围内近似发言过多"
		}
	}
	return false, ""
}

func NewHighSimilarityAndSpeedFilter(similarity float32, startCheck int) *HighSimilarityAndSpeedFilter {
	return &HighSimilarityAndSpeedFilter{
		similarity: similarity,
		startCheck: startCheck,
	}
}
