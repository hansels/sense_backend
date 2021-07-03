package ml

import tf "github.com/galeone/tensorflow/tensorflow/go"

// ObjectDetectionResponse is the response the user receives after requesting an
// object detection prediction
type ObjectDetectionResponse struct {
	Detections    []detection `json:"detections"`
	NumDetections int         `json:"numDetections"`
}

type detection struct {
	Score float32 `json:"score"`
	Label string  `json:"label"`
}

const threshold = 0.50

// NewObjectDetectionResponse creates an ObjectDetectionResponse
func NewObjectDetectionResponse(output []*tf.Tensor, labels []string) *ObjectDetectionResponse {
	detectionsAboveThreshold := 0

	detections := []detection{}

	// Use type assertion to get the values of the output tensor.
	outputDetection := output[0].Value().([][]float32)

	for i, element := range outputDetection[0] {
		if element < threshold {
			continue
		}

		detectionsAboveThreshold++

		detection := detection{
			Score: element,
			Label: labels[i],
		}
		detections = append(detections, detection)
	}

	return &ObjectDetectionResponse{
		Detections:    detections,
		NumDetections: detectionsAboveThreshold,
	}
}
