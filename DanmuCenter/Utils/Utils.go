package Utils

import (
	"io/ioutil"
	"net/http"
	"strconv"
	"unicode/utf8"
	"unsafe"

	"github.com/goccy/go-json"
)

type Empty struct{}

type TopRoomResponse struct {
	Code    int    `json:"code"`
	Msg     string `json:"msg"`
	Message string `json:"message"`
	Data    []Data `json:"data"`
}

type Data struct {
	Roomid int    `json:"roomid"`
	Uname  string `json:"uname"`
}

func GetTopRoom(max int) ([]int, error) {
	var roomIDs []int
	page := 1
	for max > 0 {
		resp, err := http.Get("https://api.live.bilibili.com/room/v1/Area/getListByAreaID?areaId=0&sort=online&pageSize=200&page=" + strconv.Itoa(page))
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		var result = new(TopRoomResponse)
		if err := json.Unmarshal(body, result); err != nil {
			return nil, err
		}
		for _, room := range result.Data {
			if max <= 0 {
				break
			}
			roomIDs = append(roomIDs, room.Roomid)
			max--
		}
		page++
	}
	return roomIDs, nil
}

type HotRoomResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	TTL     int    `json:"ttl"`
	Data    Data2  `json:"data"`
}
type Data2 struct {
	List []List `json:"list"`
}
type List struct {
	Uname  string `json:"uname"`
	RoomID int    `json:"room_id"`
}

func GetTop50HotRoom() ([]int, error) {
	resp, err := http.Get("https://api.live.bilibili.com/xlive/general-interface/v1/rank/getHotRank?room_id=1&ruid=1&area_id=0&page_size=50&source=1")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var result = new(HotRoomResponse)
	if err := json.Unmarshal(body, result); err != nil {
		return nil, err
	}
	var roomIDs []int
	for _, room := range result.Data.List {
		roomIDs = append(roomIDs, room.RoomID)
	}
	return roomIDs, nil
}

// get the number in a ant not in b
func DiffrenceIntArray(a []int, b []int) []int {
	mark := map[int]Empty{}
	result := make([]int, 0)
	for _, num := range b {
		mark[num] = Empty{}
	}
	for _, num := range a {
		if _, ok := mark[num]; !ok {
			result = append(result, num)
		}
	}
	return result
}

// make int array to map[int]Empty
func IntArrayToMap(a []int) map[int]Empty {
	mark := map[int]Empty{}
	for _, num := range a {
		mark[num] = Empty{}
	}
	return mark
}

func GetSimilarity(a, b string) float32 {
	// aLen, bLen, distance := utf8.RuneCountInString(a), utf8.RuneCountInString(b), GetEditDistance(a, b)
	aLen, bLen, distance := utf8.RuneCountInString(a), utf8.RuneCountInString(b), Distance(a, b) // 切换至不安全的编辑距离

	return maxFloat32(1 - float32(distance)/maxFloat32(float32(aLen), float32(bLen)))
}

func GetEditDistance(a, b string) int {
	if a == b {
		return 0
	}
	ar, br := []rune(a), []rune(b)
	aLen, bLen := len(ar), len(br)
	if aLen == 0 {
		return bLen
	}
	if bLen == 0 {
		return aLen
	}
	if aLen > bLen {
		aLen, bLen, ar, br = bLen, aLen, br, ar
	}
	distance := make([]int, aLen+1)
	for i := 0; i <= aLen; i++ {
		distance[i] = i
	}
	_ = distance[aLen]
	for i := 1; i <= bLen; i++ {
		left, up := i, i-1
		for j := 1; j <= aLen; j++ {
			if br[i-1] != ar[j-1] {
				up++
			}
			left, up = min(distance[j]+1, left+1, up), distance[j]
			distance[j] = left
		}
	}
	return distance[aLen]
}

func min(a int, args ...int) int {
	for _, b := range args {
		if b < a {
			a = b
		}
	}
	return a
}

func max(a int, args ...int) int {
	for _, b := range args {
		if b > a {
			a = b
		}
	}
	return a
}

func maxFloat32(a float32, args ...float32) float32 {
	for _, b := range args {
		if b > a {
			a = b
		}
	}
	return a
}

type Stack struct {
	size      int
	max       int
	container []interface{}
}

func NewStack(max int) *Stack {
	return &Stack{
		size:      0,
		max:       max,
		container: make([]interface{}, max),
	}
}

func (s *Stack) MustPush(item interface{}) {
	s.size++
	s.container[s.size-1] = item
}
func (s *Stack) Push(item interface{}) bool {
	if s.size >= s.max {
		return false
	}
	s.MustPush(item)
	return true
}

func (s *Stack) MustPeek() interface{} {
	return s.container[s.size-1]
}

func (s *Stack) Peek() (interface{}, bool) {
	if s.size <= 0 {
		return nil, false
	}
	return s.MustPeek(), true
}

func (s *Stack) MustPop() interface{} {
	s.size--
	return s.container[s.size]
}

func (s *Stack) Pop() (interface{}, bool) {
	if s.size <= 0 {
		return nil, false
	}
	return s.MustPop(), true
}

func (s *Stack) Size() int {
	return s.size
}

func (s *Stack) Empty() {
	s.size = 0
}

func str2bytes(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	b := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&b))
}

func bytes2str(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
