package impl

import (
	"DBProject/internal/models"
	"DBProject/pkg/queries"
	"time"

	"github.com/jackc/pgx"
	_ "github.com/lib/pq"
)

type PostRepositoryImpl struct {
	db *pgx.ConnPool
}

func NewPostRepository(db *pgx.ConnPool) *PostRepositoryImpl {
	return &PostRepositoryImpl{db: db}
}

func (pr *PostRepositoryImpl) GetPost(id int64) (post *models.Post, err error) {
	post = &models.Post{}
	timeScan := time.Time{}
	err = pr.db.QueryRow(queries.PostGet, id).
		Scan(
			&post.Id,
			&post.Parent,
			&post.Author,
			&post.Message,
			&post.IsEdited,
			&post.Forum,
			&post.Thread,
			&timeScan)
	post.Created = timeScan.Format(time.RFC3339)
	return
}

func (pr *PostRepositoryImpl) UpdatePost(post *models.Post) (err error) {
	_, err = pr.db.Exec(queries.PostUpdate, post.Message, post.IsEdited, post.Id)
	return
}
