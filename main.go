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

	wsSpotify := new(restful.WebService)
	wsSpotify.Path("/spotify").Produces(restful.MIME_XML, restful.MIME_JSON)
	wsSpotify.Route(wsSpotify.PUT("/pause").To(controler.PUTPauseSpotifyControler))
	wsSpotify.Route(wsSpotify.GET("/sleep").To(controler.GETSleep))
	wsSpotify.Route(wsSpotify.POST("/sleep").To(controler.POSTSleep))
	wsSpotify.Route(wsSpotify.PUT("/sleep").To(controler.PUTSleep))
	wsSpotify.Route(wsSpotify.DELETE("/sleep").To(controler.DELETESleep))
	restful.Add(wsSpotify)

	wsAccount := new(restful.WebService)
	wsAccount.Path("/account").Produces(restful.MIME_XML, restful.MIME_JSON)
	wsAccount.Route(wsAccount.GET("/isConnected").To(controler.IsConnected))
	restful.Add(wsAccount)

	cors := restful.CrossOriginResourceSharing{
		AllowedHeaders: []string{"Content-Type", "Accept"},
		AllowedDomains: []string{config.Angular},
		CookiesAllowed: true,
		Container:      restful.DefaultContainer,
	}
	restful.DefaultContainer.Filter(cors.Filter)
	restful.DefaultContainer.Filter(restful.DefaultContainer.OPTIONSFilter)

	http.HandleFunc("/callback", controler.CallbackSpotifyControler)
	http.HandleFunc("/login", controler.LoginSpotifyControler)

	go initCron()

	err := http.ListenAndServe(config.DomainName, nil)
	if err != nil {
		log.Fatalf("Listen and Serve : %s", err)
	}

	log.Println("END OF SleepSpotify")
}
