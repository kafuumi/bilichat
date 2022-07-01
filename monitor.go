package bilichat

import (
	"io"
	"sync"
	"time"

	"github.com/Hami-Lemon/bilichat/logger"
	"gopkg.in/yaml.v3"
)

const (
	logFileSize    = 1024 * 512
	danMuMsgBufCap = 256
)

var (
	logLevel                    = logger.Debug
	logAppender logger.Appender = logger.NewConsoleAppender()
	mainLogger                  = logger.New("main", logger.Info, logger.NewConsoleAppender())
)

// Config 配置信息
type Config struct {
	Rooms    []int `yaml:"rooms"` //监控的房间号
	Database struct {
		Name     string `yaml:"name"`     //选择的数据库，mysql或mongodb
		User     string `yaml:"user"`     //用户名
		Password string `yaml:"password"` //密码
		Address  string `yaml:"address"`  //数据库地址
		Port     int    `yaml:"port"`     //端口号
		Dbname   string `yaml:"dbname"`   //数据库名称
	} `yaml:"database"`
	Log struct {
		Level    string `yaml:"level"` //日志级别
		Appender string `yaml:"appender"`
	} `yaml:"log"`
}

// ReadConfig 读取配置，需要是 yaml 格式的输入流
func ReadConfig(reader io.Reader) (Config, error) {
	var con Config

	err := yaml.NewDecoder(reader).Decode(&con)
	if err != nil {
		mainLogger.Error("读取配置信息失败，%v", err)
		return Config{}, err
	}
	mainLogger.Info("监控的房间：%v", con.Rooms)
	mainLogger.Info("database: user=%s, url=%s:%d, dbname=%s",
		con.Database.User, con.Database.Address, con.Database.Port, con.Database.Dbname)
	mainLogger.Info("logger: level=%s, appender=%s", con.Log.Level, con.Log.Appender)
	return con, nil
}

type Monitor struct {
	servers []*ChatServer
	group   sync.WaitGroup
	logger  *logger.Logger
	dao     dao
}

func NewMonitor(c Config) *Monitor {
	switch c.Log.Level {
	case "Debug", "debug":
		logLevel = logger.Debug
	case "Info", "info":
		logLevel = logger.Info
	case "Warn", "warn":
		logLevel = logger.Warn
	case "Error", "error":
		logLevel = logger.Error
	default:
		mainLogger.Warn("read log level fail, default level: Debug")
	}
	switch c.Log.Appender {
	case "file":
		logAppender = logger.NewFileAppender(logFileSize)
	case "console":
	default:
		mainLogger.Warn("read log append fail, default appender: console")
	}

	m := &Monitor{
		servers: make([]*ChatServer, 0),
		logger:  logger.New("monitor", logLevel, logAppender),
	}
	for _, room := range c.Rooms {
		chatServer, err := GetChatServer(room)
		if err != nil {
			mainLogger.Error("获取弹幕服务器失败：roomId=%d, %v", room, err)
			return nil
		}
		m.servers = append(m.servers, chatServer)
	}

	database := c.Database
	var err error
	m.dao, err = newMysqlDao(database.User, database.Password,
		database.Address, database.Port, database.Dbname)
	if err != nil {
		mainLogger.Error("连接数据库失败：%v", err)
		return nil
	}
	return m
}

func work(chat *ChatServer, d dao, group *sync.WaitGroup) {
	mainLogger.Info("监控【%s】的直播间，开播=%t, roomId=%d, title=%s",
		chat.room.Liver.Uname, chat.room.IsLive, chat.room.Id, chat.room.Title)
	out := make(chan Message, chanBufSize*2) //两倍缓冲
	go chat.ReceiveMsg(out)

	l := chat.logger
	ifInsertError := func(err error) {
		if err != nil {
			l.Error("插入数据失败：%v", err)
		}
	}

	r := &(chat.room)
	danMuMsgBuf := newBuffer[*DanMuMessage](danMuMsgBufCap, time.Minute, true,
		func(items []*DanMuMessage) {
			ifInsertError(d.insertDanMuMsg(*r, items))
		})
	for {
		msg, ok := <-out
		if !ok {
			danMuMsgBuf.MustFlush()
			danMuMsgBuf.Free()
			group.Done()
			return
		}
		switch m := msg.(type) {
		case *DanMuMessage:
			danMuMsgBuf.Put(m)
		case *SuperChatMessage:
			ifInsertError(d.insertScMsg(*r, m))
		case *GiftMessage:
			ifInsertError(d.insertGiftMsg(*r, m))
		case *GuardMessage:
			ifInsertError(d.insertGuardMsg(*r, m))
		case *EntryMessage:
			ifInsertError(d.insertEntryMsg(*r, m))
		case *RoomFansMessage:
			ifInsertError(d.insertFansMsg(*r, m))
		case *RankCountMessage:
			ifInsertError(d.insertRankCountMsg(*r, m))
		case *HotRankMessage:
			ifInsertError(d.insertHotRankMsg(*r, m))
		case *LiveStatusMessage:
			r.IsLive = m.Status
			if r.IsLive {
				l.Info("[%s] 开播", r.Liver.Uname)
			} else {
				l.Info("[%s] 下播", r.Liver.Uname)
			}
		case *RoomChangeMessage:
			ifInsertError(d.insertRoomChangeMsg(*r, m))
			r.Title = m.Title
		case *WatchedChangeMessage:
			ifInsertError(d.insertWatchedChangeMsg(*r, m))
		}
	}
}

func (m *Monitor) Start() {
	for _, c := range m.servers {
		err := c.Connect()
		if err != nil {
			mainLogger.Error("连接直播间失败，roomId=%d, %v", c.room.Id, err)
			return
		}
		m.group.Add(1)
		go work(c, m.dao, &(m.group))
		time.Sleep(time.Second)
	}
}

func (m *Monitor) Stop() {
	mainLogger.Info("程序退出...")
	for _, c := range m.servers {
		c.Disconnect()
	}
	logAppender.Close()
	m.group.Wait()
}
