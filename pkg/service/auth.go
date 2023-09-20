package service

import (
	"crypto/sha1"
	"fmt"
	"task_manager"
	"task_manager/pkg/repository"
)

const salt = "fjadksf3ir20iao"

type AuthService struct {
	repo repository.Authorization
}

func NewAuthService(repo repository.Authorization) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) CreateUser(user task_manager.User) (int, error) {
	user.Password = s.generatePasswordHash(user.Password)
	return s.repo.CreateUser(user)
}

func (s *AuthService) GetUserId(username, password string) (int, error) {
	user, err := s.repo.GetUser(username, s.generatePasswordHash(password))
	return user.Id, err
}

func (s *AuthService) generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}
