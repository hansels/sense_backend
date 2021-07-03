package ml

import (
	"bytes"
	"fmt"
	tf "github.com/galeone/tensorflow/tensorflow/go"
	"github.com/galeone/tensorflow/tensorflow/go/op"
	tg "github.com/galeone/tfgo"
	"image"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

type Coco struct {
	model  *tg.Model
	labels []string
}

const path = "files/models/my_model/"

// NewCoco returns a Coco object
func NewCoco() *Coco {
	return &Coco{}
}

func readLabels(labelsFile string) ([]string, error) {
	fileBytes, err := ioutil.ReadFile(labelsFile)
	if err != nil {
		return nil, fmt.Errorf("Unable to read labels file: %v", err)
	}

	return strings.Split(string(fileBytes), "\n"), nil
}

// Load loads the my_model SavedModel.
func (c *Coco) Load() error {
	model := tg.LoadModel(path, []string{"serve"}, nil)
	c.model = model

	var err error
	c.labels, err = readLabels(strings.Join([]string{path, "labels.txt"}, ""))
	if err != nil {
		return fmt.Errorf("Error loading labels file: %v", err)
	}
	return nil
}

// Predict predicts.
func (c *Coco) Predict(data []byte) *ObjectDetectionResponse {
	tensor, _ := makeTensorFromBytes(data)

	output := c.model.Exec(
		[]tf.Output{
			c.model.Op("StatefulPartitionedCall", 0),
		},
		map[tf.Output]*tf.Tensor{
			c.model.Op("serving_default_input_1", 0): tensor,
		},
	)

	outcome := NewObjectDetectionResponse(output, c.labels)
	return outcome
}

// Convert the image in filename to a Tensor suitable as input
func makeTensorFromBytes(bytes []byte) (*tf.Tensor, error) {
	// bytes to tensor
	tensor, err := tf.NewTensor(string(bytes))
	if err != nil {
		return nil, err
	}

	// create batch
	graph, input, output, err := makeTransformImageGraph("jpeg")
	if err != nil {
		return nil, err
	}

	// Execute that graph create the batch of that image
	session, err := tf.NewSession(graph, nil)
	if err != nil {
		return nil, err
	}

	defer session.Close()

	batch, err := session.Run(
		map[tf.Output]*tf.Tensor{input: tensor},
		[]tf.Output{output},
		nil)
	if err != nil {
		return nil, err
	}
	return batch[0], nil
}

// makeBatch uses ExpandDims to convert the tensor into a batch of size 1.
func makeBatch() (graph *tf.Graph, input, output tf.Output, err error) {
	s := op.NewScope()
	input = op.Placeholder(s, tf.String)

	decode := op.DecodeJpeg(s, input, op.DecodeJpegChannels(3))

	output = op.ExpandDims(s,
		op.Cast(s, decode, tf.Float),
		op.Const(s.SubScope("make_batch"), int32(0)))
	graph, err = s.Finalize()
	return graph, input, output, err
}

func makeTensorFromImage(imageBuffer *bytes.Buffer, imageFormat string) (*tf.Tensor, error) {
	tensor, err := tf.NewTensor(imageBuffer.String())
	if err != nil {
		return nil, err
	}
	graph, input, output, err := makeTransformImageGraph(imageFormat)
	if err != nil {
		return nil, err
	}
	session, err := tf.NewSession(graph, nil)
	if err != nil {
		return nil, err
	}
	defer session.Close()
	normalized, err := session.Run(
		map[tf.Output]*tf.Tensor{input: tensor},
		[]tf.Output{output},
		nil)
	if err != nil {
		return nil, err
	}
	return normalized[0], nil
}

func makeTransformImageGraph(imageFormat string) (graph *tf.Graph, input, output tf.Output, err error) {
	const (
		H, W  = 224, 224
		Mean  = float32(0)
		Scale = float32(255)
	)
	s := op.NewScope()
	input = op.Placeholder(s, tf.String)
	// Decode PNG or JPEG
	var decode tf.Output
	if imageFormat == "png" {
		decode = op.DecodePng(s, input, op.DecodePngChannels(3))
	} else {
		decode = op.DecodeJpeg(s, input, op.DecodeJpegChannels(3))
	}
	// Div and Sub perform (value-Mean)/Scale for each pixel
	output = op.Div(s,
		op.Sub(s,
			// Resize to 224x224 with bilinear interpolation
			op.ResizeBilinear(s,
				// Create a batch containing a single image
				op.ExpandDims(s,
					// Use decoded pixel values
					op.Cast(s, decode, tf.Float),
					op.Const(s.SubScope("make_batch"), int32(0))),
				op.Const(s.SubScope("size"), []int32{H, W})),
			op.Const(s.SubScope("mean"), Mean)),
		op.Const(s.SubScope("scale"), Scale))
	graph, err = s.Finalize()
	return graph, input, output, err
}

func normalizeImage(body io.ReadCloser) (*tf.Tensor, error) {
	var buf bytes.Buffer
	io.Copy(&buf, body)

	tensor, err := tf.NewTensor(buf.String())
	if err != nil {
		return nil, err
	}

	graph, input, output, err := getNormalizedGraph()
	if err != nil {
		return nil, err
	}

	session, err := tf.NewSession(graph, nil)
	if err != nil {
		return nil, err
	}

	normalized, err := session.Run(
		map[tf.Output]*tf.Tensor{
			input: tensor,
		},
		[]tf.Output{
			output,
		},
		nil)
	if err != nil {
		return nil, err
	}

	return normalized[0], nil
}

// Creates a graph to decode, rezise and normalize an image
func getNormalizedGraph() (graph *tf.Graph, input, output tf.Output, err error) {
	s := op.NewScope()
	input = op.Placeholder(s, tf.String)
	// 3 return RGB image
	decode := op.DecodeJpeg(s, input, op.DecodeJpegChannels(3))

	// Sub: returns x - y element-wise
	output = op.Sub(s,
		// make it 224x224: inception specific
		op.ResizeBilinear(s,
			// inserts a dimension of 1 into a tensor's shape.
			op.ExpandDims(s,
				// cast image to float type
				op.Cast(s, decode, tf.Float),
				op.Const(s.SubScope("make_batch"), int32(0))),
			op.Const(s.SubScope("size"), []int32{224, 224})),
		// mean = 117: inception specific
		op.Const(s.SubScope("mean"), float32(117)))
	graph, err = s.Finalize()

	return graph, input, output, err
}

func makeTensorFromImageForInception(filename string) (*tf.Tensor, error) {
	const (
		// Some constants specific to the pre-trained model at:
		// https://storage.googleapis.com/download.tensorflow.org/models/inception5h.zip
		//
		// - The model was trained after with images scaled to 224x224 pixels.
		// - The colors, represented as R, G, B in 1-byte each were converted to
		//   float using (value - Mean)/Std.
		//
		// If using a different pre-trained model, the values will have to be adjusted.
		H, W = 224, 224
		Mean = 117
		Std  = float32(1)
	)
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	sz := img.Bounds().Size()
	if sz.X != W || sz.Y != H {
		return nil, fmt.Errorf("input image is required to be %dx%d pixels, was %dx%d", W, H, sz.X, sz.Y)
	}
	// 4-dimensional input:
	// - 1st dimension: Batch size (the model takes a batch of images as
	//                  input, here the "batch size" is 1)
	// - 2nd dimension: Rows of the image
	// - 3rd dimension: Columns of the row
	// - 4th dimension: Colors of the pixel as (B, G, R)
	// Thus, the shape is [1, 224, 224, 3]
	var ret [1][H][W][3]float32
	for y := 0; y < H; y++ {
		for x := 0; x < W; x++ {
			px := x + img.Bounds().Min.X
			py := y + img.Bounds().Min.Y
			r, g, b, _ := img.At(px, py).RGBA()
			ret[0][y][x][0] = float32(int(b>>8)-Mean) / Std
			ret[0][y][x][1] = float32(int(g>>8)-Mean) / Std
			ret[0][y][x][2] = float32(int(r>>8)-Mean) / Std
		}
	}
	return tf.NewTensor(ret)
}
