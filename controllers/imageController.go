package controllers

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	// "strings"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"

	model "github.com/gweinert/cms_scratch/models"
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

	ctx := context.Background()

	// Creates a client.
	client, err := storage.NewClient(ctx,
		option.WithServiceAccountFile("/Users/Garrett/Desktop/react-cms-e5dc3890c619.json"))
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// Sets the name for the new bucket.
	bucketName := "garrett-react-cms-test"

	reg, _ := regexp.Compile("^data:image/(png|jpg|jpeg);base64,")
	base64DataURI := reg.ReplaceAllString(r.DataURI, "")
	buf := bytes.NewBufferString(base64DataURI)
	dec := base64.NewDecoder(base64.StdEncoding, buf)
	fileName := r.FileName

	// upload object
	wc := client.Bucket(bucketName).Object(fileName).NewWriter(ctx)
	if _, err = io.Copy(wc, dec); err != nil {
		return "error copying", err
	}
	if err := wc.Close(); err != nil {
		return "error closing", err
	}

	//make url public
	acl := client.Bucket(bucketName).Object(fileName).ACL()
	if err := acl.Set(ctx, storage.AllUsers, storage.RoleReader); err != nil {
		return "", err
	}

	imageURL := strings.Join([]string{
		"https://storage.googleapis.com/",
		bucketName,
		"/",
		fileName,
	}, "")

	return imageURL, nil
}

type deleteImageReq struct {
	ImageURL string `json:"imageURL"`
	ID       int    `json:"id"`
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

	_, err = GoogleCloudDelete(req.ImageURL)
	if err != nil {
		fmt.Println("error google cloud delete")
		http.Error(w, err.Error(), 500)
	}

	ids := []int{req.ID}
	_, err = model.DeleteElements(ids)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	res := struct {
		Success int `json:"success"`
		ID      int `json:"id"`
	}{
		Success: 1,
		ID:      req.ID,
	}

	b, err := json.Marshal(res)
	if err != nil {
		http.Error(w, "server Broke", 500)
		return
	}

	fmt.Fprint(w, string(b))
}

func GoogleCloudDelete(imageURL string) (string, error) {
	ctx := context.Background()

	// Creates a client.
	client, err := storage.NewClient(ctx,
		option.WithServiceAccountFile("/Users/Garrett/Desktop/react-cms-e5dc3890c619.json"))
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
		return "", err
	}

	// Sets the name for the new bucket.
	bucketName := "garrett-react-cms-test"

	// Gets object name from end og image URL
	fileURL, err := url.Parse(imageURL)
	if err != nil {
		return "", err
	}
	filePath := fileURL.Path
	filePathArr := strings.Split(filePath, "/")
	fileName := filePathArr[len(filePathArr)-1]

	o := client.Bucket(bucketName).Object(fileName)
	if err := o.Delete(ctx); err != nil {
		return "", err
	}

	return imageURL, nil
}
