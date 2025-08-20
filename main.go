package main

import (
	"fmt"
	"html/template"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/centrifugal/centrifuge"
)

func main(){
	app := chi.NewRouter()
	tmpl, err := template.ParseGlob("templates/*.html")

	if err != nil{
		panic(err)
	}

	node, err := centrifuge.New(centrifuge.Config{})

	if err != nil {
		panic(err)
	}

	fmt.Println(node.Config().Name)

	app.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	app.Get("/", func(res http.ResponseWriter, req *http.Request) {
		err := tmpl.ExecuteTemplate(res, "index.html", nil)

		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
		}
	})

	fmt.Println("Listening on port 8080")
	http.ListenAndServe(":8080", app)
}

func auth(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		// Put authentication Credentials into request Context.
		// Since we don't have any session backend here we simply
		// set user ID as empty string. Users with empty ID called
		// anonymous users, in real app you should decide whether
		// anonymous users allowed to connect to your server or not.
		credentials := &centrifuge.Credentials{
			UserID: "",
		}
		newCtx := centrifuge.SetCredentials(ctx, credentials)
		h.ServeHTTP(w, r.WithContext(newCtx))
	})
}