package firebase

import (
	"bytes"
	"cloud.google.com/go/firestore"
	"cloud.google.com/go/storage"
	"context"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"fmt"
	"github.com/google/uuid"
	"github.com/hansels/sense_backend/common/log"
	"github.com/hansels/sense_backend/config"
	"google.golang.org/api/option"
	"io"
	"strings"
)

var (
	App *firebase.App
)

type Storage struct {
	FirebaseStorage *storage.BucketHandle
}

type UploadImageData struct {
	Ctx      context.Context
	FileName string
	File     []byte
}

func init() {
	opt := option.WithCredentialsFile("./files/firebase/firebase.json")
	cfg := &firebase.Config{
		ProjectID:     config.FirebaseProjectId,
		StorageBucket: config.FirebaseStorageBucket,
	}

	var err error
	App, err = firebase.NewApp(context.Background(), cfg, opt)
	if err != nil {
		log.Fatalf("error initializing firebase: %v\n", err)
	}

	log.Infoln("Firebase Initialization Success 🔥")
}

func InitStorage() *Storage {
	client, err := App.Storage(context.Background())
	if err != nil {
		log.Fatalln("error initializing storage: %v\n", err)
	}

	bucket, err := client.DefaultBucket()
	if err != nil {
		log.Fatalln("error initializing default bucket: %v\n", err)
	}

	return &Storage{FirebaseStorage: bucket}
}

func InitAuth() *auth.Client {
	client, err := App.Auth(context.Background())
	if err != nil {
		log.Fatalln("error initializing authentication: %v\n", err)
	}

	return client
}

func InitFirestore() *firestore.Client {
	fs, err := App.Firestore(context.Background())
	if err != nil {
		log.Fatalf("error initializing firestore: %v\n", err)
	}
	log.Infoln("Firestore Roll-on 🔥🔥🔥")
	return fs
}

func (s *Storage) UploadImage(data UploadImageData) (string, error) {
	id := uuid.New()

	object := s.FirebaseStorage.Object(data.FileName)
	writer := object.NewWriter(data.Ctx)
	writer.ObjectAttrs.Metadata = map[string]string{"firebaseStorageDownloadTokens": id.String()}
	defer writer.Close()

	if _, err := io.Copy(writer, bytes.NewReader(data.File)); err != nil {
		return "", err
	}

	url := s.GenerateURL(data, id.String())
	return url, nil
}

func (s *Storage) GenerateURL(data UploadImageData, uuid string) string {
	url := fmt.Sprintf("https://firebasestorage.googleapis.com/v0/b/%s/o/%s?alt=media&token=%s", config.FirebaseStorageBucket, data.FileName, uuid)
	strings.ReplaceAll(url, "@", "%40")
	strings.ReplaceAll(url, ":", "%3A")

	return url
}
