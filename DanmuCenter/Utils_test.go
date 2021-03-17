package DanmuCenter

import (
	"fmt"
	"testing"
)

func TestGetTopRoom(t *testing.T) {
	fmt.Println(GetTopRoom(100))
}

func TestTry(t *testing.T) {
	fmt.Println(GetSimilarity(ReplaceSimilarAndNumberRune("片姐姐就是我，照片！"), ReplaceSimilarAndNumberRune("上船解锁舰长群60多样好听好看的浮力，点播免费，还可以加学姐私人QQ，")))
	fmt.Println(GetSimilarity(ReplaceSimilarAndNumberRune("  勉的钱输 Β t y  。Pw"), ReplaceSimilarAndNumberRune("上船解锁舰长群60多样好听好看的浮力，点播免费，还可以加学姐私人QQ，")))
}

func TestMap(t *testing.T) {
	m := make(map[rune]int)
	fmt.Println(m['t'])
}

func TestGetEditDistance(t *testing.T) {
	fmt.Println(GetEditDistance("abc", "cba"))
}
