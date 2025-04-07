package response

const (
	SuccessCode         = 2001
	ErrorCode           = 4001
	ForbiddenCode       = 4031
	TooManyRequestsCode = 4291

	InvalidToken              = 3001
	InvalidRequestPayloadCode = 4000

	ErrorCodeEmailExist = 4002
	ErrorOtpNotExist    = 4004

	ErrorOtpFail            = 4003
	ErrorCodeInvalidRequest = 4005
	ErrorUserNotExist       = 4006
	ErrorPasswordNotMatch   = 4007

	ServerErrorCode = 5001

	ErrorCreateCode = 6002
	ErrorDeleteCode = 6001
	ErrorUpdateCode = 6000

	// update
	ErrorUpdateLoginCode = 6003
)

const (
	SuccessMessage         = "Success"
	ErrorMessage           = "Error"
	ForbiddenMessage       = "Forbidden"
	TooManyRequestsMessage = "Too Many Requests"

	InvalidTokenMessage            = "Invalid Token"
	ErrorEmailExistMessage         = "Email Exist"
	ServerErrorMessage             = "Server Error"
	ErrorOtpFailMessage            = "Otp Fail"
	ErrorOtpNotExistMessage        = "Otp Not Exist"
	ErrorCodeInvalidRequestMessage = "Invalid Request"
	InvalidRequestPayloadtMessage  = "Invalid Request Payload"
	ErrorCreateMessage             = "Error Create Record"
	ErrorUserNotExistMessage       = "User Not Exist"
	ErrorPasswordNotMatchMessage   = "Password Not Match"
	ErrorUpdateLoginMessage        = "Error Update Login Time"
	ErrorDeleteMessage             = "Error Delete Record"
	ErrorUpdateMessage             = "Error Update Record"
)

var CodeMessage = map[int]string{
	SuccessCode:               SuccessMessage,
	ErrorCode:                 ErrorMessage,
	ServerErrorCode:           ServerErrorMessage,
	InvalidToken:              InvalidTokenMessage,
	ErrorCodeEmailExist:       ErrorEmailExistMessage,
	ErrorOtpFail:              ErrorOtpFailMessage,
	ErrorOtpNotExist:          ErrorOtpNotExistMessage,
	ErrorCodeInvalidRequest:   ErrorCodeInvalidRequestMessage,
	InvalidRequestPayloadCode: InvalidRequestPayloadtMessage,
	ErrorCreateCode:           ErrorCreateMessage,
	ErrorUserNotExist:         ErrorUserNotExistMessage,
	ErrorPasswordNotMatch:     ErrorPasswordNotMatchMessage,
	ErrorUpdateLoginCode:      ErrorUpdateLoginMessage,
	ForbiddenCode:             ForbiddenMessage,
	ErrorDeleteCode:           ErrorDeleteMessage,
	ErrorUpdateCode:           ErrorUpdateMessage,
	TooManyRequestsCode:       TooManyRequestsMessage,
}
