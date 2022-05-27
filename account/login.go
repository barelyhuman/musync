package account

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/twinj/uuid"
	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"
)

var (
	state          = generateCrypticState()
	serverListener chan bool
)

type Auth struct {
	authenticator spotify.Authenticator
	Client        spotify.Client
	RedirectURL   string
	Token         *oauth2.Token
	Storage       *TokenStorage
}

type AuthConfig struct {
	RedirectURL string
	token       *oauth2.Token
}

type PersistDetails struct {
	RedirectURL string
	Token       *oauth2.Token
}

// generateCrypticState - generate a random string to re-evaluate spotify login
func generateCrypticState() string {
	id := uuid.NewV4()
	return id.String()
}

func (auth *Auth) getClientFromRequest(w http.ResponseWriter, r *http.Request) {
	token, err := auth.authenticator.Token(state, r)
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusNotFound)
		return
	}
	auth.Client = auth.authenticator.NewClient(token)
	auth.Token = token

}

func (auth *Auth) StartVerificationServer(port string, onComplete chan bool) {
	serverListener = make(chan bool)
	server := &http.Server{Addr: ":" + port}

	fmt.Println("Starting verification server:", port)

	http.HandleFunc("/callback", func(rw http.ResponseWriter, r *http.Request) {
		auth.getClientFromRequest(rw, r)
		rw.Write([]byte("Logged In!"))
		serverListener <- true
		onComplete <- true
	})

	go func(listener chan bool) {
		if <-listener {
			err := server.Shutdown(context.Background())
			if err != nil {
				log.Fatal(err)
			}
		}
	}(serverListener)

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

func (auth *Auth) PingUser() bool {
	user, err := auth.Client.CurrentUser()
	if err != nil {
		return false
	}
	if len(user.DisplayName) > 0 {
		return true
	}
	return false
}

func (auth *Auth) SaveToken() error {
	persist := PersistDetails{
		Token:       auth.Token,
		RedirectURL: auth.RedirectURL,
	}
	auth.Storage.SaveToken(persist)
	return nil
}

func (auth *Auth) GetAuthURL(clientId string, secretKey string) string {
	auth.authenticator.SetAuthInfo(clientId, secretKey)
	return auth.authenticator.AuthURLWithDialog(state)
}

func NewAuth(params AuthConfig) *Auth {
	auth := &Auth{}
	redirectURL := params.RedirectURL
	auth.authenticator = spotify.NewAuthenticator(redirectURL, spotify.ScopeUserReadPrivate, spotify.ScopeUserLibraryRead, spotify.ScopePlaylistModifyPublic)
	auth.Storage = NewTokenStorage()
	if params.token != nil && len(params.token.AccessToken) > 0 {
		auth.Client = auth.authenticator.NewClient(params.token)
	}
	return auth
}

func AuthFromToken() *Auth {
	auth := &Auth{}
	var persist PersistDetails
	auth.Storage = NewTokenStorage()
	auth.Storage.ReadToken(&persist)
	auth.authenticator = spotify.NewAuthenticator(persist.RedirectURL, spotify.ScopeUserReadPrivate, spotify.ScopeUserLibraryRead, spotify.ScopePlaylistModifyPublic)
	auth.Client = auth.authenticator.NewClient(persist.Token)
	return auth
}
