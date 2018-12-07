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
	HTTPHandlerfunc http.HandlerFunc
}

var pingroutes = []route{
	route{
		Name:            "index",
		Method:          "GET",
		Path:            "/",
		HTTPHandlerfunc: template.Index,
	},
	route{
		Name:            "login",
		Method:          "GET,POST",
		Path:            "/login",
		HTTPHandlerfunc: template.Login,
	},
	route{
		Name:            "signup",
		Method:          "GET,POST",
		Path:            "/signup",
		HTTPHandlerfunc: template.SignUp,
	},
	route{
		Name:            "panic",
		Method:          "GET",
		Path:            "/panic",
		HTTPHandlerfunc: template.Panic,
	},
}

//GetGorilaMuxRouter  fsjfacjs
func GetGorilaMuxRouter() *mux.Router {
	muxroute := mux.NewRouter()
	for _, route := range pingroutes {
		//var handle http.Handler
		//handle = route.HTTPHandlerfunc
		muxroute.Methods(route.Method).Path(route.Path).Handler(middleware.Logger(route.HTTPHandlerfunc, route.Name)).Name(route.Name)
		//muxroute.Methods(route.Method).Path(route.Path).HandlerFunc(middleware.Logger(route.HTTPHandlerfunc, route.Name))
		//MuxRouter.Methods(route.MethodType).Path(route.URLPattern).Name(route.FuncName).Handler(logger.Logger(handler, route.FuncName))
	}
	fmt.Println("Muxroute : ", muxroute)
	return muxroute
}
