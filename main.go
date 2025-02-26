package main

import (
	. "fmt"
	"html/template"
	"mime"
	"net/http"
	"os"
)

// Top-level just lists all songs and allows the latest version of each
// to be played from there

// Clicking on a song navigates to that song's specific page, where
// different versions can be selected and played

// Pages specific to a version of a song also exist for sharing
// purposes, but aren't a main feature otherwise

var songs []Song
var audio_dir string = "audio"
var port string = ":8080"

func main() {
	mime.AddExtensionType(".css", "text/css")
	mime.AddExtensionType(".js", "text/javascript")
	mime.AddExtensionType(".mjs", "text/javascript")

	mime.AddExtensionType(".wav", "audio/wav")
	mime.AddExtensionType(".mp3", "audio/mpeg")
	mime.AddExtensionType(".flac", "audio/flac")
	mime.AddExtensionType(".aac", "audio/aac")

	for i := 1; i < len(os.Args); i++ {
		switch os.Args[i] {
		case "-a", "--audio-dir":
			if 1+i < len(os.Args) {
				i++
				audio_dir = os.Args[i]
			} else {
				Printf("Expected audio directory path following %s", os.Args[i])
				return
			}
		case "-p", "--port":
			if 1+i < len(os.Args) {
				i++
				port = os.Args[i]
			} else {
				Error("Expected network port following %s", os.Args[i])
				return
			}
		default:
			Error("Unknown parameter %s", os.Args[i])
			return
		}
	}

	var err error
	if songs, err = scan_all_songs(audio_dir); err != nil {
		Error("%s", err.Error())
		return
	}

	Info("Discovered %d songs", len(songs))

	mux := http.NewServeMux()
	mux.Handle("GET /audio/", http.StripPrefix("/audio/", http.FileServer(http.Dir(audio_dir))))
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	mux.HandleFunc("GET /song/{song}", serve_song)
	mux.HandleFunc("GET /", serve_all_songs)
	Info("Listening on port %s", port)
	http.ListenAndServe(Sprintf(":%s", port), mux)
}

func serve_song(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseGlob("template/*.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	if err = tmpl.ExecuteTemplate(w, "song", songs); err != nil {
		Println(err)
	}
}

func serve_all_songs(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseGlob("template/*.html")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	if err = tmpl.ExecuteTemplate(w, "main", songs); err != nil {
		Println(err)
	}
}
