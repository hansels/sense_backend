package model

type PredictionResult struct {
	Verdict    string `json:"verdict" structs:"verdict"`
	IsDetected bool   `json:"is_detected" structs:"is_detected"`
	Image      string `json:"image" structs:"image"`
}
