package session

import (
	"github.com/SleepSpotify/SleepSpotify/config"
	"github.com/gorilla/sessions"
)

var store sessions.Store

// InitStore func to init the session with the config session secret
func InitStore(config config.Config) {
	store = sessions.NewCookieStore([]byte(config.SessionSecret))
}
