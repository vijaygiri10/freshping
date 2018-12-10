package main

import (
	"fmt"
	"freshping/monitor"
	"freshping/route"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	addr := "127.0.0.1:9999"

	gressfullstop := make(chan os.Signal)

	go func() {
		Signal := <-gressfullstop
		signal.Notify(gressfullstop, syscall.SIGQUIT, syscall.SIGKILL, syscall.SIGABRT)
		fmt.Printf("\t\n Stoping Backed Server \n\t Caught Signal : %s", Signal.String())
	}()

	go monitor.StartURLParser()

	route := route.GetGorilaMuxRouter()
	// This will serve files under http://IP:PORT/assets/<filename>
	route.PathPrefix("/").Handler(http.FileServer(http.Dir(".")))

	http.Handle("/", route)
	fmt.Println(http.ListenAndServe(addr, nil))
}
