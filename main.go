package main

import (
	"fmt"
	"log"
	"net/http"
	"net/smtp"
	"strings"
	"text/template"
)

const (
	yandex      = "smtp.yandex.ru"
	yandex_port = "465"
	google      = "smtp.gmail.com"
	google_port = "587"
)

type Email struct {
	Email_field     string
	Password_field  string
	Text_field      string
	Email_field_rec string
}

func (e Email) SendEmail() {
	var port string
	var host string
	switch (strings.Split(e.Email_field, "@"))[1] {
	case "yandex.ru":
		port = yandex_port
		host = yandex
	case "gmail.com":
		port = google_port
		host = google

	}

	auth := smtp.PlainAuth("", e.Email_field, e.Password_field, host)
	err := smtp.SendMail(host+":"+port, auth, e.Email_field, []string(e.Email_field_rec), []byte(e.Text_field))
	if err != nil {
		log.Println(err)
		
		fmt.Println("Sending email...")
	}
}
func SendForm(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" && r.URL.Path == "/action" {
		err := r.ParseForm()
		if err != nil {
			log.Println(err)
		}
		data := Email{
			Email_field:    r.FormValue("email"),
			Password_field: r.FormValue("password"),
			Text_field:     r.FormValue("message"),
      Email_field_rec: r.FormValue("email2"),
		}
		
		data.SendEmail()
		fmt.Fprint(w, data)
	} else {
		ts, err := template.ParseFiles("static/form.html")
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Internal Server Error", 500)
			return
		}
		ts.Execute(w, nil)

	}
}

func main() {
	// 1. Create a new http server listening on port 8080
	http.Handle("/static/", http.StripPrefix("/static", http.FileServer(http.Dir("./static"))))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		//http.ServeFile(w, r, "static/form.html")
		ts, err := template.ParseFiles("static/form.html")
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Internal Server Error", 500)
			return
		}
		err = ts.Execute(w, nil)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Internal Server Error", 500)
		}
	})
	http.HandleFunc("/action", SendForm)
	log.Fatal(http.ListenAndServe(":80", nil))
}
