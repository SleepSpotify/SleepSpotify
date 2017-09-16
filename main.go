package main

import (
	"log"
	"net/http"

	"github.com/SleepSpotify/SleepSpotify/config"
	"github.com/SleepSpotify/SleepSpotify/controler"
	"github.com/SleepSpotify/SleepSpotify/session"
	"github.com/SleepSpotify/SleepSpotify/spotify"
	restful "github.com/emicklei/go-restful"
)

func main() {
	log.Println("START OF SleepSpotify")

	config, errConfig := config.ReadConfig()
	if errConfig != nil {
		log.Fatalf("Problem with the config file : %s", errConfig)
	}

	session.InitStore(config)
	spotify.InitAuth(config)

	ws := new(restful.WebService)
	ws.Path("/spotify").Produces(restful.MIME_JSON, restful.MIME_XML)
	ws.Route(ws.PUT("/pause").To(controler.PUTPauseSpotifyControler))
	restful.Add(ws)

	http.HandleFunc("/callback", controler.CallbackSpotifyControler)
	http.HandleFunc("/login", controler.LoginSpotifyControler)

	err := http.ListenAndServe(config.DomainName, nil)
	if err != nil {
		log.Fatalf("Listen and Serve : %s", err)
	}

	log.Println("END OF SleepSpotify")
}
