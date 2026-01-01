package user

import (
	"database/sql"

	"github.com/manaswini1869/interview-prep/go-api/ecom/types"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) GetUserByEmail(email string) (*types.User, error) {
	query := "SELECT * FROM users WHERE email = ?"
	row, err := s.db.Query(query, email)
	if err != nil {
		return nil, err
	}
	user := new(types.User)
	for row.Next() {
		user, err = scanRowIntoUser(row)
		if err != nil {
			return nil, err
		}
	}
	if user.ID == 0 {
		return nil, sql.ErrNoRows
	}
	return user, nil

}

func (s *Store) CreateUser(user *types.User) error {

	return nil
}

func (s *Store) GetUserById(id int) (*types.User, error) {
	return nil, nil
}

func scanRowIntoUser(row *sql.Rows) (*types.User, error) {
	u := new(types.User)
	err := row.Scan(&u.ID, &u.FirstName, &u.LastName, &u.Email, &u.Password, &u.CreatedAt)
	if err != nil {
		return nil, err
	}
	return u, nil
}
