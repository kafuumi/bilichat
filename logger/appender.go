package logger

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"time"
)

const bufSize = 1024 //日志文件的缓冲区大小
var l = log.New(os.Stdout, "[Error]logger ", log.LstdFlags)

//Appender 负责将日志内容写入指定的目的地，目的地可以是标准输出，也可以是文件
type Appender interface {
	io.Writer
	WriteMsg(msg string) //写入日志内容
	Close()              //关闭日志输出，因为具体实现会涉及到缓冲区，在程序结束时应该调用此方法，确保日志完全写入
}

//FileAppender 将日志内容写入文件中，并根据指定的最大文件大小，自动创建新文件
type FileAppender struct {
	file    *os.File      //文件句柄
	writer  *bufio.Writer //文件写入缓冲流，写入日志是一个频繁操作，所以增加一个缓冲区，减少io操作
	maxSize int           //单个日志文件最大大小
	nowSize int           //记录写入的日志数据大小
	isClose bool          //输出流是否已经关闭
	lock    sync.Mutex
}

//NewFileAppender 创建一个 FileAppender 对象，并指定单个日志文件最大大小为 maxSize
func NewFileAppender(maxSize int) *FileAppender {
	return &FileAppender{
		file:    nil,
		writer:  nil,
		maxSize: maxSize,
		nowSize: 0,
		isClose: false,
	}
}

// Close 关闭输出流
func (f *FileAppender) Close() {
	f.lock.Lock()
	defer f.lock.Unlock()

	f.isClose = true
	if f.file != nil {
		_ = f.writer.Flush()
		_ = f.file.Close()
	}
}

// Write 写入数据
func (f *FileAppender) Write(p []byte) (int, error) {
	if f.isClose {
		fmt.Println("log file already close!")
		return 0, io.ErrClosedPipe
	}
	f.lock.Lock()
	defer f.lock.Unlock()

	//再次检查
	if f.isClose {
		fmt.Println("log file already close!")
		return 0, io.ErrClosedPipe
	}
	//写入的日志数据达到指定值，创建新的日志文件
	if f.file == nil || f.nowSize >= f.maxSize {
		f.logFile()
	}
	wn, err := f.writer.Write(p)
	if err != nil {
		l.Printf("write log msg fail! %v\n", err)
	}
	f.nowSize += wn
	return wn, err
}

// WriteMsg 向日志文件中写入日志信息 msg
func (f *FileAppender) WriteMsg(msg string) {
	_, _ = f.Write([]byte(msg))
}

//创建新的日志文件
func (f *FileAppender) logFile() {
	//文件名格式：年月日时分秒
	name := fmt.Sprintf("logs%c%s.log", os.PathSeparator,
		time.Now().Format("20060102150405"))
	file, err := os.Create(name)
	//logs 目录不存在，则创建
	if err != nil && os.IsNotExist(err) {
		err = os.Mkdir("./logs", os.ModePerm)
		if err != nil {
			l.Printf("create log dir fail! %v\n", err)
			return
		}
		file, err = os.Create(name)
	}
	if err != nil {
		l.Printf("create log file fail! %v\n", err)
		return
	}
	if f.file != nil {
		//刷新缓冲区并关闭
		_ = f.writer.Flush()
		_ = f.file.Close()
	}
	f.file = file
	f.writer = bufio.NewWriterSize(file, bufSize)
	f.nowSize = 0
}

//ConsoleAppender 向标准输出写日志
type ConsoleAppender struct{}

func NewConsoleAppender() *ConsoleAppender {
	return &ConsoleAppender{}
}

func (c *ConsoleAppender) Write(p []byte) (int, error) {
	return os.Stdout.Write(p)
}

func (c *ConsoleAppender) WriteMsg(msg string) {
	_, _ = os.Stdout.WriteString(msg)
}

func (c *ConsoleAppender) Close() {
	//ignore
}
