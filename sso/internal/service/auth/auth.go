package auth

import (
	"context"
	"crypto/sha1"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"log/slog"
	"sso/internal/domain/model"
	"sso/internal/lib/jwt"
	"sso/internal/storage"
	"strings"
	"time"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrPasswordIncorrect  = errors.New("password incorrect")
	ErrUserExists         = errors.New("user already exists")
)

type Auth struct {
	log      *slog.Logger
	storage  Storage
	salt     string
	tokenTTL time.Duration
}

type Storage interface {
	SaveUser(
		ctx context.Context,
		username string,
		passwordHash string,
	) (userId uuid.UUID, err error)
	GetUser(
		ctx context.Context,
		username string,
	) (model.User, error)
	GetApp(ctx context.Context, appID int) (model.App, error)
	IsAdmin(ctx context.Context, userID uuid.UUID) (bool, error)
}

func New(log *slog.Logger, storage Storage, tokenTTL time.Duration) *Auth {
	return &Auth{
		log:      log,
		storage:  storage,
		tokenTTL: tokenTTL,
	}
}

func (a *Auth) Login(
	ctx context.Context,
	username string,
	password string,
	appID int,
) (string, error) {
	const op = "server.Login"
	log := a.log.With(slog.String("op", op), slog.String("username", username))

	user, err := a.storage.GetUser(ctx, username)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			a.log.Warn("User not found")
			return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}
		return "", fmt.Errorf("%s: %w", op, err)
	}

	if err := comparePasswordWithHash(user.PasswordHash, password, a.salt); err != nil {
		if errors.Is(err, ErrPasswordIncorrect) {
			a.log.Warn("Password incorrect")
			return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}
		return "", fmt.Errorf("%s: %w", op, err)
	}

	app, err := a.storage.GetApp(ctx, appID)
	if err != nil {
		if errors.Is(err, storage.ErrAppNotFound) {
			a.log.Warn("App not found")
			return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}
		return "", fmt.Errorf("%s: %w", op, err)
	}

	log.Info("User logged successfully")

	token, err := jwt.NewToken(user, app, a.tokenTTL)
	if err != nil {
		a.log.Error("Failed to create token")
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return token, nil
}

func (a *Auth) Register(ctx context.Context, username string, password string) (uuid.UUID, error) {
	const op = "server.Register"
	log := a.log.With(slog.String("op", op), slog.String("username", username))

	passwordHash := generatePasswordHash(password, a.salt)
	id, err := a.storage.SaveUser(ctx, username, passwordHash)
	if err != nil {
		if errors.Is(err, storage.ErrUserExists) {
			return uuid.Nil, ErrUserExists
		}
		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}
	log.Info("User registered")
	return id, nil
}

func (a *Auth) IsAdmin(ctx context.Context, userID uuid.UUID) (bool, error) {
	const op = "server.IsAdmin"
	log := a.log.With(slog.String("op", op), slog.String("user_id", fmt.Sprint(userID)))

	isAdmin, err := a.storage.IsAdmin(ctx, userID)
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			return false, err
		}
		return false, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("Checked if user is admin", slog.Bool("is_admin", isAdmin))
	return isAdmin, nil
}

func comparePasswordWithHash(hash, password, salt string) error {
	if strings.Compare(hash, generatePasswordHash(password, salt)) == 0 {
		return nil
	} else {
		return ErrPasswordIncorrect
	}
}

func generatePasswordHash(password, salt string) string {
	hash := sha1.New()
	hash.Write([]byte(password))
	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}
