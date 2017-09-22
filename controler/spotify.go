package controler

import (
	"fmt"
	"log"
	"net/http"

	"github.com/SleepSpotify/SleepSpotify/config"
	"github.com/SleepSpotify/SleepSpotify/session"
	"github.com/SleepSpotify/SleepSpotify/spotify"
)

// CallbackSpotifyControler Controler for save tokkens from Spotify
func CallbackSpotifyControler(w http.ResponseWriter, r *http.Request) {

	ses, errSes := session.GetSessionSpotify(r)
	if errSes != nil {
		log.Println("SESSION FAILURE : ", errSes)
		http.Error(w, "Server Error", http.StatusInternalServerError)
		return
	}

	checkTok := session.GetToken(ses)
	if checkTok != nil {
		http.Redirect(w, r, fmt.Sprintf("%s/timer", config.GetConfig().Angular), http.StatusTemporaryRedirect)
		return
	}

	state := session.GetState(ses)
	if state == "" {
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		return
	}

	tok, errTkn := spotify.GetToken(state, r)
	if errTkn != nil {
		http.Error(w, errTkn.Error(), http.StatusBadRequest)
		return
	}

	session.SetToken(ses, tok)
	errSave := ses.Save(r, w)
	if errSave != nil {
		log.Println("SESSION FAILURE : ", errSave)
		http.Error(w, "Server Error", http.StatusInternalServerError)
		return
	}

	if !spotify.IsTokenValid(tok) {
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		return
	}
	http.Redirect(w, r, fmt.Sprintf("%s/timer", config.GetConfig().Angular), http.StatusTemporaryRedirect)
}

// LoginSpotifyControler controler to save a state to the session and to redirect to the spotify website
func LoginSpotifyControler(w http.ResponseWriter, r *http.Request) {
	ses, errSes := session.GetSessionSpotify(r)
	if errSes != nil {
		log.Println("SESSION FAILURE : ", errSes)
		http.Error(w, "Server Error", http.StatusInternalServerError)
		return
	}

	tok := session.GetToken(ses)
	if tok != nil {
		http.Redirect(w, r, fmt.Sprintf("%s/timer", config.GetConfig().Angular), http.StatusTemporaryRedirect)
		return
	}

	state := session.SetSRandomState(ses)

	errSave := ses.Save(r, w)
	if errSave != nil {
		log.Println("SESSION FAILURE : ", errSave)
		http.Error(w, "Server Error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, spotify.GetAuthURL(state), http.StatusTemporaryRedirect)
}
