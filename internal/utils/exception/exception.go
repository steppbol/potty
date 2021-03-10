package exception

const (
	Success             = 200
	BadRequest          = 400
	Unauthorized        = 401
	Forbidden           = 403
	NotFound            = 404
	Conflict            = 409
	InternalServerError = 500
)

var messageFlags = map[int]string{
	Success:             "ok",
	BadRequest:          "bad request",
	Unauthorized:        "unauthorized",
	Forbidden:           "forbidden",
	NotFound:            "record not found",
	Conflict:            "conflict",
	InternalServerError: "fail",
}

func GetMessage(code int) string {
	msg, ok := messageFlags[code]
	if ok {
		return msg
	}

	return messageFlags[InternalServerError]
}
