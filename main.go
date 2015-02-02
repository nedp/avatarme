package main

import (
	"bitbucket.org/nedp/avatarme/identicon"
	"net/http"
	"os"
	"strings"
	"strconv"
	"image/color"
	"log"
)

const defaultSize = 512
const nBlocks = 6

var bgColour = color.NRGBA{
	R: 0xF0,
	G: 0xF0,
	B: 0xF0,
	A: 0xFF,
}

var hashKey = []byte{
	0x11, 0xBB, 0x22, 0xAA,
	0x33, 0x00, 0xEE, 0x66,
	0x99, 0x44, 0x77, 0x88,
	0xCC, 0xFF, 0x55, 0xDD,
}

func handler(w http.ResponseWriter, r *http.Request) {
	log.Println("Recieved request", r.URL.Path)
	args := strings.Split(r.URL.Path, "/")
	path := args[1:]

	var size int
	if s, ok := r.URL.Query()["s"]; ok {
		tmp, err := strconv.ParseInt(s[0], 10, 0)
		if err != nil {
			http.Error(w, "", http.StatusBadRequest)
			return
		}
		size = int(tmp)
	} else {
		size = defaultSize
	}

	if len(path) != 1 {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	item := path[0]

	if !strings.HasSuffix(item, ".png") {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	info := strings.TrimSuffix(item, ".png")
	id := identicon.FromInfo([]byte(info))

	log.Println("DEBUG-BEFORE")

	pngData, err := id.Render(size, nBlocks, hashKey, bgColour)
	if err != nil {
		http.Error(w, "", http.StatusInternalServerError)
		return
	}

	log.Println("Created identicon for", info)
	w.Header().Set("Content-Type", "image/png")
	w.Write(pngData)

	return
}

func main() {
	log.Println("avatarme - Identicon Generator")
	port := ":8080"
	if p := os.Getenv("PORT"); p != "" {
		port = ":" + p
	}
	log.Println("Listening on port", port)

	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(port, nil))
}
