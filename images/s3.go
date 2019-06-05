package images

import (
	"dropzone-s3/respond"
	"dropzone-s3/views"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

var (
	svc    *s3.S3
	folder = "blah"
	bucket = os.Getenv("BUCKET")
)

func init() {

	creds := credentials.NewStaticCredentials(os.Getenv("AWS_ACCESS_KEY"), os.Getenv("AWS_SECRET_KEY"), "")

	_, err := creds.Get()
	if err != nil {
		log.Fatal(err)
	}

	cfg := aws.NewConfig().WithRegion(os.Getenv("REGION")).WithCredentials(creds)
	sess, err := session.NewSession(cfg)
	if err != nil {
		log.Fatal(err)
	}
	svc = s3.New(sess)
}

type PutRequest struct {
	FileName    string `json:"filename,omitempty"`
	ContentType string `json:"contentType,omitempty"`
	Extension   string `json:"extension,omitempty"`
	SignedURL   string `json:"signed_url,omitempty"` // Pre Signed URL
	Path        string `json:"path,omitempty"`
	Size        int    `json:"size,omitempty"`
}

func PreSignRequest(w http.ResponseWriter, r *http.Request) {
	putRequest := new(PutRequest)
	if err := json.NewDecoder(r.Body).Decode(putRequest); err != nil {
		respond.Json(w, http.StatusInternalServerError, nil)
	}
	fmt.Println(putRequest)
	key := fmt.Sprintf("%s/%s", folder, putRequest.FileName)
	req, _ := svc.PutObjectRequest(&s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(key),
		ContentType: aws.String(putRequest.ContentType),
		Body:        strings.NewReader(``),
	})
	urlStr, err := req.Presign(15 * time.Minute)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	data := map[string]interface{}{
		"method": req.Operation.HTTPMethod,
		"url":    urlStr,
		"fields": make(map[string]interface{}), // has to stay here
		"bucket": bucket,
		"key":    key,
	}
	respond.Json(w, http.StatusCreated, data)
}

type GetRequest struct {
	FileName string `json:"filename,omitempty"`
}

func GetSignRequest(w http.ResponseWriter, r *http.Request) {
	putRequest := new(PutRequest)
	if err := json.NewDecoder(r.Body).Decode(putRequest); err != nil {
		respond.Json(w, http.StatusInternalServerError, nil)
	}
	urlStr, _ := createPreSignURL(folder + "wolf.jpg")
	data := map[string]interface{}{
		"url": urlStr,
	}
	respond.Json(w, http.StatusCreated, data)
}
func createPreSignURL(key string) (string, error) {
	req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	return req.Presign(60 * time.Minute)
}

type Image struct {
	URL  string // Pre Signed URL
	Path string
	Size int64
}

func ListBuckets(w http.ResponseWriter, r *http.Request) {
	buckets, err := svc.ListBuckets(&s3.ListBucketsInput{})
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(*buckets.Owner)
	for _, b := range buckets.Buckets {
		fmt.Println(*b.Name)
	}
}
func GetImages(w http.ResponseWriter, r *http.Request) {
	params := &s3.ListObjectsInput{
		Bucket: aws.String(bucket),
	}
	resp, err := svc.ListObjects(params)
	if err != nil {
		fmt.Println(err)
	}

	var images []Image

	for _, obj := range resp.Contents {
		k := *obj.Key
		signed, err := createPreSignURL(k)
		if err != nil {
			fmt.Println(err)
		}

		images = append(images, Image{URL: signed, Path: url.QueryEscape(k), Size: *obj.Size})
	}

	views.Render(w, "templates/images.html", images)
}

func DeleteImage(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	k := q.Get("key")
	params := &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(k),
	}
	_, err := svc.DeleteObject(params)
	if err != nil {
		fmt.Println(err)
	}
	respond.Json(w, 200, nil)
}
