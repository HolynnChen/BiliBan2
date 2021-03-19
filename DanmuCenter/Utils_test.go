package DanmuCenter

import (
	"fmt"
	"testing"
)

func TestGetTopRoom(t *testing.T) {
	fmt.Println(GetTopRoom(100))
}

func TestTry(t *testing.T) {
	s := "我可以吗? 喜欢你 我可以吗! 喜欢你"
	fmt.Println(CompressRepeatGroup(3)(s))
}

func TestMap(t *testing.T) {
	m := make(map[rune]int)
	fmt.Println(m['t'])
}

func TestGetEditDistance(t *testing.T) {
	fmt.Println(GetEditDistance("abc", "cba"))
}
