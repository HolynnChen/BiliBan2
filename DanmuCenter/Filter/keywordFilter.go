package Filter

import (
	"log"

	"github.com/Holynnchen/BiliBan2/DanmuCenter"
	"github.com/TheFutureIsOurs/ahocorasick"
)

type KeyWordFilter struct {
	*ahocorasick.Ac
}

func (filter *KeyWordFilter) Check(center *DanmuCenter.DanmuCenter, danmu *DanmuCenter.Danmu) (bool, string) {
	return filter.MultiPatternHit([]rune(danmu.Content)), "关键词匹配"
}

func NewKeyWordFilter(keywords []string) *KeyWordFilter {
	ac, err := ahocorasick.Build(keywords)
	if err != nil {
		log.Fatal(err)
	}
	return &KeyWordFilter{ac}
}
