package util

import (
  "BeatBox/dao"
  "BeatBox/models"
  "log"
)
// Search returns a UserDisplaySong matching the keyword from any of the following
// Albums, Artists or song names
func Search(keyword string) ([]models.UserDisplaySong, error) {
    var searchResults []models.UserDisplaySong
    seenSongs := make(map[string]bool)

    // Search for matches in Song table
    songs, err := dao.SearchInSongTable(keyword)
    if err != nil {
        log.Printf("Error creating DynamoDB session: %v", err)
        return nil, err
    }

    for _, song := range songs {
        if _, ok := seenSongs[song.SongID]; ok {
            continue
        }
        // set this songID as already processed
        seenSongs[song.SongID] = true

        // query the database to get album/artist name
        artistName, err := dao.GetArtistName(song.ArtistID)
        if err != nil {
            log.Printf("Error getting artist name: %v", err)
            return nil, err
        }
        albumName, err := dao.GetAlbumName(song.AlbumID)
        if err != nil {
            log.Printf("Error getting album name: %v", err)
            return nil, err
        }

        searchResults = append(searchResults, models.UserDisplaySong{
            SongID: song.SongID,
            SongName:song.SongName,
            ArtistName:artistName,
            AlbumName:albumName,
        })
    }

    // Search for matches in Album table
    albums, err := dao.SearchInAlbumTable(keyword)
    if err != nil {
        log.Printf("Error searching album table: %v", err)
        return nil, err
    }

    // for albums found, get songs
    for _,album := range albums {
        for _,songId := range album.SongIDs {
            if _, ok := seenSongs[songId]; ok {
                continue
            }
            seenSongs[songId] = true

            // query the database to get song/artistName
            song, err := dao.GetSong(songId)
            if err != nil {
                log.Printf("Error getting song: %v", err)
                return nil, err
            }

            artistName, err := dao.GetArtistName(song.ArtistID)
            if err != nil {
                log.Printf("Error getting artist name: %v", err)
                return nil, err
            }

            searchResults = append(searchResults, models.UserDisplaySong{
                SongID:song.SongID,
                SongName:song.SongName,
                ArtistName:artistName,
                AlbumName:album.AlbumName,
            })
        }
    }

    // Search for matches in Artist table
    artists, err := dao.SearchInArtistTable(keyword)
    if err != nil {
        log.Printf("Error searching artist table: %v", err)
        return nil, err
    }

    // for each artist, get all songs
    for _,artist := range artists {
        for _,albumId := range artist.AlbumIDs {
            album, err := dao.GetAlbum(albumId)
            if err != nil {
                log.Printf("Error getting album: %v", err)
                return nil, err
            }
            // for each song in the album, add it to our search results if not already seen
            for _,songId := range album.SongIDs {
                if _, ok := seenSongs[songId]; ok {
                    continue
                }
                seenSongs[songId] = true

                song, err := dao.GetSong(songId)
                if err != nil {
                    log.Printf("Error getting song: %v", err)
                    return nil, err
                }

                searchResults = append(searchResults, models.UserDisplaySong{
                    SongID:song.SongID,
                    SongName:song.SongName,
                    ArtistName:artist.ArtistName,
                    AlbumName:album.AlbumName,
                })
            }
        }
    }

    return searchResults, nil
}
