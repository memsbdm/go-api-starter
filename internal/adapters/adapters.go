package adapters

import (
	"database/sql"
	"go-starter/internal/adapters/storage/postgres/repositories"
	"go-starter/internal/adapters/token"
	"go-starter/internal/domain/ports"
)

type Adapters struct {
	UserRepository  ports.UserRepository
	TokenRepository ports.TokenRepository
	CacheRepository ports.CacheRepository
}

func New(db *sql.DB, timeGenerator ports.TimeGenerator, cache ports.CacheRepository) *Adapters {
	return &Adapters{
		UserRepository:  repositories.NewUserRepository(db),
		TokenRepository: token.NewTokenRepository(timeGenerator),
		CacheRepository: cache,
	}
}
