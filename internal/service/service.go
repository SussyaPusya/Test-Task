package service

import (
	"context"
	"test_task/internal/dto"
)

type Repository interface {
	AddPeople(ctx context.Context, person *dto.Person) error
}

type Service struct {
	repo Repository
}

func NewRepository(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) AddPeople(ctx context.Context, person *dto.Person) error {

	err := s.repo.AddPeople(ctx, person)
	if err != nil {

	}
	return nil
}
