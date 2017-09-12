package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	model "github.com/gweinert/cms_scratch/models"
	"github.com/julienschmidt/httprouter"
)

//ShowSiteDetailFunc needs comment
func ShowSiteDetailFunc(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	s, err := model.GetSiteByUserID(1)
	if err != nil {
		log.Fatal(err)
	}

	b, err := json.Marshal(s)
	if err != nil {
		fmt.Println("json err:", err)
	}

	fmt.Fprint(w, string(b))

}

func PublishSite(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var req = make(map[string]string)

	if r.Body == nil {
		http.Error(w, "Please send a request body", 400)
		return
	}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	fileURL, err := model.BuildStaticJsonAndUpload(req["sessionId"])
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
