package DanmuCenter

import (
	"fmt"
	"testing"
)

func TestGetTopRoom(t *testing.T) {
	fmt.Println(GetTopRoom(100))
}

func TestGetHotRoom(t *testing.T) {
	fmt.Println(GetTop50HotRoom())
}

func TestTry(t *testing.T) {
	s := `\花丸/\花丸/\花丸/\花丸/\花丸/\花丸/\花丸/\花丸/`
	fmt.Println(CompressRepeatGroup(3)(s))
}

func TestMap(t *testing.T) {
	m := make(map[rune]int)
	fmt.Println(m['t'])
}

func TestGetEditDistance(t *testing.T) {
	fmt.Println(GetEditDistance("abc", "cba"))
}
