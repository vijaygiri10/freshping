package template

import (
	"fmt"
	"html/template"
	"net/http"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("html/*.gohtml"))
	fmt.Println("tpl: \n", tpl.Tree)
}

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Index")
	index := struct {
		Tittle string
	}{
		Tittle: "Index",
	}
	tpl.ExecuteTemplate(w, "index.gohtml", index)
}

func Login(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Login")
	login := struct {
		Tittle string
	}{
		Tittle: "Login",
	}
	tpl.ExecuteTemplate(w, "login.gohtml", login)
}
