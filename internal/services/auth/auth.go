package auth

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"ginhub.com/Aller101/sso/internal/domain/models"
	"ginhub.com/Aller101/sso/internal/lib/jwt"
	"ginhub.com/Aller101/sso/internal/storage"
	"golang.org/x/crypto/bcrypt"
)

// service
type Auth struct {
	log         *slog.Logger
	usrSaver    UserSaver
	usrProvider UserProvider
	appProvider AppProvider
	tokenTTL    time.Duration
}

type UserSaver interface {
	UserSave(ctx context.Context, email string, passHash []byte) (uid int64, err error)
}

type UserProvider interface {
	User(ctx context.Context, email string) (models.User, error)
	IsAdmin(ctx context.Context, userID int64) (bool, error)
}

type AppProvider interface {
	App(ctx context.Context, appID int) (models.App, error)
}

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
)

// New retn a new inst of the Auth serv
func New(log *slog.Logger, userSaver UserSaver, userProvider UserProvider, appProvider AppProvider, tokenTTL time.Duration) *Auth {
	return &Auth{
		usrSaver:    userSaver,
		usrProvider: userProvider,
		log:         log,
		appProvider: appProvider,
		tokenTTL:    tokenTTL,
	}
}

func (a *Auth) RegisterNewUser(ctx context.Context, email string, pass string) (int64, error) {
	const op = "auth.RegisterNewUser"

	// мы можем выводить настроенный в main *logger или
	// добавить еще параметров в тот же *logger, как это сделано ниже:
	log := a.log.With(
		slog.String("op", op),
		slog.String("email", email), //убрать
	)
	log.Info("registering user")

	passHash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to generate password hash")
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := a.usrSaver.UserSave(ctx, email, passHash)
	if err != nil {
		log.Error("failed to save user", slog.String("error", err.Error()))
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	return id, nil
}

func (a *Auth) Login(ctx context.Context, email string, password string, appID int) (string, error) {
	const op = "auth.Login"

	log := a.log.With(
		slog.String("op", op),
		slog.String("username", email), //убрать
	)
	a.log.Info("attempting to login user")

	user, err := a.usrProvider.User(ctx, email)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			a.log.Warn("user not found", slog.String("error", err.Error()))

			// Не стоит возвращать подробности об ошбкие: пользователь не узнает,
			// что он ввел неправильно(пароль или логин)
			// - ###нельзя будет составить базу тех, кто зарег. в приложении.

			// или, например, если попробовать восстановить акк через email,
			// приложение напишет: "письмо было отправленно на указанную почту"
			// - ###нельзя будет составить базу тех, кто зарег. в приложении.

			return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}

		// log.Error("failed to login", slog.String("error", err.Error())) -- надо ли?
		a.log.Error("failed to login", slog.String("error", err.Error()))
		return "", fmt.Errorf("%s: %w", op, err)
	}

	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		a.log.Info("invalid credentials", slog.String("error", err.Error()))

		return "", fmt.Errorf("%s: %w", op, err)
	}

	app, err := a.appProvider.App(ctx, appID)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	log.Info("user logged in successfully")

	token, err := jwt.NewToken(user, app, a.tokenTTL)
	if err != nil {
		a.log.Error("failed to generate token", slog.String("error", err.Error()))
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return token, nil

}
