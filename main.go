package main

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main(){
	app := chi.NewRouter()
	tmpl, err := template.ParseGlob("templates/*.html")

	if err != nil{
		panic(err)
	}

	app.Get("/", func(res http.ResponseWriter, req *http.Request) {
		err := tmpl.ExecuteTemplate(res, "index.html", nil)

		if err != nil {
			fmt.Println("err:", err)
		}
	})

	fmt.Println("Listening on port 8080")
	http.ListenAndServe(":8080", app)
}