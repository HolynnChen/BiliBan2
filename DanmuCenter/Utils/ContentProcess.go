package Utils

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/Holynnchen/BiliBan2/DanmuCenter/Utils/static"
	"github.com/dlclark/regexp2"
)

var replaceDict = static.ReplaceDict
var replaceSpecial = regexp.MustCompile(`[.|/\\@\-*&^ +-]`)

func ReplaceSimilarAndNumberRune(rawData string) string {
	result := strings.Map(func(code rune) rune {
		if result, ok := replaceDict[code]; ok {
			code = result
		}
		if code == 12288 {
			return -1
		}
		if code > 65280 && code < 65375 {
			return code - 65248
		}
		if code >= 0x0030 && code <= 0x0039 ||
			code >= 0x2460 && code <= 0x249b ||
			code >= 0x3220 && code <= 0x3229 ||
			code >= 0x3248 && code <= 0x324f ||
			code >= 0x3251 && code <= 0x325f ||
			code >= 0x3280 && code <= 0x3289 ||
			code >= 0x32b1 && code <= 0x32bf ||
			code >= 0xff10 && code <= 0xff19 {
			return '#'
		}
		return code
	}, rawData)
	result = replaceSpecial.ReplaceAllString(strings.ToLower(result), "")
	return result
}

func SimpleReplaceSimilar(rawData string) string {
	result := strings.Map(func(code rune) rune {
		if result, ok := replaceDict[code]; ok {
			code = result
		}
		if code == 12288 {
			return -1
		}
		if code > 65280 && code < 65375 {
			return code - 65248
		}
		if code >= 0x2460 && code <= 0x249b ||
			code >= 0x3220 && code <= 0x3229 ||
			code >= 0x3248 && code <= 0x324f ||
			code >= 0x3251 && code <= 0x325f ||
			code >= 0x3280 && code <= 0x3289 ||
			code >= 0x32b1 && code <= 0x32bf ||
			code >= 0xff10 && code <= 0xff19 {
			return '#'
		}
		return code
	}, rawData)
	result = replaceSpecial.ReplaceAllString(strings.ToLower(result), "")
	return result
}

func CompressRepeatGroup(repeatGroupMinLen int) func(string) string {
	replaceGroup := regexp2.MustCompile(`(.{`+strconv.Itoa(repeatGroupMinLen)+`,})(?=.*\1)+`, 0)
	return func(s string) string {
		m, _ := replaceGroup.FindStringMatch(s)
		for m != nil {
			word := m.Capture.String()
			split := strings.Split(s, word)
			if len(split) > 1 {
				s = split[0] + word + strings.Join(split[1:], "")
			}
			m, _ = replaceGroup.FindNextMatch(m)
		}
		return s
	}
}
