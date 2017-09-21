package controler

import (
	"log"
	"net/http"

	"github.com/SleepSpotify/SleepSpotify/session"
	"github.com/SleepSpotify/SleepSpotify/spotify"
	restful "github.com/emicklei/go-restful"
)

// IsConnected controler to know if the user is connected
func IsConnected(req *restful.Request, resp *restful.Response) {
	ses, err := session.GetSessionSpotify(req.Request)
	if err != nil {
		log.Println("Session Failure : ", err)
		resp.WriteHeaderAndEntity(http.StatusUnauthorized, jsonRep{"Server Error"})
		return
	}

	tok := session.GetToken(ses)
	if tok == nil {
		resp.WriteEntity(jsonConnected{false, ""})
		return
	}

	client := spotify.GetClient(tok)
	User, errSpot := client.CurrentUser()
	if errSpot != nil {
		log.Println("Spotify Failure : ", errSpot)
		resp.WriteHeaderAndEntity(http.StatusUnauthorized, jsonRep{"Server Error"})
		return
	}

	resp.WriteEntity(jsonConnected{true, User.DisplayName})
}

type jsonConnected struct {
	IsConnected bool
	UserName    string
}
