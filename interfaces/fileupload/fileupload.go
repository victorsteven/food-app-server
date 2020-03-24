package fileupload

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/minio/minio-go/v6"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
)

func NewFileUpload() *fileUpload {
	return &fileUpload{}
}

type fileUpload struct{}

type UploadFileInterface interface {
	UploadFile(file *multipart.FileHeader) (string, error)
}

//So what is exposed is Uploader
var _ UploadFileInterface = &fileUpload{}

//func (fu *fileUpload) UploadFilex(file *multipart.FileHeader) (string, error) {
//	f, err := file.Open()
//	if err != nil {
//		return "", errors.New("cannot open file")
//	}
//	defer f.Close()
//
//	size := file.Size
//	//The image should not be more than 500KB
//	fmt.Println("the size: ", size)
//	if size > int64(512000) {
//		return "", errors.New("sorry, please upload an Image of 500KB or less")
//	}
//	//only the first 512 bytes are used to sniff the content type of a file,
//	//so, so no need to read the entire bytes of a file.
//	buffer := make([]byte, size)
//	f.Read(buffer)
//	fileType := http.DetectContentType(buffer)
//	//if the image is valid
//	if !strings.HasPrefix(fileType, "image") {
//		return "", errors.New("please upload a valid image")
//	}
//	fileBytes := bytes.NewReader(buffer)
//	filePath := FormatFile(file.Filename)
//	path := "/profile-photos/" + filePath
//
//	params := &s3.PutObjectInput{
//		Bucket:        aws.String("chodapi"), //this is the name i saved the bucket that contains the image
//		Key:           aws.String(path),
//		Body:          fileBytes,
//		ContentLength: aws.Int64(size),
//		ContentType:   aws.String(fileType),
//		ACL:           aws.String("public-read"),
//	}
//	s3Config := &aws.Config{
//		Credentials: credentials.NewStaticCredentials(
//			os.Getenv("DO_SPACES_KEY"), os.Getenv("DO_SPACES_SECRET"), os.Getenv("DO_SPACES_TOKEN")),
//		Endpoint: aws.String(os.Getenv("DO_SPACES_ENDPOINT")),
//		Region:   aws.String(os.Getenv("DO_SPACES_REGION")),
//	}
//	newSession := session.New(s3Config)
//	s3Client := s3.New(newSession)
//
//	_, err = s3Client.PutObject(params)
//	if err != nil {
//		fmt.Println("actual error: ", err)
//		return "", errors.New("something went wrong uploading image")
//	}
//	return filePath, nil
//}

func (fu *fileUpload) UploadFile(file *multipart.FileHeader) (string, error) {
	f, err := file.Open()
	if err != nil {
		return "", errors.New("cannot open file")
	}
	defer f.Close()

	size := file.Size
	//The image should not be more than 500KB
	fmt.Println("the size: ", size)
	if size > int64(512000) {
		return "", errors.New("sorry, please upload an Image of 500KB or less")
	}
	//only the first 512 bytes are used to sniff the content type of a file,
	//so, so no need to read the entire bytes of a file.
	buffer := make([]byte, size)
	f.Read(buffer)
	fileType := http.DetectContentType(buffer)
	//if the image is valid
	if !strings.HasPrefix(fileType, "image") {
		return "", errors.New("please upload a valid image")
	}
	filePath := FormatFile(file.Filename)

	accessKey := os.Getenv("DO_SPACES_KEY")
	secKey := os.Getenv("DO_SPACES_SECRET")
	endpoint := os.Getenv("DO_SPACES_ENDPOINT")
	ssl := true

	// Initiate a client using DigitalOcean Spaces.
	client, err := minio.New(endpoint, accessKey, secKey, ssl)
	if err != nil {
		log.Fatal(err)
	}
	fileBytes := bytes.NewReader(buffer)
	cacheControl := "max-age=31536000"
	// make it public
	userMetaData := map[string]string{"x-amz-acl": "public-read"}
	n, err := client.PutObject("chodapi", filePath, fileBytes, size, minio.PutObjectOptions{ContentType: fileType, CacheControl: cacheControl, UserMetadata: userMetaData})
	if err != nil {
		fmt.Println("the error", err)
		return "", errors.New("something went wrong")
	}
	fmt.Println("Successfully uploaded bytes: ", n)
	return filePath, nil
}
