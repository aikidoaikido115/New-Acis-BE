package constants

const (
	// UrineType valid values
	// UrineOutput วัดเป็น มิลลิลิตร (ml) — ใช้กับผู้ป่วยที่ต้องติดตามปริมาณปัสสาวะอย่างละเอียด เช่น โรคไต
	UrineTypeML = "ml"
	// UrineOutput วัดเป็น ครั้ง (times) — ใช้กับผู้ป่วยทั่วไปที่วัดแค่จำนวนครั้ง
	UrineTypeTimes = "times"

	// BloodGlucose limits (mg/dL)
	// < 70 = Hypoglycemia, > 180 = Hyperglycemia (postprandial threshold)
	MinBloodGlucose = 1.0
	MaxBloodGlucose = 1000.0

	NormalBloodGlucoseLow  = 70.0
	NormalBloodGlucoseHigh = 180.0

	// FluidIn limits (mL) — ปริมาณน้ำที่รับเข้าทั้งหมดต่อบันทึก/กะ
	MinFluidIn = 0.0
	MaxFluidIn = 10000.0

	// FluidOut limits (mL) — ปริมาณน้ำที่ออกจากร่างกายทั้งหมดต่อบันทึก/กะ
	MinFluidOut = 0.0
	MaxFluidOut = 10000.0

	// UrineOutput when UrineType = "ml" (mL)
	MinUrineOutputML = 0.0
	MaxUrineOutputML = 5000.0

	// UrineOutput when UrineType = "times" (ครั้ง)
	MinUrineOutputTimes = 0.0
	MaxUrineOutputTimes = 50.0

	// Stool limits (times/day or per shift)
	MinStool = 0
	MaxStool = 30

	// DiaperChange limits (times/day or per shift)
	MinDiaperChange = 0
	MaxDiaperChange = 30

	// Normal ranges (per shift/record) for abnormality detection

	// FluidIn (mL/shift) — < 300 = dehydration risk, > 2500 = overhydration risk
	NormalFluidInLow  = 300.0
	NormalFluidInHigh = 2500.0

	// FluidOut (mL/shift) — < 200 = fluid retention risk, > 2000 = excessive loss
	NormalFluidOutLow  = 200.0
	NormalFluidOutHigh = 2000.0

	// UrineOutput when UrineType = "ml" (mL/shift) — < 150 = oliguria risk
	NormalUrineOutputMLLow  = 150.0
	NormalUrineOutputMLHigh = 2000.0

	// UrineOutput when UrineType = "times" (ครั้ง/shift) — < 2 = oliguria risk, > 8 = polyuria risk
	NormalUrineOutputTimesLow  = 2.0
	NormalUrineOutputTimesHigh = 8.0

	// Stool (times/shift) — > 3 = possible diarrhea
	NormalStoolHigh = 3

	// DiaperChange (times/shift) — > 6 = excessive
	NormalDiaperChangeHigh = 6
)
