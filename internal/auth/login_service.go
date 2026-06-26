package auth

import (
	"context"
	"errors"
	"fmt"

	"github.com/vgadzh/telegram-message-collector/internal/domain/user"
	"github.com/vgadzh/telegram-message-collector/internal/observability/logctx"
	"github.com/vgadzh/telegram-message-collector/internal/repo"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type LoginService struct {
	users *repo.UserRepo
}

func NewLoginService(users *repo.UserRepo) *LoginService {
	return &LoginService{
		users: users,
	}
}

func (s *LoginService) Authenticate(ctx context.Context, login, password string) (*user.User, error) {
	logger := logctx.FromContext(ctx)

	u, err := s.users.GetByLogin(ctx, login)
	if err != nil {
		logger.Debug("users.GetByLogin error", zap.Error(err))
		if errors.Is(err, repo.ErrUserNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, fmt.Errorf("get user by login: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)); err != nil {
		logger.Debug("bcrypt error", zap.Error(err))
		return nil, ErrInvalidCredentials
	}

	return u, nil
}
