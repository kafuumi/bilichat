package main

import (
	"os"
	"os/signal"

	"github.com/Hami-Lemon/bilichat"
)

func main() {
	//读取设置文件
	configReader, err := os.Open("./setting.yaml")
	if err != nil {
		panic(err)
	}
	con, err := bilichat.ReadConfig(configReader)
	if err != nil {
		panic(err)
	}
	_ = configReader.Close()
	monitor := bilichat.NewMonitor(con)
	monitor.Start()
	defer monitor.Stop()

	ch := make(chan os.Signal)
	signal.Notify(ch, os.Kill, os.Interrupt)
	<-ch
}
