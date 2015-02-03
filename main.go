package main

import (
	"github.com/nedp/avatarme/identicon"
	"net/http"
	"os"
	"strings"
	"strconv"
	"image/color"
	"log"
)

const (
	defaultSize = 130
	defaultBorder = 1
	defaultNBlocks = 11
)

var bgColour = color.NRGBA{
	R: 0xF0,
	G: 0xF0,
	B: 0xF0,
	A: 0xFF,
}

func handler(w http.ResponseWriter, r *http.Request) {
	log.Println("Recieved request", r.URL.Path)
	args := strings.Split(r.URL.Path, "/")
	path := args[1:]

	var size int
	if q, ok := r.URL.Query()["s"]; ok {
		tmp, err := strconv.ParseInt(q[0], 10, 0)
		if err != nil {
			size = defaultSize
		} else {
			size = int(tmp)
		}
	} else {
		size = defaultSize
	}

	var border int
	if q, ok := r.URL.Query()["b"]; ok {
		tmp, err := strconv.ParseInt(q[0], 10, 0)
		if err != nil {
			border = defaultBorder
		} else {
			border = int(tmp)
		}
	} else {
		border = defaultBorder
	}

	var nBlocks int
	if q, ok := r.URL.Query()["n"]; ok {
		tmp, err := strconv.ParseInt(q[0], 10, 0)
		if err != nil {
			nBlocks = defaultNBlocks
		} else {
			nBlocks = int(tmp)
		}
	} else {
		nBlocks = defaultNBlocks
	}

	size -= size % (nBlocks+2*border)

	item := path[0]

	if !strings.HasSuffix(item, ".png") {
		http.Error(w, "", http.StatusBadRequest)
		return
	}

	info := strings.TrimSuffix(item, ".png")
	id := identicon.FromInfo([]byte(info))

	pngData := id.Hash([]byte(info)).Design(nBlocks).Render(size, border, bgColour)

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
