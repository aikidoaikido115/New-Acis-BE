package usecases

type IntakeLabelInput struct {
    LabelName string `json:"label_name" binding:"required"`
    NoteText  *string `json:"note_text"`
}
