package main

import (
  "context"
  "flag"
  "fmt"
  "log"
  "net/http"

  "golang.ngrok.com/ngrok"
  "golang.ngrok.com/ngrok/config"

  "BeatBox/handler"
)

func runNgrok(ctx context.Context, hostname string) error {
    l, err := ngrok.Listen(ctx,
        config.HTTPEndpoint(
            config.WithDomain(hostname),
        ),
        ngrok.WithAuthtokenFromEnv(),
        )
    if err != nil {
        return err
    }
    fmt.Println("Listening at ", l.Addr())

    return http.Serve(l, nil)
}

func runLocalhost(port string) error {
    log.Println("Attempting to start server on ", port)
    // this is a blocking function that will run indefinitely until stopped
    err := http.ListenAndServe(":" + port, nil)
    if err != nil {
        return err
    }
    return nil
}

func setHandlers() {
    // Endpoint handlers
    http.HandleFunc("/api/music/stream", handler.FetchSongById)
    http.HandleFunc("/api/music/search", handler.HandleKeywordSearch)
    http.HandleFunc("/api/user/insert", handler.InsertUserHandler)
    http.HandleFunc("/api/music/songs", handler.HandleFetchSongs)
}

func main() {
    // Define a flag named "port" with a default value of "8080"
    port := flag.String("port", "8080", "The port to start the server on if running locally (not useable under ngrok mode)")

    // Define a flag named "ngrok" with a default value of "false"
    useNgrok := flag.Bool("ngrok", false, "Enable running on ngrok. Default random domain")

    // Define a flag named "hostname" with a default value of ""
    hostname := flag.String("hostname", "", "The hostname to use if running ngrok (only useable if you already own the domain and are registered with ngrok)")

    // Parse the flags
    flag.Parse()
    // set http handlers
    setHandlers()

    // if ngrok flag run on ngrok
    // else localhost
    if *useNgrok {
        if err := runNgrok(context.Background(), *hostname); err != nil {
            log.Fatal(err)
        }
    } else {
        if err := runLocalhost(*port); err != nil {
            log.Fatal(err)
        }
    }
}
