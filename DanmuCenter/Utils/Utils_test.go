package Utils

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
	s := `\秋絵/♫\秋絵/♫mp\秋絵/♫\秋絵/♫\秋絵/♫qoc\秋絵/♫\秋絵/♫usm\秋絵/♫\秋絵/♫\秋絵/`
	fmt.Println(s)
	fmt.Println(CompressRepeatGroup(2)(s))
	fmt.Println(CompressRepeatGroup2(2)(s))
}

func BenchmarkCompress1(b *testing.B) {
	s := `\秋絵/♫\秋絵/♫mp\秋絵/♫\秋絵/♫\秋絵/♫qoc\秋絵/♫\秋絵/♫usm\秋絵/♫\秋絵/♫\秋絵/`
	compressor := CompressRepeatGroup(2)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		compressor(s)
	}
}

func BenchmarkCompress2(b *testing.B) {
	s := `\秋絵/♫\秋絵/♫mp\秋絵/♫\秋絵/♫\秋絵/♫qoc\秋絵/♫\秋絵/♫usm\秋絵/♫\秋絵/♫\秋絵/`
	compressor := CompressRepeatGroup2(2)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		compressor(s)
	}
}

func TestMap(t *testing.T) {
	m := make(map[rune]int)
	fmt.Println(m['t'])
}

func TestGetEditDistance(t *testing.T) {
	fmt.Println(GetEditDistance("##########星⼩智实属⽜逼", "............"))
	fmt.Println(GetSimilarity("##########星⼩智实属⽜逼", "............"))
}
