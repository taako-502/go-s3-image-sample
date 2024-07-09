package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

var templates = template.Must(template.ParseFiles("templates/upload.html"))

func main() {
	http.HandleFunc("/", uploadFormHandler)
	http.HandleFunc("/upload", uploadHandler)
	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func uploadFormHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "upload.html", nil)
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		log.Println("Invalid request method")
		renderTemplate(w, "upload.html", "Invalid request method")
		return
	}

	file, header, err := r.FormFile("image")
	if err != nil {
		log.Printf("Failed to read file: %v", err)
		renderTemplate(w, "upload.html", fmt.Sprintf("Failed to read file: %v", err))
		return
	}
	defer file.Close()

	region := os.Getenv("AWS_REGION")
	bucket := os.Getenv("S3_BUCKET")

	log.Printf("AWS_REGION: %s, S3_BUCKET: %s", region, bucket)

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(region)},
	)
	if err != nil {
		log.Printf("Failed to create AWS session: %v", err)
		renderTemplate(w, "upload.html", fmt.Sprintf("Failed to create AWS session: %v", err))
		return
	}

	uploader := s3.New(sess)

	_, err = uploader.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(header.Filename),
		Body:   file,
	})
	if err != nil {
		log.Printf("Failed to upload file to S3: %v", err)
		renderTemplate(w, "upload.html", fmt.Sprintf("Failed to upload file to S3: %v", err))
		return
	}

	url := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", bucket, region, header.Filename)
	log.Printf("Successfully uploaded file to S3: %s", url)
	renderTemplate(w, "upload.html", fmt.Sprintf("Successfully uploaded file to S3: %s", url))
}

func renderTemplate(w http.ResponseWriter, tmpl string, message interface{}) {
	data := struct {
		Message interface{}
	}{
		Message: message,
	}

	if err := templates.ExecuteTemplate(w, tmpl, data); err != nil {
		log.Printf("Failed to render template: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}