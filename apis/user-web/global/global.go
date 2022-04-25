package global

import (
	ut "github.com/go-playground/universal-translator"
	"luke544187758/user-web/proto"
)

var (
	Trans ut.Translator

	UserServiceClient proto.UserClient
)
