package bootstrap

import (
	"context"

	"github.com/vgadzh/telegram-message-collector/internal/config"
	"github.com/vgadzh/telegram-message-collector/internal/domain/user"
	"github.com/vgadzh/telegram-message-collector/internal/repo"
)

func BootstrapAdmin(ctx context.Context, repo *repo.UserRepo, cfg config.Admin) error {
	exists, err := repo.ExistsByLogin(ctx, cfg.Login)
	if err != nil {
		return err
	}

	if exists {
		return nil
	}

	_, err = repo.Create(ctx, &user.User{
		Login:        cfg.Login,
		PasswordHash: cfg.PasswordHash,
		Role:         user.RoleAdmin,
	})
	return err
}
