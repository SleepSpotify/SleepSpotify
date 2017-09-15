package controler

import (
	"log"
	"net/http"

	"github.com/Preovaleo/SleepSpotify/session"
	"github.com/Preovaleo/SleepSpotify/spotify"
)

// CallbackSpotifyControler Controler for save tokkens from Spotify
func CallbackSpotifyControler(w http.ResponseWriter, r *http.Request) {

	ses, errSes := session.GetSessionSpotify(r)
	if errSes != nil {
		http.Error(w, jsonErrMessage("Server Error"), http.StatusInternalServerError)
		log.Println("Session Failure : ", errSes)
		return
	}

	checkTok := session.GetToken(ses)
	if checkTok != nil {
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	state := session.GetState(ses)
	if state == "" {
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		return
	}

	tok, errTkn := spotify.GetToken(state, r)
	if errTkn != nil {
		http.Error(w, jsonErrMessage(errTkn.Error()), http.StatusBadRequest)
		return
	}

	session.SetToken(ses, tok)
	errSave := ses.Save(r, w)
	if errSave != nil {
		http.Error(w, jsonErrMessage(errSave.Error()), http.StatusInternalServerError)
		return
	}

	if !spotify.IsTokenValid(tok) {
		http.Redirect(w, r, "/login", http.StatusTemporaryRedirect)
		return
	}
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}

// LoginSpotifyControler controler to save a state to the session and to redirect to the spotify website
func LoginSpotifyControler(w http.ResponseWriter, r *http.Request) {
	ses, errSes := session.GetSessionSpotify(r)
	if errSes != nil {
		http.Error(w, jsonErrMessage("Server Error"), http.StatusInternalServerError)
		log.Println("Session Failure : ", errSes)
		return
	}

	state := session.SetSRandomState(ses)

	errSave := ses.Save(r, w)
	if errSave != nil {
		http.Error(w, jsonErrMessage(errSave.Error()), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, spotify.GetAuthURL(state), http.StatusTemporaryRedirect)
}
