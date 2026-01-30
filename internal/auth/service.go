package auth

import (
	"context"
	"errors"
	"finance-app/internal/repository"
	"finance-app/lib"
	"finance-app/utils"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type SignUpParams struct {
	Name       string `json:"name"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	ProviderId string `json:"provider_id"`
}

type SignInParams struct {
	Email      string `json:"email"`
	Password   string `json:"password"`
	ProviderId string `json:"provider_id"`
}

type SignInUserMeta struct {
	UserAgent string
	IpAddress string
}

type Service struct {
	repository *repository.Queries
}

type ServiceHandler interface {
	SignUp(ctx context.Context, args SignUpParams) (*repository.CreateAccountRow, error)
	SignIn(ctx context.Context, args SignInParams, meta SignInUserMeta) (*CreateSessionRowWithCSRF, error)
	Me(ctx context.Context, token string) (*repository.FindUserByIdRow, error)
	Csrf(ctx context.Context) (string, error)
}

func NewService(db *repository.Queries) ServiceHandler {
	return &Service{
		repository: db,
	}
}

func (s *Service) SignUp(ctx context.Context, args SignUpParams) (*repository.CreateAccountRow, error) {
	hash, err := lib.HashPassword(args.Password)
	if err != nil {
		return nil, err
	}
	fmt.Println(args.Email)
	fmt.Println(args.Password)

	// FIND USER
	fUser, err := s.repository.FindUserByEmail(ctx, args.Email)
	if err != nil {
		// CREATE NEW USER AND ACCOUNT IF ROWS NOT FOUND
		if errors.Is(err, pgx.ErrNoRows) {
			uPayload := repository.CreateUserParams{
				Name:  args.Name,
				Email: args.Email,
			}

			nUser, err := s.repository.CreateUser(ctx, uPayload)
			if err != nil {
				return nil, err
			}

			accPayload := repository.CreateAccountParams{
				UserID:     nUser.ID,
				Password:   hash,
				ProviderID: args.ProviderId,
			}

			nAcc, err := s.repository.CreateAccount(ctx, accPayload)
			if err != nil {
				return nil, err
			}

			return &nAcc, nil
		}

		// THROW ANYWAY
		return nil, utils.InternalServerError
	}

	// FIND ACCOUNT
	fAccPayload := repository.FindAccountByUserIdAndProviderIDParams{
		UserID:     fUser.ID,
		ProviderID: args.ProviderId,
	}

	_, err = s.repository.FindAccountByUserIdAndProviderID(ctx, fAccPayload)
	if err != nil {
		// CREATE NEW ACCOUNT WITH SAME ID
		if errors.Is(err, pgx.ErrNoRows) {
			accPayload := repository.CreateAccountParams{
				UserID:     fUser.ID,
				Password:   hash,
				ProviderID: args.ProviderId,
			}

			nAcc, err := s.repository.CreateAccount(ctx, accPayload)
			if err != nil {
				return nil, err
			}

			return &nAcc, nil
		}

		// THROW ANYWAY
		return nil, utils.InternalServerError
	}

	// ACCOUNT CONFLICT
	return nil, utils.ErrorAccountConflict
}

func (s *Service) SignIn(ctx context.Context, args SignInParams, meta SignInUserMeta) (*CreateSessionRowWithCSRF, error) {
	// FIND USER
	fUser, err := s.repository.FindUserByEmail(ctx, args.Email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, utils.ErrorEmailNotFound
		}
		return nil, err
	}

	// SESSION PARAMS FUNCTION
	cSessionPayload := CSessionParams{
		UserID:    fUser.ID,
		IpAddress: meta.IpAddress,
		UserAgent: meta.UserAgent,
	}

	// FIND ACCOUNT
	fAccPayload := repository.FindAccountByUserIdAndProviderIDParams{
		UserID:     fUser.ID,
		ProviderID: args.ProviderId,
	}

	fAcc, err := s.repository.FindAccountByUserIdAndProviderID(ctx, fAccPayload)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, utils.ErrorAccountNotFound
		}
		return nil, err
	}

	// COMPARE PASSWORD
	err = lib.ComparePassword(fAcc.Password, args.Password)
	if err != nil {
		return nil, utils.ErrorAccountWrongPassword
	}

	// FIND SESSION
	_, err = s.repository.FindSessionByUserId(ctx, fUser.ID)
	if err != nil {
		// CREATE NEW SESSION IF NOT FOUND
		if errors.Is(err, pgx.ErrNoRows) {
			cSession, err := CreateSession(ctx, s, cSessionPayload)
			if err != nil {
				return nil, err
			}
			return cSession, nil
		}
		return nil, err
	}

	// DROP/DELETE SESSION
	err = s.repository.DeleteSessionByUserId(ctx, fAcc.UserID)
	if err != nil {
		return nil, err
	}

	// CREATE NEW SESSION
	cSession, err := CreateSession(ctx, s, cSessionPayload)
	if err != nil {
		return nil, err
	}

	return cSession, nil
}

func (s *Service) Me(ctx context.Context, token string) (*repository.FindUserByIdRow, error) {
	fSession, err := s.repository.FindSessionByToken(ctx, token)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, utils.ErrorSessionNotFound
		}
		return nil, err
	}

	fUser, err := s.repository.FindUserById(ctx, fSession.UserID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, utils.ErrorNoUserFound
		}
	}

	return &fUser, nil
}

func (s *Service) Csrf(ctx context.Context) (string, error) {
	csrfToken, err := lib.GenerateCsrfToken(32)
	if err != nil {
		return "", err
	}

	return csrfToken, nil
}
