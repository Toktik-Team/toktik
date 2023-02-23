package main

import (
	"context"
	"toktik/constant/biz"
	"toktik/kitex_gen/douyin/auth"
	"toktik/logging"
	"toktik/repo"
	commonModel "toktik/repo/model"

	"github.com/segmentio/ksuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// AuthServiceImpl implements the last service interface defined in the IDL.
type AuthServiceImpl struct{}

// Authenticate implements the AuthServiceImpl interface.
func (s *AuthServiceImpl) Authenticate(ctx context.Context, req *auth.AuthenticateRequest) (resp *auth.AuthenticateResponse, err error) {
	logger := logging.Logger.WithFields(logrus.Fields{
		"method": "Authenticate",
		"token":  req.Token,
	})
	logger.Debugf("Process start")
	if req == nil {
		logger.Warningf("request is nil")
		return &auth.AuthenticateResponse{
			StatusCode: biz.RequestIsNil,
			StatusMsg:  biz.InternalServerErrorStatusMsg,
		}, nil
	}
	userToken := repo.Q.UserToken
	first, err := userToken.WithContext(ctx).Where(userToken.Token.Eq(req.Token)).First()
	if err != nil {
		logger.Warningf(biz.TokenNotFoundMessage)
		return &auth.AuthenticateResponse{
			StatusCode: biz.TokenNotFound,
			StatusMsg:  biz.TokenNotFoundMessage,
		}, nil
	}
	resp = &auth.AuthenticateResponse{
		StatusCode: 0,
		StatusMsg:  "success",
		UserId:     first.UserID,
	}
	return
}

// Register implements the AuthServiceImpl interface.
func (s *AuthServiceImpl) Register(ctx context.Context, req *auth.RegisterRequest) (resp *auth.RegisterResponse, err error) {
	resp = &auth.RegisterResponse{}
	user := repo.Q.User
	userToken := repo.Q.UserToken
	dbUser, err := user.WithContext(ctx).Where(user.Username.Eq(req.Username)).Select().Find()
	if err != nil {
		resp = &auth.RegisterResponse{
			StatusCode: biz.ServiceNotAvailable,
			StatusMsg:  biz.InternalServerErrorStatusMsg,
		}
		return
	}
	if len(dbUser) > 0 {
		resp = &auth.RegisterResponse{
			StatusCode: biz.UserNameExist,
			StatusMsg:  "用户名已存在",
		}
		return
	}
	hashedPwd, err := HashPassword(req.Password)
	if err != nil {
		resp = &auth.RegisterResponse{
			StatusCode: biz.ServiceNotAvailable,
			StatusMsg:  biz.InternalServerErrorStatusMsg,
		}
		return
	}

	newUser := commonModel.User{
		Username: req.Username,
		Password: &hashedPwd,
	}
	err = user.WithContext(ctx).Save(
		&newUser)
	if err != nil {
		resp = &auth.RegisterResponse{
			StatusCode: biz.ServiceNotAvailable,
			StatusMsg:  biz.InternalServerErrorStatusMsg,
		}
		return
	}
	token := ksuid.New().String()
	err = userToken.WithContext(ctx).Save(&commonModel.UserToken{
		Token:    token,
		Username: newUser.Username,
		UserID:   newUser.ID,
	})

	if err != nil {
		resp = &auth.RegisterResponse{
			StatusCode: biz.ServiceNotAvailable,
			StatusMsg:  biz.InternalServerErrorStatusMsg,
		}
		return
	}
	resp.Token = token
	resp.UserId = newUser.ID
	resp.StatusCode = 0
	resp.StatusMsg = "success"
	return
}

// Login implements the AuthServiceImpl interface.
func (s *AuthServiceImpl) Login(ctx context.Context, req *auth.LoginRequest) (resp *auth.LoginResponse, err error) {
	resp = &auth.LoginResponse{}
	user := repo.Q.User
	userToken := repo.Q.UserToken
	dbUser, err := user.WithContext(ctx).Where(user.Username.Eq(req.Username)).Select().Find()
	if err != nil {
		resp = &auth.LoginResponse{
			StatusCode: biz.ServiceNotAvailable,
			StatusMsg:  biz.InternalServerErrorStatusMsg,
		}
		return
	}
	if len(dbUser) != 1 {
		resp = &auth.LoginResponse{
			StatusCode: biz.UserNotFound,
			StatusMsg:  "用户不存在",
		}
		return
	}
	if !CheckPasswordHash(req.Password, *dbUser[0].Password) {
		resp = &auth.LoginResponse{
			StatusCode: biz.PasswordIncorrect,
			StatusMsg:  "密码错误",
		}
		return
	}

	tokens, err := userToken.WithContext(ctx).Where(userToken.UserID.Eq(dbUser[0].ID)).Find()
	if len(tokens) == 1 {
		resp.Token = tokens[0].Token
		resp.UserId = dbUser[0].ID
		resp.StatusCode = 0
		resp.StatusMsg = "success"
		return
	}
	if err != nil {
		return nil, err
	}
	token := ksuid.New().String()
	err = userToken.WithContext(ctx).Save(&commonModel.UserToken{
		Token:    token,
		Username: dbUser[0].Username,
		UserID:   dbUser[0].ID,
	})

	if err != nil {
		resp = &auth.LoginResponse{
			StatusCode: biz.ServiceNotAvailable,
			StatusMsg:  biz.InternalServerErrorStatusMsg,
		}
		return
	}
	resp.Token = token
	resp.UserId = dbUser[0].ID
	resp.StatusCode = 0
	resp.StatusMsg = "success"
	return
}
