package constants

const (
	// Temperature limits (Celsius)
	MinTemperature = 30.0
	MaxTemperature = 45.0

	// Heart rate limits (bpm)
	MinHeartRate = 20
	MaxHeartRate = 250

	// Breathing rate limits (breaths per minute)
	MinBreathingRate = 5
	MaxBreathingRate = 60

	// Blood pressure systolic limits (mmHg)
	MinBloodPressureSystolic = 50
	MaxBloodPressureSystolic = 300

	// Blood pressure diastolic limits (mmHg)
	MinBloodPressureDiastolic = 30
	MaxBloodPressureDiastolic = 200

	// Oxygen saturation limits (%)
	MinOxygenSaturation = 50
	MaxOxygenSaturation = 100

	// Temperature (°C) - ปกติ 36.1-37.2
	// < 35.0 = Hypothermia (ตัวเย็นผิดปกติ)
	// > 37.5 = Fever (มีไข้)
	// > 38.0 = High Fever (ไข้สูง)
	NormalTempLow  = 35.0
	NormalTempHigh = 37.5

	// Heart Rate (bpm) - ผู้ใหญ่ปกติ 60-100
	// < 60 = Bradycardia (หัวใจเต้นช้า)
	// > 100 = Tachycardia (หัวใจเต้นเร็ว)
	NormalHeartRateLow  = 60
	NormalHeartRateHigh = 100

	// Breathing Rate (breaths/min) - ผู้ใหญ่ปกติ 12-20
	// < 12 = Bradypnea (หายใจช้า)
	// > 20 = Tachypnea (หายใจเร็ว)
	NormalBreathingRateLow  = 12
	NormalBreathingRateHigh = 20

	// Blood Pressure Systolic (mmHg)
	// < 90 = Hypotension (ความดันต่ำ)
	// > 140 = Hypertension Stage 2 (ความดันสูง)
	NormalSystolicLow  = 90
	NormalSystolicHigh = 140

	// Blood Pressure Diastolic (mmHg)
	// < 60 = Hypotension
	// > 90 = Hypertension Stage 2
	NormalDiastolicLow  = 60
	NormalDiastolicHigh = 90

	// Oxygen Saturation (%)
	// < 95 = Low (ออกซิเจนต่ำ - ควรเฝ้าระวัง)
	// < 90 = Critical (วิกฤต - ต้องการออกซิเจน)
	NormalOxygenSaturationLow = 95
)
