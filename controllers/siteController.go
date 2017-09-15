package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	model "github.com/gweinert/cms_scratch/models"
	"github.com/julienschmidt/httprouter"
)

//ShowSiteDetailFunc needs comment
func ShowSiteDetailFunc(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	sessionIDCookie, err := r.Cookie("sessionId")
	if err != nil {
		log.Fatal(err)
		http.Error(w, err.Error(), 400)
		return
	}

	// if sessionIDCookie {
	fmt.Println("session id cookie", sessionIDCookie)
	sessionID := sessionIDCookie.Value
	// }

	user, err := model.GetUserFromSessionID(sessionID)
	if err != nil {
		log.Fatal(err)
		http.Error(w, err.Error(), 400)
		return
	}

	if user == nil {
		log.Fatal(err)
		http.Error(w, err.Error(), 400)
		return
	}

	s, err := model.GetSiteByUserID(user.ID)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	b, err := json.Marshal(s)
	if err != nil {
		fmt.Println("json err:", err)
		http.Error(w, err.Error(), 400)
		return
	}

	fmt.Fprint(w, string(b))

}

func PublishSite(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	sessionIDCookie, err := r.Cookie("sessionId")
	if err != nil {
		log.Fatal(err)
		http.Error(w, err.Error(), 400)
	}

	sessionID := sessionIDCookie.Value

	bucketName := strings.Split(r.Host, ":")[0]
	bucketName = strings.Join([]string{"garrett-react-cms", bucketName}, "-")

	fileURL, err := model.BuildStaticJsonAndUpload(sessionID, bucketName)
	if err != nil {
		log.Fatal(err)
	}

	resSuccess := struct {
		Success int    `json:"success"`
		FileURL string `json:"fileURL"`
	}{
		Success: 1,
		FileURL: fileURL,
	}

	b, err := json.Marshal(resSuccess)
	if err != nil {
		fmt.Println("json err:", err)
	}

	fmt.Fprint(w, string(b))
}
