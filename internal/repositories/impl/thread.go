package impl

import (
	"DBProject/internal/models"
	"DBProject/pkg/errors"
	"DBProject/pkg/handlerows"
	"DBProject/pkg/queries"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx"
	_ "github.com/lib/pq"
)

type ThreadRepositoryImpl struct {
	db *pgx.ConnPool
}

func NewThreadRepository(db *pgx.ConnPool) *ThreadRepositoryImpl {
	return &ThreadRepositoryImpl{db: db}
}

func (tr *ThreadRepositoryImpl) GetBySlug(slug string) (thread *models.Thread, err error) {
	thread = &models.Thread{}
	err = tr.db.QueryRow(queries.ThreadGetSlug, slug).
		Scan(&thread.Id, &thread.Title, &thread.Author, &thread.Forum, &thread.Message, &thread.Votes, &thread.Slug, &thread.Created)
	return
}

func (tr *ThreadRepositoryImpl) GetById(id int64) (thread *models.Thread, err error) {
	thread = &models.Thread{}
	err = tr.db.QueryRow(queries.ThreadGetId, id).
		Scan(&thread.Id, &thread.Title, &thread.Author, &thread.Forum, &thread.Message, &thread.Votes, &thread.Slug, &thread.Created)
	return
}

func (tr *ThreadRepositoryImpl) CreateThread(thread *models.Thread) (err error) {
	err = tr.db.QueryRow(queries.ThreadCreate, thread.Title, thread.Author, thread.Forum, thread.Message, thread.Slug, thread.Created).
		Scan(
			&thread.Id,
			&thread.Created)
	return
}

func (tr *ThreadRepositoryImpl) GetThread(slugOrId interface{}) (*models.Thread, error) {
	thread := &models.Thread{}
	var err error
	switch slugOrId.(type) {
	case string:
		err = tr.db.QueryRow(queries.ThreadGetSlug, slugOrId).
			Scan(&thread.Id, &thread.Title, &thread.Author, &thread.Forum, &thread.Message, &thread.Votes, &thread.Slug, &thread.Created)
	case int64:
		id, _ := strconv.Atoi(slugOrId.(string))
		err = tr.db.QueryRow(queries.ThreadGetId, int64(id)).
			Scan(&thread.Id, &thread.Title, &thread.Author, &thread.Forum, &thread.Message, &thread.Votes, &thread.Slug, &thread.Created)
	}
	return thread, err
}

func (tr *ThreadRepositoryImpl) GetThreadVotes(id int64) (int32, error) {
	var votes int32
	err := tr.db.QueryRow(queries.ThreadVotes, id).Scan(&votes)
	return votes, err
}

func (tr *ThreadRepositoryImpl) UpdateThread(thread *models.Thread) error {
	_, err := tr.db.Exec(queries.ThreadUpdate, thread.Title, thread.Message, thread.Id)
	return err
}

func (tr *ThreadRepositoryImpl) createPartPosts(thread *models.Thread, posts *models.Posts, from, to int, created time.Time, createdFormatted string) (err error) {

	args := make([]interface{}, 0, 0)
	query := queries.PostPart
	j := 0
	for i := from; i < to; i++ {
		(*posts)[i].Forum = thread.Forum
		(*posts)[i].Thread = thread.Id
		(*posts)[i].Created = createdFormatted
		query += fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d),", j*6+1, j*6+2, j*6+3, j*6+4, j*6+5, j*6+6)
		if (*posts)[i].Parent != 0 {
			args = append(args, (*posts)[i].Parent, (*posts)[i].Author, (*posts)[i].Message, thread.Forum, thread.Id, created)
		} else {
			args = append(args, nil, (*posts)[i].Author, (*posts)[i].Message, thread.Forum, thread.Id, created)
		}
		j++
	}
	query = query[:len(query)-1]
	query += " RETURNING id;"

	isSuccess := false
	k := 0

	for !isSuccess {

		resultRows, err := tr.db.Query(query, args...)
		if err != nil {
			fmt.Println(err)
			return errors.ErrParentPostNotExist
		}
		defer resultRows.Close()

		for i := from; resultRows.Next(); i++ {
			isSuccess = true
			var id int64
			if err = resultRows.Scan(&id); err != nil {
				return err
			}
			(*posts)[i].Id = id
		}
		k++
		if k >= 3 {
			break
		}
	}
	return
}

func (tr *ThreadRepositoryImpl) CreateThreadPosts(thread *models.Thread, posts *models.Posts) (err error) {
	created := time.Now()
	createdFormatted := created.Format(time.RFC3339)

	parts := len(*posts) / 20
	for i := 0; i < parts+1; i++ {
		if i == parts {
			if i*20 != len(*posts) {
				err = tr.createPartPosts(thread, posts, i*20, len(*posts), created, createdFormatted)
				if err != nil {
					return err
				}
			}
		} else {
			err = tr.createPartPosts(thread, posts, i*20, i*20+20, created, createdFormatted)
			if err != nil {
				return err
			}
		}
	}
	return
}

func (tr *ThreadRepositoryImpl) GetThreadPostsTree(id int64, limit, since int, desc bool) (*[]models.Post, error) {
	var rows *pgx.Rows
	var err error
	if since == -1 {
		if desc {
			rows, err = tr.db.Query(strings.Join([]string{queries.ThreadTreeBase, queries.ThreadTreeSinceDesc}, ""), id, limit)
		} else {
			rows, err = tr.db.Query(strings.Join([]string{queries.ThreadTreeBase, queries.ThreadTreeSince}, ""), id, limit)
		}
	} else {
		if desc {
			rows, err = tr.db.Query(strings.Join([]string{queries.ThreadTreeBase, queries.ThreadTreeDesc}, ""), id, since, limit)
		} else {
			rows, err = tr.db.Query(strings.Join([]string{queries.ThreadTreeBase, queries.ThreadTree}, ""), id, since, limit)
		}
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return handlerows.Post(rows)
}

func (tr *ThreadRepositoryImpl) GetThreadPostsParentTree(threadID int64, limit, since int, desc bool) (*[]models.Post, error) {
	var rows *pgx.Rows
	var err error
	if since == -1 {
		if desc {
			rows, err = tr.db.Query(strings.Join([]string{queries.ThreadParentBase, queries.ThreadParentTreeSinceDesc}, ""), threadID, limit)
		} else {
			rows, err = tr.db.Query(strings.Join([]string{queries.ThreadParentBase, queries.ThreadParentTreeSince}, ""), threadID, limit)
		}
	} else {
		if desc {
			rows, err = tr.db.Query(strings.Join([]string{queries.ThreadParentBase, queries.ThreadParentTreeDesc}, ""), threadID, since, limit)
		} else {
			rows, err = tr.db.Query(strings.Join([]string{queries.ThreadParentBase, queries.ThreadParentTree}, ""), threadID, since, limit)
		}
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return handlerows.Post(rows)
}

func (tr *ThreadRepositoryImpl) GetThreadPostsFlat(id int64, limit, since int, desc bool) (*[]models.Post, error) {
	var rows *pgx.Rows
	var err error
	if since == -1 {
		if desc {
			rows, err = tr.db.Query(strings.Join([]string{queries.ThreadFlatBase, queries.ThreadFlatSinceDesc}, ""), id, limit)
		} else {
			rows, err = tr.db.Query(strings.Join([]string{queries.ThreadFlatBase, queries.ThreadFlatSince}, ""), id, limit)
		}
	} else {
		if desc {
			rows, err = tr.db.Query(strings.Join([]string{queries.ThreadFlatBase, queries.ThreadFlatDesc}, ""), id, since, limit)
		} else {
			rows, err = tr.db.Query(strings.Join([]string{queries.ThreadFlatBase, queries.ThreadFlat}, ""), id, since, limit)
		}
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return handlerows.Post(rows)
}
