package repository

import (
	"context"
	"fmt"
	"test_task/internal/dto"
	"test_task/pkg/logger"

	sq "github.com/Masterminds/squirrel"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type Repository struct {
	pg *pgxpool.Pool
}

func NewRepository(pg *pgxpool.Pool) *Repository {

	return &Repository{pg: pg}

}

func (r *Repository) AddPeople(ctx context.Context, person *dto.Person) (string, error) {
	query := sq.Insert("people").
		Columns("name", "surname", "patronymic", "gender", "age", "nationality").
		Values(person.Name, person.Surname, person.Patronymic, person.Gender, person.Age, person.Nationality).
		Suffix("RETURNING id").
		PlaceholderFormat(sq.Dollar)
	if person.Patronymic == "" {
		query = sq.Insert("people").
			Columns("name", "surname", "gender", "age", "nationality").
			Values(person.Name, person.Surname, person.Gender, person.Age, person.Nationality).
			Suffix("RETURNING id").
			PlaceholderFormat(sq.Dollar)
	}

	sql, args, err := query.ToSql()
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Info(ctx, "Failed to build SQL:", zap.Error(err))

		return "", err
	}

	var id string
	err = r.pg.QueryRow(ctx, sql, args...).Scan(&id)
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Info(ctx, ":INSERT failed:", zap.Error(err))

		return "", err
	}

	return id, nil
}

func (r *Repository) DeletePerson(ctx context.Context, id string) error {
	query := sq.Delete("people").
		Where(sq.Eq{"id": id}).
		PlaceholderFormat(sq.Dollar)

	sqlStr, args, err := query.ToSql()
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Info(ctx, "Failed to build DELETE SQL:", zap.Error(err))
		return err
	}

	_, err = r.pg.Exec(ctx, sqlStr, args...)
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Info(ctx, "DELETE failed:", zap.Error(err))
		return err
	}

	return nil
}

func (r *Repository) GetPeople(ctx context.Context, filter *dto.PersonFilter, limit, offset int) ([]dto.Person, error) {

	fmt.Println(filter)

	filetrMap := detectFilter(filter)

	query := sq.Select("id", "name", "surname", "patronymic", "gender", "age", "nationality").
		From("people").
		PlaceholderFormat(sq.Dollar)

	if len(filetrMap) == 0 {
		query = sq.Select("id", "name", "surname", "patronymic", "gender", "age", "nationality").
			From("people").
			Where(filetrMap).
			PlaceholderFormat(sq.Dollar)

	}

	sqlStr, args, err := query.ToSql()
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Debug(ctx, "Failed to build SELECT SQL:", zap.Error(err))
		return nil, err
	}

	logger.GetLoggerFromCtx(ctx).Debug(ctx, "Executing query: "+sqlStr)

	rows, err := r.pg.Query(ctx, sqlStr, args...)
	if err != nil {
		logger.GetLoggerFromCtx(ctx).Debug(ctx, "SELECT query failed:", zap.Error(err))
		return nil, err
	}
	defer rows.Close()

	var result []dto.Person

	for rows.Next() {
		var p dto.Person
		err := rows.Scan(&p.ID, &p.Name, &p.Surname, &p.Patronymic, &p.Gender, &p.Age, &p.Nationality)
		if err != nil {
			logger.GetLoggerFromCtx(ctx).Debug(ctx, "Row scan failed:", zap.Error(err))
			continue
		}
		result = append(result, p)
	}

	logger.GetLoggerFromCtx(ctx).Debug(ctx, fmt.Sprintf("Found %d people", len(result)))

	return result, nil
}

func detectFilter(filter *dto.PersonFilter) sq.Eq {
	hash := sq.Eq{}
	if filter.Name != "" {
		fmt.Println("provberka", filter.Name)

		hash["name"] = filter.Name
	}
	if filter.Surname != "" {

		hash["surname"] = filter.Surname
	}
	if filter.Patronym != "" {

		hash["patronymic"] = filter.Patronym
	}
	if filter.Gender != "" {

		hash["gender"] = filter.Gender
	}
	if filter.Age != "" {

		hash["age"] = filter.Age
	}
	if filter.Country != "" {
		hash["nationality"] = filter.Country
	}

	return hash
}
