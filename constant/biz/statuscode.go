package biz

// OK
const (
	OkStatusCode = 0 + iota
)

// Service Not Available
const (
	ServiceNotAvailable = 503000
)

// Internal Server Error
const (
	RedisError = 500000 + iota
	ProtoMarshalError
	ProtoUnmarshalError
	UnableToCreateComment
	UnableToDeleteComment
	UnableToQueryVideo
	UnableToQueryComment
	UnableToQueryUser
	SQLQueryErrorStatusCode
	Unable2GenerateUUID
	Unable2CreateThumbnail
	Unable2UploadVideo
	Unable2UploadCover
	Unable2CreateDBEntry
	RequestIsNil
	UnableToQueryFollowList
	UnableToQueryFollowerList
	UnableToCreateRelation
	UnableToDeleteRelation
	UnableToLike
	UnableToCancelLike
)

// Bad Request
const (
	UserNameExist                    = 400000 + iota
	UserNotFound                     = 400002
	InvalidCommentActionType         = 400001
	VideoNotFound                    = 400002
	Unable2ParseLatestTimeStatusCode = 400001
	InvalidContentType               = 400101
	TargetUserNotFound               = 400101
	RelationNotFound                 = 400102
	RelationAlreadyExists            = 400103
	InvalidToUserId                  = 400104
)

// Unauthorized
const (
	PasswordIncorrect = 401003
	TokenNotFound     = 401001
)

// Forbidden
const (
	ActorIDNotMatch = 403001
)
