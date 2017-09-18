package session

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/gorilla/sessions"
	"golang.org/x/oauth2"
)

// GetSessionSpotify function to get the session from the spotify part
func GetSessionSpotify(r *http.Request) (*sessions.Session, error) {
	return store.Get(r, "spotify")
}

// GetState get the State from Spotify
func GetState(session *sessions.Session) string {
	state, ok := session.Values["SpotifyState"]
	if !ok {
		return ""
	}
	return state.(string)
}

// SetSRandomState generate and state random state for spotify
func SetSRandomState(session *sessions.Session) string {
	rnd := randStringBytesMaskImprSrc(25)
	session.Values["SpotifyState"] = rnd
	return rnd
}

// GetToken set the token from Spotify
func GetToken(session *sessions.Session) *oauth2.Token {
	tokjson, ok := session.Values["SpotifyToken"]
	if !ok {
		return nil
	}
	var tok oauth2.Token
	errUnmar := json.Unmarshal(tokjson.([]byte), &tok)
	if errUnmar != nil {
		log.Println("TOKEN UNMARSHAL error : ", errUnmar)
		return nil
	}
	return &tok
}

// SetToken set the token from Spotify
func SetToken(session *sessions.Session, token *oauth2.Token) {
	ret, err := json.Marshal(token)
	if err != nil {
		ret = []byte("{}")
		log.Println("TOKEN MARSHAL error : ", err)
	}
	session.Values["SpotifyToken"] = ret
}

// GetIDPause function to retrieve the ID pause from the session
func GetIDPause(session *sessions.Session) (int64, error) {
	i, ok := session.Values["SpotifyID"]
	if !ok {
		return 0, nil
	}

	iInt, ok := i.(int64)
	if !ok {
		return 0, errors.New("Bad shapped ID")
	}

	return iInt, nil
}

// SetIDPause function to set the ID pause in the session
func SetIDPause(session *sessions.Session, i int64) {
	session.Values["SpotifyID"] = i
}
