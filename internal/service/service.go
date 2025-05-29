package service

import (
	"context"
	"fmt"
	"test_task/internal/clients"
	"test_task/internal/dto"
	"test_task/pkg/logger"
)

type Repository interface {
	AddPeople(ctx context.Context, person *dto.Person) (string, error)
	DeletePerson(ctx context.Context, id string) error
	GetPeople(ctx context.Context, filter *dto.PersonFilter, limit, offset int) ([]dto.Person, error)
	UpdatePerson(ctx context.Context, person *dto.Person) error
}

type Service struct {
	clent *clients.ExternalAPI
	repo  Repository
}

func NewService(api *clients.ExternalAPI, repo Repository) *Service {
	return &Service{clent: api, repo: repo}
}

func (s *Service) AddPeople(ctx context.Context, person *dto.Person) (string, error) {
	logger.GetLoggerFromCtx(ctx).Info(ctx, "Start AddPeople service")

	logger.GetLoggerFromCtx(ctx).Debug(ctx, "Calling enrichment service (GetAll)")
	person, err := s.clent.GetAll(ctx, person)
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Debug(ctx, "Enrichment service failed: "+err.Error())
		return "", err
	}

	logger.GetLoggerFromCtx(ctx).Info(ctx, "Person enriched successfully, calling repository to insert")

	id, err := s.repo.AddPeople(ctx, person)
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Debug(ctx, "Repository failed to insert person: "+err.Error())
		return "", err
	}

	logger.GetLoggerFromCtx(ctx).Info(ctx, "Person added successfully with ID: "+id)

	return id, nil
}

func (s *Service) DeletePerson(ctx context.Context, id string) error {
	err := s.repo.DeletePerson(ctx, id)
	logger.GetLoggerFromCtx(ctx).Info(ctx, "Start DeletePerson service for ID: "+id)
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Debug(ctx, "Failed to delete person: "+err.Error())
		return err
	}

	logger.GetLoggerFromCtx(ctx).Info(ctx, "Person deleted successfully")
	return nil

}

func (s *Service) GetPeople(ctx context.Context, filter *dto.PersonFilter, limit, offset int) ([]dto.Person, error) {
	logger.GetLoggerFromCtx(ctx).Info(ctx, "Start GetPeople service")

	logger.GetLoggerFromCtx(ctx).Debug(ctx, fmt.Sprintf("Filter: %+v, Limit: %d, Offset: %d", filter, limit, offset))
	people, err := s.repo.GetPeople(ctx, filter, limit, offset)
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Debug(ctx, "Failed to get people from repository: "+err.Error())
		return nil, err
	}

	logger.GetLoggerFromCtx(ctx).Info(ctx, fmt.Sprintf("Successfully fetched %d people", len(people)))
	return people, nil
}

func (s *Service) UpdatePerson(ctx context.Context, person *dto.Person) error {

	logger.GetLoggerFromCtx(ctx).Info(ctx, "Start UpdatePerson service")

	logger.GetLoggerFromCtx(ctx).Debug(ctx, "Calling enrichment service (GetAll)")
	person, err := s.clent.GetAll(ctx, person)
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Debug(ctx, "Enrichment service failed: "+err.Error())
		return err
	}

	logger.GetLoggerFromCtx(ctx).Info(ctx, "Person enriched successfully, updating in repository")

	err = s.repo.UpdatePerson(ctx, person)

	if err != nil {
		logger.GetLoggerFromCtx(ctx).Debug(ctx, "Repository failed to update person: "+err.Error())
		return err
	}
	logger.GetLoggerFromCtx(ctx).Info(ctx, "Person updated successfully")
	return nil
}
