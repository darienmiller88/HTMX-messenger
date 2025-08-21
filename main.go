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

	node.OnConnect(func(client *centrifuge.Client) {
		// In our example transport will always be Websocket but it can be different.
		transportName := client.Transport().Name()

		// In our example clients connect with JSON protocol but it can also be Protobuf.
		transportProto := client.Transport().Protocol()
		fmt.Printf("client connected via %s (%s)", transportName, transportProto)

		client.OnSubscribe(func(e centrifuge.SubscribeEvent, cb centrifuge.SubscribeCallback) {
			fmt.Printf("client subscribes on channel %s", e.Channel)
			cb(centrifuge.SubscribeReply{}, nil)
		})

		client.OnPublish(func(e centrifuge.PublishEvent, cb centrifuge.PublishCallback) {
			fmt.Printf("client publishes into channel %s: %s", e.Channel, string(e.Data))
			cb(centrifuge.PublishReply{}, nil)
		})

		// Set Disconnect handler to react on client disconnect events.
		client.OnDisconnect(func(e centrifuge.DisconnectEvent) {
			fmt.Printf("client disconnected")
		})
	})

	// Run node. This method does not block. See also node.Shutdown method
	// to finish application gracefully.
	if err := node.Run(); err != nil {
		panic(err)
	}

	// Serve Websocket connections using WebsocketHandler.
	wsHandler := centrifuge.NewWebsocketHandler(node, centrifuge.WebsocketConfig{})
	
	app.Handle("/connection/websocket", auth(wsHandler))
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