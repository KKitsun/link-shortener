package response

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

func ErrorResponse(errText string, err error) Response {
	return Response{
		Status: errText,
		Error:  err.Error(),
	}
}
