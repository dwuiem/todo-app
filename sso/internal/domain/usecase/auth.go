package usecase

import (
	"context"
	"crypto/sha1"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"log/slog"
	"sso/internal/adapter/repository"
	"sso/internal/domain/entity"
	"sso/internal/jwt"
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
	repo     Repository
	salt     string
	tokenTTL time.Duration
}

type Repository interface {
	SaveUser(
		ctx context.Context,
		username string,
		passwordHash string,
	) (userId uuid.UUID, err error)
	GetUser(
		ctx context.Context,
		username string,
	) (entity.User, error)
	GetApp(ctx context.Context, appID int) (entity.App, error)
	IsAdmin(ctx context.Context, userID uuid.UUID) (bool, error)
}

func New(log *slog.Logger, repository Repository, tokenTTL time.Duration) *Auth {
	return &Auth{
		log:      log,
		repo:     repository,
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

	user, err := a.repo.GetUser(ctx, username)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
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

	app, err := a.repo.GetApp(ctx, appID)
	if err != nil {
		if errors.Is(err, repository.ErrAppNotFound) {
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
	id, err := a.repo.SaveUser(ctx, username, passwordHash)
	if err != nil {
		if errors.Is(err, repository.ErrUserExists) {
			return uuid.Nil, ErrUserExists
		}
		return uuid.Nil, fmt.Errorf("%s: %w", op, err)
	}
	log.Info("User registered")
	return id, nil
}

func (a *Auth) IsAdmin(ctx context.Context, userID uuid.UUID) (bool, error) {
	const op = "server.IsAdmin"
	log := a.log.With(slog.String("op", op), slog.String("userID", fmt.Sprint(userID)))

	isAdmin, err := a.repo.IsAdmin(ctx, userID)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
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
