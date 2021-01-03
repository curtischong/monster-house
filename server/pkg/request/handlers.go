package request

import (
	"../config"
	"../database"
	"../storage"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"mime/multipart"
	"net/http"
	"strings"
	"github.com/google/uuid"
	)

// TODO: Store s3Client
type RequestHandler struct{
	s3Client *storage.S3Client
	postgresClient *database.PostgresClient
}

func NewRequestHandler(
	config *config.Config	,
) *RequestHandler{
	return &RequestHandler{
		s3Client:       storage.NewS3Client(config),
		postgresClient: database.NewPostgresClient(config),
	}
}

func (handler *RequestHandler) HandleGetAllPhotos(
	w http.ResponseWriter, r *http.Request,
){
	fileUrls , err := handler.s3Client.GetAllFileURLs()
	if err != nil{
		handler.sendInternalServerError(w, err)
	}
	fileUrlsBytes, _ := json.Marshal(fileUrls)
	handler.sendStandardHeaders(w)
	w.Write(fileUrlsBytes)
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
	}

	// 2. upload the file to S3
	// TODO(Curtis): consider validating the fileType using: https://tika.apache.org/
	log.Infof("storing file with name=%s", fileName)
	photoID, err := handler.s3Client.UploadFile(file, fileType)
	if err != nil{
		handler.sendInternalServerError(w, err)
	}

	// 3. Insert the Photo metadata into the DB
	err = handler.postgresClient.InsertPhoto(photoID, fileName, fileType)
	if err != nil{
		handler.sendInternalServerError(w, err)
	}

	// 3. Store the tags
	tagIDs, err := handler.parseAndStoreTags(r)
	if err != nil{
		handler.sendInternalServerError(w, err)
	}

	// 4. store the photo-tag associations
	err = handler.postgresClient.InsertPhotoTags(photoID, tagIDs, false)
	if err != nil{
		handler.sendInternalServerError(w, err)
	}
	return
}

func (handler *RequestHandler) parseFile(
	r *http.Request,
) (file multipart.File, fileName, fileType string, err error){

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
)(tagIDs []uuid.UUID, err error){
	if formError := r.ParseForm(); formError != nil {
		err = formError
		return
	}
	tagsFormValue := r.PostForm["tags"]
	if len(tagsFormValue) == 0{
		// No tags
		return nil, nil
	}

	tags := []string{}
	json.Unmarshal([]byte(tagsFormValue[0]), &tags)

	tagIDs = make([]uuid.UUID, 0, len(tags))
	for _, tagName := range tags{
		tagID, insertErr := handler.postgresClient.InsertTagIfNotExist(tagName)
		if insertErr != nil{
			return nil, fmt.Errorf("cannot insert tagName=%s, err=%s", tagName, insertErr)
		}
		tagIDs = append(tagIDs, tagID)
	}
	return
}

func (handler *RequestHandler) sendStandardHeaders(
	w http.ResponseWriter,
){
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")
}

func (handler *RequestHandler) sendInternalServerError(
	w http.ResponseWriter,
	err error,
){
	handler.sendStandardHeaders(w)
	log.Error(err)
	w.WriteHeader(http.StatusInternalServerError)
	return
}