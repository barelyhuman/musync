package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"
	"gopkg.in/yaml.v2"
)

type Config struct {
	ClientID       string `yaml:"client_id"`
	ClientSecret   string `yaml:"client_secret"`
	PlaylistTarget string `yaml:"playlist"`
	Port           string `yaml:"port"`
}

var (
	ServerListener chan bool
	redirectURL    string
	auth           spotify.Authenticator
)

var (
	state           = GenerateCrypticState()
	terminalGreen   = "\033[32m"
	terminalRed     = "\033[31m"
	terminalReset   = "\033[0m"
	globalConfig    *Config
	skipConfimation *bool
)

func main() {
	configPath := flag.String("c", "./musync.yaml", "point to config")
	skipConfimation = flag.Bool("skip", false, "skip confirmation question")

	flag.Parse()

	globalConfig = parseConfig(*configPath)

	fmt.Println("config received:", globalConfig)

	if globalConfig.Port == "" {
		globalConfig.Port = "3000"
	}

	redirectURL = "http://localhost" + ":" + globalConfig.Port + "/callback"
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

	token, err := readToken()

	if err != nil {
		log.Fatal(err)
	}

	authPassed := true

	clientInstance = auth.NewClient(token)
	_, err = clientInstance.CurrentUser()

	if err != nil {
		authPassed = false
	}

	auth.SetAuthInfo(globalConfig.ClientID, globalConfig.ClientSecret)
	url := auth.AuthURL(state)

	if !authPassed {
		server := &http.Server{Addr: ":" + globalConfig.Port}

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
	}
	// Start the track shift process
	handleTrackShift(clientInstance)
}

func getClientFromRequest(w http.ResponseWriter, r *http.Request) spotify.Client {
	// use the same state string here that you used to generate the URL
	token, err := auth.Token(state, r)

	saveToken(token)
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusNotFound)
		return spotify.Client{}
	}

	// create a client using the specified token
	return auth.NewClient(token)
}

type Token struct {
	AccessToken  string    `yaml:"access_token"`
	TokenType    string    `yaml:"type"`
	RefreshToken string    `yaml:"refresh_token"`
	Expiry       time.Time `yaml:"expiry"`
}

func readToken() (*oauth2.Token, error) {
	token := &oauth2.Token{}
	data, err := ioutil.ReadFile("token.yml")

	if err != nil {
		return token, nil
	}

	configToken := &Token{}
	yaml.Unmarshal(data, configToken)

	token.AccessToken = configToken.AccessToken
	token.RefreshToken = configToken.RefreshToken
	token.TokenType = configToken.TokenType
	token.Expiry = configToken.Expiry
	return token, nil
}

func saveToken(token *oauth2.Token) {
	var sb strings.Builder
	sb.WriteString("refresh_token: " + token.RefreshToken + "\n")
	sb.WriteString("type: " + token.TokenType + "\n")
	sb.WriteString("access_token: " + token.AccessToken + "\n")
	sb.WriteString("expiry: " + token.Expiry.String() + "\n")

	ioutil.WriteFile("token.yml", []byte(sb.String()), os.ModePerm)
}

func handleTrackShift(client spotify.Client) {
	playlistId := globalConfig.PlaylistTarget
	user, err := client.CurrentUser()

	if err != nil {
		log.Fatal("Failed to get user...\n Error:", err)
	}

	fmt.Println("Logged in as: ",
		colorString(terminalGreen, user.DisplayName),
	)

	fmt.Println("====\nProcessing Library\n====")

	tracks, err := client.CurrentUsersTracks()
	totalCount := tracks.Total
	fmt.Println(
		"Total Tracks in Library:" + colorString(terminalGreen, fmt.Sprint(totalCount)),
	)
	playlist, err := client.GetPlaylist(spotify.ID(playlistId))
	if err != nil {
		log.Fatal("Failed while trying to get playlist, Error:", err)
	}

	fmt.Println("Comparing Library and Playlist: ",
		colorString(terminalGreen, playlist.Name),
	)

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

	if !*skipConfimation && !askForConfirmation("Do you want to continue ?") {
		fmt.Println(colorString(terminalRed, "Cancelled"))
		return
	}

	if len(uniqueTrackIds) < 1 {
		fmt.Println(colorString(terminalRed, "Nothing to move!"))
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

func parseConfig(path string) *Config {
	fileData, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(`Error reading config, make sure you have config file 
		named musync.yaml or point 
		to another config using the -c flag`)
	}

	config := &Config{}
	yaml.Unmarshal([]byte(fileData), config)

	return config
}
