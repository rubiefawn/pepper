const audio = document.querySelector("audio");
const progress = document.getElementById("playback-progress");
const btn_theme = document.getElementById("btn-theme-toggle");
const time = document.getElementById("playback-time");
const btn_main = document.getElementById("btn-main");
const name_main = document.getElementById("name-main");

function set_song(name, url) {
	audio.src = url;
	progress.max = audio.duration;
	name_main.innerHTML = name;
}

function toggle_playback(force_play) {
	if (!audio.src) { return; }
	if (force_play ?? audio.paused) {
		audio.play();
	} else {
		audio.pause();
	}
}

function audio_progress_readable() {
	const bm = `${Math.floor(audio.currentTime / 60)}`.padStart(2, '0');
	const bs = `${Math.floor(audio.currentTime % 60)}`.padStart(2, '0');
	const em = `${Math.floor(audio.duration / 60)}`.padStart(2, '0');
	const es = `${Math.floor(audio.duration % 60)}`.padStart(2, '0');
	return `${bm}:${bs}/${em}:${es}`;
}

audio.addEventListener("timeupdate", _ => {
	progress.value = audio.currentTime;
	progress.max = audio.duration;
	time.innerHTML = audio_progress_readable();
});

btn_theme.addEventListener("click", _ => {
	const root_html = document.querySelector("html");
	const theme = root_html.getAttribute("theme");
	root_html.setAttribute("theme", theme == "dark" ? "light" : "dark");
});
btn_main.addEventListener("click", _ => toggle_playback());
progress.addEventListener("input", _ => audio.fastSeek(progress.value));
document.querySelectorAll(".song").forEach(song => {
	const version_selector = song.querySelector(".song-revision");
	const name = song.querySelector(".song-name").innerHTML;
	const btn = song.querySelector(".btn-play-song");
	const btn_check = btn.querySelector("input[type=checkbox]");
	btn.addEventListener("click", e => {
		// Pause if this revision is already playing
		const already_playing_this_rev = version_selector.value == audio.getAttribute("src");
		if (!already_playing_this_rev) {
			set_song(`${name} @ ${version_selector.selectedOptions[0].innerHTML}`, version_selector.value);
			document.querySelectorAll("button>input[type=checkbox]").forEach(it => it.checked = false);
			btn_check.checked = true;
			toggle_playback(true);
		} else {
			toggle_playback();
			btn_check.checked = !audio.paused;
		}
	});
});
