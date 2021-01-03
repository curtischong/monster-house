package request

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"strings"
	"../config"
	"../storage"
	"../database"
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
	r.ParseMultipartForm(32 << 20) // limit your max input length!
	file, header, err := r.FormFile("file")
	if err != nil {
		handler.sendInternalServerError(w, err)
	}
	defer file.Close()
	name := strings.Split(header.Filename, ".")

	// TODO(Curtis): consider validating the fileType using: https://tika.apache.org/
	fmt.Printf("File name %s\n", name[0])
	fileType := name[1]
	err = handler.s3Client.UploadFile(file, fileType)

	if err != nil{
		handler.sendInternalServerError(w, err)
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