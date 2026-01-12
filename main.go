package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

var (
	storage_path string = "storage"
	serve_port   int    = 8080
	db_port      int    = 5432
	encoder      string = "libopus"
	hide_login   bool   = false // If true, login page must be navigated to manually
)

func main() {
	for i := 1; i < len(os.Args); i++ {
		switch os.Args[i] {
		case "--hide-login":
			hide_login = true
		case "--show-login":
			hide_login = false
		case "-s", "--storage-path":
			if 1+i >= len(os.Args) {
				fmt.Printf("Expected user content storage path following %s\n", os.Args[i])
				os.Exit(2)
			}
			i++
			storage_path = os.Args[i]
			// TODO: Check that storage_path exists and permissions are set properly
		case "-p", "--port":
			if 1+i >= len(os.Args) {
				fmt.Printf("Expected network port to serve on following %s\n", os.Args[i])
				os.Exit(2)
			}
			i++
			if v, err := strconv.Atoi(os.Args[i]); err != nil || v < 0 {
				fmt.Printf("%q is not a valid network port\n", os.Args[i])
				os.Exit(2)
			} else {
				serve_port = v
			}
		case "-d", "--db-port":
			if 1+i >= len(os.Args) {
				fmt.Printf("Expected network port of database following %s\n", os.Args[i])
				os.Exit(2)
			}
			i++
			if v, err := strconv.Atoi(os.Args[i]); err != nil || v < 0 {
				fmt.Printf("%q is not a valid network port\n", os.Args[i])
				os.Exit(2)
			} else {
				db_port = v
			}
		// TODO: TLS certificate and key files parameters
		default:
			fmt.Printf("Unknown parameter %q\n", os.Args[i])
			os.Exit(2)
		}
	}
	// TODO: Config file so params don't have to be passed every single time? Params should override config file settings

	find_ffmpeg := exec.Command("ffmpeg", "-encoders")
	var find_ffmpeg_output strings.Builder
	find_ffmpeg.Stdout = &find_ffmpeg_output
	if err := find_ffmpeg.Run(); err != nil {
		log.Fatalf("Could not find ffmpeg: %s", err.Error())
	} else if encoders := find_ffmpeg_output.String(); !strings.Contains(encoders, "libopus Opus") {
		if !strings.Contains(encoders, "Opus") {
			log.Fatalf("Found ffmpeg, but it has no Opus encoder")
		} else {
			log.Printf("Found ffmpeg, but libopus Opus encoder is missing; falling back to ffmpeg's Opus encoder")
			encoder = "opus"
		}
	}

	compressed_audio_dir := filepath.Join(storage_path, "audio") // Opus-encoded copies of audio assets go here
	raw_audio_dir := filepath.Join(storage_path, "raw")          // Original audio assets go here
	user_images_dir := filepath.Join(storage_path, "img")        // Album art & avatars go here
	for _, path := range []string{compressed_audio_dir, raw_audio_dir, user_images_dir} {
		if err := os.MkdirAll(path, os.ModePerm); err != nil {
			log.Fatalf("Could not create %q: %s", path, err.Error())
		}
	}

	mux := http.NewServeMux()
	mux.Handle("GET /static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	mux.Handle("GET /a/", http.StripPrefix("/a/", http.FileServer(http.Dir(compressed_audio_dir))))
	mux.Handle("GET /r/", http.StripPrefix("/r/", http.FileServer(http.Dir(raw_audio_dir))))
	mux.Handle("GET /i/", http.StripPrefix("/i/", http.FileServer(http.Dir(user_images_dir))))
	// mux.HandleFunc("GET /{username}/{song}", TODO_FUNC_SERVE_SONG) // TODO: Links to individual revisions accomplished via URL query parameter
	// mux.HandleFunc("GET /{username}", TODO_FUNC_SERVE_USER_SONGS)
	// mux.HandleFunc("GET /{$}", TODO_FUNC_SERVE_EVERY_SONG)
	log.Printf("Listening on port %d", serve_port)
	// TODO: Use TLS
	if err := http.ListenAndServe(fmt.Sprintf(":%d", serve_port), mux); err != nil {
		log.Fatalf("Could not serve on port %d: %s", serve_port, err.Error())
	}
}
