package request

import (
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"strings"

	"../common"
	"../config"
	"../database"
	"../storage"
	"../tagger"
	"../utils"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

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

// HandleGetPhotos fetches the photos that have at least 1 matching tag in the query
func (handler *RequestHandler) HandleGetPhotos(
	w http.ResponseWriter, r *http.Request,
) {
	query := r.URL.Query()["query"]
	log.Infof("Fetching photos with query=%s", query)

	// validate that the query is valid
	if len(query) != 1 {
		handler.sendStatusBadRequest(w, fmt.Errorf("multiple queries found in query"))
		return
	}

	tags := strings.Split(query[0], ",")
	// Now find the photos that have at least 1 matching tag
	matchingPhotoIDs, err := handler.findPhotosIDsWithTags(tags)
	if err != nil {
		handler.sendInternalServerError(w, err)
		return
	}
	handler.writePhotoResponseDataFromPhotoIDs(w, matchingPhotoIDs)
}

// findPhotosIDsWithTags returns a slice of photoIDs that have at least 1 tag
// in the slice of tags (this is an OR operation not AND)
func (handler *RequestHandler) findPhotosIDsWithTags(
	tags []string,
) ([]uuid.UUID, error) {
	// Note: Since a photo might have multiple tags that match,
	// so we have to be careful to not double count each photoID.
	// To account for this, we are using a map to act as a "set"
	photosIDsFound := make(map[uuid.UUID]bool, 0)
	for _, tag := range tags {
		trimmedTag := strings.TrimSpace(tag)
		photoIDs, err := handler.postgresClient.QueryAllPhotosWithTag(trimmedTag)
		if err != nil {
			return nil, err
		}
		for _, photoID := range photoIDs {
			// If the photoID was already in the photosIDsFound map, this will do nothing
			// since the photo was already matched!
			photosIDsFound[photoID] = true
		}
	}
	return utils.GetArrayOfUUIDFromMapOfUUID(photosIDsFound), nil
}

// HandleGetAllPhotos fetches all photos stored in the repository
func (handler *RequestHandler) HandleGetAllPhotos(
	w http.ResponseWriter, r *http.Request,
) {
	log.Info("Fetching all photos")
	allPhotoIDs, err := handler.postgresClient.QueryAllPhotoIDs()
	if err != nil {
		handler.sendInternalServerError(w, err)
		return
	}
	handler.writePhotoResponseDataFromPhotoIDs(w, allPhotoIDs)
}

// writePhotoResponseDataFromPhotoIDs fetches the response data for each photoID
// and marshals it into a response for the client
func (handler *RequestHandler) writePhotoResponseDataFromPhotoIDs(
	w http.ResponseWriter,
	photoIDs []uuid.UUID,
) {
	allPhotos := make([]common.PhotoReponseData, 0, len(photoIDs))
	for _, photoID := range photoIDs {
		photoReponseData, err := handler.getPhotoResponseData(photoID)
		if err != nil {
			handler.sendInternalServerError(w, err)
			return
		}
		allPhotos = append(allPhotos, photoReponseData)
	}
	handler.sendStatusOK(w)
	fileUrlsBytes, _ := json.Marshal(allPhotos)
	w.Write(fileUrlsBytes)
}

// getPhotoResponseData fetches metadata for the photoID from the DB and returns the data as a PhotoReponseData object
func (handler *RequestHandler) getPhotoResponseData(
	photoID uuid.UUID,
) (common.PhotoReponseData, error) {
	fileUrl := fmt.Sprintf("%s/%s/%s", handler.config.AWSConfig.S3Endpoint,
		handler.config.AWSConfig.S3BucketName, photoID)

	allTags, err := handler.postgresClient.QueryAllTagsForPhoto(photoID)
	if err != nil {
		return common.PhotoReponseData{}, err
	}

	return common.PhotoReponseData{
		ID:   photoID.String(),
		Url:  fileUrl,
		Tags: allTags,
	}, nil
}

// HandleUpload saves the uploaded file into s3, generates tags for the photo, and saves the tags into the DB
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

	// 4. Calculate the tags from the image
	userTags, err := handler.parseTagsFromRequest(r)
	if err != nil {
		handler.sendInternalServerError(w, err)
		return
	}

	generatedTags, err := handler.imaggaClient.GetTagsForPhoto(photoID)
	if err != nil {
		handler.sendInternalServerError(w, err)
		return
	}

	// 5. Insert the tag information into the DB
	userTagIDs, err := handler.insertTagsAndGetTagIDs(userTags)
	if err != nil {
		handler.sendInternalServerError(w, err)
		return
	}

	generatedTagIDs, err := handler.insertTagsAndGetTagIDs(generatedTags)
	if err != nil {
		handler.sendInternalServerError(w, err)
		return
	}

	// 6. store the photo-tag associations in the DB
	err = handler.postgresClient.InsertPhotoTags(photoID, userTagIDs, false)
	if err != nil {
		handler.sendInternalServerError(w, err)
		return
	}
	err = handler.postgresClient.InsertPhotoTags(photoID, generatedTagIDs, true)
	if err != nil {
		handler.sendInternalServerError(w, err)
		return
	}

	handler.sendStatusOK(w)
	return
}

func (handler *RequestHandler) parseFile(
	r *http.Request,
) (file multipart.File, fileName, fileType string, err error) {
	// from: https://stackoverflow.com/questions/40684307/how-can-i-receive-an-uploaded-file-using-a-golang-net-http-server
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

func (handler *RequestHandler) parseTagsFromRequest(
	r *http.Request,
) ([]string, error) {
	if formError := r.ParseForm(); formError != nil {
		return nil, formError
	}
	tagsFormValue := r.PostForm["tags"]
	if len(tagsFormValue) == 0 {
		// No tags
		return nil, nil
	}

	tags := []string{}
	json.Unmarshal([]byte(tagsFormValue[0]), &tags)
	return tags, nil
}

func (handler *RequestHandler) insertTagsAndGetTagIDs(
	tags []string,
) (tagIDs []uuid.UUID, err error) {

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
