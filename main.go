package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

var templates = template.Must(template.ParseFiles("templates/upload.html"))

const (
    S3_REGION = "us-west-2" // 変更が必要
    S3_BUCKET = "your-bucket-name" // 変更が必要
)

func main() {
    http.HandleFunc("/", uploadFormHandler)
    http.HandleFunc("/upload", uploadHandler)
    log.Println("Server started on :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}

func uploadFormHandler(w http.ResponseWriter, r *http.Request) {
    if err := templates.ExecuteTemplate(w, "upload.html", nil); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
        return
    }

    file, header, err := r.FormFile("image")
    if err != nil {
        http.Error(w, "Failed to read file", http.StatusBadRequest)
        return
    }
    defer file.Close()

    sess, err := session.NewSession(&aws.Config{
        Region: aws.String(S3_REGION)},
    )
    if err != nil {
        http.Error(w, "Failed to create AWS session", http.StatusInternalServerError)
        return
    }

    uploader := s3.New(sess)

    _, err = uploader.PutObject(&s3.PutObjectInput{
        Bucket: aws.String(S3_BUCKET),
        Key:    aws.String(header.Filename),
        Body:   file,
        ACL:    aws.String("public-read"),
    })
    if err != nil {
        http.Error(w, "Failed to upload file to S3", http.StatusInternalServerError)
        return
    }

    url := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", S3_BUCKET, S3_REGION, header.Filename)
    fmt.Fprintf(w, "Successfully uploaded file to S3: %s\n", url)
}