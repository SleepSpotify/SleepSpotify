package controler

import (
	"log"
	"net/http"

	"github.com/SleepSpotify/SleepSpotify/session"
	"github.com/SleepSpotify/SleepSpotify/spotify"
	restful "github.com/emicklei/go-restful"
)

// PUTPauseSpotifyControler controler to pause spotify
func PUTPauseSpotifyControler(req *restful.Request, resp *restful.Response) {

	ses, err := session.GetSessionSpotify(req.Request)
	if err != nil {
		log.Println("Session Failure : ", err)
		resp.WriteHeaderAndEntity(http.StatusUnauthorized, JSONError{"Server Error"})
		return
	}

	tok := session.GetToken(ses)
	if tok == nil {
		resp.WriteHeaderAndEntity(http.StatusUnauthorized, JSONError{"You are not connected"})
		return
	}

	client := spotify.GetClient(tok)
	client.Pause()

	resp.WriteEntity(JSONActionDone{"Musique Paused"})

}
