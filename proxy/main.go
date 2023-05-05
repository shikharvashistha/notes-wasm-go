package main

import (
	"io"
	"net/http"
	"net/url"

	"github.com/iris-contrib/middleware/cors"
	iris "github.com/kataras/iris/v12"
)

func main() {
	app := iris.New()
	app.Use(iris.Compression)
	crs := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		AllowedHeaders:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "OPTIONS", "PUT", "DELETE", "PATCH"},
	})

	app.UseRouter(crs)

	app.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		url, err := url.Parse(r.URL.RawQuery)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		req := &http.Request{
			URL:    url,
			Method: r.Method,
			Body:   r.Body,
			Header: r.Header,
		}

		resp, err := http.DefaultTransport.RoundTrip(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusServiceUnavailable)
			return
		}

		defer resp.Body.Close()

		w.WriteHeader(resp.StatusCode)
		io.Copy(w, resp.Body)
	})

	http.ListenAndServe(":8081", app)
}
