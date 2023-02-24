package model

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"toktik/constant/config"
	"toktik/logging"
	"toktik/rpc"
)

// User 用户表 /*
type User struct {
	Model                   // 基础模型
	Username        string  `gorm:"not null;unique;size: 32;index"` // 用户名
	Password        *string `gorm:"not null;size: 32"`              // 密码
	Avatar          *string // 用户头像
	BackgroundImage *string // 背景图片
	Signature       *string // 个人简介

	updated bool
}

func (u *User) IsUpdated() bool {
	return u.updated
}

func (u *User) isEmail() bool {
	parts := strings.Split(u.Username, "@")
	if len(parts) != 2 {
		return false
	}

	userPart := parts[0]
	if len(userPart) == 0 {
		return false
	}

	domainPart := parts[1]
	if len(domainPart) < 3 {
		return false
	}
	if !strings.Contains(domainPart, ".") {
		return false
	}
	if strings.HasPrefix(domainPart, ".") || strings.HasSuffix(domainPart, ".") {
		return false
	}

	return true
}

func getEmailMD5(email string) (md5String string) {
	lowerEmail := strings.ToLower(email)
	hasher := md5.New()
	hasher.Write([]byte(lowerEmail))
	md5Bytes := hasher.Sum(nil)
	md5String = hex.EncodeToString(md5Bytes)
	return
}

type unsplashResponse []struct {
	Urls struct {
		Regular string `json:"regular"`
	} `json:"urls"`
}

func getImageFromUnsplash(query string) (url string, err error) {
	unsplashUrl := fmt.Sprintf("https://api.unsplash.com/photos/random?query=%s&count=1", query)

	resp, err := rpc.HttpRequest("GET", unsplashUrl, nil, rpc.WithAuthorizationHeader("Client-ID "+config.EnvConfig.UNSPLASH_ACCESS_KEY))
	if err != nil {
		return "", err
	}

	if resp.Body != nil {
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				logging.Logger.Errorf("getImageFromUnsplash: %v", err)
			}
		}(resp.Body)
	}
	body, _ := io.ReadAll(resp.Body)

	var response unsplashResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", err
	}

	url = response[0].Urls.Regular

	if url == "" {
		return "", fmt.Errorf("getImageFromUnsplash: url is empty")
	}
	return
}

func getCravatarUrl(email string) string {
	return `https://cravatar.cn/avatar/` + getEmailMD5(email) + `?d=` + "identicon"
}

func (u *User) GetBackgroundImage() (url string) {
	if u.BackgroundImage != nil && *u.BackgroundImage != "" {
		return *u.BackgroundImage
	}

	defer func() {
		u.BackgroundImage = &url
		u.updated = true
	}()

	unsplashURL, err := getImageFromUnsplash(u.Username)
	if err != nil {
		catURL, err := getImageFromUnsplash("cat")
		if err != nil {
			return getCravatarUrl(u.Username)
		}
		return catURL
	}
	return unsplashURL
}

func (u *User) GetUserAvatar() (url string) {
	if u.Avatar != nil && *u.Avatar != "" {
		return *u.Avatar
	}

	defer func() {
		u.Avatar = &url
		u.updated = true
	}()

	if u.isEmail() {
		return getCravatarUrl(u.Username)
	}

	unsplashURL, err := getImageFromUnsplash(u.Username)
	if err != nil || unsplashURL == "" {
		catURL, err := getImageFromUnsplash("cat")
		if err != nil || catURL == "" {
			return getCravatarUrl(u.Username)
		}
		return catURL
	}

	return unsplashURL
}

func (u *User) GetSignature() (signature string) {
	if u.Signature != nil &&
		*u.Signature != "" &&
		*u.Signature != u.Username /* For compatibility */ {
		return *u.Signature
	}

	defer func() {
		u.Signature = &signature
		u.updated = true
	}()

	resp, err := rpc.HttpRequest("GET", "https://v1.hitokoto.cn/?encode=text", nil)
	if err != nil {
		logging.Logger.Errorf("GetSignature: %v", err)
		signature = u.Username
		return
	}

	if resp.StatusCode != http.StatusOK {
		logging.Logger.Errorf("GetSignature: %v", err)
		signature = u.Username
		return
	}

	if resp.Body != nil {
		defer func(Body io.ReadCloser) {
			err := Body.Close()
			if err != nil {
				logging.Logger.Errorf("GetSignature: %v", err)
			}
		}(resp.Body)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logging.Logger.Errorf("GetSignature: %v", err)
		signature = u.Username
		return
	}

	signature = string(body)

	return
}
