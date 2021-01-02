package request

import (
	"../config"
	"../storage"
	"fmt"
	"net/http"
	"strings"
	log "github.com/sirupsen/logrus"
)

// TODO: Store s3client
type RequestHandler struct{
	s3client *storage.S3Client
}

func NewRequestHandler(
	config *config.Config	,
) *RequestHandler{
	return &RequestHandler{
		s3client: storage.NewS3Client(config),
	}
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
	fmt.Printf("File name %s\n", name[0])
	err = handler.s3client.UploadFile(file)

	if err != nil{
		handler.sendInternalServerError(w, err)
	}
	return
}

func (handler *RequestHandler) sendInternalServerError(
	w http.ResponseWriter,
	err error,
){
	log.Error(err)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(http.StatusInternalServerError)
	return
}