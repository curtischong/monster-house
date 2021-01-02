package request

import (
	"fmt"
	"net/http"
	"strings"
)

// TODO: Store s3client
type RequestHandler struct{}

func NewRequestHandler() *RequestHandler{
	return &RequestHandler{}
}

// from: https://stackoverflow.com/questions/40684307/how-can-i-receive-an-uploaded-file-using-a-golang-net-http-server
func (handler* RequestHandler) HandleUpload(
	w http.ResponseWriter, r *http.Request,
) {
	r.ParseMultipartForm(32 << 20) // limit your max input length!
	file, header, err := r.FormFile("file")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	name := strings.Split(header.Filename, ".")
	fmt.Printf("File name %s\n", name[0])
	return
}