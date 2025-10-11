package stores

import (
	"app/internal/models"
	"app/internal/models/dto"
	"context"
	"database/sql"
	"fmt"
	"log"

	"firebase.google.com/go/v4/auth"
)

type UserStore struct {
	db *sql.DB
}

func NewUserStore(db *sql.DB) *UserStore {
	return &UserStore{
		db: db,
	}
}

func (us *UserStore) CreateUser(ctx context.Context, user *auth.UserRecord, req *dto.CreateUserRequest) (bool, error) {
	if us == nil || us.db == nil {
		return false, fmt.Errorf("contest store: db is not initialized")
	}

	const q = `
	INSERT INTO users (id, name, email, usn, mobile_number, current_year, department)
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	ON CONFLICT (id) DO NOTHING
	`

	res, err := us.db.ExecContext(ctx, q, user.UID, req.Name, user.Email, req.USN, req.MobileNumber, req.CurrentYear, req.Department)
	if err != nil {
		log.Printf("user-store: insert error %v", err)
		return false, fmt.Errorf("insert error: %w", err)
	}

	// If rows affected is 0, then the user already exists
	rows, err := res.RowsAffected()
	if err != nil {
		log.Printf("user-store: update error %v", err)
		return false, fmt.Errorf("update error: %w", err)
	}

	if rows == 0 {
		return false, nil
	}

	return true, nil
}

func (us *UserStore) GetUserProfile(ctx context.Context, userID string) (*models.User, error) {
	if us == nil || us.db == nil {
		return nil, fmt.Errorf("db is not initialized")
	}

	const q = `
	SELECT id, name, email, usn, mobile_number, current_year, department
	FROM users
	WHERE id = $1
	`

	row := us.db.QueryRowContext(ctx, q, userID)

	var user models.User
	if err := row.Scan(&user.ID, &user.Name, &user.Email, &user.USN, &user.MobileNumber, &user.CurrentYear, &user.Department); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		log.Printf("user-store: row error %v", err)
		return nil, fmt.Errorf("row error: %w", err)
	}

	return &user, nil
}

func (us *UserStore) UpdateUserProfile(ctx context.Context, userID string, req *dto.UpdateUserProfileRequest) (bool, error) {
	if us == nil || us.db == nil {
		return false, fmt.Errorf("contest store: db is not initialized")
	}

	const q = `
	UPDATE users
	SET name = $2, usn = $3, mobile_number = $4, current_year = $5, department = $6
	WHERE id = $1
	`

	res, err := us.db.ExecContext(ctx, q, userID, req.Name, req.USN, req.MobileNumber, req.CurrentYear, req.Department)
	if err != nil {
		log.Printf("user-store: update error %v", err)
		return false, fmt.Errorf("update error: %w", err)
	}

	// If rows affected is 0, then the user does not exist
	rows, err := res.RowsAffected()
	if err != nil {
		log.Printf("user-store: update error %v", err)
		return false, fmt.Errorf("update error: %w", err)
	}

	if rows == 0 {
		return false, nil
	}

	return true, nil
}
