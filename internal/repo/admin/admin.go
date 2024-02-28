package admin

import (
	"errors"
	"fmt"
	"time"
	"wabot/internal/database"
	"wabot/internal/helpers"

	"github.com/alexedwards/argon2id"
)

type Admin struct {
	Id           int       `db:"id" json:"id"`
	Username     string    `db:"username" json:"username"`
	Password     string    `db:"password" json:"password,omitempty"`
	LastactiveTs time.Time `db:"lastactive_ts" json:"lastactive_ts"`
	CreatedTs    time.Time `db:"created_ts" json:"created_ts"`
	IsDeleted    bool      `db:"is_deleted" json:"is_deleted"`
}

type AdminRepo struct {
	db *database.DB
}

type AdminRepository interface {
	CreateAdmin(*Admin) (*Admin, error)
	Login(username string, password string) (*Admin, error)
}

func NewAdminRepo(db *database.DB) AdminRepository {
	return &AdminRepo{db}
}

func (repo *AdminRepo) CreateAdmin(admin *Admin) (*Admin, error) {
	if !helpers.ValidUsername(admin.Username) {
		return nil, errors.New("username is not valid, must be at least 4 characters long and contain only lowercase letters and numbers")
	}
	if !helpers.StrongPassword(admin.Password) {
		return nil, errors.New("password is not strong enough, must be at least 8 characters long and contain at least one uppercase letter, one lowercase letter, one number and one special character")
	}
	hash, err := setPassword(admin.Password)
	if err != nil {
		return nil, err
	}
	admin.Password = hash

	admin.CreatedTs = time.Now()

	if repo.isUsernameExists(admin.Username) {
		return nil, errors.New("username already exists")
	}

	query := `INSERT INTO tbl_admin (username, password, created_ts, is_deleted) VALUES ($1, $2, $3, $4)`
	_, err = repo.db.Exec(query, admin.Username, admin.Password, admin.CreatedTs, 0)
	if err != nil {
		return nil, err
	}
	return admin, nil
}

func (repo *AdminRepo) Login(username, password string) (*Admin, error) {
	query := `SELECT id, username, password FROM tbl_admin WHERE username = $1 AND is_deleted = $2 LIMIT 1`
	row, err := repo.db.Query(query, username, false)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer row.Close()

	admin := Admin{}
	if row.Next() {
		err = row.Scan(&admin.Id, &admin.Username, &admin.Password)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
	} else {
		return nil, errors.New("Invalid username or password")
	}

	match, err := checkPassword(admin.Password, password)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	if !match {
		return nil, errors.New("Invalid username or password")
	}

	update := `UPDATE tbl_admin SET lastactive_ts = $1 WHERE id = $2`
	_, err = repo.db.Exec(update, time.Now(), admin.Id)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &admin, nil
}

func (repo *AdminRepo) isUsernameExists(username string) bool {
	query := `SELECT id FROM tbl_admin WHERE username = $1 AND is_deleted = $2 LIMIT 1`
	row, err := repo.db.Query(query, username, 0)
	if err != nil {
		return false
	}
	defer row.Close()
	return row.Next()
}

func setPassword(password string) (string, error) {
	hash, err := argon2id.CreateHash(password, argon2id.DefaultParams)
	if err != nil {
		return "", err
	}
	return hash, nil
}

func checkPassword(hash, password string) (bool, error) {
	match, err := argon2id.ComparePasswordAndHash(password, hash)
	if err != nil {
		return false, err
	}
	return match, nil
}
