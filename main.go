package main

import (
	"fmt"
	"log"
	"net/http"
	"text/template"
)

func SendForm(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" && r.URL.Path == "/action" {
		fmt.Fprint(w, r.URL.Path)
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
