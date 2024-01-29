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

func main() {
	// Set up AWS S3 client
	awsConfig := aws.NewConfig().WithRegion(os.Getenv("AWS_REGION"))
	awsSession := session.Must(session.NewSession(awsConfig))
	s3Client := s3.New(awsSession)

	// Set up Gin router
	router := gin.Default()

	// Serve static files (HTML, CSS, JS)
	router.LoadHTMLGlob("static/*")

	// Route to upload image
	router.POST("/upload", func(c *gin.Context) {
		file, err := c.FormFile("image")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Image not provided"})
			return
		}

		// Open the file
		uploadedFile, err := file.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error opening file"})
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
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error uploading file to S3"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Image uploaded successfully"})
	})

	// Route to serve the HTML file
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{})
	})

	// Run the server
	router.Run(":8080")
}
