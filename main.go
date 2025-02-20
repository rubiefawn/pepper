package main

import (
	. "fmt"
	"html/template"
	"mime"
	"net/http"
	"os"
	fpath "path/filepath"
	"sort"
	"time"

	"github.com/BurntSushi/toml"
)

// Top-level just lists all songs and allows the latest version of each
// to be played from there

// Clicking on a song navigates to that song's specific page, where
// different versions can be selected and played

// Pages specific to a version of a song also exist for sharing
// purposes, but aren't a main feature otherwise

// TODO: Support for album art
// TODO: Revision comments

type Revision struct {
	Path     template.URL
	Modified time.Time
	// Comments string
}

type SongInfo struct {
	Name              string
	NameIsPlaceholder bool `toml:"name_is_placeholder"`
	Emoji             string
	IsReleased        bool `toml:"released"`
}

type Song struct {
	SongInfo
	Revisions []Revision
}

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
				Printf("error: expected audio directory path following %s\n", os.Args[i])
				return
			}
		case "-p", "--port":
			if 1+i < len(os.Args) {
				i++
				port = os.Args[i]
			} else {
				Printf("error: expected network port following %s\n", os.Args[i])
				return
			}
		default:
			Printf("error: unknown parameter %s\n", os.Args[i])
			return
		}
	}

	var err error
	if songs, err = scan_all_songs(audio_dir); err != nil {
		Printf("error: %s\n", err.Error())
		return
	}

	for _, song := range songs {
		Printf("%s %s %v %v\n", song.Emoji, song.Name, song.IsReleased, song.NameIsPlaceholder)
		for _, rev := range song.Revisions {
			Printf("   %s\n", string(rev.Path))
		}
	}

	mux := http.NewServeMux()
	mux.Handle("GET /audio/", http.StripPrefix("/audio/", http.FileServer(http.Dir(audio_dir))))
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	mux.HandleFunc("GET /song/{song}", serve_song)
	mux.HandleFunc("GET /", serve_all_songs)
	http.ListenAndServe(Sprintf(":%s", port), mux)
}

func scan_all_songs(in_path string) (songs []Song, err error) {
	in_path = fpath.Clean(in_path)
	toml_path := fpath.ToSlash(fpath.Join(in_path, "pepper.toml"))
	var data []byte
	if data, err = os.ReadFile(toml_path); err != nil {
		return
	}
	var pepper_toml struct{ Songs []SongInfo }
	if err = toml.Unmarshal(data, &pepper_toml); err != nil {
		return
	}

	for _, song_info := range pepper_toml.Songs {
		if song_info.Emoji == "" {
			song_info.Emoji = "🌶️"
		}
		song_path := fpath.ToSlash(fpath.Join(in_path, song_info.Name))
		var r []Revision
		if r, err = scan_revisions(song_path); err != nil {
			return
		}
		songs = append(songs, Song{song_info, r})
	}
	sort.Slice(songs, func(a int, b int) bool { return songs[a].Revisions[0].Modified.After(songs[b].Revisions[0].Modified) })
	return
}

func scan_revisions(song_path string) (revisions []Revision, err error) {
	var files []os.DirEntry
	if files, err = os.ReadDir(song_path); err != nil {
		return
	}
	for _, rev := range files {
		var i os.FileInfo
		if i, err = rev.Info(); err != nil {
			return
		}
		ext := fpath.Ext(rev.Name())
		if !(ext == ".mp3" || ext == ".wav" || ext == ".flac" || ext == ".aac") {
			continue
		}
		revisions = append(revisions, Revision{template.URL(rev.Name()), i.ModTime()})
	}
	sort.Slice(revisions, func(a int, b int) bool { return revisions[a].Modified.After(revisions[b].Modified) })
	return
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
