package api

import (
	"context"
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"github.com/fatih/structs"
	"github.com/google/uuid"
	"github.com/hansels/sense_backend/common/log"
	"github.com/hansels/sense_backend/common/response"
	"github.com/hansels/sense_backend/config"
	"github.com/hansels/sense_backend/src/firebase"
	"github.com/hansels/sense_backend/src/ml"
	"github.com/hansels/sense_backend/src/model"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"net/http"
	"time"
)

type MyClaims struct {
	jwt.StandardClaims
	Email string `json:"email"`
}

func (a *API) CheckUser(w http.ResponseWriter, r *http.Request) *response.JSONResponse {
	ctx := context.Background()

	var req model.CheckUserData
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Errorln("CheckUserData Json Decode Error : %+v", err)
		return response.NewJSONResponse().SetData(false)
	}

	_, err = a.Module.Firestore.Collection("users").Doc(req.Email).Get(ctx)
	if err != nil {
		return response.NewJSONResponse().SetData(false)
	}

	return response.NewJSONResponse().SetData(true)
}

func (a *API) Login(w http.ResponseWriter, r *http.Request) *response.JSONResponse {
	ctx := context.Background()

	var req model.LoginData
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Errorln("LoginData Json Decode Error : %+v", err)
		return response.NewJSONResponse().SetError(response.ErrBadRequest).SetMessage("Bad Request")
	}

	user, err := a.getUser(ctx, req.Email)
	if err != nil {
		return response.NewJSONResponse().SetError(response.ErrNoValidUserFound).SetMessage("Login Unsuccessful!")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(config.PasswordSalt+req.Password))
	if err != nil {
		return response.NewJSONResponse().SetError(response.ErrNoValidUserFound).SetMessage("Login Unsuccessful!")
	}

	claims := MyClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Duration(720) * time.Hour).Unix(),
			Issuer:    "Sense",
		},
		Email: user.Email,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString(config.SignatureKey)
	if err != nil {
		log.Errorln("JWT Token Signing error : %+v", err)
		return response.NewJSONResponse().SetError(response.ErrBadRequest).SetMessage("Bad Request")
	}

	// Never send password to FE (it's dangerous)
	user.Password = ""
	return response.NewJSONResponse().SetData(map[string]interface{}{"token": signedToken, "user": structs.Map(user)})
}

func (a *API) RegisterUser(w http.ResponseWriter, r *http.Request) *response.JSONResponse {
	ctx := context.Background()

	var user model.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Errorln("RegisterData Json Decode Error : %+v", err)
		return response.NewJSONResponse().SetError(response.ErrBadRequest).SetMessage("Bad Request")
	}

	doc := a.Module.Firestore.Collection("users").Doc(user.Email)

	_, err = doc.Get(ctx)
	if err == nil {
		return response.NewJSONResponse().SetError(response.ErrAlreadyRegistered).SetMessage("User Already Registered!")
	}

	password, err := bcrypt.GenerateFromPassword([]byte(config.PasswordSalt+user.Password), bcrypt.DefaultCost)
	if err != nil {
		return response.NewJSONResponse().SetError(response.ErrBadRequest).SetMessage("Bad Request")
	}

	user.Password = string(password)
	user.Type = "Member"

	_, err = doc.Set(ctx, structs.Map(user))
	if err != nil {
		log.Errorln("Write to Firestore error : %+v", err)
		return response.NewJSONResponse().SetError(response.ErrBadRequest).SetMessage("Bad Request")
	}

	return response.NewJSONResponse().SetData("OK")
}

func (a *API) Predict(w http.ResponseWriter, r *http.Request) *response.JSONResponse {
	ctx := context.Background()

	// ML Prediction
	mlModel := a.Module.Model
	file, _, err := r.FormFile("data")
	if err != nil {
		return response.NewJSONResponse().SetError(response.ErrInternalServerError).SetMessage("Internal Server Error")
	}
	defer file.Close()

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return response.NewJSONResponse().SetError(response.ErrInternalServerError).SetMessage("Internal Server Error")
	}

	outcome := mlModel.Predict(fileBytes)
	result, err := a.generateResultFromML(outcome)
	if err != nil {
		return response.NewJSONResponse().SetError(response.ErrBadRequest).SetMessage("Bad Request")
	}

	// Upload Image to Storage Firebase
	id, err := uuid.NewUUID()
	if err != nil {
		return response.NewJSONResponse().SetError(response.ErrInternalServerError).SetMessage("Internal Server Error")
	}

	uploadData := firebase.UploadImageData{
		Ctx:      ctx,
		FileName: id.String(),
		File:     fileBytes,
	}
	url, err := a.Module.Storage.UploadImage(uploadData)
	if err != nil {
		return response.NewJSONResponse().SetError(response.ErrInternalServerError).SetMessage("Internal Server Error")
	}

	// URL for Download Image, currently log for health checking
	log.Infoln(url)
	result.Image = url
	return response.NewJSONResponse().SetData(structs.Map(result))
}

func (a *API) Ping(w http.ResponseWriter, r *http.Request) *response.JSONResponse {
	log.Println("PING Called!")
	return response.NewJSONResponse().SetData("Ping!!!")
}

func (a *API) InsertResort(w http.ResponseWriter, r *http.Request) *response.JSONResponse {
	ctx := context.Background()

	var resort model.Resort
	err := json.NewDecoder(r.Body).Decode(&resort)
	if err != nil {
		log.Errorln("Resort Json Decode Error : %+v", err)
		return response.NewJSONResponse().SetError(response.ErrBadRequest).SetMessage("Bad Request")
	}

	doc := a.Module.Firestore.Collection("resorts").Doc(resort.Name)

	_, err = doc.Set(ctx, structs.Map(resort))
	if err != nil {
		log.Errorln("Write to Firestore error : %+v", err)
		return response.NewJSONResponse().SetError(response.ErrBadRequest).SetMessage("Bad Request")
	}

	return response.NewJSONResponse().SetData("OK")
}

func (a *API) generateResultFromML(outcome *ml.ObjectDetectionResponse) (*model.PredictionResult, error) {
	if outcome.NumDetections == 0 {
		return &model.PredictionResult{IsDetected: false}, nil
	}

	result := &model.PredictionResult{IsDetected: true}
	result.Verdict = outcome.Detections[0].Label

	return result, nil
}

func (a *API) getUser(ctx context.Context, id string) (*model.User, error) {
	ds, err := a.Module.Firestore.Collection("users").Doc(id).Get(ctx)
	if err != nil {
		return nil, err
	}

	jsonString, err := json.Marshal(ds.Data())
	if err != nil {
		log.Errorln("User Marshal Error : %+v", err)
		return nil, err
	}

	user := &model.User{}
	err = json.Unmarshal(jsonString, &user)
	if err != nil {
		log.Errorln("User Unmarshal Error : %+v", err)
		return nil, err
	}
	return user, nil
}
