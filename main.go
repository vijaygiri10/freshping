package main

import (
	"fmt"
	"freshping/monitor"
	"freshping/route"
	"net/http"
	"net/http/pprof"
	"os"
	"os/signal"
	"syscall"

	//	"github.com/pkg/profile"

	_ "net/http/pprof"
)

func main() {
	// CPU profiling by default
	//defer profile.Start().Stop()
	// Memory profiling by default
	//defer profile.Start(profile.MemProfile).Stop()
	addr := "127.0.0.1:9999"

	gressfullstop := make(chan os.Signal)

	go func() {
		Signal := <-gressfullstop
		signal.Notify(gressfullstop, syscall.SIGQUIT, syscall.SIGKILL, syscall.SIGABRT)
		fmt.Printf("\t\n Stoping Backed Server \n\t Caught Signal : %s", Signal.String())
	}()
	//context.Background()
	go monitor.StartURLParser()

	route := route.GetGorilaMuxRouter()
	// This will serve files under http://IP:PORT/assets/<filename>
	route.PathPrefix("/").Handler(http.FileServer(http.Dir(".")))

	// Register pprof handlers
	route.HandleFunc("/debug/pprof/", pprof.Index)
	route.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	route.HandleFunc("/debug/pprof/profile", pprof.Profile)
	route.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	route.HandleFunc("/debug/pprof/trace", pprof.Trace)

	// Register All Mux Route to Http Handler
	http.Handle("/", route)
	fmt.Println(http.ListenAndServe(addr, nil))
}
