package Filter

import "github.com/Holynnchen/BiliBan2/DanmuCenter"

// 被封禁了就直接忽略
type haveBeenBanFilter struct{}

func (filter *haveBeenBanFilter) SaveCheck(center *DanmuCenter.DanmuCenter, danmu *DanmuCenter.Danmu) (bool, string) {
	_, ok := center.BanDB.Load(danmu.UserID)
	return ok, ""
}

// Filter -> haveBeenBan
func NewHaveBeenBanFilter() *haveBeenBanFilter {
	return &haveBeenBanFilter{}
}