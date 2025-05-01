package main

import (
	"html/template"
	"os"
	fpath "path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
)

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

func (i *SongInfo) Uri() string {
	return strings.ToLower(strings.ReplaceAll(i.Name, " ", "-"))
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
