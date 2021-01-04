package tagger

const (
	tagAPIURL     = "https://api.imagga.com/v2/tags"
	successStatus = "success"
)

type TagResponse struct {
	Result struct {
		Tags []struct {
			Confidence float64 `json:"confidence"`
			Tag        struct {
				En string `json:"en"`
			} `json:"tag"`
		} `json:"tags"`
	} `json:"result"`
	Status struct {
		Text string `json:"text"`
		Type string `json:"type"`
	} `json:"status"`
}
