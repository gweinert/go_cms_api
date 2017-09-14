package services

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"strings"

	"google.golang.org/api/option"

	"cloud.google.com/go/storage"
)

func GoogleCloudUpload(file io.Reader, bucketName string, fileName string) (string, error) {

	fmt.Println("google bucketname", bucketName)
	ctx := context.Background()

	// oauthConfig := oauth2.ConfigFromJSON(os.Getenv("GOOGLE_CLOUD_KEYFILE_JSON"))
	// Creates a client.
	client, err := storage.NewClient(ctx,
		option.WithCredentialsFile(os.Getenv("GOOGLE_CREDENTIALS")),
	)
	// option.WithServiceAccountFile("/Users/Garrett/Desktop/react-cms-e5dc3890c619.json"))
	// option.WithTokenSource(oauth.tokenSource))

	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// upload object
	wc := client.Bucket(bucketName).Object(fileName).NewWriter(ctx)
	wc.ContentType = getContentType(fileName)
	if _, err = io.Copy(wc, file); err != nil {
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

	fileLink := strings.Join([]string{
		"https://storage.googleapis.com/",
		bucketName,
		"/",
		fileName,
	}, "")

	return fileLink, nil
}

func GoogleCloudDelete(bucketName string, fileURLs []string) ([]string, error) {
	ctx := context.Background()

	// Creates a client.
	client, err := storage.NewClient(ctx,
		option.WithCredentialsFile(os.Getenv("GOOGLE_CREDENTIALS")),
		// option.WithServiceAccountFile("/Users/Garrett/Desktop/react-cms-e5dc3890c619.json"))

	)
	// option.WithTokenSource(os.Getenv("REACT_CMS_GOOGLE_CREDENTIALS1")))

	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
		return nil, err
	}

	// Sets the name for the new bucket.
	// bucketName := "garrett-react-cms-test"

	// goes through array of imageURLs and deletes all of them
	for _, fileURL := range fileURLs {

		// Gets object name from end og image URL
		file, err := url.Parse(fileURL)
		if err != nil {
			return nil, err
		}
		filePath := file.Path
		filePathArr := strings.Split(filePath, "/")
		fileName := filePathArr[len(filePathArr)-1]

		// deletes image
		o := client.Bucket(bucketName).Object(fileName)
		if err := o.Delete(ctx); err != nil {
			return nil, err
		}
	}

	return fileURLs, nil
}

func getContentType(fileName string) string {
	fileType := strings.Split(fileName, ".")[1]

	switch strings.ToLower(fileType) {
	case "jpg":
		return "image/jpeg"
	case "png":
		return "image/png"
	case "json":
		return "application/json"
	default:
		return "image/jpeg"
	}
}
