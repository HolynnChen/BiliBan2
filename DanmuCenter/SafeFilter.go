package DanmuCenter

import "unicode/utf8"

//单重复率大于repeatTarget视作正常弹幕
type highRepeatWordFilter struct {
	repeatTarget float32
}

func (filter *highRepeatWordFilter) Check(center *DanmuCenter, danmu *Danmu) (bool, string) {
	countMap := make(map[rune]int)
	max := 0
	for _, data := range danmu.Content {
		countMap[data]++
		if countMap[data] > max {
			max = countMap[data]
		}
	}
	return float32(max)/float32(utf8.RuneCountInString(danmu.Content)) > filter.repeatTarget, ""
}

func NewHighReatWordFilter(repeatTarget float32) *highRepeatWordFilter {
	return &highRepeatWordFilter{
		repeatTarget: repeatTarget,
	}
}

//小于lenTarget视作正常弹幕
type lenFilter struct {
	repeatGroupCompress func(string) string
	lenTarget           int
}

type lenFilterOption func(*lenFilter)

func (filter *lenFilter) Check(center *DanmuCenter, danmu *Danmu) (bool, string) {
	content := danmu.Content
	if filter.repeatGroupCompress != nil {
		content = filter.repeatGroupCompress(content)
	}
	return utf8.RuneCountInString(content) < filter.lenTarget, ""
}

func NewLenFilter(lenTarget int, options ...lenFilterOption) *lenFilter {
	filter := &lenFilter{
		lenTarget: lenTarget,
	}
	for _, option := range options {
		option(filter)
	}
	return filter
}

func SetLenFilterCompressRepeatGroup(minLen int) lenFilterOption {
	return func(lf *lenFilter) {
		lf.repeatGroupCompress = CompressRepeatGroup(minLen)
	}
}

// 被封禁了就直接忽略
type haveBeenBanFilter struct{}

func (filter *haveBeenBanFilter) Check(center *DanmuCenter, danmu *Danmu) (bool, string) {
	_, ok := center.BanDB.Load(danmu.UserID)
	return ok, ""
}

// Filter -> haveBeenBan
func NewHaveBeenBanFilter() *haveBeenBanFilter {
	return &haveBeenBanFilter{}
}

// uid小于该区域视作安全
type uidFilter struct {
	uidTarget int64
}

func (filter *uidFilter) Check(center *DanmuCenter, danmu *Danmu) (bool, string) {
	return danmu.UserID < filter.uidTarget, ""
}

// Filter-> uid<uidTarget
func NewUIDFilter(uidTarget int64) *uidFilter {
	return &uidFilter{
		uidTarget: uidTarget,
	}
}

// userLevel >= levelTarget视为正常
type levelFilter struct {
	levelTarget int
}

func (filter *levelFilter) Check(center *DanmuCenter, danmu *Danmu) (bool, string) {
	return danmu.UserLevel >= filter.levelTarget, ""
}

// Filter-> userLever >= levelTarget
func NewUserLevelFilter(levelTarget int) *levelFilter {
	return &levelFilter{
		levelTarget: levelTarget,
	}
}

// fansMedalLever >=levelTarget视为正常
type fansMedalFilter struct {
	levelTarget int
}

func (filter *fansMedalFilter) Check(center *DanmuCenter, danmu *Danmu) (bool, string) {
	return danmu.MedalLevel >= filter.levelTarget, ""
}

// Filter-> fansMedalLevel >= leverTarget
func NewFansMedalFilter(levelTarget int) *fansMedalFilter {
	return &fansMedalFilter{
		levelTarget: levelTarget,
	}
}
