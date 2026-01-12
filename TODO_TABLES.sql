-- TODO: Users, sessions, auth stuff

-- Users may create and grant other users access to multiple artists. Songs and revisions belong to a particular artist.
CREATE TABLE artists (
	id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
	name VARCHAR (80) UNIQUE NOT NULL,
	url VARCHAR (80) UNIQUE NOT NULL
);

CREATE TABLE users_artists (
	user_id BIGINT FOREIGN KEY REFERENCES users NOT NULL,
	artist BIGINT FOREIGN KEY REFERENCES artists NOT NULL,
);
COMMENT ON TABLE users_artists IS 'What users have administrative access to what artists.';

CREATE TABLE songs (
	id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
	artist BIGINT FOREIGN KEY REFERENCES artists NOT NULL,
	name VARCHAR (80) NOT NULL,
	name_is_placeholder BOOLEAN NOT NULL DEFAULT FALSE,
	emoji VARCHAR (8),
	UNIQUE (artist, name)
);
COMMENT ON COLUMN songs.name_is_placeholder IS 'Whether or not to place quote marks around the song name, indicating the name is only a placeholder.';
COMMENT ON COLUMN songs.emoji IS 'An emoji, if any, to display alongside the song name.';

CREATE TABLE revisions (
	id BIGINT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
	song BIGINT FOREIGN KEY REFERENCES songs NOT NULL,
	when DATE NOT NULL,
	audio_compressed VARCHAR (250) NOT NULL,
	audio_original VARCHAR (250) NOT NULL,
	album_art VARCHAR (250),
	description TEXT
);
COMMENT ON COLUMN revisions.when IS 'The date on which this revision was rendered.';
COMMENT ON COLUMN revisions.audio_compressed IS 'The file name of the Opus-encoded audio asset.';
COMMENT ON COLUMN revisions.audio_original IS 'The file name of the original audio asset.';
COMMENT ON COLUMN revisions.description IS 'Text to display when viewing this particular revision.';

CREATE TABLE revision_links (
	revision BIGINT FOREIGN KEY REFERENCES revisions NOT NULL,
	type ENUM ('Other', 'Bandcamp', 'Soundcloud', 'Spotify', 'Apple Music') NOT NULL,
	url VARCHAR (250) NOT NULL,
	PRIMARY KEY (revision, type)
);
COMMENT ON TABLE revision_links IS 'Optional links to display when viewing a particular revision, such as when that revision represents a released version of a song.';
COMMENT ON COLUMN revision_links.type IS 'The destination service of the outgoing link; used to determine what validation to perform on url and what icon to display.';
