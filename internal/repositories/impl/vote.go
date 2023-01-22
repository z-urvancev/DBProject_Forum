package impl

import (
	"DBProject/internal/models"
	"DBProject/pkg/queries"
	"github.com/jackc/pgx"
	_ "github.com/lib/pq"
)

type VoteRepositoryImpl struct {
	db *pgx.ConnPool
}

func NewVoteRepository(db *pgx.ConnPool) *VoteRepositoryImpl {
	return &VoteRepositoryImpl{db: db}
}

func (vr *VoteRepositoryImpl) VoteForThread(id int64, vote *models.Vote) error {
	_, err := vr.db.Exec(queries.Vote, vote.Nickname, id, vote.Voice)
	return err
}
