// Package handler/handler.go
package handler

import (
  "BeatBox/dao"
  "BeatBox/models"
  "BeatBox/util"
  "encoding/json"
  "io"
  "log"
  "net/http"
  "os"
)

type File struct {
    Name string `json:"name"`
    Id string `json:"id"`
}

func InsertUserHandler(w http.ResponseWriter, r *http.Request) {
    var user models.User
    doBadStuff(w)
    err := json.NewDecoder(r.Body).Decode(&user)
    if err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    err = dao.InsertUser(user.Email, user.Name)
    if err != nil {
      http.Error(w, "Error inserting user, user may already exist", http.StatusConflict)
      return
    }
    w.WriteHeader(http.StatusOK)
}

// FetchSongById Fetch song by ID from database
func FetchSongById(w http.ResponseWriter, r *http.Request) {
    // Retrieve ethe song ID from the query param
    // E.G. MYURL.com/song?id={songID}
    songId := r.URL.Query().Get("id")

    // Call the getSong function to retrieve the file
    file, err := dao.GetSongFromS3ById(songId)
    if err != nil {
        http.Error(w, "Error retrieving file from S3", http.StatusInternalServerError)
        return
    }
    defer file.Close()

    // Set the appropriate Content-Type header
    w.Header().Set("Content-Type", "audio/mpeg")

    doBadStuff(w)

    // Stream the file to the client
    io.Copy(w, file)
}

// HandleFetchSongs Fetch songs from databaas
func HandleFetchSongs(w http.ResponseWriter, r *http.Request) {
    // Get the songs
    userDisplaySongs, err := dao.GetSongsAsUserDisplaySong()
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Convert the userDisplaySongs slice to JSON
    userDisplaySongsJson, err := json.Marshal(userDisplaySongs)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Send the JSON as the response
    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

    doBadStuff(w)

    w.Write(userDisplaySongsJson)
}

// HandleKeywordSearch Search the database for matching keyword
func HandleKeywordSearch(w http.ResponseWriter, r *http.Request) {
    // Since this is a read heavy request we're just going to ignore preflight
    if isPreflightRequest(w, r) {
        return
    }

    // Retrieve the search from the query
    searchTerm := r.URL.Query().Get("searchTerm")

    // get search results
    userDisplaySongs, err := util.Search(searchTerm)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Convert the userDisplaySongs slice to JSON
    userDisplaySongsJson, err := json.Marshal(userDisplaySongs)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // Send the JSON as the response
    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

    doBadStuff(w)

    w.Write(userDisplaySongsJson)
}

func isPreflightRequest(w http.ResponseWriter, r *http.Request) bool {
  if r.Method == http.MethodOptions {
    // This is a preflight request
    handlePreflight(w, r)
    return true
  } else {
    return false
  }
}

// if preflight just return OK
func handlePreflight(w http.ResponseWriter, r *http.Request) {
  // Set CORS headers to allow the actual request
  w.Header().Set("Access-Control-Allow-Origin", "*")
  w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
  w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
  w.Header().Set("Access-Control-Max-Age", "3600")

  // Return a "200 OK" response
  w.WriteHeader(http.StatusOK)
}

// Set CORS to allow any requestor (this is bad practice in prod but fine since it's just for development)
func doBadStuff(w http.ResponseWriter) {
    // Set the Access-Control-Allow-Origin header to allow requests from the react client
    w.Header().Set("Access-Control-Allow-Origin", "*")
}

// ================================================== LOCALLY HOSTED RETRIEVAL CODE ===================================================

// Fetch song by ID from local
func FetchSongByIdLocal(w http.ResponseWriter, r *http.Request) {
    // Retrieve ethe song ID from the query param
    // E.G. MYURL.com/song?id={songID}
    songId := r.URL.Query().Get("id")

    // Open the mp3 file
    f, err := os.Open("./scripts/sampledata/songs/" + songId)
    if err != nil {
        http.Error(w, "Error opening file", http.StatusInternalServerError)
        return
    }
    defer f.Close()

    // Set the appropriate Content-Type header
    w.Header().Set("Content-Type", "audio/mpeg")

    doBadStuff(w)

    // Read and send the audio file in chunks of 300kb each
    const chunkSize = 3000000
    buf := make([]byte, chunkSize)
    for {
        // Read a chunk
        n, err := f.Read(buf)
        if err != nil && err != io.EOF {
            http.Error(w, "Error reading file", http.StatusInternalServerError)
            return
        }
        if n == 0 {
            break
        }

        // Send the chunk
        if _, err := w.Write(buf[:n]); err != nil {
            http.Error(w, "Error sending file", http.StatusInternalServerError)
            return
        }

        // Flush the buffer to send the chunk to the client
        if flusher, ok := w.(http.Flusher); ok {
            flusher.Flush()
        } else {
            // Flusher not supported, send the file in one shot
            io.Copy(w, f)
        }
    }
}
// Retrieve song mp3s from local file system
func HandleFilesLocal(w http.ResponseWriter, r *http.Request) {
    dir := "./scripts/sampledata/songs"

    doBadStuff(w)
    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

    d, err := os.Open(dir)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer d.Close()

    files, err := d.Readdir(-1)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    var fileList []File
    for _, file := range files {
        fileList = append(fileList, File{file.Name(), file.Name()})
    }

    json.NewEncoder(w).Encode(fileList)
}

// Retrieve songs from local folder
func HandleSearchLocal(w http.ResponseWriter, r *http.Request) {
    dir := "./scripts/sampledata/songs"
    // Retrieve ethe song ID from the query param
    // E.G. MYURL.com/song?id={songID}
    searchTerm := r.URL.Query().Get("searchTerm")

    doBadStuff(w)
    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

    d, err := os.Open(dir)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    defer d.Close()

    files, err := d.Readdir(-1)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    var fileList []File
    //TODO: Remove when actual search is working
    if (searchTerm == "nothing") {
        log.Print("Nothing")
    } else {
        for _, file := range files {
            fileList = append(fileList, File{file.Name(), file.Name()})
        }
    }

    json.NewEncoder(w).Encode(fileList)
}