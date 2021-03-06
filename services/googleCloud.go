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

	"golang.org/x/oauth2/google"

	"cloud.google.com/go/storage"
)

func GoogleCloudUpload(file io.Reader, bucketName string, fileName string) (string, error) {

	fmt.Println("google bucketname", bucketName)

	ctx := context.Background()
	json := os.Getenv("REACT_CMS_GOOGLE_CREDENTIALS") // `{"type": "service_account", "project_id": "my-project", ...}`

	jwtConfig, err := google.JWTConfigFromJSON([]byte(json), storage.ScopeFullControl)
	if err != nil {
		log.Fatal(err)
	}
	ts := jwtConfig.TokenSource(ctx)

	client, err := storage.NewClient(ctx, option.WithTokenSource(ts)) //option.WithCredentialsFile("https://00e9e64bac710723f321221e0d775a3da461c2308a36a2f0e3-apidata.googleusercontent.com/download/storage/v1/b/react-cms-utilities/o/react-cms-e5dc3890c619.json?qk=AD5uMEuweJUR19VCe2NURhMvYBT1G_4pWiJXZQcZUbXdblu9T4hJp3x9xY_QhrHA_ckNIVug0ZaQ-fNapoGifba3CI05ilDg17nSkHnpIuZdpBZppB8SdwGoIn3_vwLqcxKNjl1Nl8ClMhsbzPWxlsDDtO9xOaz7CQofvw5FPpPOWEWRmZ9E1qOAPITtphiifmBR85W4y8jNmQLoYEQhyvimcVcQ1zJvKob47xyS2Kf5y1rp27cJzOlnIaB-5XmB-R3UteMgn4eYelS1PjP-cgxbjfNxaMfyfFn3mfPf9djlsJ1ayUY_MdwSsy34u_KPbaLgScGGYcB1uWb0lrHcQVr0Bq1GVmLQ-PK5Bd93jmlGFKmIDS-ZNHFAJahICBXUlzi2z786WXI0rOJJ-0HVFV3_oc8eiq1EEROfHffdKrYGrxdFpuPIL0dusZvhUIGTpaET_ebSmcmWLuRnLLzxX66lFTClArqQ_1Ejk8aP3isSKhimolOasGvhuMOEY3I2N6l0n4OpJ3n1c9RAFzdhTfn6lBnlu6dWWCIt7HW7h-MWdnT8Rw5TW8siIrxOaxA47Q646DrCSbB9udLNrtqoLtMd-Z5vOYfGQ8BPU4YDX4HI1jbWwfyeJiwvVIUSLYNWjijqyzBEMOtHio7e5oYSQa4NJzQZlo5V79pqMgfgNGcGmQlDDruoDkJdgao9KcmwkQKBuTecOKNfxwXlX8dHcvwejOYFx3yBski23U7rYHr2GgU8Oz_vMBEGLZGMyu3soFhSXkhEa88D"),

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

	fmt.Println("google bucketname", bucketName)

	json := os.Getenv("REACT_CMS_GOOGLE_CREDENTIALS") // `{"type": "service_account", "project_id": "my-project", ...}`

	jwtConfig, err := google.JWTConfigFromJSON([]byte(json), storage.ScopeFullControl)
	if err != nil {
		log.Fatal(err)
	}
	ts := jwtConfig.TokenSource(ctx)

	client, err := storage.NewClient(ctx, option.WithTokenSource(ts)) //option.WithCredentialsFile("https://00e9e64bac710723f321221e0d775a3da461c2308a36a2f0e3-apidata.googleusercontent.com/download/storage/v1/b/react-cms-utilities/o/react-cms-e5dc3890c619.json?qk=AD5uMEuweJUR19VCe2NURhMvYBT1G_4pWiJXZQcZUbXdblu9T4hJp3x9xY_QhrHA_ckNIVug0ZaQ-fNapoGifba3CI05ilDg17nSkHnpIuZdpBZppB8SdwGoIn3_vwLqcxKNjl1Nl8ClMhsbzPWxlsDDtO9xOaz7CQofvw5FPpPOWEWRmZ9E1qOAPITtphiifmBR85W4y8jNmQLoYEQhyvimcVcQ1zJvKob47xyS2Kf5y1rp27cJzOlnIaB-5XmB-R3UteMgn4eYelS1PjP-cgxbjfNxaMfyfFn3mfPf9djlsJ1ayUY_MdwSsy34u_KPbaLgScGGYcB1uWb0lrHcQVr0Bq1GVmLQ-PK5Bd93jmlGFKmIDS-ZNHFAJahICBXUlzi2z786WXI0rOJJ-0HVFV3_oc8eiq1EEROfHffdKrYGrxdFpuPIL0dusZvhUIGTpaET_ebSmcmWLuRnLLzxX66lFTClArqQ_1Ejk8aP3isSKhimolOasGvhuMOEY3I2N6l0n4OpJ3n1c9RAFzdhTfn6lBnlu6dWWCIt7HW7h-MWdnT8Rw5TW8siIrxOaxA47Q646DrCSbB9udLNrtqoLtMd-Z5vOYfGQ8BPU4YDX4HI1jbWwfyeJiwvVIUSLYNWjijqyzBEMOtHio7e5oYSQa4NJzQZlo5V79pqMgfgNGcGmQlDDruoDkJdgao9KcmwkQKBuTecOKNfxwXlX8dHcvwejOYFx3yBski23U7rYHr2GgU8Oz_vMBEGLZGMyu3soFhSXkhEa88D"),

	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

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
