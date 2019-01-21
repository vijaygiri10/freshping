package template

import (
	"fmt"
	"html/template"
	"net/http"
)

var tpl *template.Template

type SignUpForm struct {
	FieldNames []string
	Fields     map[string]string
	Errors     map[string]string
}

func init() {
	tpl = template.Must(template.ParseGlob("html/*.gohtml"))
	//fmt.Println("tpl: \n", tpl.Tree)
}

func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Index")
	index := struct {
		Title string
	}{Title: "Index Page"}

	if err := tpl.ExecuteTemplate(w, "index.gohtml", index); err != nil {
		fmt.Println("error Execute template : ", err)
		fmt.Fprintln(w, err.Error())
	}
}

func Login(w http.ResponseWriter, r *http.Request) {
	s := SignUpForm{}
	s.FieldNames = []string{"username", "firstName", "lastName", "email"}
	s.Fields = make(map[string]string)
	s.Errors = make(map[string]string)
	switch r.Method {
	case "POST":
		if err := tpl.ExecuteTemplate(w, "validatelogin.gohtml", &s); err != nil {
			fmt.Println("error Execute template : ", err)
			fmt.Fprintln(w, err.Error())
		}
	case "GET":
		if err := tpl.ExecuteTemplate(w, "login.gohtml", &s); err != nil {
			fmt.Println("error Execute template : ", err)
			fmt.Fprintln(w, err.Error())
		}
	}

}
func StreamToBuffer(stream io.Reader) *bytes.Buffer {
	buf := new(bytes.Buffer)
	buf.ReadFrom(stream)
	return buf
}

func (c *Controller) Login(w http.ResponseWriter, req *http.Request) {
	fmt.Println(w, " ---- Login ---- ")
	fmt.Println("Method : ", req.Method)
	fmt.Println("Form : ", req.Form)
	fmt.Println("URL : ", req.URL)
	fmt.Println("Header : ", req.Header)
	fmt.Println("body : ", mux.Vars(req.Form))
	fmt.Println()
	fmt.Println("body : ", StreamToBuffer(req.Body).String())
	var user userDetails
	type userDetails struct {
		firstName, lastName, email, password string
	}
	if err := json.Unmarshal(StreamToBuffer(req.Body).Bytes(), &register.UserD); err != nil {
		fmt.Fprintf(w, err.Error())
	}
	userD

}

func SignUp(w http.ResponseWriter, r *http.Request) {
	s := SignUpForm{}
	s.FieldNames = []string{"username", "firstName", "lastName", "email"}
	s.Fields = make(map[string]string)
	s.Errors = make(map[string]string)

	switch r.Method {
	case "POST":
		tpl.ExecuteTemplate(w, "validatesignupform.gohtml", &s)
	case "GET":
		tpl.ExecuteTemplate(w, "signupform.gohtml", &s)
	default:
		tpl.ExecuteTemplate(w, "signupform.gohtml", &s)
	}

}

func Panic(w http.ResponseWriter, r *http.Request) {

	panic("Panic cretaed")
}
