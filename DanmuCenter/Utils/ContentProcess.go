package Utils

import (
	"bytes"
	"log"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"unicode/utf8"

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

func CompressRepeatGroup2(repeatGroupMinLen int) func(string) string {
	operator := NewUTF8StringSuffixArray()
	var buffer bytes.Buffer
	return func(s string) string {
		buffer.Reset()
		buffer.WriteString(s)
		for {
			value := buffer.Bytes()
			operator.Init(bytes2str(value))
			word := operator.MaxAreaString(repeatGroupMinLen)
			if word == "" {
				return bytes2str(value)
			}
			sWord := str2bytes(word)
			buffer.Reset()
			start := 0
			for {
				if p := bytes.Index(value[start:], sWord); p > -1 {

					buffer.Write(value[start : start+p])
					if start == 0 {
						buffer.Write(sWord)
					}
					start += p + len(sWord)
				} else {
					buffer.Write(value[start:])
					break
				}
			}
		}
	}
}

type UTF8StringSuffixArray struct {
	sa     []int
	rank   []int
	height []int
	source []rune
	len    int
	stack  *Stack
}

func (s *UTF8StringSuffixArray) Len() int {
	return s.len
}

func (s *UTF8StringSuffixArray) Less(i, j int) bool {
	len_i := s.len - s.sa[i]
	len_j := s.len - s.sa[j]
	len_min := min(len_i, len_j)
	for k := 0; k < len_min; k++ {
		if s.source[s.sa[i]+k] == s.source[s.sa[j]+k] {
			continue
		}
		return s.source[s.sa[i]+k] < s.source[s.sa[j]+k]
	}
	return len_i < len_j
}

func (s *UTF8StringSuffixArray) Swap(i, j int) {
	s.sa[i], s.sa[j] = s.sa[j], s.sa[i]
}

func (s *UTF8StringSuffixArray) commonLen(i, j int) int {
	if i > j {
		i, j = j, i
	}
	if i == j {
		return 0
	}
	return min(s.height[i], s.height[i+1:j+1]...)
}

func (s *UTF8StringSuffixArray) Init(str string) {
	s.len = 0
	if utf8.RuneCountInString(str) > 63 {
		log.Fatal("长度大于63")
	}
	for _, r := range str {
		s.source[s.len] = r
		s.len++
	}
	s.source[s.len] = 0
	s.len++
	for i := 0; i < s.len; i++ {
		s.sa[i] = i
	}
	sort.Sort(s)
	for i := 0; i < s.len; i++ {
		s.rank[s.sa[i]] = i
	}
	for i, k := 0, 0; i < s.len-1; i++ {
		if k > 0 {
			k--
		}
		for s.source[i+k] == s.source[s.sa[s.rank[i]-1]+k] {
			k++
		}
		s.height[s.rank[i]] = k
	}
	//方便后续计算最大权重
	s.height[s.len] = 0
}

func (s *UTF8StringSuffixArray) MaxAreaString(minLen int) string {
	if s.len <= 1 {
		return string(s.source[:1])
	}
	area := 0
	left := 0
	h := 0
	s.stack.MustPush(0)
	for i := 1; i < s.len+1; i++ { //len+1是因为额外占用掉一个height最后的空
		for s.height[s.stack.MustPeek().(int)] > s.height[i] {
			height := s.height[s.stack.MustPop().(int)]
			width := i - 1 - s.stack.MustPeek().(int)
			tmpLeft := i - width
			tmpArea := 0
			if height >= minLen && width > 1 {
				for j := tmpLeft; j < i; j++ {
					if j-tmpLeft >= height {
						tmpLeft = j
						tmpArea += height
					}
				}
				if tmpArea > area {
					h = height
					left = i - width
					area = tmpArea
				}
			}
		}
		s.stack.MustPush(i)
	}
	s.stack.Empty()
	return string(s.source[s.sa[left] : s.sa[left]+h])
}

//反正就短弹幕，直接限制最大输入长度为63
func NewUTF8StringSuffixArray() *UTF8StringSuffixArray {
	return &UTF8StringSuffixArray{
		source: make([]rune, 64),
		sa:     make([]int, 64),
		rank:   make([]int, 64),
		height: make([]int, 64),
		stack:  NewStack(64),
	}
}
