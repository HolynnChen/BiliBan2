package main

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/Holynnchen/BiliBan2/DanmuCenter"
	"github.com/goccy/go-json"
)

type SyncData struct {
	UserID    int64  `json:"UserId"`
	UserName  string `json:"UserName"`
	RoomID    int    `json:"RoomId"`
	Content   string `json:"Content"`
	CT        string `json:"ct"`
	TimeStamp int64  `json:"TimeStamp"`
	Reason    string `json:"Reason"`
}

func syncBan(banData *DanmuCenter.BanData) {
	jsonData, _ := json.Marshal(SyncData{
		UserID:    banData.UserID,
		UserName:  banData.UserName,
		RoomID:    banData.RoomID,
		Content:   banData.Content,
		CT:        banData.CT,
		TimeStamp: banData.Timestamp,
		Reason:    banData.Reason,
	})
	resp, err := http.DefaultClient.Post(env["sync_url"].(string), "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println(err)
		return
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
}

type QueryBanRequest struct {
	// CursorID int64 `json:"CursorId"`
}

type QueryBanResponse struct {
	Code int            `json:"code"`
	Msg  string         `json:"msg"`
	Data []QueryBanData `json:"data"`
}

type QueryBanData struct {
	CursorID int64 `json:"CursorId"`
	Danmaku  struct {
		Comment string `json:"Comment"`
	} `json:"Danmaku"`
}

func queryBan(CursorID int64) ([]QueryBanData, error) {
	jsonData, _ := json.Marshal(QueryBanRequest{
		// CursorID: CursorID,
	})
	resp, err := http.DefaultClient.Post(env["query_url"].(string), "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var result QueryBanResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}
	if result.Code != 0 {
		return nil, errors.New(result.Msg)
	}
	for index := 0; index < len(result.Data); index++ {
		if result.Data[index].CursorID > CursorID {
			return result.Data[index:], nil
		}
	}
	return []QueryBanData{}, nil
}

type ProxyResponse struct {
	CheckCount int    `json:"check_count"`
	FailCount  int    `json:"fail_count"`
	LastStatus int    `json:"last_status"`
	LastTime   string `json:"last_time"`
	Proxy      string `json:"proxy"`
	Region     string `json:"region"`
	Source     string `json:"source"`
	Type       string `json:"type"`
}

var reqTimeLimit = time.Now().Unix()

func maxInt(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

func getProxy() func(*http.Request) (*url.URL, error) {
	if reqTimeLimit-time.Now().Unix() <= 60 {
		reqTimeLimit = maxInt(time.Now().Unix(), reqTimeLimit) + 3
		return nil // 低频下用自己ip
	}
	resp, err := http.DefaultClient.Get(env["proxy_url"].(string))
	if err != nil {
		log.Println(err)
		return nil
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return nil
	}
	var result ProxyResponse
	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Println(err)
		return nil
	}
	proxyUrl, err := url.Parse("http://" + result.Proxy)
	if err != nil {
		log.Println(err)
		return nil
	}
	return http.ProxyURL(proxyUrl)
}
