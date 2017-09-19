package main

import (
	"log"
	"net/http"

	"github.com/SleepSpotify/SleepSpotify/config"
	"github.com/SleepSpotify/SleepSpotify/controler"
	"github.com/SleepSpotify/SleepSpotify/db"
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
	db.InitDB(config)

	errDB := db.DB.Ping()
	if errDB != nil {
		log.Fatal("DB FAILURE : ", errDB)
	}

	ws := new(restful.WebService)
	ws.Path("/spotify").Produces(restful.MIME_XML, restful.MIME_JSON)
	ws.Route(ws.PUT("/pause").To(controler.PUTPauseSpotifyControler))
	ws.Route(ws.GET("/sleep").To(controler.GETSleep))
	ws.Route(ws.POST("/sleep").To(controler.POSTSleep))
	ws.Route(ws.PUT("/sleep").To(controler.PUTSleep))
	ws.Route(ws.DELETE("/sleep").To(controler.DELETESleep))
	restful.Add(ws)

	http.HandleFunc("/callback", controler.CallbackSpotifyControler)
	http.HandleFunc("/login", controler.LoginSpotifyControler)

	go initCron()

	err := http.ListenAndServe(config.DomainName, nil)
	if err != nil {
		log.Fatalf("Listen and Serve : %s", err)
	}

	log.Println("END OF SleepSpotify")
}
