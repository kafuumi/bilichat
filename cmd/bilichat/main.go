package main

import (
	"fmt"
	"github.com/Hami-Lemon/bilichat"
	"os"
	"os/signal"
	"strings"
)

func main() {
	reader := strings.NewReader(`rooms:
  - 6
database:
  user: 'carol'
  password: 'mysqlcarol'
  address: 'localhost'
  port: 3306
  dbname: 'live_info'
log:
  level: 'debug'
  appender: 'console'
`)
	con, _ := bilichat.ReadConfig(reader)
	monitor := bilichat.NewMonitor(con)
	monitor.Start()
	defer monitor.Stop()
	ch := make(chan os.Signal)
	signal.Notify(ch, os.Kill, os.Interrupt)
	<-ch
	fmt.Printf("exit\n")
}
