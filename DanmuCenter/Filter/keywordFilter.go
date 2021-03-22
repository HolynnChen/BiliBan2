package Filter

import (
	"log"

	"github.com/Holynnchen/BiliBan2/DanmuCenter"
	"github.com/TheFutureIsOurs/ahocorasick"
)

type keyWordFilter struct {
	*ahocorasick.Ac
}

func (filter *keyWordFilter) Check(center *DanmuCenter.DanmuCenter, danmu *DanmuCenter.Danmu) (bool, string) {
	return filter.MultiPatternHit([]rune(danmu.Content)), "关键词匹配"
}
func (filter *keyWordFilter) SaveCheck(center *DanmuCenter.DanmuCenter, danmu *DanmuCenter.Danmu) (bool, string) {
	return filter.Check(center, danmu)
}
func (filter *keyWordFilter) SafeCheck(center *DanmuCenter.DanmuCenter, danmu *DanmuCenter.Danmu) (bool, string) {
	return filter.Check(center, danmu)
}
func (filter *keyWordFilter) BanCheck(center *DanmuCenter.DanmuCenter, danmu *DanmuCenter.Danmu) (bool, string) {
	return filter.Check(center, danmu)
}

func NewKeyWordFilter(keywords []string) *keyWordFilter {
	ac, err := ahocorasick.Build(keywords)
	if err != nil {
		log.Fatal(err)
	}
	return &keyWordFilter{ac}
}
