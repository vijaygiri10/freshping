package main

import (
	"fmt"
	"freshping/route"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	addr := "127.0.0.1:9999"

	gressfullstop := make(chan os.Signal)
	signal.Notify(gressfullstop, syscall.SIGQUIT, syscall.SIGKILL, syscall.SIGABRT)

	go func() {
		Signal := <-gressfullstop
		fmt.Printf("\t\n Stoping Backed Server \n\t Caught Signal : %s", Signal.String())
	}()

	route := route.GetGorilaMuxRouter()
	http.Handle("/", route)
	http.ListenAndServe(addr, nil)
}
