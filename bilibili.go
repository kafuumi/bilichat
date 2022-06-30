package bilichat

import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/andybalholm/brotli"
	"github.com/tidwall/gjson"
)

var (
	ErrVerify = errors.New("verify fail") //进入直播间失败
	reqHeader = map[string]string{
		"Accept-Language": "zh-CN,zh;q=0.9",
		"Accept-Encoding": "gzip, deflate, br",
		"User-Agent":      "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/102.0.0.0 Safari/537.36",
	}
)

// Liver 主播
type Liver struct {
	Uid   int64  //uid
	Uname string //昵称
}

// Room 直播间
type Room struct {
	Liver  Liver  //该直播所对应的主播
	Id     int    //房间号
	Rid    int    //真实房间号
	Title  string //直播间标题
	IsLive bool   //是否正在直播
}

type BiliClient struct {
	client *http.Client
}

func NewClient() *BiliClient {
	return &BiliClient{client: &http.Client{}}
}

func handleResp(resp *http.Response, err error) (*gjson.Result, error) {
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	//请求失败
	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}

	//解压
	var reader io.Reader = resp.Body
	contentEncoding := resp.Header.Get("Content-Encoding")
	switch contentEncoding {
	case "gzip":
		reader, err = gzip.NewReader(resp.Body)
		if err != nil {
			return nil, err
		}
	case "deflate":
		reader = flate.NewReader(resp.Body)
	case "br":
		reader = brotli.NewReader(resp.Body)
	default:
		fmt.Printf(contentEncoding)
	}
	buf := &bytes.Buffer{}
	_, err = io.Copy(buf, reader)
	if err != nil {
		return nil, err
	}
	r := gjson.ParseBytes(buf.Bytes())
	code := r.Get("code").Int()
	if code != 0 {
		return nil, errors.New(r.Get("message").String())
	}
	return &r, nil
}

func (b *BiliClient) get(u string) (*gjson.Result, error) {
	req, err := http.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}
	//设置请求头
	for name, value := range reqHeader {
		req.Header.Add(name, value)
	}
	return handleResp(b.client.Do(req))
}

// LiverInfo 获取主播信息
func (b *BiliClient) LiverInfo(uid int64) (Liver, error) {
	u := "https://api.bilibili.com/x/space/acc/info?mid=" + strconv.FormatInt(uid, 10)
	resp, err := b.get(u)
	if err != nil {
		return Liver{}, err
	}
	liver := Liver{
		Uid: uid,
	}
	liver.Uname = resp.Get("data.name").String()
	return liver, nil
}

// RoomInfo 获取直播间信息
func (b *BiliClient) RoomInfo(id int) (Room, error) {
	u := "https://api.live.bilibili.com/room/v1/Room/get_info?room_id=" + strconv.Itoa(id)
	resp, err := b.get(u)
	if err != nil {
		return Room{}, err
	}
	room := Room{
		Id: id,
	}
	data := resp.Get("data")
	room.Rid = int(data.Get("room_id").Int())
	room.IsLive = data.Get("live_status").Int() == 1
	room.Title = data.Get("title").String()

	liverUid := data.Get("uid").Int()
	liver, err := b.LiverInfo(liverUid)
	if err != nil {
		return Room{}, err
	}
	room.Liver = liver
	return room, nil
}
