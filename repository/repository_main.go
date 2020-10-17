package repository

import (
	"fmt"
	"github.com/google/logger"
	"github.com/jackc/pgx"
)

type Repository struct {
	DbConnections int
	pool          *pgx.ConnPool
	UsersRepo     *UsersRepo
	StudentsRepo  *StudentsRepo
	MastersRepo   *MastersRepo
	ThemesRepo    *ThemesRepo
	LanguagesRepo *LanguagesRepo
	VideosRepo    *VideosRepo
	AvatarsRepo   *AvatarsRepo
}

func (repository *Repository) Init(config pgx.ConnConfig) error {
	var err error
	repository.pool, err = pgx.NewConnPool(pgx.ConnPoolConfig{
		ConnConfig:     config,
		MaxConnections: repository.DbConnections,
	})
	if err != nil {
		return err
	}
	err = repository.dropTables()
	if err != nil {
		return err
	}
	err = repository.createTables()
	if err != nil {
		return err
	}
	err = repository.fillTables()
	if err != nil {
		return err
	}
	repository.StudentsRepo = &StudentsRepo{repository}
	repository.MastersRepo = &MastersRepo{repository}
	repository.UsersRepo = &UsersRepo{repository}
	repository.LanguagesRepo = &LanguagesRepo{repository}
	repository.ThemesRepo = &ThemesRepo{repository}
	repository.VideosRepo = &VideosRepo{repository}
	repository.AvatarsRepo = &AvatarsRepo{repository}
	return nil
}

// relation style tables creation

func (repository *Repository) createTables() error {
	_, err := repository.pool.Exec(TABLES_CREATION)
	if err != nil {
		return err
	}
	return nil
}

func (repository *Repository) dropTables() error {
	_, err := repository.pool.Exec(TABLES_DROPPING)
	if err != nil {
		return err
	}
	return nil
}

func (repository *Repository) fillTables() error {
	_, err := repository.pool.Exec(TABLES_FILLING)
	if err != nil {
		return err
	}
	return nil
}

func (repository *Repository) GetPool() *pgx.ConnPool {
	return repository.pool
}

func (repository *Repository) startTransaction() (*pgx.Tx, error) {
	db := repository.GetPool()
	transaction, err := db.Begin()
	if err != nil {
		dbError := fmt.Errorf("failed to start transaction: %v", err.Error())
		logger.Errorf(dbError.Error())
		return transaction, dbError
	}
	return transaction, err
}

func (repository *Repository) commitTransaction(transaction *pgx.Tx) error {
	err := transaction.Commit()
	if err != nil {
		dbError := fmt.Errorf("error commit: %v", err.Error())
		logger.Errorf(dbError.Error())
		return dbError
	}
	return nil
}

func (repository *Repository) rollbackTransaction(transaction *pgx.Tx) error {
	err := transaction.Rollback()
	if err != nil {
		dbError := fmt.Errorf("failed to rollback: %v", err.Error())
		logger.Errorf(dbError.Error())
		return err
	}
	return err
}
