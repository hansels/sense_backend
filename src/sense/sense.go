package sense

import (
	"cloud.google.com/go/firestore"
	"github.com/dgrijalva/jwt-go"
	"github.com/hansels/sense_backend/common/errors"
	"github.com/hansels/sense_backend/common/response"
	"github.com/hansels/sense_backend/common/router"
	"github.com/hansels/sense_backend/config"
	"github.com/hansels/sense_backend/src/firebase"
	"github.com/hansels/sense_backend/src/ml"
	"net/http"
	"strings"
)

type Opts struct {
	Firestore *firestore.Client
	Storage   *firebase.Storage
	Model     *ml.Coco
}

type Module struct {
	Firestore *firestore.Client
	Storage   *firebase.Storage
	Model     *ml.Coco
}

const authPrefix string = "Bearer "
const authPrefixLower string = "bearer "

func New(opts *Opts) *Module {
	return &Module{Firestore: opts.Firestore, Storage: opts.Storage, Model: opts.Model}
}

func getBearerToken(r *http.Request) string {
	var token string
	if authHeader := r.Header.Get("Authorization"); authHeader != "" {
		token = strings.Replace(authHeader, authPrefix, "", 1)
	}
	return token
}

func (m *Module) Authorize(h router.Handle) router.Handle {
	return func(w http.ResponseWriter, r *http.Request) *response.JSONResponse {
		userId, err := m.GetAuthorization(r)
		if err != nil {
			return response.NewJSONResponse().SetError(response.ErrForbiddenResource).SetLog("error", err).SetMessage("Unauthorized Access!")
		}
		r.Header.Set("UserID", userId)

		return h(w, r)
	}
}

func (m *Module) GetAuthorization(r *http.Request) (string, error) {
	var token string
	var authorized bool

	authToken := getBearerToken(r)
	if authToken != "" {
		authorized = true
		token = authToken
	}

	if authorized == false {
		return "", errors.New("Token is required!")
	}

	userId, err := m.GetAuthorizationFromToken(token)
	if err != nil {
		return "", err
	}

	return userId, nil
}

func (m *Module) GetAuthorizationFromToken(tokenString string) (string, error) {
	if tokenString == "" {
		return "", errors.New("Token should not be empty!")
	}

	if strings.Contains(tokenString, authPrefix) {
		tokenString = strings.Replace(tokenString, authPrefixLower, "", 1)
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if method, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Signing method invalid")
		} else if method != jwt.SigningMethodHS256 {
			return nil, errors.New("Signing method invalid")
		}
		return config.SignatureKey, nil
	})
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return "", errors.New("Token Unauthorized")
	}

	userId, err := CheckClaims(claims)
	if err != nil {
		return "", err
	}

	return userId, nil
}

func CheckClaims(claims jwt.MapClaims) (string, error) {
	email := claims["email"].(string)
	if email == "" {
		return "", errors.New("Token Unauthorized")
	}

	iss := claims["iss"].(string)
	if iss != "Sense" {
		return "", errors.New("Token Unauthorized")
	}

	// TOKEN HAS NO EXPIRATION - (Unrecommended)
	//exp := claims["exp"].(float64)
	//if int64(exp) < time.Now().Unix() {
	//	return "", errors.New("Token Expired")
	//}

	return email, nil
}
