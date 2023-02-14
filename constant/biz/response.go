package biz

const (
	OkStatusCode = 0
)

var (
	OkStatusMsg = "OK"

	BadRequestStatusMsg          = "Unable to finish request, please check your parameters. If you think this is a bug, please contact us."
	ForbiddenStatusMsg           = "You are not allowed to access this resource."
	InternalServerErrorStatusMsg = "The server had an error while processing your request. Sorry about that!"
)
