package httpx

import "fmt"

type HttpError struct {
	StatusCode int
	Resp       string
}

func (e HttpError) Error() string {
	return fmt.Sprintf("StatusCode[%+v], Resp[%+v]", e.StatusCode, e.Resp)
}
