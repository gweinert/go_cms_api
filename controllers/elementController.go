package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	model "github.com/gweinert/cms_scratch/models"
	"github.com/julienschmidt/httprouter"
)

//GetElements gets all elements of a given page
func GetElements(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	pageID, err := strconv.Atoi(ps.ByName("pageID"))
	if err != nil {
		http.Error(w, "Forbidden", 403)
	}

	els, err := model.GetElementsByPageID(pageID)
	if err != nil {
		log.Fatal(err)
	}

	b, err := json.Marshal(els)
	if err != nil {
		fmt.Println("json err:", err)
	}

	fmt.Fprint(w, "%s", string(b))
}
