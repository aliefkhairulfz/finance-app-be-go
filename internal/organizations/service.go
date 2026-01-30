package organizations

import "finance-app/internal/repository"

type Service struct {
	repository *repository.Queries
}

type ServiceHandler interface{}
