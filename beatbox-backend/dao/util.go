package dao

import (
    "BeatBox/models"
    "io"

    "github.com/aws/aws-sdk-go/aws"
    "github.com/aws/aws-sdk-go/aws/session"
    "github.com/aws/aws-sdk-go/service/dynamodb"
    "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
    "github.com/aws/aws-sdk-go/service/s3"
    "log"
)

// CreateDynamoDBSession Initialize a session that the SDK will use to load
// credentials from the shared credentials file ~/.aws/credentials
// and region from the shared configuration file ~/.aws/config.
// NOTE: By default ~/.aws/config is not used unless the AWS_SDK_LOAD_CONFIG flag is set to true
// This is the SharedConfigState as defined below
func CreateDynamoDBSession() (*dynamodb.DynamoDB, error) {
    sess, err := session.NewSessionWithOptions(session.Options{
        SharedConfigState: session.SharedConfigEnable,
    })
    if err != nil {
        return nil, err
    }

    return dynamodb.New(sess), nil
}

// CreateS3Session Initialize a session that the SDK will use to load
// credentials from the shared credentials file ~/.aws/credentials
// and region from the shared configuration file ~/.aws/config.
// NOTE: By default ~/.aws/config is not used unless the AWS_SDK_LOAD_CONFIG flag is set to true
// This is the SharedConfigState as defined below
func CreateS3Session() (*s3.S3, error) {
    sess, err := session.NewSessionWithOptions(session.Options{
        SharedConfigState: session.SharedConfigEnable,
    })
    if err != nil {
        return nil, err
    }

    return s3.New(sess), nil
}

func GetSong(songId string) (models.Song, error) {
    var song models.Song

    // Get session
    svc, err := CreateDynamoDBSession()
    if err != nil {
        log.Printf("Error creating DynamoDB session: %v", err)
        return song, err
    }

    // Build the get item input
    input := &dynamodb.GetItemInput{
        TableName: aws.String("Song"),
        Key: map[string]*dynamodb.AttributeValue{
            "SongID": {
                S: aws.String(songId),
            },
        },
    }

    // Execute the get item request
    result, err := svc.GetItem(input)
    if err != nil {
        log.Printf("Error getting item from the table song: %v", err)
        return song, err
    }

    // Unmarshal the item into a struct
    err = dynamodbattribute.UnmarshalMap(result.Item, &song)
    if err != nil {
        log.Printf("Error unmarshaling item: %v", err)
        return song, err
    }

    return song, nil
}

func GetAlbum(albumID string) (models.Album, error) {
    var album models.Album

    // Get session
    svc, err := CreateDynamoDBSession()
    if err != nil {
        log.Printf("Error creating DynamoDB session: %v", err)
        return album, err
    }

    // Build the get item input
    input := &dynamodb.GetItemInput{
        TableName: aws.String("Album"),
        Key: map[string]*dynamodb.AttributeValue{
            "AlbumID": {
                S: aws.String(albumID),
            },
        },
    }

    // Execute the get item request
    result, err := svc.GetItem(input)
    if err != nil {
        log.Printf("Error getting item from the tabl albume: %v", err)
        return album, err
    }

    // Unmarshal the item into a struct
    err = dynamodbattribute.UnmarshalMap(result.Item, &album)
    if err != nil {
        log.Printf("Error unmarshaling item: %v", err)
        return album, err
    }

    return album, nil
}

func GetArtist(artistID string) (models.Artist, error) {
    var artist models.Artist

    // Get session
    svc, err := CreateDynamoDBSession()
    if err != nil {
        log.Printf("Error creating DynamoDB session: %v", err)
        return artist, err
    }

    // we have to use QueryInput instead of GetItemInput because we don't know the sort key
    // the structure of the table is partition key: ArtistID, sort key: ArtistName
    input := &dynamodb.QueryInput{
        TableName: aws.String("Artist"),
        KeyConditionExpression: aws.String("ArtistID = :artistId"),
        ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
            ":artistId": {
                S: aws.String(artistID),
            },
        },
        Select: aws.String("SPECIFIC_ATTRIBUTES"),
        ProjectionExpression: aws.String("ArtistName"),
        Limit: aws.Int64(1),
    }

    // Execute the get item request
    result, err := svc.Query(input)
    if err != nil {
        log.Printf("Error getting item from the table artist (name): %v", err)
        return artist, err
    }

    if result.LastEvaluatedKey != nil {
        log.Printf("DATA INCONSISTENCY: There are more than one result, this means you have duplicate ArtistIDs")
    }

    // Unmarshal the item into a struct
    err = dynamodbattribute.UnmarshalMap(result.Items[0], &artist)
    if err != nil {
        log.Printf("Error unmarshaling item: %v", err)
        return artist, err
    }

    return artist, nil
}

func GetAlbumName(albumID string) (string, error) {
    var albumName string

    // Get session
    svc, err := CreateDynamoDBSession()
    if err != nil {
        log.Printf("Error creating DynamoDB session: %v", err)
        return albumName, err
    }

    // Build the get item input
    input := &dynamodb.GetItemInput{
        TableName: aws.String("Album"),
        Key: map[string]*dynamodb.AttributeValue{
            "AlbumID": {
                S: aws.String(albumID),
            },
        },
        ProjectionExpression: aws.String("AlbumName"),
    }

    // Execute the get item request
    result, err := svc.GetItem(input)
    if err != nil {
        log.Printf("Error getting item from the table album: %v", err)
        return albumName, err
    }

    // Unmarshal the item into a struct
    var album struct {
        AlbumName string `json:"AlbumName"`
    }

    // Unmarshal the item into a struct
    err = dynamodbattribute.UnmarshalMap(result.Item, &album)
    if err != nil {
        log.Printf("Error unmarshaling item: %v", err)
        return albumName, err
    }

    albumName = album.AlbumName
    return albumName, nil
}

func GetArtistName(artistID string) (string, error) {
    var artistName string

    // Get session
    svc, err := CreateDynamoDBSession()
    if err != nil {
        log.Printf("Error creating DynamoDB session: %v", err)
        return artistName, err
    }

    // Build the get item input
    input := &dynamodb.GetItemInput{
        TableName: aws.String("Artist"),
        Key: map[string]*dynamodb.AttributeValue{
            "ArtistID": {
                S: aws.String(artistID),
            },
        },
        ProjectionExpression: aws.String("ArtistName"),
    }

    // Execute the get item request
    result, err := svc.GetItem(input)
    if err != nil {
        log.Printf("Error getting item from the table artist (name): %v", err)
        return artistName, err
    }

    // Unmarshal the item into a struct
    var artist struct {
        ArtistName string `json:"ArtistName"`
    }

    // Unmarshal the item into a struct
    err = dynamodbattribute.UnmarshalMap(result.Item, &artist)
    if err != nil {
        log.Printf("Error unmarshaling item: %v", err)
        return artistName, err
    }

    artistName = artist.ArtistName

    return artistName, nil
}

func SearchInSongTable(keyword string) ([]models.Song, error) {
    // Get session
    svc, err := CreateDynamoDBSession()
    if err != nil {
        log.Printf("Error creating DynamoDB session: %v", err)
    }
    var songs []models.Song

    // Build the query input
    input := &dynamodb.QueryInput{
        TableName: aws.String("Song"),
        IndexName: aws.String("SongName-SongID-index"),
        KeyConditionExpression: aws.String("SongName = :songName"),
        ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
            ":songName": {
                S: aws.String(keyword),
            },
        },
        ReturnConsumedCapacity: aws.String("TOTAL"),
    }

    // Execute the query
    result, err := svc.Query(input)
    if err != nil {
        return nil, err
    }

    // Unmarshal the query result into the songs slice
    err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &songs)
    if err != nil {
        return nil, err
    }

    return songs, nil
}

func SearchInArtistTable(keyword string) ([]models.Artist, error) {
    // Get session
    svc, err := CreateDynamoDBSession()
    if err != nil {
        log.Printf("Error creating DynamoDB session: %v", err)
    }
    var artists []models.Artist

    // Build the query input
    input := &dynamodb.QueryInput{
        TableName: aws.String("Artist"),
        IndexName: aws.String("ArtistName-ArtistID-index"),
        KeyConditionExpression: aws.String("ArtistName = :artistName"),
        ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
            ":artistName": {
                S: aws.String(keyword),
            },
        },
      ReturnConsumedCapacity: aws.String("TOTAL"),
    }

    // Execute the query
    result, err := svc.Query(input)
    if err != nil {
        return nil, err
    }

    // Unmarshal the query result into the songs slice
    err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &artists)
    if err != nil {
        return nil, err
    }

    return artists, nil
}

func SearchInAlbumTable(keyword string) ([]models.Album, error) {
    // Get session
    svc, err := CreateDynamoDBSession()
    if err != nil {
        log.Printf("Error creating DynamoDB session: %v", err)
    }
    var albums []models.Album

    // Build the query input
    input := &dynamodb.QueryInput{
        TableName: aws.String("Album"),
        IndexName: aws.String("AlbumName-AlbumID-index"),
        KeyConditionExpression: aws.String("AlbumName = :albumName"),
        ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
            ":albumName": {
                S: aws.String(keyword),
            },
        },
        ReturnConsumedCapacity: aws.String("TOTAL"),
    }

    // Execute the query
    result, err := svc.Query(input)
    if err != nil {
        return nil, err
    }

    // Unmarshal the query result into the songs slice
    err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &albums)
    if err != nil {
        return nil, err
    }

    return albums, nil
}


// GetSongFromS3ById Retrieves the specific song by ID of the S3 bucket
// This assumes the file is at the root
// Returns io.ReadCloser to be streamed to the requestor
func GetSongFromS3ById(songId string) (io.ReadCloser, error) {
    // Get session
    svc, err := CreateS3Session()
    if err != nil {
        log.Printf("Error creating S3 session: %v", err)
    }

    key := songId + ".mp3"
    params := &s3.GetObjectInput{
        Bucket: aws.String("beatbox-songbucket"), // S3 bucket name TODO: abstract
        Key:    aws.String(key),
    }

    // Send the GetObject request to S3
    result, err := svc.GetObject(params)
    if err != nil {
        log.Printf("Error getting object from S3: %v", err)
        return nil, err
    }

    return result.Body, nil
}

func InsertUser(email, name string) error {
    // Get session
    svc, err := CreateDynamoDBSession()
    if err != nil {
        log.Fatal(err)
    }

    // Create the user object
    item := models.User{
        Email:   email,
        Name:  name,
    }

    av, err := dynamodbattribute.MarshalMap(item)
    if err != nil {
        log.Printf("Got error marshalling new user item: %s", err)
        return err
    }

    // Create item in table User
    tableName := "User"

    input := &dynamodb.PutItemInput{
        Item:      av,
        TableName: aws.String(tableName),
        ConditionExpression: aws.String("attribute_not_exists(Email)"),
    }

    _, err = svc.PutItem(input)
    if err != nil {
        log.Printf("Got error calling PutItem: %s", err)
        return err
    }

    log.Print("Successfully added '" + item.Email + "' (" + name + ") to table " + tableName)

    return nil
}

// GetSongsAsUserDisplaySong Gets songs in format of
// models.UserDisplaySong and return a slice of songs
func GetSongsAsUserDisplaySong() ([]models.UserDisplaySong, error) {
    // Initialize the slice of UserDisplaySong structs
    var userDisplaySongs []models.UserDisplaySong

    // Get session
    svc, err := CreateDynamoDBSession()
    if err != nil {
        log.Printf("Error connecting to DynamoDB: %v", err)
        return userDisplaySongs, err
    }

    // Query the Song table to get all SongID, SongName, ArtistID, and AlbumID
    songs, err := svc.Scan(&dynamodb.ScanInput{
        TableName: aws.String("Song"),
        ProjectionExpression: aws.String("SongID, SongName, ArtistID, AlbumID"),
        })
    if err != nil {
        log.Printf("Error querying Song table: %v", err)
      return userDisplaySongs, err
    }

    // Create a map to store the artist name and album name for each song
    songDetails := make(map[string]models.UserDisplaySong)

    // Query the Artist table to get all ArtistID and ArtistName
    artists, err := svc.Scan(&dynamodb.ScanInput{
        TableName: aws.String("Artist"),
        ProjectionExpression: aws.String("ArtistID, ArtistName"),
        })
    if err != nil {
        log.Printf("Error querying Artist table: %v", err)
        return userDisplaySongs, err
    }

    // Add the artist name to the songDetails map for each song
    for _, artist := range artists.Items {
        artistId := *artist["ArtistID"].S
        artistName := *artist["ArtistName"].S
        for _, song := range songs.Items {
            if *song["ArtistID"].S == artistId {
                songId := *song["SongID"].S
                songName := *song["SongName"].S
                songDetails[songId] = models.UserDisplaySong{
                    SongID: songId,
                    SongName:songName,
                    ArtistName:artistName,
                    AlbumName: "",
                    }
            }
        }
    }

    // Query the Album table to get all AlbumID and AlbumName
    albums, err := svc.Scan(&dynamodb.ScanInput{
        TableName: aws.String("Album"),
        ProjectionExpression: aws.String("AlbumID, AlbumName"),
        })
    if err != nil {
        log.Printf("Error querying Album table: %v", err)
        return userDisplaySongs, err
    }

    // Add the album name to the songDetails map for each song
    for _, album := range albums.Items {
        albumId := *album["AlbumID"].S
        albumName := *album["AlbumName"].S
        for _,song := range songs.Items {
            if *song["AlbumID"].S == albumId {
                songId := *song["SongID"].S
                songDetail := songDetails[songId]
                songDetail.AlbumName = albumName
                songDetails[songId] = songDetail
            }
        }
    }

    // Convert the songDetails map to a slice of UserDisplaySong structs
    for _, songDetail := range songDetails {
        userDisplaySongs = append(userDisplaySongs,songDetail)
    }
    return userDisplaySongs, nil
}
