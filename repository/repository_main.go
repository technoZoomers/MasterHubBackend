package repository

import "github.com/jackc/pgx"

type Repository struct {
	pool *pgx.ConnPool
	UsersRepo *UsersRepo
	StudentsRepo *StudentsRepo
	MastersRepo *MastersRepo
	ThemesRepo *ThemesRepo
	LanguagesRepo *LanguagesRepo
	VideosRepo *VideosRepo
	AvatarsRepo *AvatarsRepo
}

var repo Repository

const dbConnections = 20

func Init(config pgx.ConnConfig) error {
	var err error
	repo.pool, err = pgx.NewConnPool(pgx.ConnPoolConfig{
		ConnConfig:     config,
		MaxConnections: dbConnections,
	})
	if err != nil {
		return err
	}
	//err = repo.dropTables()
	//if err != nil {
	//	return err
	//}
	//err = repo.createTables()
	//if err != nil {
	//	return err
	//}
	//err = repo.fillTables()
	//if err != nil {
	//	return err
	//}
	repo.StudentsRepo = &StudentsRepo{}
	repo.MastersRepo = &MastersRepo{}
	repo.UsersRepo = &UsersRepo{}
	repo.LanguagesRepo = &LanguagesRepo{}
	repo.ThemesRepo = &ThemesRepo{}
	repo.VideosRepo = &VideosRepo{}
	repo.AvatarsRepo = &AvatarsRepo{}
	return nil
}

// relation style tables creation

func (repo *Repository) createTables() error {
	_, err := repo.pool.Exec(TABLES_CREATION)
	if err != nil {
		return err
	}
	return nil
}

func (repo *Repository) dropTables() error {
	_, err := repo.pool.Exec(TABLES_DROPPING)
	if err != nil {
		return err
	}
	return nil
}

func (repo *Repository) fillTables() error {
	_, err := repo.pool.Exec(TABLES_FILLING)
	if err != nil {
		return err
	}
	return nil
}

func getPool() *pgx.ConnPool {
	return repo.pool
}

func GetUsersRepo() UsersRepoI {
	return repo.UsersRepo
}

func GetMastersRepo() MastersRepoI {
	return repo.MastersRepo
}

func GetStudentsRepo() StudentsRepoI {
	return repo.StudentsRepo
}

func GetThemesRepo() ThemesRepoI {
	return repo.ThemesRepo
}

func GetLanguagesRepo() LanguagesRepoI {
	return repo.LanguagesRepo
}

func GetAvatarsRepo() AvatarsRepoI {
	return repo.AvatarsRepo
}

func GetVideosRepo() VideosRepoI {
	return repo.VideosRepo
}