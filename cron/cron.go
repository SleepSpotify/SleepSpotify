package cron

import (
	"log"
	"time"

	"github.com/SleepSpotify/SleepSpotify/db"
	sleepspotify "github.com/SleepSpotify/SleepSpotify/spotify"
	"github.com/robfig/cron"
	"github.com/zmb3/spotify"
)

// InitCron the function to call to enable the cron every second
func InitCron() {
	c := cron.New()
	c.AddFunc("* * * * * *", cronSpotify)
	go c.Run()
}

// the function that will be called every second
func cronSpotify() {
	pauses, errDb := db.GetFromUts(time.Now().Unix())
	if errDb != nil {
		log.Println("DB FAILURE : ", errDb)
	}

	for _, pause := range pauses {
		tok, errTok := pause.GetToken()
		if errTok != nil {
			log.Println("MARSHAL FAILURE : ", errTok)
			continue
		}

		client := sleepspotify.GetClient(tok)
		go pauseSpotifyRoutine(client, pause)
	}
}

func pauseSpotifyRoutine(client spotify.Client, pause db.Sleep) {
	err := client.Pause()
	if err != nil {
		log.Println("SPOTIFY FAILURE : ", err)
	} else {
		pause.Delete()
	}
}
