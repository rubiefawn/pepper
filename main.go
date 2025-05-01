package main

import (
	. "fmt"
	"html/template"
	"mime"
	"net/http"
	"os"
	"strings"
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
var tmpl *template.Template

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

	tmpl, err = template.ParseGlob("template/*.html")
	if err != nil {
		Error("%s", err.Error())
		return
	}

	mux := http.NewServeMux()
	mux.Handle("GET /audio/", http.StripPrefix("/audio/", http.FileServer(http.Dir(audio_dir))))
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	mux.HandleFunc("GET /song/{name}", serve_song)
	mux.HandleFunc("GET /rescan", rescan_songs)
	mux.HandleFunc("GET /reparse", reparse_templates)
	mux.HandleFunc("GET /{$}", serve_all_songs)
	Info("Listening on port %s", port)
	http.ListenAndServe(Sprintf(":%s", port), mux)
}

// HACK: Individual song pages are currently just the home page but
// filtered by song name. This is done since the individual song
// template has no <head>, JavaScript, or playback controls. This needs
// to be changed in the future to an actual dedicated song page.
func serve_song(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	var song []Song
	for i, s := range songs {
		if strings.EqualFold(name, s.Uri()) {
			song = songs[i : i+1]
			break
		}
	}

	if song == nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if err := tmpl.ExecuteTemplate(w, "main", song); err != nil {
		Error("%s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}
}

func serve_all_songs(w http.ResponseWriter, r *http.Request) {
	if err := tmpl.ExecuteTemplate(w, "main", songs); err != nil {
		Error("%s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}
}

func rescan_songs(w http.ResponseWriter, _ *http.Request) {
	var err error
	if songs, err = scan_all_songs(audio_dir); err != nil {
		Error("%s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	msg := Sprintf("Discovered %d songs", len(songs))
	Info("%s", msg)
	w.Write([]byte(msg))
}

func reparse_templates(w http.ResponseWriter, _ *http.Request) {
	var err error
	var t *template.Template
	if t, err = template.ParseGlob("template/*.html"); err != nil {
		Error("%s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	tmpl = t
}
