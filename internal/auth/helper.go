package auth

import (
	"context"
	"finance-app/internal/repository"
	"finance-app/lib"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type CSessionParams struct {
	UserID    string
	IpAddress string
	UserAgent string
}

type CreateSessionRowWithCSRF struct {
	CreateSessionRow *repository.CreateSessionRow
	CsrfToken        string
}

func CreateSession(ctx context.Context, s *Service, payload CSessionParams) (*CreateSessionRowWithCSRF, error) {
	tExpiresAt := time.Now().Add(time.Hour * 1)
	t, err := lib.GenerateSessionToken(32)
	if err != nil {
		return nil, err
	}

	csrfT, err := lib.GenerateCsrfToken(32)
	if err != nil {
		return nil, err
	}

	cSessionPayload := repository.CreateSessionParams{
		Token: t,
		ExpiresAt: pgtype.Timestamptz{
			Time:             tExpiresAt,
			InfinityModifier: pgtype.Finite,
			Valid:            true,
		},
		UserID: payload.UserID,
		IpAddress: pgtype.Text{
			String: payload.IpAddress,
			Valid:  true,
		},
		UserAgent: pgtype.Text{
			String: payload.UserAgent,
			Valid:  true,
		},
	}

	cSession, err := s.repository.CreateSession(ctx, cSessionPayload)
	if err != nil {
		return nil, err
	}

	cSessionResponsePayload := &CreateSessionRowWithCSRF{
		CreateSessionRow: &cSession,
		CsrfToken:        csrfT,
	}

	return cSessionResponsePayload, nil
}
