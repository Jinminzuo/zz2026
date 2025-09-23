package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

// JsonResult 返回结构
type JsonResult struct {
	Code     int         `json:"code"`
	ErrorMsg string      `json:"errorMsg,omitempty"`
	Data     interface{} `json:"data"`
}

// HelloWorldHandler 返回 hello world
func HelloWorldHandler(w http.ResponseWriter, r *http.Request) {
	// 仅支持 GET
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(JsonResult{
			Code:     -1,
			ErrorMsg: "请求方法不支持",
		})
		return
	}

	// 构建返回结果
	result := JsonResult{
		Code: 0,
		Data: "helloworld",
	}

	// 设置响应头
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

// 微信开放接口服务前缀（云托管内部调用）
const openAPIHost = "https://api.weixin.qq.com" // 云托管内部直接调用原接口即可

// ---------------------- 标签列表 ----------------------
type Tag struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Count int    `json:"count"`
}

type TagsResp struct {
	Tags []Tag `json:"tags"`
}

// handler 获取标签列表
func TagsHandler(w http.ResponseWriter, r *http.Request) {
	url := fmt.Sprintf("%s/cgi-bin/tags/get", openAPIHost)

	resp, err := http.Get(url)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	// 可直接转发给前端
	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}

// ---------------------- 按标签群发 ----------------------
type Filter struct {
	IsToAll bool `json:"is_to_all"`
	TagID   int  `json:"tag_id"`
}

type MpNews struct {
	MediaID string `json:"media_id"`
}

type MassSendReq struct {
	Filter  Filter `json:"filter"`
	MpNews  MpNews `json:"mpnews"`
	MsgType string `json:"msgtype"`
}

func SendHandler(w http.ResponseWriter, r *http.Request) {
	// 从前端接收 JSON
	var req MassSendReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	url := fmt.Sprintf("%s/cgi-bin/message/mass/sendall", openAPIHost)
	payload, _ := json.Marshal(req)

	resp, err := http.Post(url, "application/json;charset=utf-8", bytes.NewReader(payload))
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)

	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.Write(body)
}
