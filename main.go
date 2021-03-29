package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"strings"

	"github.com/zmb3/spotify"
)

var (
	ServerListener chan bool
	clientID       *string
	clientSecret   *string
	playlistTarget *string
	port           *string
	redirectURL    string
	auth           spotify.Authenticator
)

var (
	state         = GenerateCrypticState()
	terminalGreen = "\033[32m"
	terminalRed   = "\033[31m"
	terminalReset = "\033[0m"
)

func main() {
	clientID = flag.String("clientId", "", "spotify client id")
	clientSecret = flag.String("clientSecret", "", "spotify client secret")
	playlistTarget = flag.String("targetPlaylist", "", "playlist to sync the library to")
	port = flag.String("port", "3000", "port to run the initial oauth server on")

	flag.Parse()

	redirectURL = "http://localhost" + ":" + *port + "/callback"
	auth = spotify.NewAuthenticator(redirectURL, spotify.ScopeUserReadPrivate, spotify.ScopeUserLibraryRead, spotify.ScopePlaylistModifyPublic)

	if runtime.GOOS == "windows" {
		terminalGreen = ""
		terminalRed = ""
		terminalReset = ""
	}

	initiateMusicSync()
}

func initiateMusicSync() {
	var clientInstance spotify.Client
	ServerListener = make(chan bool)

	auth.SetAuthInfo(*clientID, *clientSecret)
	url := auth.AuthURL(state)

	server := &http.Server{Addr: ":" + *port}

	fmt.Println("Created server instance")

	http.HandleFunc("/callback", func(rw http.ResponseWriter, r *http.Request) {
		clientInstance = getClientFromRequest(rw, r)
		ServerListener <- true
	})

	fmt.Println("Please open the below URL in a browser:\n" + url)

	go func(listener chan bool) {
		select {
		case <-listener:
			{
				err := server.Shutdown(context.Background())
				if err != nil {
					log.Fatal(err)
				}
			}
		}
	}(ServerListener)

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}

	// Start the track shift process
	handleTrackShift(clientInstance)
}

func getClientFromRequest(w http.ResponseWriter, r *http.Request) spotify.Client {
	// use the same state string here that you used to generate the URL
	token, err := auth.Token(state, r)
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusNotFound)
		return spotify.Client{}
	}

	// create a client using the specified token
	return auth.NewClient(token)
}

func handleTrackShift(client spotify.Client) {
	playlistId := *playlistTarget
	user, err := client.CurrentUser()

	if err != nil {
		log.Fatal("Failed to get user...\n Error:", err)
	}

	fmt.Println("Logged in as: ", user.DisplayName)

	fmt.Println("Getting library tracks...")

	tracks, err := client.CurrentUsersTracks()
	totalCount := tracks.Total
	fmt.Println(
		"Total Tracks in Library:" + colorString(terminalGreen, fmt.Sprint(totalCount)),
	)
	playlist, err := client.GetPlaylist(spotify.ID(playlistId))
	if err != nil {
		log.Fatal("Failed while trying to get playlist, Error:", err)
	}

	fmt.Println("Starting sync to move to playlist: ", playlist.Name)

	trackIdsInLibrary := getAllTrackIds(tracks, client)
	trackIdsInPlaylist := getAllPlaylistIds(playlist, client)

	uniqueTrackIds := []spotify.ID{}

	for _, idFromLib := range trackIdsInLibrary {
		if !doesSliceContain(trackIdsInPlaylist, idFromLib) {
			uniqueTrackIds = append(uniqueTrackIds, spotify.ID(idFromLib))
		}
	}

	fmt.Println(
		"Total Tracks to be Moved: " + colorString(terminalGreen, fmt.Sprint(len(uniqueTrackIds))),
	)

	if !askForConfirmation("Do you want to continue ?") {
		fmt.Println(colorString(terminalRed, "Cancelled"))
		return
	}

	batches := createPlaylistIDBatches(uniqueTrackIds, 100)

	for _, batch := range batches {
		_, err = client.AddTracksToPlaylist(spotify.ID(playlistId), batch...)
		if err != nil {
			log.Fatal("Failed to sync library with playlist,Error:", err)
		}
	}

	fmt.Println(colorString(terminalGreen, "Done!"))
}

func getAllTrackIds(tracks *spotify.SavedTrackPage, client spotify.Client) []string {
	storeRef := []string{}
	for page := 1; ; page++ {
		for _, track := range tracks.Tracks {
			storeRef = append(storeRef, track.ID.String())
		}
		err := client.NextPage(tracks)
		if err == spotify.ErrNoMorePages {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
	}
	return storeRef
}

func getAllPlaylistIds(playlist *spotify.FullPlaylist, client spotify.Client) []string {
	storeRef := []string{}
	for page := 1; ; page++ {
		for _, track := range playlist.Tracks.Tracks {
			storeRef = append(storeRef, string(track.Track.ID))
		}
		err := client.NextPage(&playlist.Tracks)
		if err == spotify.ErrNoMorePages {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
	}
	return storeRef
}

func doesSliceContain(dataSlice []string, toCompare string) bool {
	for _, value := range dataSlice {
		if value == toCompare {
			return true
		}
	}

	return false
}

func createPlaylistIDBatches(slice []spotify.ID, batch int) [][]spotify.ID {
	var batches [][]spotify.ID
	for i := 0; i < len(slice); i += batch {
		end := i + batch

		if end > len(slice) {
			end = len(slice)
		}

		batches = append(batches, slice[i:end])
	}

	return batches
}

func colorString(color string, toWrite string) string {
	var sb strings.Builder
	sb.WriteString(color)
	sb.WriteString(toWrite)
	sb.WriteString(terminalReset)
	return sb.String()
}

func askForConfirmation(s string) bool {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("%s [y/n]: ", s)

		response, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}

		response = strings.ToLower(strings.TrimSpace(response))

		if response == "y" || response == "yes" {
			return true
		} else if response == "n" || response == "no" {
			return false
		}
	}
}
