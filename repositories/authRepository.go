package repositories

import (
	"context"
	"e-commerce-API/models"

	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthRepository struct {
	db *pgxpool.Pool
}

func NewAuthRepository(conn *pgxpool.Pool) *AuthRepository {
	return &AuthRepository{db: conn}
}

func (repository *AuthRepository) Create(ctx context.Context, user models.User) (int, error) {
	var id int

	err := repository.db.QueryRow(ctx, "insert into users(name, email, password_hash) values($1, $2, $3) returning id", user.Name, user.Email, user.PasswordHash).Scan(&id)

	if err != nil {
		return 0, err
	}
	return id, err
}

func (repository *AuthRepository) FindByEmail(ctx context.Context, email string) (models.User, error) {
	row := repository.db.QueryRow(ctx, "select id, name, email, password_hash from users where email = $1", email)

	var user models.User
	err := row.Scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash)

	if err != nil {
		return models.User{}, err
	}
	return user, err
}
