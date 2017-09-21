package controler

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/SleepSpotify/SleepSpotify/db"
	"github.com/SleepSpotify/SleepSpotify/session"
	restful "github.com/emicklei/go-restful"
)

// GETSleep controler to get the uts when the music will stop
func GETSleep(req *restful.Request, resp *restful.Response) {
	ses, errSes := session.GetSessionSpotify(req.Request)
	if errSes != nil {
		log.Println("Session Failure : ", errSes)
		resp.WriteHeaderAndEntity(http.StatusInternalServerError, jsonRep{"Server Error"})
		return
	}

	tok := session.GetToken(ses)
	if tok == nil {
		resp.WriteHeaderAndEntity(http.StatusUnauthorized, jsonRep{"You are not connected"})
		return
	}

	id, errID := session.GetIDPause(ses)
	if errID != nil {
		resp.WriteHeaderAndEntity(http.StatusInternalServerError, jsonRep{"Server Error"})
		return
	}

	if id == 0 {
		resp.WriteEntity(jsonSleepNbr{0, nil})
		return
	}

	pause, errDB := db.GetFromID(id)
	if errDB != nil {
		log.Println("DB failure : ", errDB)
		resp.WriteHeaderAndEntity(http.StatusInternalServerError, jsonRep{"Server Error"})
		return
	}

	ret := []db.Sleep{pause}
	resp.WriteEntity(jsonSleepNbr{1, ret})
}

// POSTSleep controler to set a sleep uts
func POSTSleep(req *restful.Request, resp *restful.Response) {
	ses, errSes := session.GetSessionSpotify(req.Request)
	if errSes != nil {
		log.Println("Session Failure : ", errSes)
		resp.WriteHeaderAndEntity(http.StatusInternalServerError, jsonRep{"Server Error"})
		return
	}

	tok := session.GetToken(ses)
	if tok == nil {
		resp.WriteHeaderAndEntity(http.StatusUnauthorized, jsonRep{"You are not connected"})
		return
	}

	uts, errPar := getUTS(req)
	if errPar != nil {
		resp.WriteHeaderAndEntity(http.StatusBadRequest, jsonRep{"Bad shaped uts or unset"})
		return
	}

	if uts < time.Now().Unix() {
		resp.WriteHeaderAndEntity(http.StatusBadRequest, jsonRep{"UTS allready passed"})
		return
	}

	pause, errSlp := db.NewSleep(0, tok, uts)
	if errSlp != nil {
		log.Println("UNMARSHAL Failure : ", errSlp)
		resp.WriteHeaderAndEntity(http.StatusInternalServerError, jsonRep{"Server Error"})
		return
	}

	errDB := pause.Insert()
	if errDB != nil {
		log.Println("DB Failure : ", errDB)
		resp.WriteHeaderAndEntity(http.StatusInternalServerError, jsonRep{"Server Error"})
		return
	}

	session.SetIDPause(ses, pause.ID)
	resp.WriteEntity(pause)
}

// PUTSleep controler to update a sleep uts
func PUTSleep(req *restful.Request, resp *restful.Response) {
	ses, errSes := session.GetSessionSpotify(req.Request)
	if errSes != nil {
		log.Println("Session Failure : ", errSes)
		resp.WriteHeaderAndEntity(http.StatusInternalServerError, jsonRep{"Server Error"})
		return
	}

	tok := session.GetToken(ses)
	if tok == nil {
		resp.WriteHeaderAndEntity(http.StatusUnauthorized, jsonRep{"You are not connected"})
		return
	}

	ID, errID := session.GetIDPause(ses)
	if errID != nil {
		log.Println("SESSION Failure : ", errID)
		resp.WriteHeaderAndEntity(http.StatusInternalServerError, jsonRep{"Server Error"})
		return
	}

	pause, errGET := db.GetFromID(ID)
	if errGET != nil {
		log.Println("DB Failure : ", errGET)
		resp.WriteHeaderAndEntity(http.StatusInternalServerError, jsonRep{"Server Error"})
		return
	}

	uts, errPar := getUTS(req)
	if errPar != nil {
		resp.WriteHeaderAndEntity(http.StatusBadRequest, jsonRep{"Bad shaped uts"})
		return
	}

	pause.Uts = uts

	errUp := pause.Update()
	if errUp != nil {
		log.Println("DB Failure : ", errUp)
		resp.WriteHeaderAndEntity(http.StatusInternalServerError, jsonRep{"Server Error"})
		return
	}

	resp.WriteEntity(pause)

}

// DELETESleep controler to delete a sleep uts
func DELETESleep(req *restful.Request, resp *restful.Response) {
	ses, errSes := session.GetSessionSpotify(req.Request)
	if errSes != nil {
		log.Println("Session Failure : ", errSes)
		resp.WriteHeaderAndEntity(http.StatusInternalServerError, jsonRep{"Server Error"})
		return
	}

	tok := session.GetToken(ses)
	if tok == nil {
		resp.WriteHeaderAndEntity(http.StatusUnauthorized, jsonRep{"You are not connected"})
		return
	}

	ID, errID := session.GetIDPause(ses)
	if errID != nil {
		log.Println("SESSION Failure : ", errID)
		resp.WriteHeaderAndEntity(http.StatusInternalServerError, jsonRep{"Server Error"})
		return
	}

	pause, errGET := db.GetFromID(ID)
	if errGET != nil {
		log.Println("DB Failure : ", errGET)
		resp.WriteHeaderAndEntity(http.StatusInternalServerError, jsonRep{"Server Error"})
		return
	}

	errDel := pause.Delete()
	if errDel != nil {
		log.Println("DB Failure : ", errDel)
		resp.WriteHeaderAndEntity(http.StatusInternalServerError, jsonRep{"Server Error"})
		return
	}

	resp.WriteEntity(jsonRep{"Done"})
}

func getUTS(req *restful.Request) (int64, error) {
	utsString, errPar := req.BodyParameter("uts")
	if errPar != nil {
		return 0, errPar
	}

	utsInt, errInt := strconv.ParseInt(utsString, 10, 64)
	if errInt != nil {
		return 0, errInt
	}

	return utsInt, nil
}
