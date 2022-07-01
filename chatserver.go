package bilichat

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Hami-Lemon/bilichat/logger"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

const (
	chanBufSize = 64
)

type ChatServer struct {
	room   Room //对应的直播间
	host   string
	port   int
	token  string
	conn   *websocket.Conn //websocket链接
	msgCh  chan []byte     //收到的数据包，已经过解压、拆包
	logger *logger.Logger
}

// getter

func (c *ChatServer) Room() Room {
	return c.room
}

func (c *ChatServer) Host() string {
	return c.host
}

func (c *ChatServer) Port() int {
	return c.port
}

// Connect 连接弹幕服务器
func (c *ChatServer) Connect() error {
	dialer := websocket.Dialer{
		HandshakeTimeout: 5 * time.Second,
		ReadBufferSize:   4 * 1024,
		WriteBufferSize:  512,
	}
	u := fmt.Sprintf("wss://%s:%d/sub", c.host, c.port)
	//请求头
	h := http.Header{}
	for name, value := range reqHeader {
		h.Add(name, value)
	}
	h.Add("Origin", "https://live.bilibili.com")
	h.Add("Cache-Control", "no-cache")

	conn, _, err := dialer.Dial(u, h)
	if err != nil {
		return err
	}
	c.conn = conn

	//进房验证
	err = c.verify()
	if err != nil {
		return err
	}
	//发送心跳包
	go c.heartbeat()
	//读取数据包并处理
	go c.handle()
	return nil
}

func (c *ChatServer) Disconnect() {
	_ = c.conn.Close()
}

//流水线模型 handle ==> unpackMsg ==> ReceiveMsg
func (c *ChatServer) handle() {
	unpackCh := make(chan []byte, chanBufSize)
	go c.unpackMsg(unpackCh)
	for {
		_, buf, err := c.conn.ReadMessage()
		if err != nil {
			if buf == nil {
				close(unpackCh)
				break
			} else {
				c.logger.Error("读取websocket消息失败, %v", err)
			}
		}
		select {
		case unpackCh <- buf:
		default:
			c.logger.Warn("读取消息 ==> 解包数据包，阻塞！")
		}
	}
}

func (c *ChatServer) unpackMsg(in <-chan []byte) {
	for {
		msg, ok := <-in
		if !ok {
			close(c.msgCh)
			return
		}
		for _, packet := range unpack(msg) {
			select {
			case c.msgCh <- packet:
			default:
				c.logger.Warn("unpackMsg ==> ReceiveMsg，阻塞！")
			}
		}
	}
}

// ReceiveMsg 解析消息,将获取到的消息写入到 out 中
func (c *ChatServer) ReceiveMsg(out chan<- Message) {
	for {
		srcMsg, ok := <-c.msgCh
		if !ok {
			close(out)
			return
		}
		msg := parseMsg(srcMsg)
		if msg != nil {
			select {
			case out <- msg:
			default:
				c.logger.Warn("ReceiveMsg ==> out，阻塞！type: %s", msg.MsgType())
			}
		}
	}
}

//发送验证消息
func (c *ChatServer) verify() error {
	verifyMsg := map[string]interface{}{
		"platform": "web",
		"protover": 3,
		"uid":      0,
		"roomid":   c.room.Rid,
		"type":     2,
		"key":      c.token,
	}
	body, _ := json.Marshal(verifyMsg)
	err := c.conn.WriteMessage(websocket.BinaryMessage, pack(verPlain, opEnterRoom, body))
	if err != nil {
		c.logger.Error("发送验证信息失败！%v", err)
		return err
	}

	//读取服务端回传的消息，判断是否成功进入直播间，如果进入失败，服务端会断开连接
	_, buf, err := c.conn.ReadMessage()
	if err != nil {
		c.logger.Error("读取验证信息回响失败,进入失败！%v", err)
		return err
	}
	op, body := unpackPacket(buf)
	if op != opEnterRoomReply {
		return errors.New(string(body))
	}
	return nil
}

//周期性发送心跳包，间隔为30秒
func (c *ChatServer) heartbeat() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()
	//心跳包内容，可以是任意内容，空数据也可以
	heartbeatPacket := []byte{0x52, 0x33, 0x52, 0x33, 0x52, 0x33, 0x52, 0x33, 0x52, 0x33, 0x52, 0x33, 0x52, 0x33}
	err := c.conn.WriteMessage(websocket.BinaryMessage, pack(verInt, opHeartbeat, heartbeatPacket))
	if err != nil {
		c.logger.Error("发送心跳包失败！%v", err)
		return
	}
	for range ticker.C {
		err = c.conn.WriteMessage(websocket.BinaryMessage, pack(verInt, opHeartbeat, heartbeatPacket))
		if err != nil {
			c.logger.Error("发送心跳包失败！%v", err)
			break
		}
	}
}

// GetChatServer 获取弹幕服务器地址
func GetChatServer(roomId int) (*ChatServer, error) {
	b := NewClient()
	r, err := b.RoomInfo(roomId)
	if err != nil {
		return nil, err
	}
	v := url.Values{}
	v.Add("id", strconv.Itoa(r.Rid))
	v.Add("type", "0")
	u := "https://api.live.bilibili.com/xlive/web-room/v1/index/getDanmuInfo?" + v.Encode()

	resp, err := b.get(u)
	if err != nil {
		return nil, err
	}
	c := &ChatServer{
		room:   r,
		msgCh:  make(chan []byte, chanBufSize),
		logger: logger.New("chat-"+r.Liver.Uname, logLevel, logAppender),
	}
	data := resp.Get("data")
	c.token = data.Get("token").String()

	host := data.Get("host_list.0")
	c.host = host.Get("host").String()
	c.port = int(host.Get("wss_port").Int())
	return c, nil
}
