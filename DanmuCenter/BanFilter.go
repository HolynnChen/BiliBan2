package DanmuCenter

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
	banWindow     []string
	banWindowSize int
	writeMark     int
	nowSize       int
	similarity    float32
}

func (filter *banWindowFilter) Check(center *DanmuCenter, danmu *Danmu) (bool, string) {
	content := ReplaceSimilarAndNumberRune(danmu.Content)
	for i := 1; i < filter.nowSize; i++ {
		if GetSimilarity(filter.banWindow[(filter.writeMark-i+filter.banWindowSize)%filter.banWindowSize], content) > filter.similarity {
			return true, "匹配封禁窗口"
		}
	}
	return false, ""
}

func (filter *banWindowFilter) Ban(banData *BanData) {
	content := ReplaceSimilarAndNumberRune(banData.Content)
	for i := 1; i < filter.nowSize; i++ {
		if GetSimilarity(filter.banWindow[(filter.writeMark-i+filter.banWindowSize)%filter.banWindowSize], content) > filter.similarity {
			return
		}
	}
	filter.banWindow[filter.writeMark] = content
	filter.writeMark = (filter.writeMark + 1) % filter.banWindowSize
	if filter.nowSize < filter.banWindowSize {
		filter.nowSize++
	}
}

func NewBanWindowFilter(banWindowSize int, similarity float32) *banWindowFilter {
	return &banWindowFilter{
		banWindow:     make([]string, banWindowSize, banWindowSize),
		banWindowSize: banWindowSize,
		similarity:    similarity,
		writeMark:     0,
		nowSize:       0,
	}
}
