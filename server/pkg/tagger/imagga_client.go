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
	photoURL := fmt.Sprintf("%s/%s/%s",
		client.config.ServerConfig.NgrokURL,
		client.config.AWSConfig.S3BucketName,
		photoID.String())

	req, err := http.NewRequest("GET", tagAPIURL, nil)
	req.Header.Add("Authorization", "Basic "+client.basicAuth)

	q := req.URL.Query()
	q.Add("image_url", photoURL)
	req.URL.RawQuery = q.Encode()
	log.Infof("fetch imagga image tags at URL=%s", req.URL.String())

	httpClient := http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("imagga tag query for photoId=%s returned statusCode=%s", photoID, resp.StatusCode)
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response TagResponse
	if err := json.Unmarshal(bodyBytes, &response); err != nil {
		return nil, err
	}

	if response.Status.Type != successStatus {
		return nil, fmt.Errorf("cannot fetch tags for image. Status.Type=%s, Status.Text=%s",
			response.Status.Type, response.Status.Text)
	}

	// Now fetch the n best generated tags
	bestGeneratedTags := make([]string, 0)
	// Note: the API already sorts the tag confidence in descending order
	maxNumberOfTags := client.config.ServerConfig.UseNBestGeneratedTags
	for i := 0; i < maxNumberOfTags && i < len(response.Result.Tags); i++ {
		tagData := response.Result.Tags[i]
		tag := tagData.Tag.En
		bestGeneratedTags = append(bestGeneratedTags, tag)
	}
	log.Infof("the generated tags for photoID=%s are: %s", photoID, strings.Join(bestGeneratedTags, ", "))
	return bestGeneratedTags, nil
}
