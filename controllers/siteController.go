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
