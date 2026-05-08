package entities

type ImageURL struct {
	URL string `json:"url"`
}

type Participation struct {
	ResidentID       string           `json:"resident_id" gorm:"primaryKey"`
	ASID             string           `json:"as_id" gorm:"primaryKey"`
	IsParticipating  bool             `json:"is_participating" gorm:"not null"`
	ImgURLs          []ImageURL       `json:"img_urls" gorm:"type:jsonb;not null;default:'[]';serializer:json"`
	Resident         Resident         `json:"-" gorm:"foreignKey:ResidentID;references:ID"`
	ActivitySchedule ActivitySchedule `json:"activity_schedule" gorm:"foreignKey:ASID;references:ID"`
}
