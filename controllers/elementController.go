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

	fmt.Fprint(w, string(b))
}

type deleteReq struct {
	GroupID int   `json:"groupId"`
	IDs     []int `json:"ids"`
}

// DeleteElement given an id from body, deletes and returns id of deleted element
func DeleteElements(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	// var req = make(map[string]string)
	// ids := make([]int, 0)
	req := new(deleteReq)

	if r.Body == nil {
		http.Error(w, "Please send a request body", 400)
		return
	}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		fmt.Println("error json")
		http.Error(w, err.Error(), 400)
		return
	}

	dIds, err := model.DeleteElements(req.IDs)
	if err != nil {
		http.Error(w, "server Broke", 500)
		return
	}

	resData := struct {
		Success int   `json:"success"`
		IDs     []int `json:"ids"`
		GroupID int   `json:"groupId"`
	}{
		Success: 1,
		IDs:     dIds,
		GroupID: req.GroupID,
	}
	b, err := json.Marshal(resData)
	if err != nil {
		http.Error(w, "server Broke", 500)
		return
	}

	fmt.Fprint(w, string(b))
}

// func DeleteGroupSlide(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
// 	var req = make(map[string]int)
// 	var res = make(map[string]int)

// 	if r.Body == nil {
// 		http.Error(w, "Please send a request body", 400)
// 		return
// 	}
// 	err := json.NewDecoder(r.Body).Decode(&req)
// 	if err != nil {
// 		fmt.Println("error json")
// 		http.Error(w, err.Error(), 400)
// 		return
// 	}

// 	delGroupID, delSlideIndex, err := model.DeleteElementGroupSlide(req["groupId"], req["sortOrder"])
// 	if err != nil {
// 		http.Error(w, "server Broke", 500)
// 		return
// 	}

// 	res["success"] = 1
// 	res["groupId"] = delGroupID
// 	res["slideIndex"] = delSlideIndex
// 	b, err := json.Marshal(res)
// 	if err != nil {
// 		http.Error(w, "server Broke", 500)
// 		return
// 	}

// 	fmt.Fprint(w, string(b))
// }
