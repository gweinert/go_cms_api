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

//GetPages gets pages for a specific site
func GetPages(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	id, err := strconv.Atoi(ps.ByName("siteID"))
	if err != nil {
		http.Error(w, "Forbidden", 403)
	}

	p, err := model.GetPagesBySiteID(id)
	if err != nil {
		log.Fatal(err)
	}

	b, err := json.Marshal(p)
	if err != nil {
		fmt.Println("json err:", err)
	}

	fmt.Fprint(w, string(b))
}

//CreateNewPage creates a new page
func CreateNewPage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	r.ParseForm()
	fmt.Println("form", r.Form)
}

type pageSuccess struct {
	Success int         `json:"success"`
	Data    *model.Page `json:"data"`
}

//UpdatePage updates a page
func UpdatePage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	p := new(model.Page)

	if r.Body == nil {
		http.Error(w, "Please send a request body", 400)
		return
	}
	err := json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	fmt.Printf("%+v\n", p)

	_, err = model.SavePage(p)
	if err != nil {
		http.Error(w, "server Broke", 500)
		return
	}

	res := pageSuccess{Success: 1, Data: p}
	js, err := json.Marshal(res)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(js)

}
