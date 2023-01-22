package impl

import (
	"DBProject/internal/models"
	"DBProject/pkg/handlerows"
	"DBProject/pkg/queries"
	"github.com/jackc/pgx"
	_ "github.com/lib/pq"
)

type ForumRepositoryImpl struct {
	db *pgx.ConnPool
}

func NewForumRepository(db *pgx.ConnPool) *ForumRepositoryImpl {
	return &ForumRepositoryImpl{db: db}
}

func (fr *ForumRepositoryImpl) CreateForum(forum *models.Forum) (err error) {
	_, err = fr.db.Exec(queries.ForumCreate, forum.Title, forum.User, forum.Slug)
	return err
}

func (fr *ForumRepositoryImpl) GetInfoAboutForum(slug string) (forum *models.Forum, err error) {
	forum = new(models.Forum)
	err = fr.db.QueryRow(queries.ForumGetBySlug, slug).Scan(&forum.Title, &forum.User, &forum.Slug, &forum.Posts, &forum.Threads)
	return forum, err
}

func (fr *ForumRepositoryImpl) GetForumUsers(slug string, limit int, since string, desc bool) (*[]models.User, error) {
	var query string

	var result *pgx.Rows
	var innerError error

	if since != "" {
		if desc {
			query = queries.ForumGetUsersSinceDesc
		} else {
			query = queries.ForumGetUsersSince
		}
		result, innerError = fr.db.Query(query, slug, since, limit)
		if innerError != nil {
			return nil, innerError
		}
	} else {
		if desc {
			query = queries.ForumGetUsersDesc
		} else {
			query = queries.ForumGetUsers
		}
		result, innerError = fr.db.Query(query, slug, limit)
		if innerError != nil {
			return nil, innerError
		}
	}
	defer result.Close()
	return handlerows.User(result)
}

func (fr *ForumRepositoryImpl) GetForumThreads(slug string, limit int, since string, desc bool) (threads *[]models.Thread, err error) {
	var query string

	var result *pgx.Rows
	var innerError error

	if since != "" {
		if desc {
			query = queries.ForumGetThreadsSinceDesc
		} else {
			query = queries.ForumGetThreadsSince
		}
		result, innerError = fr.db.Query(query, slug, since, limit)
		if innerError != nil {
			return
		}
	} else {
		if desc {
			query = queries.ForumGetThreadsDesc
		} else {
			query = queries.ForumGetThreads
		}
		result, innerError = fr.db.Query(query, slug, limit)
		if innerError != nil {
			return
		}
	}

	defer result.Close()
	return handlerows.Thread(result)
}
