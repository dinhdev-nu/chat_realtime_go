package response

const (
	SuccessCode = 2001
	ErrorCode   = 4001

	ServerErrorCode = 5001
)

const (
	SuccessMessage = "Success"
	ErrorMessage   = "Error"

	ServerErrorMessage = "Server Error"
)

var CodeMessage = map[int]string{
	SuccessCode:     SuccessMessage,
	ErrorCode:       ErrorMessage,
	ServerErrorCode: ServerErrorMessage,
}