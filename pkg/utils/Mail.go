package utils

import (
	"bytes"
	"errors"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"

	"github.com/aikidoaikido115/New-Acis-BE/configs"
	"github.com/aikidoaikido115/New-Acis-BE/modules/entities"

	"gopkg.in/gomail.v2"
)

func SendMail(templatePath string, user *entities.User, otp string, config configs.Mail) error {
	var body bytes.Buffer

	// Check if file exists
	_, err := os.Stat(templatePath)
	if os.IsNotExist(err) {
		// Try with absolute path based on working directory
		workingDir, _ := os.Getwd()
		templatePath = filepath.Join(workingDir, templatePath)
	}

	t, err := template.ParseFiles(templatePath)
	if err != nil {
		log.Printf("Template parsing error: %v", err)
		return err
	}

	if templatePath == "" {
		return errors.New("template is empty")
	}

	err = t.Execute(&body, struct {
		Username string
		OTP      string
	}{
		Username: user.Username,
		OTP:      otp,
	})

	if err != nil {
		return err
	}

	m := gomail.NewMessage()
	m.SetHeader("From", config.Sender)
	m.SetHeader("To", user.Email)
	m.SetHeader("Subject", "Recovery Your Password")
	m.SetBody("text/html", body.String())
	port, err := strconv.Atoi(config.Port)
	if err != nil {
		return err
	}

	d := gomail.NewDialer(config.Host, port, config.Sender, config.Key)
	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}

func NormalizeEmail(email string) (string, error) {
	email = strings.ToLower(email)
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return "", errors.New("invalid Email")
	}

	localPart, domain := parts[0], parts[1]
	localPart = strings.ReplaceAll(localPart, ".", "")
	email = localPart + "@" + domain
	return email, nil
}
