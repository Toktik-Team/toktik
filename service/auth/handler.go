package main

import (
	"context"
	"fmt"
	"github.com/segmentio/ksuid"
	"golang.org/x/crypto/bcrypt"
	"toktik/kitex_gen/douyin/auth"
	"toktik/repo"
	commonModel "toktik/repo/model"
	authModel "toktik/service/auth/model"
	"toktik/service/web/mw"
)

var StatusCodes = struct {
	UserNameExist       uint32
	ServiceNotAvailable uint32
	UserNotFound        uint32
	PasswordIncorrect   uint32
}{
	400101,
	503101,
	400102,
	401103,
}

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
	if req == nil {
		return nil, fmt.Errorf("failed")
	}
	userToken := repo.Q.UserToken
	first, err := userToken.WithContext(ctx).Where(userToken.Token.Eq(req.Token)).First()
	if err != nil {
		return nil, fmt.Errorf("failed")
	}
	resp = &auth.AuthenticateResponse{
		StatusCode: 0,
		StatusMsg:  string(mw.AUTH_RESULT_SUCCESS),
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
			StatusCode: StatusCodes.ServiceNotAvailable,
			StatusMsg:  "数据库查询失败",
		}
		return
	}
	if len(dbUser) > 0 {
		resp = &auth.RegisterResponse{
			StatusCode: StatusCodes.UserNameExist,
			StatusMsg:  "用户名已存在",
		}
		return
	}
	hashedPwd, err := HashPassword(req.Password)
	if err != nil {
		resp = &auth.RegisterResponse{
			StatusCode: StatusCodes.ServiceNotAvailable,
			StatusMsg:  "密码加密失败",
		}
		return
	}

	newUser := commonModel.User{
		Username:      req.Username,
		Password:      &hashedPwd,
		FollowCount:   0,
		FollowerCount: 0,
		Name:          req.Username,
		Role:          "0",
	}
	err = user.WithContext(ctx).Save(
		&newUser)
	if err != nil {
		resp = &auth.RegisterResponse{
			StatusCode: StatusCodes.ServiceNotAvailable,
			StatusMsg:  "数据库保存失败",
		}
		return
	}
	token := ksuid.New().String()
	err = userToken.WithContext(ctx).Save(&authModel.UserToken{
		Token:    token,
		Username: newUser.Username,
		UserID:   newUser.ID,
		Role:     newUser.Role,
	})

	if err != nil {
		resp = &auth.RegisterResponse{
			StatusCode: StatusCodes.ServiceNotAvailable,
			StatusMsg:  "数据库保存失败",
		}
		return
	}
	resp.Token = token
	resp.UserID = newUser.ID,
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
			StatusCode: StatusCodes.ServiceNotAvailable,
			StatusMsg:  "数据库查询失败",
		}
		return
	}
	if len(dbUser) != 1 {
		resp = &auth.LoginResponse{
			StatusCode: StatusCodes.UserNotFound,
			StatusMsg:  "用户不存在",
		}
		return
	}
	if !CheckPasswordHash(req.Password, *dbUser[0].Password) {
		resp = &auth.LoginResponse{
			StatusCode: StatusCodes.PasswordIncorrect,
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
	err = userToken.WithContext(ctx).Save(&authModel.UserToken{
		Token:    token,
		Username: dbUser[0].Username,
		UserID:   dbUser[0].ID,
		Role:     dbUser[0].Role,
	})

	if err != nil {
		resp = &auth.LoginResponse{
			StatusCode: StatusCodes.ServiceNotAvailable,
			StatusMsg:  "数据库保存失败",
		}
		return
	}
	resp.Token = token
	resp.UserId = dbUser[0].ID
	resp.StatusCode = 0
	resp.StatusMsg = "success"
	return
}
