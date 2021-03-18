package DanmuCenter

import (
	"fmt"
	"testing"
)

func TestGetTopRoom(t *testing.T) {
	fmt.Println(GetTopRoom(100))
}

func TestTry(t *testing.T) {
	text1 := ReplaceSimilarAndNumberRune("y22c１▪∁ＯM→全是倮躰的校花们哟e")
	text2 := ReplaceSimilarAndNumberRune("搜<y9h2.∁n>带好纸巾去鲁个痛筷e")
	fmt.Println(text1)
	fmt.Println(text2)
	fmt.Println(GetSimilarity(text1, text2))
}

func TestMap(t *testing.T) {
	m := make(map[rune]int)
	fmt.Println(m['t'])
}

func TestGetEditDistance(t *testing.T) {
	fmt.Println(GetEditDistance("abc", "cba"))
}
