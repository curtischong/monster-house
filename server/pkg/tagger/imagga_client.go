package tagger

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"

	"../config"
	"github.com/google/uuid"
)

// This client is used to call the Imagga API to generate our photo tags

type ImaggaClient struct {
	config    *config.Config
	basicAuth string
}

func NewImaggaClient(
	config *config.Config,
) *ImaggaClient {
	basicAuthString := fmt.Sprintf("%s:%s",
		config.SecretsConfig.ImaggaUser, config.SecretsConfig.ImaggaSecret)
	basicAuth := base64.StdEncoding.EncodeToString([]byte(basicAuthString))

	return &ImaggaClient{
		config:    config,
		basicAuth: basicAuth,
	}
}

func (client *ImaggaClient) GetTagsForPhoto(
	photoID uuid.UUID,
) ([]string, error) {
	// 1. query for the tags from the API
	resp, err := client.queryTagsForPhoto(photoID)
	if err != nil {
		return nil, fmt.Errorf("cannot queryTagsForPhoto err=%s", err)
	}

	// 2. Parse the response and validate that the request was successful
	tagResponse, err := client.getTagResponseAndValidateResponseSuccess(resp, photoID)
	if err != nil {
		return nil, err
	}
	// 3. Fetch the n best generated tags from the response
	return client.getTopTagsInResponse(tagResponse, photoID), nil
}

// queryTagsForPhoto makes an http request to imagga to tag the photo
// TODO(Curtis): implement retry logic
func (client *ImaggaClient) queryTagsForPhoto(
	photoID uuid.UUID,
) (*http.Response, error) {
	// Since the photos are stored in the local s3, we need to expose them to the internet using ngrok
	// This builds the photoURL to send to imagga
	photoURL := fmt.Sprintf("%s/%s/%s",
		client.config.ServerConfig.NgrokURL,
		client.config.AWSConfig.S3BucketName,
		photoID.String())

	// Now make the request and return the response
	req, err := http.NewRequest("GET", tagAPIURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Basic "+client.basicAuth)

	q := req.URL.Query()
	q.Add("image_url", photoURL)
	req.URL.RawQuery = q.Encode()
	log.Infof("fetch imagga image tags at URL=%s", req.URL.String())

	httpClient := http.Client{}
	return httpClient.Do(req)
}

// getTagResponseAndValidateResponseSuccess unmarshalls the response into a TagResponse object and also vli
func (client *ImaggaClient) getTagResponseAndValidateResponseSuccess(
	resp *http.Response,
	photoID uuid.UUID,
) (TagResponse, error) {
	if resp.StatusCode != http.StatusOK {
		return TagResponse{}, fmt.Errorf("imagga tag query for photoId=%s returned statusCode=%d",
			photoID, resp.StatusCode)
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return TagResponse{}, err
	}

	var response TagResponse
	if err := json.Unmarshal(bodyBytes, &response); err != nil {
		return TagResponse{}, err
	}

	if response.Status.Type != successStatus {
		return TagResponse{}, fmt.Errorf("cannot fetch tags for image. Status.Type=%s, Status.Text=%s",
			response.Status.Type, response.Status.Text)
	}
	return response, nil
}

// getTopTagsInResponse returns the top N tags from the response from the Imagga API
// Note: The top N tags is determined by a config variable
func (client *ImaggaClient) getTopTagsInResponse(
	tagResponse TagResponse,
	photoID uuid.UUID,
) (bestGeneratedTags []string) {
	bestGeneratedTags = make([]string, 0)
	// Note: the API already sorts the tag confidence in descending order
	maxNumberOfTags := client.config.ServerConfig.UseNBestGeneratedTags
	for i := 0; i < maxNumberOfTags && i < len(tagResponse.Result.Tags); i++ {
		tagData := tagResponse.Result.Tags[i]
		tag := tagData.Tag.En
		bestGeneratedTags = append(bestGeneratedTags, tag)
	}
	log.Infof("the generated tags for photoID=%s are: %s", photoID, strings.Join(bestGeneratedTags, ", "))
	return
}
