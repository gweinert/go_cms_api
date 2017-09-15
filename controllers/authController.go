package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	model "github.com/gweinert/cms_scratch/models"
	"github.com/julienschmidt/httprouter"
)

// func GetGUID(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
// 	res := map[string]int
// 	res["GUID"] = time.Now().Unix()

// 	b, err := json.Marshal(res)
// 	if err != nil {
// 		fmt.Println("json err:", err)
// 	}

// 	fmt.Fprint(w, string(b))
// }

func BasicAuth(h httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		// Get the Basic Authentication credentials
		// user, password, hasAuth := r.BasicAuth()
		sessionIDCookie, err := r.Cookie("sessionId")
		if err != nil {
			log.Fatal(err)
			http.Error(w, err.Error(), 400)
		}

		sessionID := sessionIDCookie.Value
		user, err := model.GetUserFromSessionID(sessionID)
		if err != nil {
			log.Fatal(err)
			http.Error(w, err.Error(), 400)
		}

		if user != nil {
			// Delegate request to the given handle
			h(w, r, ps)
		} else {
			// Request Basic Authentication otherwise
			w.Header().Set("WWW-Authenticate", "Basic realm=Restricted")
			http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
		}
	}
}

type loginReq struct {
	Email string `json:"email"`
	Hash  string `json:"hash"`
}

type loginRes struct {
	Success   int         `json:"success"`
	User      *model.User `json:"user"`
	SessionID string      `json:"sessionId"`
}

func Login(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	loginReq := new(loginReq)
	loginRes := new(loginRes)

	if r.Body == nil {
		http.Error(w, "Please send a request body", 400)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&loginReq)
	if err != nil {
		fmt.Println("error json")
		http.Error(w, err.Error(), 400)
		return
	}

	user, sessionID, err := model.LoginUser(loginReq.Email, loginReq.Hash)
	if err != nil {
		fmt.Println("error user")
		http.Error(w, err.Error(), 400)
		return
	}

	if user != nil {
		loginRes.Success = 1
		loginRes.User = user
		loginRes.SessionID = sessionID
	} else {
		loginRes.Success = 0
		loginRes.User = nil
		loginRes.SessionID = ""
	}

	b, err := json.Marshal(loginRes)
	if err != nil {
		http.Error(w, "server Broke", 500)
		return
	}
	cookie := http.Cookie{Name: "sessionId", Value: sessionID}
	http.SetCookie(w, &cookie)

	fmt.Fprint(w, string(b))

}

func GetUserFromSessionID(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var req = make(map[string]string)
	var loginRes = new(loginRes)

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

	user, err := model.GetUserFromSessionID(req["id"])
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	if user != nil {
		user.PasswordHash = ""
		user.PasswordSalt = ""
		loginRes.Success = 1
		loginRes.User = user
		loginRes.SessionID = req["id"]
	} else {
		loginRes.Success = 0
		loginRes.User = nil
		loginRes.SessionID = ""
	}

	b, err := json.Marshal(loginRes)
	if err != nil {
		http.Error(w, "server Broke", 500)
		return
	}

	fmt.Fprint(w, string(b))
}
