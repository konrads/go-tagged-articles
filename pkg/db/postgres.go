package db

import (
	"database/sql"
	"log"
	"time"

	"github.com/konrads/go-tagged-articles/pkg/model"
	"github.com/lib/pq"
)

type PostgresDB struct {
	db *sql.DB
}

func NewPostgresDB(connUri *string) *PostgresDB {
	db, err := sql.Open("postgres", *connUri)
	if err != nil {
		log.Fatalf("failed to connect to postgres on uri: %s", *connUri)
	}
	return &PostgresDB{db: db}
}

func (db *PostgresDB) GetArticle(id string) (*model.Article, error) {
	rows := db.db.QueryRow("SELECT title, date, body, tags FROM article WHERE id = $1", id)
	var title string
	var date time.Time
	var body string
	var tags []sql.NullString

	err := rows.Scan(&title, &date, &body, pq.Array(&tags))
	if err == sql.ErrNoRows {
		return nil, err
	}
	if err != nil {
		log.Fatalf("failed to unmarshall postgres row due to %v", err)
		return nil, err
	}
	res := model.Article{
		ID: id,
		ArticleReq: &model.ArticleReq{
			Title: title,
			Date:  model.Date(date),
			Body:  body,
			Tags:  toStringArr(tags),
		},
	}
	return &res, nil
}

func (db *PostgresDB) SaveArticle(a *model.Article) (int, error) {
	res, err := db.db.Exec(
		"INSERT INTO article (id, title, date, body, tags) VALUES ($1, $2, $3, $4, $5) ON CONFLICT(id) DO NOTHING",
		a.ID, a.Title, time.Time(a.Date), a.Body, pq.Array(a.Tags),
	)
	if err != nil {
		return 0, err
	}
	ra, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return int(ra), nil
}

func (db *PostgresDB) GetArticles(tag string, date model.Date) ([]*model.Article, error) {
	rows, err := db.db.Query("SELECT id, title, date, body, tags FROM article WHERE date = $1 and tags @> ARRAY[$2]::text[]", time.Time(date), tag)
	defer rows.Close()
	if err != nil {
		return []*model.Article{}, err
	}

	var id string
	var title string
	var body string
	var tags []sql.NullString

	res := []*model.Article{}
	for rows.Next() {
		err := rows.Scan(&id, &title, &date, &body, pq.Array(&tags))
		if err != nil {
			return []*model.Article{}, err
		}
		a := model.Article{
			ID: id,
			ArticleReq: &model.ArticleReq{
				Title: title,
				Date:  model.Date(date),
				Body:  body,
				Tags:  toStringArr(tags),
			},
		}
		res = append(res, &a)
	}
	return res, nil
}

func (db *PostgresDB) Close() error {
	return db.db.Close()
}

func toStringArr(arr []sql.NullString) []string {
	res := make([]string, len(arr))
	for i, x := range arr {
		res[i] = x.String
	}
	return res
}
