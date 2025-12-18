package utils

import (
	"errors"
	"mime/multipart"
	"net/http"
)

// DetectFileType - ตรวจสอบ file type จาก file content (รองรับ images และ PDF)
func DetectFileType(file multipart.File) (string, error) {
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
	case "application/pdf":
		return ".pdf", nil
	default:
		return "", errors.New("unsupported file format: " + mimeType)
	}
}
