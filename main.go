package main

import (
	chanpprof "go-pprof-test/chan-pprof"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
)

var (
	ExitChan chan os.Signal
)

func init() {
	ExitChan = make(chan os.Signal)
	signal.Notify(ExitChan, os.Interrupt)
}

func main() {
	go http.ListenAndServe("127.0.0.1:9527", nil)
	chanpprof.UncloseChanPPROF(ExitChan, true)
}
