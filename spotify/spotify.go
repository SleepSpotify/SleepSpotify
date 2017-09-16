package spotify

import (
	"fmt"
	"net/http"

	"golang.org/x/oauth2"

	"github.com/SleepSpotify/SleepSpotify/config"
	"github.com/zmb3/spotify"
)

var auth spotify.Authenticator

// InitAuth function to init the lib Spotify with the config
func InitAuth(config config.Config) {
	auth = spotify.NewAuthenticator(fmt.Sprintf("http://%s/callback", config.DomainName), spotify.ScopeUserModifyPlaybackState)
	auth.SetAuthInfo(config.Spotify.ClientID, config.Spotify.ClientSecret)
}

// IsTokenValid function to know if the token is valid
func IsTokenValid(token *oauth2.Token) bool {
	client := auth.NewClient(token)
	_, err := client.CurrentUser()
	if err != nil {
		return false
	}
	return true
}

// GetToken function to get the oauthToken from a request
func GetToken(state string, r *http.Request) (*oauth2.Token, error) {
	return auth.Token(state, r)
}

// GetAuthURL function to get the authentification url on spotify website
func GetAuthURL(state string) string {
	return auth.AuthURL(state)
}

// GetClient Get the client from spotify API
func GetClient(tok *oauth2.Token) spotify.Client {
	return auth.NewClient(tok)
}
