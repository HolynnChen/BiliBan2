package main

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"sync"
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
	CursorID int64 `json:"CursorId"`
}

type QueryBanResponse struct {
	Code int            `json:"code"`
	Msg  string         `json:"msg"`
	Data []QueryBanData `json:"data"`
}

type QueryBanData struct {
	CursorID int64  `json:"CursorId"`
	Reason   string `json:"Reason"`
	Danmaku  struct {
		Comment string `json:"Comment"`
	} `json:"Danmaku"`
}

func queryBan(CursorID int64) ([]QueryBanData, error) {
	jsonData, _ := json.Marshal(QueryBanRequest{
		CursorID: CursorID,
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

func httpGet(url string, out interface{}) error {
	resp, err := http.DefaultClient.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, out)
	if err != nil {
		return err
	}
	return nil
}

func getProxy() func(*http.Request) (*url.URL, error) {
	if reqTimeLimit-time.Now().Unix() <= 60 {
		reqTimeLimit = maxInt(time.Now().Unix(), reqTimeLimit) + 3
		return nil // 低频下用自己ip
	}
	var result ProxyResponse
	if err := httpGet(env["proxy_url"].(string), &result); err != nil {
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

type proxyPool struct {
	proxys []string
}

type getProxyResp struct {
	Proxy string `json:"proxy"`
}

func (p *proxyPool) Get() func(*http.Request) (*url.URL, error) {
	proxyUrl, _ := url.Parse("http://" + p.proxys[rand.Intn(len(p.proxys))])
	return http.ProxyURL(proxyUrl)
}

func (p *proxyPool) Sync() {
	var resp = make([]getProxyResp, 0)
	if err := httpGet(env["proxy_all_url"].(string), &resp); err != nil {
		log.Println(err)
		return
	}
	tmp := p.proxys
	for _, i := range resp {
		tmp = append(tmp, i.Proxy)
	}
	result := make([]string, 0)
	var wg sync.WaitGroup
	for _, i := range tmp {
		wg.Add(1)
		go func(proxy string) {
			if p.Check(proxy) {
				tmp = append(tmp, proxy)
			}
			defer wg.Done()
		}(i)
	}
	wg.Wait()
	if len(result) > 0 {
		result = result[:100]
	}
	p.proxys = result
}

type biliIpResp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (p *proxyPool) Check(proxy string) bool {
	var rsp biliIpResp
	proxyUrl, _ := url.Parse("http://" + proxy)
	client := &http.Client{Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		Proxy:           http.ProxyURL(proxyUrl),
	}, Timeout: 10 * time.Second}
	resp, err := client.Get("https://api.live.bilibili.com/xlive/web-room/v1/index/getIpInfo")
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false
	}
	err = json.Unmarshal(body, &rsp)
	if err != nil {
		return false
	}
	if rsp.Code != 0 {
		return false
	}
	return true
}
