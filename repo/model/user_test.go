package model

import (
	"github.com/cloudwego/hertz/pkg/common/test/assert"
	"strings"
	"testing"
)

func TestUser_GetUserAvatar(t *testing.T) {
	exampleAvatar := "https://example.image"
	// should get default avatar
	var userHasAvatar = User{
		Username: "nyanki",
		Avatar:   &exampleAvatar,
	}
	avatar := userHasAvatar.GetUserAvatar()
	assert.Assert(t, avatar == exampleAvatar)
	assert.Assert(t, !userHasAvatar.updated)

	// should get cravatar url when username is email
	var userHasEmail = User{
		Username: "example@email.com",
	}
	avatar = userHasEmail.GetUserAvatar()
	assert.Assert(t, avatar == getCravatarUrl(userHasEmail.Username))
	assert.Assert(t, userHasEmail.updated)

	// should get unsplash url when username is not email
	var userHasUsername = User{
		Username: "example",
	}
	avatar = userHasUsername.GetUserAvatar()
	assert.Assert(t, strings.HasPrefix(avatar, "https://"))
	assert.Assert(t, userHasUsername.updated)
}
