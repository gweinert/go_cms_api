package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	model "github.com/gweinert/cms_scratch/models"
	"github.com/julienschmidt/httprouter"
)

type groupedPageElements struct {
	Group    *model.ElementGroup `json:"group"`
	Elements []*model.Element    `json:"elements"`
}

type newGroupResponse struct {
	Success int                  `json:"success"`
	Data    *groupedPageElements `json:"data"`
}

// CreateNewGroup turns json into group struct and sends to model
func CreateNewGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	g := new(model.ElementGroup)
	ng := new(model.ElementGroup)
	pes := make([]*model.Element, 0)
	data := new(groupedPageElements)

	if r.Body == nil {
		http.Error(w, "Please send a request body", 400)
		return
	}
	err := json.NewDecoder(r.Body).Decode(&g)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	pes, ng, err = model.CreateNewGroup(g)
	if err != nil {
		log.Fatal(err)
		http.Error(w, err.Error(), 500)

	}

	data.Group = ng
	data.Elements = pes
	res := newGroupResponse{Success: 1, Data: data}

	b, err := json.Marshal(res)
	if err != nil {
		fmt.Println("json err:", err)
		http.Error(w, err.Error(), 500)
	}

	fmt.Fprint(w, string(b))
}
