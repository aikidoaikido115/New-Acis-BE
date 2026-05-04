package utils

import (
	"bytes"
	"fmt"
	"log"
	"mime/multipart"

	"github.com/aikidoaikido115/New-Acis-BE/configs"

	storage_go "github.com/supabase-community/storage-go"
)

// UploadFile - รับ multipart.File โดยตรง
func UploadFile2Supa(file multipart.File, fileName string, dir string, config configs.Supabase) (string, error) {
	log.Printf("=== SUPABASE CONFIG CHECK ===")
	log.Printf("URL: '%s' (len=%d)", config.URL, len(config.URL))
	log.Printf("Bucket: '%s' (len=%d)", config.Bucket, len(config.Bucket))
	log.Printf("ServiceKey: '%s...' (len=%d)", config.ServiceKey[:20], len(config.ServiceKey))
	log.Printf("=== END CONFIG CHECK ===")

	if config.URL == "" || config.ServiceKey == "" || config.Bucket == "" {
		log.Printf("Config validation failed - missing required fields")
		return "", fmt.Errorf("invalid Supabase config")
	}

	// Supabase Storage URL should be: https://PROJECT_REF.supabase.co/storage/v1
	storageURL := config.URL + "/storage/v1"
	log.Printf("Creating storage client with URL: %s", storageURL)

	storageClient := storage_go.NewClient(storageURL, config.ServiceKey, nil)
	if storageClient == nil {
		return "", fmt.Errorf("failed to create storage client")
	}

	contentType := getContentType(fileName)
	options := storage_go.FileOptions{
		ContentType: &contentType,
	}

	fullFileName := dir + fileName

	log.Printf("About to upload - Bucket: %s, FileName: %s, ContentType: %s", config.Bucket, fullFileName, contentType)

	resp, err := storageClient.UploadFile(config.Bucket, fullFileName, file, options)
	if err != nil {
		log.Printf("Upload error details: %v", err)
		log.Printf("Error type: %T", err)
		if storageErr, ok := err.(*storage_go.StorageError); ok {
			log.Printf("StorageError - Status: %d, Message: %s", storageErr.Status, storageErr.Message)
		}
		return "", fmt.Errorf("failed to upload file: %v", err)
	}

	log.Printf("Upload successful! Response: %+v", resp)

	url := fmt.Sprintf("%s/storage/v1/object/public/%s/%s", config.URL, config.Bucket, fullFileName)
	return url, nil
}

// UploadFileFromBytes - รับ byte array โดยตรง
func UploadFileFromBytes(data []byte, fileName string, dir string, config configs.Supabase) (string, error) {
	log.Printf("Supabase Config - URL: %s, Bucket: %s", config.URL, config.Bucket)
	log.Printf("Supabase Config - ServiceKey: %s", config.ServiceKey)
	if config.URL == "" || config.ServiceKey == "" || config.Bucket == "" {
		return "", fmt.Errorf("invalid Supabase config")
	}

	storageClient := storage_go.NewClient(config.URL, config.ServiceKey, nil)
	if storageClient == nil {
		return "", fmt.Errorf("failed to create storage client")
	}

	contentType := getContentType(fileName)
	options := storage_go.FileOptions{
		ContentType: &contentType,
	}

	fullFileName := dir + fileName
	reader := bytes.NewReader(data)

	log.Printf("About to upload bytes - Bucket: %s, FileName: %s, ContentType: %s, Size: %d", config.Bucket, fullFileName, contentType, len(data))

	_, err := storageClient.UploadFile(config.Bucket, fullFileName, reader, options)
	if err != nil {
		log.Printf("Upload error details: %v", err)
		log.Printf("Error type: %T", err)
		if storageErr, ok := err.(*storage_go.StorageError); ok {
			log.Printf("StorageError - Status: %d, Message: %s", storageErr.Status, storageErr.Message)
		}
		return "", fmt.Errorf("failed to upload file: %v", err)
	}

	url := fmt.Sprintf("%s/storage/v1/object/public/%s/%s", config.URL, config.Bucket, fullFileName)
	return url, nil
}

// Helper function สำหรับกำหนด content type
func getContentType(fileName string) string {
	switch {
	case contains(fileName, ".jpg") || contains(fileName, ".jpeg"):
		return "image/jpeg"
	case contains(fileName, ".png"):
		return "image/png"
	case contains(fileName, ".gif"):
		return "image/gif"
	case contains(fileName, ".webp"):
		return "image/webp"
	case contains(fileName, ".pdf"):
		return "application/pdf"
	default:
		return "image/jpeg"
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && s[len(s)-len(substr):] == substr
}
