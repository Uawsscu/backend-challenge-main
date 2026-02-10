package mongodb

import (
	"time"

	"github.com/backend-challenge/user-api/internal/domain"
)

type userDoc struct {
	ID        string    `bson:"_id,omitempty"`
	Name      string    `bson:"name"`
	Email     string    `bson:"email"`
	Password  string    `bson:"password"`
	CreatedAt time.Time `bson:"created_at"`
}

func fromDomain(u *domain.User) *userDoc {
	if u == nil {
		return nil
	}
	return &userDoc{
		ID:        u.ID,
		Name:      u.Name,
		Email:     u.Email,
		Password:  u.Password,
		CreatedAt: u.CreatedAt,
	}
}

func (d *userDoc) toDomain() *domain.User {
	if d == nil {
		return nil
	}
	return &domain.User{
		ID:        d.ID,
		Name:      d.Name,
		Email:     d.Email,
		Password:  d.Password,
		CreatedAt: d.CreatedAt,
	}
}
