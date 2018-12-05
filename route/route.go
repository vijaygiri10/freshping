package route

import (
	"fmt"
	"freshping/middleware"
	"freshping/template"
	"net/http"

	"github.com/gorilla/mux"
)

var routes []route

type route struct {
	Name            string
	Method          string
	Path            string
	HttpHandlerfunc http.HandlerFunc
}

var pingroutes = []route{
	route{
		Name:            "index",
		Method:          "GET",
		Path:            "/index",
		HttpHandlerfunc: template.Index,
	},
	route{
		Name:            "login",
		Method:          "GET,POST",
		Path:            "/login",
		HttpHandlerfunc: template.Login,
	},
}

//GetGorilaMuxRouter
func GetGorilaMuxRouter() *mux.Router {
	muxroute := mux.NewRouter().StrictSlash(true)
	for _, route := range pingroutes {
		muxroute.Methods(route.Method).Path(route.Path).Handler(middleware.Logger(route.HttpHandlerfunc, route.Name)).Name(route.Name)
		//MuxRouter.Methods(route.MethodType).Path(route.URLPattern).Name(route.FuncName).Handler(logger.Logger(handler, route.FuncName))
	}
	fmt.Println("Muxroute : ", muxroute)
	return muxroute
}
