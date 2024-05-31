package mixinapi

type ResponseType int

const (
	JsonResponseType ResponseType = iota + 1
	StringResponseType
	HTMLResponseType
	FileResponseType
	RedirectResponseType
	StreamResponseType
	AnyResponseType
)

type Response struct {
	BizCode     int          `json:"-"; description:"Business code"`
	StatusCode  int          `json:"-"; description:"HTTP status code"`
	Content     any          `json:"-"; description:"HTTP response body"`
	ContentType string       `json:"-"; description:"HTTP response content type"`
	Type        ResponseType `json:"-"; description:"HTTP response type"`
}