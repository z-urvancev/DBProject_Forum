package impl

import (
	"DBProject/internal/models"
	"DBProject/pkg/handlerows"
	"DBProject/pkg/queries"
	"github.com/jackc/pgx"
	_ "github.com/lib/pq"
)

type UserRepositoryImpl struct {
	db *pgx.ConnPool
}

func NewUserRepository(db *pgx.ConnPool) *UserRepositoryImpl {
	return &UserRepositoryImpl{db: db}
}

func (ur *UserRepositoryImpl) CreateUser(user *models.User) error {
	_, err := ur.db.Exec(queries.UserCreate, user.Nickname, user.Fullname, user.About, user.Email)
	return err
}

func (ur *UserRepositoryImpl) UpdateUser(user *models.User) error {
	return ur.db.QueryRow(queries.UserUpdate, user.Fullname, user.About, user.Email, user.Nickname).Scan(&user.Fullname, &user.About, &user.Email)
}

func (ur *UserRepositoryImpl) GetInfoAboutUser(nickname string) (*models.User, error) {
	user := new(models.User)
	err := ur.db.QueryRow(queries.UserGet, nickname).Scan(&user.Nickname, &user.Fullname, &user.About, &user.Email)
	return user, err
}

func (ur *UserRepositoryImpl) GetSimilarUsers(user *models.User) (*[]models.User, error) {
	resultRows, err := ur.db.Query(queries.UserGetSimilar, user.Nickname, user.Email)
	if err != nil {
		return nil, err
	}
	defer resultRows.Close()
	return handlerows.User(resultRows)
}
