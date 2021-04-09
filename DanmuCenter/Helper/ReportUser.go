package Helper

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/Holynnchen/BiliBan2/DanmuCenter"
)

const reportUrl = "https://api.live.bilibili.com/xlive/web-ucenter/v1/dMReport/Report"

func report(cookie, csrf string, banData *DanmuCenter.BanData) {
	params := url.Values{}
	params.Add("roomid", strconv.Itoa(banData.RoomID))
	params.Add("tuid", strconv.FormatInt(banData.UserID, 10))
	params.Add("msg", banData.Content)
	params.Add("ts", strconv.FormatInt(banData.Timestamp, 10))
	params.Add("reason", "垃圾广告")
	params.Add("reason_id", "3")
	params.Add("csrf", csrf)
	params.Add("csrf_token", csrf)
	req, err := http.NewRequest("POST", reportUrl, strings.NewReader(params.Encode()))
	if err != nil {
		log.Println(err)
		return
	}
	req.Header.Add("content-type", "application/x-www-form-urlencoded")
	req.Header.Add("cookie", cookie)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
}

func Reporter(cookie, csrf string) func(*DanmuCenter.BanData) {
	return func(bd *DanmuCenter.BanData) {
		report(cookie, csrf, bd)
	}
}
