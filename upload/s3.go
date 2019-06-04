package upload

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

var (
	svc    *s3.S3
	folder = "main"
)

func init() {

	token := ""
	creds := credentials.NewStaticCredentials(os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"), token)

	_, err := creds.Get()
	if err != nil {
		log.Fatal(err)
	}

	cfg := aws.NewConfig().WithRegion(os.Getenv("REGION")).WithCredentials(creds)

	svc = s3.New(session.New(), cfg)
}

type UppyRequest struct {
	ContentType string `json:"contentType,omitempty"`
	FileName    string `json:"filename,omitempty"`
}

func PreSignRequest(w http.ResponseWriter, r *http.Request) {
	uppyRequest := new(UppyRequest)
	if err := json.NewDecoder(r.Body).Decode(uppyRequest); err != nil {
		jsonRespond(w, http.StatusInternalServerError, nil)
	}
	req, _ := svc.PutObjectRequest(&s3.PutObjectInput{
		Bucket:      aws.String(os.Getenv("BUCKET")),
		Key:         aws.String(fmt.Sprintf("%s/%s", folder, uppyRequest.FileName)),
		ContentType: aws.String(uppyRequest.ContentType),
		Body:        strings.NewReader(``),
	})
	urlStr, err := req.Presign(15 * time.Minute)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	data := map[string]interface{}{
		"method": req.Operation.HTTPMethod,
		"url":    urlStr,
		"fields": make(map[string]interface{}),
	}
	jsonRespond(w, http.StatusCreated, data)
}
func jsonRespond(w http.ResponseWriter, statusCode int, v interface{}) {
	var err error
	b := []byte(``)
	if v != nil {
		b, err = json.Marshal(v)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(b); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
