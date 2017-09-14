package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gweinert/cms_scratch/services"
	"github.com/julienschmidt/httprouter"
)

type contactReq struct {
	ToEmail string `json:"toEmail"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Subject string `json:"subject"`
	Message string `json:"message"`
}

// SendContactMail sends an email using mailgun
func SendContactMail(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	req := new(contactReq)

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

	body := strings.Join([]string{
		"From: \n\t Name:",
		req.Name,
		"\n\t Email: ",
		req.Email,
		"\n\n Message: ",
		req.Message,
	}, " ")

	mail := services.Mail{
		Name:    req.Name,
		ToEmail: req.ToEmail,
		Subject: req.Subject,
		Body:    body,
	}

	_, err = services.SendSimpleMessage(&mail)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	resSuccess := struct {
		Success int `json:"success"`
	}{
		Success: 1,
	}

	b, err := json.Marshal(resSuccess)
	if err != nil {
		fmt.Println("json err:", err)
	}

	fmt.Fprint(w, string(b))

}

type bookingReq struct {
	ToEmail        string `json:"toEmail"`
	Name           string `json:"name"`
	Email          string `json:"email"`
	Subject        string `json:"subject"`
	TattooLocation string `json:"tattooLocation"`
	Size           string `json:"size"`
	Description    string `json:"description"`
	ImageURI       string `json:"imageURI"`
}

// SendContactMail sends an email using mailgun
func SendBookingMail(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	req := new(bookingReq)

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

	body := strings.Join([]string{
		"From: \n\t Name:",
		req.Name,
		"\n\t Email: ",
		req.Email,
		"\n\n Location of Tattoo: ",
		req.TattooLocation,
		"\n\n Size of Tattoo: ",
		req.Size,
		"\n\n Description: ",
		req.Description,
	}, " ")
	html := ""

	if req.ImageURI != "" {
		html = strings.Join([]string{
			`<html>
				<body>
					<div>
						<h2>Name: `, req.Name, `</h2>
						<h2>Email: `, req.Email, `</h2>
						<h2>Location of Tattoo: </h2>
						<h3>`, req.TattooLocation, `</h3>
						<h2>Size  of Tattoo: `, req.Size, `</h2>
						<h2>Description: </h2>
						<h3>`, req.Description, `</h3>
							<img width="300" height="auto" src="`, req.ImageURI, `">
					</div>
				</body>
			</html`,
		}, "")
	}

	mail := services.Mail{
		Name:    req.Name,
		ToEmail: req.ToEmail,
		Subject: req.Subject,
		Body:    body,
		Html:    html,
	}

	_, err = services.SendSimpleMessage(&mail)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	resSuccess := struct {
		Success int `json:"success"`
	}{
		Success: 1,
	}

	b, err := json.Marshal(resSuccess)
	if err != nil {
		fmt.Println("json err:", err)
	}

	fmt.Fprint(w, string(b))

}
