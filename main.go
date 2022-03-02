package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/zmb3/spotify"
)

var (
	state            = generateCrypticState()
	ServerListener   chan bool
	redirectURL      string
	auth             spotify.Authenticator
	defaultTokenPath string
	prompt           = colorString(terminalGreen, "\n    >> ")
)

var (
	globalConfig    *Config
	skipConfimation *bool
)

func bail(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	homeDir, err := os.UserHomeDir()
	bail(err)

	defaultConfigPath := path.Join(homeDir, ".config", "musync", "musync.yml")
	defaultTokenPath = path.Join(homeDir, ".config", "musync", "token.yml")

	configPath := flag.String("c", defaultConfigPath, "point to config")
	skipConfimation = flag.Bool("skip", false, "skip the confirmation")

	// enable skipping the confirmation
	// when running inside a CI
	if os.Getenv("IS_CI") == "true" {
		*skipConfimation = true
	}

	flag.Parse()

	// parse the global config based on the flag paths
	config := &Config{}
	config.parseConfig(*configPath)

	globalConfig = config

	if globalConfig.Port == "" {
		globalConfig.Port = "3000"
	}

	redirectURL = "http://localhost" + ":" + globalConfig.Port + "/callback"
	auth = spotify.NewAuthenticator(redirectURL, spotify.ScopeUserReadPrivate, spotify.ScopeUserLibraryRead, spotify.ScopePlaylistModifyPublic)

	initiateMusicSync()
}

func initiateMusicSync() {
	var clientInstance spotify.Client
	var authPassed bool

	ServerListener = make(chan bool)

	token, err := readToken(defaultTokenPath)

	if err != nil {
		authPassed = false
	} else if token == nil {
		authPassed = false
	} else {
		authPassed = true
	}

	clientInstance = auth.NewClient(token)
	clientInstance.AutoRetry = true

	_, err = clientInstance.CurrentUser()

	if err != nil {
		authPassed = false
	}

	auth.SetAuthInfo(globalConfig.ClientID, globalConfig.ClientSecret)

	url := auth.AuthURL(state)

	if !authPassed {
		server := &http.Server{Addr: ":" + globalConfig.Port}

		fmt.Printf("\nStarting local verification server...\n\n")

		http.HandleFunc("/callback", func(rw http.ResponseWriter, r *http.Request) {
			clientInstance = getClientFromRequest(rw, r)
			rw.Write([]byte("Logged In!"))
			ServerListener <- true
		})

		fmt.Println("Please open the below URL in a browser:\n" + colorString(terminalGreen, url))

		go func(listener chan bool) {
			if <-listener {
				err := server.Shutdown(context.Background())
				if err != nil {
					log.Fatal(err)
				}
			}
		}(ServerListener)

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}
	// Start the track shift process
	handleTrackShift(clientInstance)
}

func getClientFromRequest(w http.ResponseWriter, r *http.Request) spotify.Client {
	// use the same state string here that you used to generate the URL
	token, err := auth.Token(state, r)

	saveToken(token, defaultTokenPath)
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusNotFound)
		return spotify.Client{}
	}

	// create a client using the specified token
	return auth.NewClient(token)
}

func handleTrackShift(client spotify.Client) {
	playlistId := globalConfig.PlaylistTarget
	user, err := client.CurrentUser()

	if err != nil {
		log.Fatal("Failed to get user...\n Error:", err)
	}

	fmt.Println("\nLogged in as:",
		colorString(terminalGreen, user.DisplayName),
	)

	fmt.Println(prompt + "Processing Library")

	tracks, err := client.CurrentUsersTracks()

	bail(err)

	totalCount := tracks.Total
	fmt.Println(prompt+
		"Total Tracks in Library:", colorString(terminalGreen, fmt.Sprint(totalCount)),
	)
	playlist, err := client.GetPlaylist(spotify.ID(playlistId))
	if err != nil {
		log.Fatal("Failed while trying to get playlist, Error:", err)
	}

	fmt.Println(prompt+"Comparing Library and Playlist:",
		colorString(terminalGreen, playlist.Name),
	)

	fmt.Println(prompt + "Please wait...")

	trackIdsInLibrary := getAllTrackIds(tracks, client)
	trackIdsInPlaylist := getAllPlaylistIds(playlist, client)

	uniqueTrackIds := []spotify.ID{}

	for _, idFromLib := range trackIdsInLibrary {
		if !doesSliceContain(trackIdsInPlaylist, idFromLib) {
			uniqueTrackIds = append(uniqueTrackIds, spotify.ID(idFromLib))
		}
	}

	fmt.Println(
		prompt + "Total Tracks to be Moved: " + colorString(terminalGreen, fmt.Sprint(len(uniqueTrackIds))),
	)

	if !*skipConfimation && !confirmSyncPrompt("Do you want to continue ?") {
		fmt.Println(colorString(terminalRed, "\n    Cancelled"))
		return
	}

	if len(uniqueTrackIds) < 1 {
		fmt.Println(colorString(terminalGreen, "\n    Nothing to move!"))
	}

	batches := getPlaylistIDChunks(uniqueTrackIds, 100)

	for _, batch := range batches {
		_, err = client.AddTracksToPlaylist(spotify.ID(playlistId), batch...)
		if err != nil {
			log.Fatal("Failed to sync library with playlist,Error:", err)
		}
	}

	fmt.Println(colorString(terminalGreen, "\n    Done!"))
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

func getPlaylistIDChunks(slice []spotify.ID, batch int) [][]spotify.ID {
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

func confirmSyncPrompt(s string) bool {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Printf("\n    %s [y/n]: ", s)

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
