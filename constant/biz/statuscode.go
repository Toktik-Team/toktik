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
	UnableToQueryIsFollow
	UnableToDeleteRelation
	UnableToLike
	UnableToCancelLike
	UnableToQueryFavorite
	UnableToQueryTotalFavorited
)

// Bad Request
const (
	UserNameExist = 400000 + iota
	UserNotFound
	InvalidCommentActionType
	VideoNotFound
	Unable2ParseLatestTimeStatusCode
	InvalidContentType
	RelationNotFound
	RelationAlreadyExists
	InvalidToUserId
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
