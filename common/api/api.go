package api

//not used as of now
import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

var client *http.Client
var api *API

func init() {
	client = &http.Client{}
	api = &API{}
}

type API struct{}

type APIError struct {
	Code         string `json:"code"`
	ErrorMessage string `json:"error"`
}

func (m *APIError) Error() string {
	return fmt.Sprintf("Error making api call: Code %s %s", m.Code, m.ErrorMessage)
}

func decodeResponse(body io.Reader, to interface{}) error {
	// b, _ := ioutil.ReadAll(body)
	// fmt.Println("Body:",string(b))
	// err := json.Unmarshal(b, to)
	err := json.NewDecoder(body).Decode(to)

	if err != nil {
		return fmt.Errorf("error decoding body; %s", err.Error())
	}
	return nil
}

type APIResponse interface {
	ErrorMessage() string
	ErrorCode() string
}

func apiError(resp *http.Response, r APIResponse) error {

	if err := decodeResponse(resp.Body, r); err != nil {
		return err
	}

	var err APIError
	if r.ErrorMessage() != "" {
		err = APIError{Code: r.ErrorCode(), ErrorMessage: r.ErrorMessage()}
	} else {
		err = APIError{Code: strconv.Itoa(resp.StatusCode), ErrorMessage: resp.Status}
	}
	return &err
}

func Do(ctx context.Context, req *http.Request, r APIResponse) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second*1)
	defer cancel()
	req = req.WithContext(ctx)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		io.CopyN(ioutil.Discard, resp.Body, 512)
		resp.Body.Close()
	}()

	if resp.StatusCode != 200 {
		return apiError(resp, r)
	}

	return decodeResponse(resp.Body, r)
}
