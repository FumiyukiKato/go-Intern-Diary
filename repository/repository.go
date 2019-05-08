package repository

import (
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/hatena/go-Intern-Diary/model"
	"github.com/jmoiron/sqlx"
)

type Repository interface {
	CreateNewUser(name string, passwordHash string) error
	FindUserByName(name string) (*model.User, error)
	FindUserByID(id uint64) (*model.User, error)
	ListUsersByIDs(userIDs []uint64) ([]*model.User, error)
	FindPasswordHashByName(name string) (string, error)
	CreateNewToken(userID uint64, token string, expiresAt time.Time) error
	FindUserByToken(token string) (*model.User, error)

	CreateNewDiary(userID uint64, name string) (*model.Diary, error)
	ListDiariesByUserID(userID, limit, offset uint64) ([]*model.Diary, error)
	DeleteDiary(userID, diaryID uint64) error
	ListArticlesByDiaryID(diaryID, limit, offset uint64) ([]*model.Article, error)
	FindDiaryByID(diaryID, userID uint64) (*model.Diary, error)
	CreateNewArticle(diaryID uint64, title string, content string) (*model.Article, error)
	FindArticleByID(articleID, diaryID uint64) (*model.Article, error)

	Close() error
}

func New(dsn string) (Repository, error) {
	db, err := sqlx.Open("mysql", dsn+"?parseTime=true&loc=Asia%2FTokyo")
	if err != nil {
		return nil, fmt.Errorf("Opening mysql failed: %v", err)
	}
	return &repository{db: db}, nil
}

type repository struct {
	db *sqlx.DB
}

func (r *repository) generateID() (uint64, error) {
	var id uint64
	err := r.db.Get(&id, "SELECT UUID_SHORT()")
	return id, err
}

func (r *repository) Close() error {
	return r.db.Close()
}
