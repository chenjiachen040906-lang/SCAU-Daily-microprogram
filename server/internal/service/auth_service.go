package service

import (
	"errors"
	"time"

	"scau-daily/internal/config"
	"scau-daily/internal/dto"
	"scau-daily/internal/jwxt"
	"scau-daily/internal/middleware"
	"scau-daily/internal/model"
	"scau-daily/internal/repository"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// AuthService handles authentication and user session management.
type AuthService struct{}

// Login authenticates against JWXT, finds or creates the user, and returns tokens.
func (s *AuthService) Login(studentID, password string) (*dto.AuthResponse, error) {
	// Authenticate via JWXT
	client := jwxt.NewClient(config.AppConfig.JWXTBaseURL)
	result, err := client.Login(studentID, password)
	if err != nil {
		return nil, errors.New("教务系统连接失败: " + err.Error())
	}
	if !result.Success {
		return nil, errors.New(result.Message)
	}

	// Find or create user in database
	user, err := repository.FindUserByStudentID(studentID)
	if err != nil {
		// User does not exist, create a new record
		user = &model.User{
			StudentID: studentID,
		}
		user, err = repository.CreateUser(user)
		if err != nil {
			return nil, errors.New("创建用户失败: " + err.Error())
		}
	}

	// Hash and store the JWXT password for offline verification
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err == nil {
		_ = repository.UpdateUserPassword(user.ID, string(hash))
	}

	// Update last login timestamp
	now := time.Now()
	_ = repository.UpdateUserLastLogin(user.ID, now)

	// Generate token pair
	accessToken, err := middleware.GenerateAccessToken(user.ID, user.StudentID)
	if err != nil {
		return nil, errors.New("生成访问令牌失败: " + err.Error())
	}
	refreshToken, err := middleware.GenerateRefreshToken(user.ID, user.StudentID)
	if err != nil {
		return nil, errors.New("生成刷新令牌失败: " + err.Error())
	}

	return &dto.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User: dto.UserDTO{
			ID:         user.ID,
			StudentID:  user.StudentID,
			Name:       user.Name,
			Department: user.Department,
			Major:      user.Major,
			Grade:      user.Grade,
			AvatarURL:  user.AvatarURL,
		},
	}, nil
}

// WxLogin is a placeholder for WeChat mini-program login.
// Returns an empty AuthResponse, a need_bind flag (true), and no error.
// The need_bind flag indicates the WeChat account must be bound to a student ID.
func (s *AuthService) WxLogin(code string) (*dto.AuthResponse, bool, error) {
	// TODO: implement WeChat code-to-openid exchange via wx API
	// For now, always require binding
	return &dto.AuthResponse{}, true, nil
}

// BindStudent verifies credentials via JWXT and binds a student ID to an existing user.
func (s *AuthService) BindStudent(userID uuid.UUID, studentID, password string) (*dto.AuthResponse, error) {
	// Verify credentials through JWXT
	client := jwxt.NewClient(config.AppConfig.JWXTBaseURL)
	result, err := client.Login(studentID, password)
	if err != nil {
		return nil, errors.New("教务系统连接失败: " + err.Error())
	}
	if !result.Success {
		return nil, errors.New(result.Message)
	}

	// Fetch the user
	user, err := repository.FindUserByID(userID)
	if err != nil {
		return nil, errors.New("用户不存在")
	}

	// Hash the password
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("密码加密失败: " + err.Error())
	}

	// Update user with student info and hashed password
	user.StudentID = studentID
	user.PasswordHash = string(hash)
	if err := repository.UpdateUser(user); err != nil {
		return nil, errors.New("更新用户信息失败: " + err.Error())
	}

	// Generate fresh token pair
	accessToken, err := middleware.GenerateAccessToken(user.ID, user.StudentID)
	if err != nil {
		return nil, errors.New("生成访问令牌失败: " + err.Error())
	}
	refreshToken, err := middleware.GenerateRefreshToken(user.ID, user.StudentID)
	if err != nil {
		return nil, errors.New("生成刷新令牌失败: " + err.Error())
	}

	return &dto.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User: dto.UserDTO{
			ID:         user.ID,
			StudentID:  user.StudentID,
			Name:       user.Name,
			Department: user.Department,
			Major:      user.Major,
			Grade:      user.Grade,
			AvatarURL:  user.AvatarURL,
		},
	}, nil
}

// RefreshToken validates a refresh token and issues a new access+refresh token pair.
func (s *AuthService) RefreshToken(refreshToken string) (*dto.AuthResponse, error) {
	claims, err := middleware.ParseRefreshToken(refreshToken)
	if err != nil {
		return nil, errors.New("无效的刷新令牌")
	}

	user, err := repository.FindUserByID(claims.UserID)
	if err != nil {
		return nil, errors.New("用户不存在")
	}

	accessToken, err := middleware.GenerateAccessToken(user.ID, user.StudentID)
	if err != nil {
		return nil, errors.New("生成访问令牌失败: " + err.Error())
	}
	newRefreshToken, err := middleware.GenerateRefreshToken(user.ID, user.StudentID)
	if err != nil {
		return nil, errors.New("生成刷新令牌失败: " + err.Error())
	}

	return &dto.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: newRefreshToken,
		User: dto.UserDTO{
			ID:         user.ID,
			StudentID:  user.StudentID,
			Name:       user.Name,
			Department: user.Department,
			Major:      user.Major,
			Grade:      user.Grade,
			AvatarURL:  user.AvatarURL,
		},
	}, nil
}

// GetMe returns the current user's profile information.
func (s *AuthService) GetMe(userID uuid.UUID) (*dto.UserDTO, error) {
	user, err := repository.FindUserByID(userID)
	if err != nil {
		return nil, errors.New("用户不存在")
	}

	return &dto.UserDTO{
		ID:         user.ID,
		StudentID:  user.StudentID,
		Name:       user.Name,
		Department: user.Department,
		Major:      user.Major,
		Grade:      user.Grade,
		AvatarURL:  user.AvatarURL,
	}, nil
}
