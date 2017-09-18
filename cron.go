package main

import (
	"log"
	"time"

	"github.com/SleepSpotify/SleepSpotify/db"
	sleepspotify "github.com/SleepSpotify/SleepSpotify/spotify"
	"github.com/robfig/cron"
	"github.com/zmb3/spotify"
)

func initCron() {
	c := cron.New()
	c.AddFunc("* * * * * *", cronSpotify)
	c.Run()
}

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
		go func(client spotify.Client, pause db.Sleep) {
			err := client.Pause()
			if err != nil {
				log.Println("SPOTIFY FAILURE : ", err)
			} else {
				pause.Delete()
			}
		}(client, pause)
	}
}
