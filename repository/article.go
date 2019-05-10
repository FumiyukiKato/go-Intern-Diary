package repository

import (
	"database/sql"
	"time"

	"github.com/hatena/go-Intern-Diary/model"
	"github.com/jmoiron/sqlx"
)

var articleNotFoundError = model.NotFoundError("article")

func (r *repository) ListArticlesByDiaryID(diaryID, limit, offset uint64) ([]*model.Article, error) {
	articles := make([]*model.Article, 0, limit)
	err := r.db.Select(&articles,
		`SELECT id, title, content, diary_id, updated_at FROM article
			WHERE diary_id = ?
			ORDER BY updated_at DESC LIMIT ? OFFSET ?`,
		diaryID, limit, offset,
	)
	return articles, err
}

func (r *repository) CreateNewArticle(diaryID uint64, title string, content string) (*model.Article, error) {
	id, err := r.generateID()
	if err != nil {
		return nil, err
	}
	now := time.Now()
	_, err = r.db.Exec(
		`INSERT INTO article 
			(id, title, content, diary_id, created_at, updated_at)
			VALUES(?,?,?,?,?,?)`,
		id, title, content, diaryID, now, now,
	)
	if err != nil {
		return nil, err
	}
	return &model.Article{ID: id, Title: title, DiaryID: diaryID, UpdatedAt: now}, nil
}

func (r *repository) FindArticleByID(articleID, diaryID uint64) (*model.Article, error) {
	var article model.Article
	err := r.db.Get(&article,
		`SELECT id, title, content, diary_id, updated_at FROM article
			WHERE id = ? AND diary_id = ? LIMIT 1`,
		articleID, diaryID,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, articleNotFoundError
		}
		return nil, err
	}
	return &article, nil
}

func (r *repository) DeleteArticle(articleID, diaryID uint64) (err error) {
	_, err = r.db.Exec(
		`DELETE FROM article
			WHERE id = ? AND diary_id = ?`,
		articleID, diaryID,
	)
	return
}

func (r *repository) ListArticlesByIDs(articleIDs []uint64) ([]*model.Article, error) {
	if len(articleIDs) == 0 {
		return nil, nil
	}
	articles := make([]*model.Article, len(articleIDs))
	query, args, err := sqlx.In(
		`SELECT id, title, content, diary_id, updated_at FROM article
			WHERE id IN (?)
			ORDER BY created_at DESC`, articleIDs,
	)
	if err != nil {
		return nil, err
	}
	err = r.db.Select(&articles, query, args...)
	return articles, err
}

func (r *repository) ListArticlesByDiaryIDs(diaryIDs []uint64) (map[uint64][]*model.Article, error) {
	if len(diaryIDs) == 0 {
		return nil, nil
	}
	query, args, err := sqlx.In(
		`SELECT id, title, content, diary_id, updated_at FROM article
			WHERE diary_id IN (?)`,
		diaryIDs,
	)
	if err != nil {
		return nil, err
	}
	rows, err := r.db.Queryx(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	articles := make(map[uint64][]*model.Article)
	for rows.Next() {
		var article model.Article
		rows.Scan(
			&article.ID,
			&article.Title,
			&article.Content,
			&article.DiaryID,
			&article.UpdatedAt,
		)
		articles[article.DiaryID] = append(articles[article.DiaryID], &article)
	}
	return articles, nil
}
