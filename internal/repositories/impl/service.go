package impl

import (
	"DBProject/internal/models"
	"DBProject/pkg/queries"
	"github.com/jackc/pgx"
	_ "github.com/lib/pq"
)

type ServiceRepositoryImpl struct {
	db *pgx.ConnPool
}

func NewServiceRepository(db *pgx.ConnPool) *ServiceRepositoryImpl {
	return &ServiceRepositoryImpl{db: db}
}

func (sr *ServiceRepositoryImpl) ClearService() (err error) {
	_, err = sr.db.Exec(queries.ServiceClear)
	return
}

func (sr *ServiceRepositoryImpl) GetService() (status *models.Status, err error) {
	status = &models.Status{}
	err = sr.db.QueryRow(queries.ServiceGet).
		Scan(
			&status.User,
			&status.Forum,
			&status.Thread,
			&status.Post)
	return
}
