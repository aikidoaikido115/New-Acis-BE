package models

import (
	"github.com/aikidoaikido115/New-Acis-BE/pkg/ai"
)

// AllergyCheckError is a custom error that carries allergy check response data
type AllergyCheckError struct {
	Message  string
	Status   string
	Response interface{}
}

type AllergyCheckSummary struct {
	MainMenuPassed   bool                     `json:"main_menu_passed"`
	BackupMenuPassed bool                     `json:"backup_menu_passed"`
	MainMenuResult   *ai.CheckAllergyResponse `json:"main_menu_result,omitempty"`
	BackupMenuResult *ai.CheckAllergyResponse `json:"backup_menu_result,omitempty"`
}

func (e *AllergyCheckError) Error() string {
	return e.Message
}

func NewAllergyCheckError(message string, status string, response interface{}) *AllergyCheckError {
	return &AllergyCheckError{
		Message:  message,
		Status:   status,
		Response: response,
	}
}
