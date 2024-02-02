package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	// Load environment variables from .env file
	if err := godotenv.Load(); err != nil {
		panic("Error loading .env file")
	}
}

func setupAWSS3Client() s3.S3 {
	// Set up AWS S3 client
	awsConfig := aws.NewConfig().WithRegion(os.Getenv("AWS_REGION"))
	awsSession := session.Must(session.NewSession(awsConfig))
	s3Client := s3.New(awsSession)
	return *s3Client
}

func main() {

	// Set up Gin router
	router := gin.Default()

	// Serve static files (HTML, CSS, JS)
	router.LoadHTMLGlob("static/**/*") // Because /static/ has a one more depth subdirectory
	router.Static("/static", "./static")

	// Route to upload image
	router.POST("/upload", serverImageUploadHandler)

	// Route to view(fetch) images
	router.GET("/images", serverImageReceiveHandler)

	// Route to serve the HTML file
	router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusSeeOther, "/static/index.html")
	})

	// Run the server
	router.Run(fmt.Sprintf(":%s", os.Getenv("CLIENT_WEB_ACCESS_PORT"))) // e.g., ":8080"
}

// upload an image from the webserver
func serverImageUploadHandler(context *gin.Context) {
	s3Client := setupAWSS3Client()

	file, err := context.FormFile("image")
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Image not provided"})
		return
	}

	uploader := context.Request.FormValue("uploader")

	// Open the file
	uploadedFile, err := file.Open()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Error opening file"})
		return
	}
	defer uploadedFile.Close()

	// Upload file to S3 with metadata
	_, err = s3Client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(os.Getenv("S3_BUCKET_NAME")),
		Key:    aws.String(file.Filename),
		Body:   uploadedFile,
		Metadata: map[string]*string{
			"uploader": aws.String(uploader),
		},
	})

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Error uploading file to S3"})
		return
	}

	// Redirect to the "upload_success.html" page
	context.Redirect(http.StatusSeeOther, "/static/upload_success.html")

}

// Fetch images
func serverImageReceiveHandler(context *gin.Context) {
	s3Client := setupAWSS3Client()

	// List objects in S3 bucket
	result, err := s3Client.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket: aws.String(os.Getenv("S3_BUCKET_NAME")),
	})

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Error listing objects in S3 bucket"})
		return
	}

	// Extract object keys and additional information from the result
	var images []gin.H
	for _, item := range result.Contents {
		// Get additional information (e.g., file size) for each object
		headOutput, err := s3Client.HeadObject(&s3.HeadObjectInput{
			Bucket: aws.String(os.Getenv("S3_BUCKET_NAME")),
			Key:    item.Key,
		})

		if err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving object details from S3"})
			return
		}

		// Add object details to the list
		images = append(images, gin.H{
			"key":          *item.Key,
			"uploader":     headOutput.Metadata["Uploader"], // Use capitalized character
			"size":         *headOutput.ContentLength,
			"lastModified": headOutput.LastModified,
			"objectAccessURL": fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s",
				os.Getenv("S3_BUCKET_NAME"),
				os.Getenv("AWS_REGION"),
				url.QueryEscape(*item.Key)),
		})
	}

	// Return the list of images with additional information as JSON
	context.JSON(http.StatusOK, gin.H{"images": images})
}
