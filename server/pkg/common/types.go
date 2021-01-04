package common

type TagResponseData struct {
	Name        string
	IsGenerated bool
}

type PhotoReponseData struct {
	ID   string
	Url  string
	Tags []TagResponseData
}
