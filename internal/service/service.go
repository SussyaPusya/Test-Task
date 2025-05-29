package service

import (
	"context"
	"test_task/internal/clients"
	"test_task/internal/dto"
	"test_task/pkg/logger"
)

type Repository interface {
	AddPeople(ctx context.Context, person *dto.Person) (string, error)
	DeletePerson(ctx context.Context, id string) error
	GetPeople(ctx context.Context, filter *dto.PersonFilter, limit, offset int) ([]dto.Person, error)
}

type Service struct {
	clent *clients.ExternalAPI
	repo  Repository
}

func NewService(api *clients.ExternalAPI, repo Repository) *Service {
	return &Service{clent: api, repo: repo}
}

func (s *Service) AddPeople(ctx context.Context, person *dto.Person) (string, error) {

	person, err := s.clent.GetAll(ctx, person)
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Debug(ctx, err.Error())
		return "", err
	}

	id, err := s.repo.AddPeople(ctx, person)
	if err != nil {
		return "", err
	}

	return id, nil
}

func (s *Service) DeletePerson(ctx context.Context, id string) error {
	err := s.repo.DeletePerson(ctx, id)

	if err != nil {
		return err
	}

	return nil

}

func (s *Service) GetPeople(ctx context.Context, filter *dto.PersonFilter, limit, offset int) ([]dto.Person, error) {
	return s.repo.GetPeople(ctx, filter, limit, offset)
}
