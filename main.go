package main

import (
	"net/http"
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
	router.LoadHTMLGlob("static/*")
	router.Static("/static", "./static")

	// Route to upload image
	router.POST("/upload", serverImageUploadHandler)

	// Route to view(fetch) images
	router.GET("/images", serverImageReceiveHandler)

	// Route to serve the HTML file
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})

	// Run the server
	router.Run(":8080")
}

// upload an image from the webserver
func serverImageUploadHandler(context *gin.Context) {

	s3Client := setupAWSS3Client()

	file, err := context.FormFile("image")
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Image not provided"})
		return
	}

	// Open the file
	uploadedFile, err := file.Open()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Error opening file"})
		return
	}
	defer uploadedFile.Close()

	// Upload file to S3
	_, err = s3Client.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(os.Getenv("S3_BUCKET_NAME")),
		Key:    aws.String(file.Filename),
		Body:   uploadedFile,
	})

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Error uploading file to S3"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Image uploaded successfully"})
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

	// Extract object keys from the result
	var objectKeys []string
	for _, item := range result.Contents {
		objectKeys = append(objectKeys, *item.Key)
	}

	// Return the list of images as JSON
	context.JSON(http.StatusOK, gin.H{"images": objectKeys})
}
