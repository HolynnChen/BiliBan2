package Helper

import (
	"github.com/Holynnchen/BiliBan2/DanmuCenter"
)

type boolFilter func(*DanmuCenter.DanmuCenter, *DanmuCenter.Danmu) (bool, string)

func Safe(filter boolFilter) DanmuCenter.Filter {
	return func(center *DanmuCenter.DanmuCenter, danmu *DanmuCenter.Danmu) (DanmuCenter.Action, string) {
		result, reason := filter(center, danmu)
		if result {
			return DanmuCenter.Break, reason
		}
		return DanmuCenter.Continue, reason
	}

}

func Break(filter boolFilter) DanmuCenter.Filter {
	return Safe(filter)
}

func Ban(filter boolFilter) DanmuCenter.Filter {
	return func(center *DanmuCenter.DanmuCenter, danmu *DanmuCenter.Danmu) (DanmuCenter.Action, string) {
		result, reason := filter(center, danmu)
		if result {
			return DanmuCenter.Ban, reason
		}
		return DanmuCenter.Continue, reason
	}
}
func Continue(filter boolFilter) DanmuCenter.Filter {
	return func(center *DanmuCenter.DanmuCenter, danmu *DanmuCenter.Danmu) (DanmuCenter.Action, string) {
		filter(center, danmu)
		return DanmuCenter.Continue, ""
	}
}
