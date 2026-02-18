package repository

import (
	"context"
	"fmt"

	"github.com/aclgo/grpc-jwt/internal/models"
	"github.com/jmoiron/sqlx"
)

type postgresRepo struct {
	db *sqlx.DB
}

func NewPostgresRepo(db *sqlx.DB) *postgresRepo {
	return &postgresRepo{
		db: db,
	}
}

func (p *postgresRepo) Add(ctx context.Context, user *models.User) (*models.User, error) {
	createdUser := models.User{}
	err := p.db.QueryRowxContext(ctx, queryAddUser,
		user.UserID,
		user.Name,
		user.Lastname,
		user.Password,
		user.Email,
		user.Role,
		user.Verified,
		user.CreatedAt,
		user.UpdatedAt,
	).StructScan(&createdUser)

	if err != nil {
		return nil, fmt.Errorf("Add.QueryRowxContext: %v", err)
	}
	return &createdUser, nil
}

func (p *postgresRepo) FindByID(ctx context.Context, userID string) (*models.User, error) {
	user := models.User{}
	err := p.db.GetContext(ctx, &user, queryByID, userID)
	if err != nil {
		return nil, fmt.Errorf("FindByID.GetContext: %v", err)
	}
	return &user, nil
}
func (p *postgresRepo) FindByEmail(ctx context.Context, userEmail string) (*models.User, error) {
	user := models.User{}
	err := p.db.GetContext(ctx, &user, queryFindByEmail, userEmail)
	if err != nil {
		return nil, fmt.Errorf("FindByEmail.GetContext: %v", err)
	}
	return &user, nil
}
func (p *postgresRepo) Update(ctx context.Context, user *models.User) (*models.User, error) {
	updatedUser := models.User{}

	err := p.db.QueryRowxContext(ctx, queryUpdate,
		user.Name,
		user.Lastname,
		user.Password,
		user.Email,
		user.Role,
		user.Verified,
		user.UpdatedAt,
		user.UserID,
	).StructScan(&updatedUser)

	if err != nil {
		return nil, fmt.Errorf("Update.QueryRowxContext: %v", err)
	}

	return &updatedUser, nil
}

func (p *postgresRepo) Delete(ctx context.Context, userID string) error {
	if _, err := p.db.Exec(queryDelete, userID); err != nil {
		return err
	}

	return nil
}
