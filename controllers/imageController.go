package controllers

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"regexp"
	// "strings"

	model "github.com/gweinert/cms_scratch/models"
	"github.com/gweinert/cms_scratch/services"
	"github.com/julienschmidt/httprouter"
)

type uploadReq struct {
	ID             int    `json:"elementId"`
	TempID         string `json:"tempId"`
	DataURI        string `json:"dataUri"`
	FileName       string `json:"fileName"`
	FileType       string `json:"fileType"`
	PageID         int    `json:"pageId"`
	SortOrder      int    `json:"sortOrder"`
	GroupID        int    `json:"groupId"`
	GroupSortOrder int    `json:"groupSortOrder"`
	Name           string `json:"name"`
}

func UploadImage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	req := new(uploadReq)

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

	imageURL, err := googleCloudUpload(req)
	if err != nil {
		fmt.Println("error uploading")
		http.Error(w, err.Error(), 500)
		return
	}

	id, err := model.SaveImageURL(req.ID, imageURL, req.PageID, req.SortOrder, req.GroupID, req.GroupSortOrder, req.Name)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	resData := struct {
		Success        int    `json:"success"`
		ImageURL       string `json:"imageURL"`
		ID             int    `json:"elementId"`
		TempID         string `json:"tempId"`
		PageID         int    `json:"pageId"`
		SortOrder      int    `json:"sortOrder"`
		GroupID        int    `json:"groupId"`
		GroupSortOrder int    `json:"groupSortOrder"`
		Name           string `json:"name"`
	}{
		Success:        1,
		ImageURL:       imageURL,
		ID:             id,
		TempID:         req.TempID,
		PageID:         req.PageID,
		SortOrder:      req.SortOrder,
		GroupID:        req.GroupID,
		GroupSortOrder: req.GroupSortOrder,
		Name:           req.Name,
	}
	b, err := json.Marshal(resData)
	if err != nil {
		http.Error(w, "server Broke", 500)
		return
	}

	fmt.Fprint(w, string(b))
}

func googleCloudUpload(r *uploadReq) (string, error) {

	// Sets the name for the new bucket.
	bucketName := "garrett-react-cms-test"

	reg, _ := regexp.Compile("^data:image/(png|jpg|jpeg);base64,")
	base64DataURI := reg.ReplaceAllString(r.DataURI, "")
	buf := bytes.NewBufferString(base64DataURI)
	dec := base64.NewDecoder(base64.StdEncoding, buf)
	fileName := r.FileName

	imageURL, err := services.GoogleCloudUpload(dec, bucketName, fileName)
	if err != nil {
		log.Fatal(err)
		return "", err
	}

	return imageURL, nil
}

type deleteImageReq struct {
	ImageURLs []string `json:"imageURLs"`
	IDs       []int    `json:"ids"`
}

func DeleteImage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	req := new(deleteImageReq)

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

	_, err = googleCloudDelete(req.ImageURLs)
	if err != nil {
		fmt.Println("error google cloud delete")
		http.Error(w, err.Error(), 500)
	}

	// ids := []int{req.ID}
	// _, err = model.DeleteElements(ids)
	// if err != nil {
	// 	http.Error(w, err.Error(), 500)
	// 	return
	// }

	res := struct {
		Success int   `json:"success"`
		ID      []int `json:"id"`
	}{
		Success: 1,
		ID:      req.IDs,
	}

	b, err := json.Marshal(res)
	if err != nil {
		http.Error(w, "server Broke", 500)
		return
	}

	fmt.Fprint(w, string(b))
}

func googleCloudDelete(imageURLs []string) ([]string, error) {
	// Sets the name for the new bucket.
	bucketName := "garrett-react-cms-test"

	imageURLs, err := services.GoogleCloudDelete(bucketName, imageURLs)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}

	return imageURLs, nil
}
