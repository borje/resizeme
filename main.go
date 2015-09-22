package main

import (
	"flag"
	"fmt"
	"image"
	"image/jpeg"
	"log"
	"net/http"
	"strconv"

	"github.com/disintegration/imaging"
	"github.com/gorilla/mux"
)

var port = flag.Int("port", 8081, "Listening port")

func main() {
	flag.Parse()
	router := mux.NewRouter()
	router.HandleFunc("/thumbnail/{width:[0-9]+}x{height:[0-9]+}", resizeHandler)
	addr := fmt.Sprintf(":%d", *port)
	log.Println("Listening on port: ", addr)
	http.ListenAndServe(addr, router)
}

func resizeHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	width, _ := strconv.Atoi(vars["width"])
	height, _ := strconv.Atoi(vars["height"])
	url := r.URL.Query().Get("url")
	log.Println(width, height, url)
	if url == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	resp, err := http.Get(url)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "image/jpeg")
	w.Header().Set("Content-Disposition", "inline; filename='image.jpg'")
	if err != nil {
		log.Fatal(err)
	}
	img, _, err := image.Decode(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	m := imaging.Thumbnail(img, width, height, imaging.NearestNeighbor)

	jpeg.Encode(w, m, nil)
}
