package usecase

import (
	"context"
	"strings"

	"github.com/ziadrahmatullah/minimarket-app/apperror"
	"github.com/ziadrahmatullah/minimarket-app/appjwt"
	"github.com/ziadrahmatullah/minimarket-app/entity"
	"github.com/ziadrahmatullah/minimarket-app/hasher"
	"github.com/ziadrahmatullah/minimarket-app/repository"
	"github.com/ziadrahmatullah/minimarket-app/valueobject"
)

type AuthUsecase interface {
	Register(ctx context.Context, data *entity.User) error
	Login(ctx context.Context, user *entity.User) (string, error)
}

type authUsecase struct {
	userRepo repository.UserRepository
	jwt      appjwt.Jwt
	hash     hasher.Hasher
}

func NewAuthUsecase(
	userRepo repository.UserRepository,
	jwt appjwt.Jwt,
	hash hasher.Hasher,
) AuthUsecase {
	return &authUsecase{
		userRepo: userRepo,
		jwt:      jwt,
		hash:     hash,
	}
}

func (u *authUsecase) Register(ctx context.Context, user *entity.User) error {
	emailQuery := valueobject.NewQuery().Condition("email", valueobject.Equal, user.Email)
	fetchedUser, err := u.userRepo.FindOne(ctx, emailQuery)
	if err != nil {
		return err
	}
	if fetchedUser != nil {
		return apperror.NewResourceAlreadyExistError("user", "email", user.Email)
	}
	hashPass, err := u.hash.Hash(user.Password)
	if err != nil {
		return err
	}
	user.Password = string(hashPass)
	parts := strings.Split(user.Email, "@")
	user.Username = parts[0]
	user.Role = entity.RoleUser
	_, err = u.userRepo.Create(ctx, user)
	if err != nil {
		return err
	}
	return nil
}

func (u *authUsecase) Login(ctx context.Context, user *entity.User) (string, error) {
	emailQuery := valueobject.NewQuery().Condition("email", valueobject.Equal, user.Email)
	fetchedUser, err := u.userRepo.FindOne(ctx, emailQuery)
	if err != nil {
		return "", err
	}
	if fetchedUser == nil {
		return "", apperror.NewResourceNotFoundError("user", "email", user.Email)
	}
	if !(u.hash.Compare(fetchedUser.Password, user.Password)) {
		return "", apperror.NewInvalidCredentialsError()
	}
	token, err := u.jwt.GenerateToken(fetchedUser)
	if err != nil {
		return "", err
	}
	return token, nil
}
