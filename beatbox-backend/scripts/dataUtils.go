package main

import (
  "BeatBox/dao"
  "BeatBox/models"
  "encoding/json"
  "os"
  "path/filepath"
  "fmt"
  "io/ioutil"
  "log"

  "github.com/aws/aws-sdk-go/aws"
  "github.com/aws/aws-sdk-go/service/dynamodb"
  "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"

  "github.com/google/uuid"
  "github.com/icrowley/fake"
)

// NumData The number of songs/albums/artists to generate
// NOTE: By default there will be 1 song per artist and one album per song
// If you want to change that just edit the code below
const NumData = 43

// FilePath Path to the json objects
const FilePath = "scripts/sampledata/"

// FilePathToMp3 Path to the mp3 files to rename
const FilePathToMp3 = "scripts/sampledata/songs"

// Scripts to create and insert data for the tables listed in models/model.go
// You can simply use the existing json files in sampledata if they exist otherwise generate your own
func main() {
//    // Generate random data into JSON
//    generateArtists()
//    generateAlbumsAndSongsFromArtists()
//    // Use existing JSON to insert into DynamoDB
//    insertArtistsFromJsonToDynamoDB()
//    insertAlbumsFromJsonToDynamoDB()
//    insertSongsFromJsonIntoDynamoDB()
//   //  Use existing JSON to rename songs to match generated SongIDs
//    renameMp3FilesToSongIDsJson()
}

func generateArtists() {
    var artists []models.Artist

    for i := 0; i < NumData; i++ {
        id := uuid.New().String()
        name := fake.FullName()
        albumIDs := []string{uuid.New().String()}
        artists = append(artists, models.Artist{ArtistID: id, ArtistName: name, AlbumIDs: albumIDs})
    }

    jsonData, err := json.MarshalIndent(artists, "", "    ")
    if err != nil {
        log.Fatalf("Error marshalling artists to JSON: %v", err)
    }

    err = ioutil.WriteFile(FilePath + "artists.json", jsonData, 0644)
    if err != nil {
        log.Fatalf("Error writing JSON to file: %v", err)
    }

    fmt.Println("Successfully generated artists.json")
}

// Requires artist.json from generateArtists
func generateAlbumsAndSongsFromArtists() {
    // Read the existing artist JSON file
    artistsJSON, err := ioutil.ReadFile(FilePath + "artists.json")
    if err != nil {
        fmt.Println(err)
        return
    }

    // Unmarshal the JSON data into an array of Artist structs
    var artists []models.Artist
    err = json.Unmarshal(artistsJSON, &artists)
    if err != nil {
        fmt.Println(err)
        return
    }

    // Create an array of Album structs
    var albums []models.Album
    // Create an array of Song structs
    var songs []models.Song
    // generate albums and songs per artist (default one album and one song per artist for sample data)
    for _, artist := range artists {
        // Generate random album/song names (I did my best)
        albumName := fake.Brand()
        songName := fake.ProductName()
        songID := []string{uuid.New().String()}

        album := models.Album{
            AlbumID:       artist.AlbumIDs[0],
            AlbumName:     albumName,
            ArtistID:      artist.ArtistID,
            SongIDs:       songID,
            }

        // append to the albums struct
        albums = append(albums, album)

        song := models.Song{
            SongID:       songID[0],
            SongName:     songName,
            ArtistID:     artist.ArtistID,
            AlbumID:      artist.AlbumIDs[0],
            }

        // append to the songs struct
        songs = append(songs, song)
    }

    // Marshal the array of Album structs into JSON data
    albumsJSON, err := json.MarshalIndent(albums, "", "  ")
    if err != nil {
        fmt.Println(err)
        return
    }

    // Write the JSON data to a new file
    err = ioutil.WriteFile(FilePath + "albums.json", albumsJSON, 0644)
    if err != nil {
        fmt.Println(err)
        return
    }

    // Marshal the array of Album structs into JSON data
    songsJSON, err := json.MarshalIndent(songs, "", "  ")
    if err != nil {
        fmt.Println(err)
        return
    }

    // Write the JSON data to a new file
    err = ioutil.WriteFile(FilePath + "songs.json", songsJSON, 0644)
    if err != nil {
        fmt.Println(err)
        return
    }

    fmt.Println("Successfully generated artists.json and songs.json")
}

func insertArtistsFromJsonToDynamoDB() {
    // Get session
    svc, err := dao.CreateDynamoDBSession()
    if err != nil {
        log.Fatal(err)
    }

    // Read data from the JSON file
    jsonFile, err := ioutil.ReadFile(FilePath + "artists.json")
    if err != nil {
        log.Fatalf("Error reading JSON file: %v", err)
    }

    // Unmarshal the JSON data into a slice of Artist structs
    var artists []models.Artist
    err = json.Unmarshal(jsonFile, &artists)
    if err != nil {
        log.Fatalf("Error unmarshalling JSON data: %v", err)
    }

    // Insert each artist into the DynamoDB table
    for _, artist := range artists {
        // Convert the artist struct to a map of attribute values
        av, err := dynamodbattribute.MarshalMap(artist)
        if err != nil {
            log.Fatalf("Error marshalling artist: %v", err)
        }

        // Create input for the PutItem method
        input := &dynamodb.PutItemInput{
            Item:      av,
            TableName: aws.String("Artist"),
            }
            _, err = svc.PutItem(input)
            if err != nil {
                log.Fatalf("Error inserting artist into DynamoDB table: %v", err)
            }

            fmt.Printf("Successfully inserted artist %s into DynamoDB table\n", artist)
    }
}

// Insert artists from the generated JSON into dynamoDB
func insertAlbumsFromJsonToDynamoDB() {
    // Get session
    svc, err := dao.CreateDynamoDBSession()
    if err != nil {
        log.Fatal(err)
    }

    // Read data from the JSON file
    jsonFile, err := ioutil.ReadFile(FilePath + "albums.json")
    if err != nil {
        log.Fatalf("Error reading JSON file: %v", err)
    }

    // Unmarshal the JSON data into a slice of Artist structs
    var artists []models.Album
    err = json.Unmarshal(jsonFile, &artists)
    if err != nil {
        log.Fatalf("Error unmarshalling JSON data: %v", err)
    }

    // Insert each artist into the DynamoDB table
    for _, artist := range artists {
        // Convert the artist struct to a map of attribute values
        av, err := dynamodbattribute.MarshalMap(artist)
        if err != nil {
            log.Fatalf("Error marshalling artist: %v", err)
        }

        // Create input for the PutItem method
        input := &dynamodb.PutItemInput{
            Item:      av,
            TableName: aws.String("Album"),
            }
            _, err = svc.PutItem(input)
            if err != nil {
                log.Fatalf("Error inserting artist into DynamoDB table: %v", err)
            }

            fmt.Printf("Successfully inserted album %s into DynamoDB table\n", artist)
    }
}

// Insert songs from JSON into dynamo DB
func insertSongsFromJsonIntoDynamoDB() {
    // Get session
    svc, err := dao.CreateDynamoDBSession()
    if err != nil {
        log.Fatal(err)
    }

    // Read data from the JSON file
    jsonFile, err := ioutil.ReadFile(FilePath + "songs.json")
    if err != nil {
        log.Fatalf("Error reading JSON file: %v", err)
    }

    // Unmarshal the JSON data into a slice of Artist structs
    var artists []models.Song
    err = json.Unmarshal(jsonFile, &artists)
    if err != nil {
        log.Fatalf("Error unmarshalling JSON data: %v", err)
    }

    // Insert each artist into the DynamoDB table
    for _, artist := range artists {
        // Convert the artist struct to a map of attribute values
        av, err := dynamodbattribute.MarshalMap(artist)
        if err != nil {
            log.Fatalf("Error marshalling artist: %v", err)
        }

        // Create input for the PutItem method
        input := &dynamodb.PutItemInput{
            Item:      av,
            TableName: aws.String("Song"),
            }
            _, err = svc.PutItem(input)
            if err != nil {
                log.Fatalf("Error inserting artist into DynamoDB table: %v", err)
            }

            fmt.Printf("Successfully inserted song %s into DynamoDB table\n", artist)
    }
}

// Rename the files from the given path to match the songIDs generated above
// (files should be royalty free, can get from pixabay.com)
func renameMp3FilesToSongIDsJson() {
    // Read the JSON data from a file or a stream
    jsonData, err := os.Open(FilePath + "songs.json")
    if err != nil {
        log.Fatalf("Unable to read JSON: %v", err)
    }
    defer jsonData.Close()

    // Decode the JSON data into a struct
    var songs []models.Song
    if err := json.NewDecoder(jsonData).Decode(&songs); err != nil {
        log.Fatalf("Unable to decode JSON: %v", err)
    }


    // Create a slice of mp3 files in assets/
    var mp3Files []string
    filepath.Walk(FilePathToMp3, func(path string, info os.FileInfo, err error) error {
        if !info.IsDir() && filepath.Ext(path) == ".mp3" {
            mp3Files = append(mp3Files, path)
        }
        return nil
    })

    // rename the files
    for i, songId := range songs {
        if i < len(mp3Files) {
            newPath := filepath.Join(filepath.Dir(mp3Files[i]), songId.SongID + ".mp3")
            os.Rename(mp3Files[i], newPath)
        } else {
            log.Println("Some files not renamed, you likely haven't generated enough songIDs")
            log.Printf("SongIDs: %d", len(songs))
            log.Printf("MP3 Files: %d", len(mp3Files))
        }
    }
}