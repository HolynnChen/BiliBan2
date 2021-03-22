package Utils

import (
	"io/ioutil"
	"net/http"
	"strconv"
	"unicode/utf8"

	jsoniter "github.com/json-iterator/go"
)

type Empty struct{}

var json = jsoniter.ConfigCompatibleWithStandardLibrary

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
	resp, err := http.Get("https://api.live.bilibili.com/room/v1/Area/getListByAreaID?areaId=0&sort=online&pageSize=" + strconv.Itoa(max) + "&page=1")
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
	var roomIDs []int
	for _, room := range result.Data {
		roomIDs = append(roomIDs, room.Roomid)
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
	aLen, bLen, distance := utf8.RuneCountInString(a), utf8.RuneCountInString(b), GetEditDistance(a, b)
	return maxFloat32(1 - float32(distance)/maxFloat32(float32(aLen), float32(bLen)))
}

func GetEditDistance(a, b string) int {
	ar, br := []rune(a), []rune(b)
	aLen, bLen := len(ar), len(br)
	distance := make([]int, aLen+1)
	for i := 0; i <= aLen; i++ {
		distance[i] = i
	}
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

func maxFloat32(a float32, args ...float32) float32 {
	for _, b := range args {
		if b > a {
			a = b
		}
	}
	return a
}
