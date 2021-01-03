package request

import (
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"strings"

	"../config"
	"../database"
	"../storage"
	"../tagger"
	"../utils"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

type PhotoReponseData struct {
	ID   string
	Url  string
	Tags []string
}

// TODO: Fix up request handling
// https://stackoverflow.com/questions/43799703/how-to-do-error-http-error-handling-in-go-language
type RequestHandler struct {
	config         *config.Config
	s3Client       *storage.S3Client
	postgresClient *database.PostgresClient
	imaggaClient   *tagger.ImaggaClient
}

func NewRequestHandler(
	config *config.Config,
) *RequestHandler {
	return &RequestHandler{
		config:         config,
		s3Client:       storage.NewS3Client(config),
		postgresClient: database.NewPostgresClient(config),
		imaggaClient:   tagger.NewImaggaClient(config),
	}
}

func (handler *RequestHandler) HandleGetPhotos(
	w http.ResponseWriter, r *http.Request,
) {
	// TODO(Curtis): put this validation code into a validator function
	query := r.URL.Query()["query"]
	if len(query) != 1 {
		handler.sendStatusBadRequest(w, fmt.Errorf("multiple queries found in query"))
		return
	}
	tags := strings.Split(query[0], ",")

	photosIDsFound := make(map[uuid.UUID]bool, 0)
	for _, tag := range tags {
		photoIDs, err := handler.postgresClient.QueryAllPhotosWithTag(tag)
		if err != nil {
			handler.sendInternalServerError(w, err)
			return
		}
		for _, photoID := range photoIDs {
			photosIDsFound[photoID] = true
		}
	}
	photoIDs := utils.GetArrayOfUUIDFromMapOfUUID(photosIDsFound)
	handler.writePhotoResponseDataFromPhotoIDs(w, photoIDs)
}

func (handler *RequestHandler) HandleGetAllPhotos(
	w http.ResponseWriter, r *http.Request,
) {
	allPhotoIDs, err := handler.postgresClient.QueryAllPhotoIDs()
	if err != nil {
		handler.sendInternalServerError(w, err)
	}
	handler.writePhotoResponseDataFromPhotoIDs(w, allPhotoIDs)
}

func (handler *RequestHandler) writePhotoResponseDataFromPhotoIDs(
	w http.ResponseWriter,
	photoIDs []uuid.UUID,
) {
	allPhotos := make([]PhotoReponseData, 0, len(photoIDs))
	for _, photoID := range photoIDs {
		photoReponseData, err := handler.getPhotoReponseData(photoID)
		if err != nil {
			handler.sendInternalServerError(w, err)
		}
		allPhotos = append(allPhotos, photoReponseData)
	}
	fileUrlsBytes, _ := json.Marshal(allPhotos)
	handler.sendStandardHeaders(w)
	w.Write(fileUrlsBytes)
}

func (handler *RequestHandler) getPhotoReponseData(
	photoID uuid.UUID,
) (PhotoReponseData, error) {
	fileUrl := fmt.Sprintf("%s/%s/%s", handler.config.AWSConfig.S3Endpoint,
		handler.config.AWSConfig.S3BucketName, photoID)

	allTags, err := handler.postgresClient.QueryAllTagsForPhoto(photoID)
	if err != nil {
		return PhotoReponseData{}, err
	}

	return PhotoReponseData{
		ID:   photoID.String(),
		Url:  fileUrl,
		Tags: allTags,
	}, nil
}

// from: https://stackoverflow.com/questions/40684307/how-can-i-receive-an-uploaded-file-using-a-golang-net-http-server
func (handler *RequestHandler) HandleUpload(
	w http.ResponseWriter, r *http.Request,
) {
	// 1. Parse the file
	file, fileName, fileType, err := handler.parseFile(r)
	defer file.Close()
	if err != nil {
		handler.sendInternalServerError(w, err)
		return
	}

	// 2. upload the file to S3
	// TODO(Curtis): consider validating the fileType using: https://tika.apache.org/
	log.Infof("storing file with name=%s", fileName)
	photoID, err := handler.s3Client.UploadFile(file, fileType)
	if err != nil {
		handler.sendInternalServerError(w, err)
		return
	}

	// 3. Insert the Photo metadata into the DB
	err = handler.postgresClient.InsertPhoto(photoID, fileName, fileType)
	if err != nil {
		handler.sendInternalServerError(w, err)
		return
	}

	// 3. Store the tags
	tagIDs, err := handler.parseAndStoreTags(r)
	if err != nil {
		handler.sendInternalServerError(w, err)
		return
	}

	// 4. store the photo-tag associations
	err = handler.postgresClient.InsertPhotoTags(photoID, tagIDs, false)
	if err != nil {
		handler.sendInternalServerError(w, err)
	}
	return
}

func (handler *RequestHandler) parseFile(
	r *http.Request,
) (file multipart.File, fileName, fileType string, err error) {

	r.ParseMultipartForm(32 << 20) // limit your max input length!
	file, header, err := r.FormFile("file")
	if err != nil {
		return
	}
	name := strings.Split(header.Filename, ".")
	fileName = name[0]
	fileType = name[1]
	return
}

func (handler *RequestHandler) parseAndStoreTags(
	r *http.Request,
) (tagIDs []uuid.UUID, err error) {
	if formError := r.ParseForm(); formError != nil {
		err = formError
		return
	}
	tagsFormValue := r.PostForm["tags"]
	if len(tagsFormValue) == 0 {
		// No tags
		return nil, nil
	}

	tags := []string{}
	json.Unmarshal([]byte(tagsFormValue[0]), &tags)

	tagIDs = make([]uuid.UUID, 0, len(tags))
	for _, tagName := range tags {
		tagID, insertErr := handler.postgresClient.InsertTagIfNotExist(tagName)
		if insertErr != nil {
			return nil, fmt.Errorf("cannot insert tagName=%s, err=%s", tagName, insertErr)
		}
		tagIDs = append(tagIDs, tagID)
	}
	return
}

func (handler *RequestHandler) sendStandardHeaders(
	w http.ResponseWriter,
) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")
}

func (handler *RequestHandler) sendInternalServerError(
	w http.ResponseWriter,
	err error,
) {
	handler.sendStandardHeaders(w)
	log.Error(err)
	w.WriteHeader(http.StatusInternalServerError)
	return
}

func (handler *RequestHandler) sendStatusBadRequest(
	w http.ResponseWriter,
	err error,
) {
	handler.sendStandardHeaders(w)
	log.Error(err)
	w.WriteHeader(http.StatusBadRequest)
	return
}
