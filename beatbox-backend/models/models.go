package models


// Note to self: They're backtics not apostrophies

type User struct {
    Email string `json:"Email"`
    Name  string `json:"Name"`
}

type Artist struct {
    ArtistID string `json:"ArtistID"`
    ArtistName    string `json:"ArtistName"`
    AlbumIDs  []string `json:"AlbumIDs"`
}

type Album struct {
    AlbumID      string `json:"AlbumID"`
    AlbumName    string `json:"AlbumName"`
    ArtistID string `json:"ArtistID"`
    SongIDs  []string `json:"SongIDs"`
}

type Song struct {
    SongID      string `json:"SongID"`
    SongName    string `json:"SongName"`
    ArtistID string `json:"ArtistID"`
    AlbumID  string `json:"AlbumID"`
}

type UserDisplaySong struct {
    SongID      string `json:"SongID"`
    SongName    string `json:"SongName"`
    ArtistName  string `json:"ArtistName"`
    AlbumName   string `json:"AlbumName"`
}