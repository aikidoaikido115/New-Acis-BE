package utils

import (
	"errors"
	"mime/multipart"
	"net/http"
)

// DetectImageType - ตรวจสอบ image type จาก file content
func DetectImageType(file multipart.File) (string, error) {
	// อ่าน 512 bytes แรกเพื่อตรวจสอบ MIME type
	buffer := make([]byte, 512)
	_, err := file.Read(buffer)
	if err != nil {
		return "", err
	}

	// Reset file pointer กลับไปที่จุดเริ่มต้น
	file.Seek(0, 0)

	// ตรวจสอบ MIME type
	mimeType := http.DetectContentType(buffer)

	switch mimeType {
	case "image/jpeg":
		return ".jpg", nil
	case "image/png":
		return ".png", nil
	case "image/gif":
		return ".gif", nil
	case "image/webp":
		return ".webp", nil
	default:
		return "", errors.New("unsupported image format: " + mimeType)
	}
}
