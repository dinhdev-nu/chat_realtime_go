package response

const (
	SuccessCode = 2001
	ErrorCode   = 4001

	InvalidToken              = 3001
	InvalidRequestPayloadCode = 4000

	ErrorCodeEmailExist     = 4002
	ErrorOtpFail            = 4003
	ErrorOtpNotExist        = 4004
	ErrorCodeInvalidRequest = 4005

	ServerErrorCode = 5001
)

const (
	SuccessMessage = "Success"
	ErrorMessage   = "Error"

	InvalidTokenMessage            = "Invalid Token"
	ErrorEmailExistMessage         = "Email Exist"
	ServerErrorMessage             = "Server Error"
	ErrorOtpFailMessage            = "Otp Fail"
	ErrorOtpNotExistMessage        = "Otp Not Exist"
	ErrorCodeInvalidRequestMessage = "Invalid Request"
	InvalidRequestPayloadtMessage  = "Invalid Request Payload"
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
}
