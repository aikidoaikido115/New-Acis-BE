package ai

// CheckAllergyRequest is the payload sent to the model API.
type CheckAllergyRequest struct {
	MenuData       MenuData        `json:"menu_data"`
	AllergyDetails []AllergyDetail `json:"allergy_details"`
}

type MenuData struct {
	MenuName        string `json:"menu_name"`
	MenuDescription string `json:"menu_description"`
}

type AllergyDetail struct {
	AllergyID   string `json:"allergy_id"`
	AllergyName string `json:"allergy_name"`
	Count       int    `json:"count"`
}
