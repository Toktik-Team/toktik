package main

import (
	"context"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"toktik/kitex_gen/douyin/auth"
	gen "toktik/repo"
	"toktik/service/web/mw"
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
	userToken := gen.Q.UserToken
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
	// TODO: Your code here...
	return
}

// Login implements the AuthServiceImpl interface.
func (s *AuthServiceImpl) Login(ctx context.Context, req *auth.LoginRequest) (resp *auth.LoginResponse, err error) {
	// TODO: Your code here...
	return
}
