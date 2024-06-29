package infra

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/BernardoDenkvitts/MySQLApp/internal/types"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
)

const (
	username      = "root"
	MySQLpassword = "root"
	hostname      = "127.0.0.1:3306"
	dbname        = "mysqluser"
	parseTime     = "true"
)

func dsn() string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=%s", username, MySQLpassword, hostname, dbname, parseTime)
}

func scanIntoUser(rows *sql.Rows) (*types.User, error) {
	user := new(types.User)
	if err := rows.Scan(
		&user.Id,
		&user.FirstName,
		&user.LastName,
		&user.CreatedAt,
	); err != nil {
		return nil, err
	}

	return user, nil
}

type Storage interface {
	Init() error
	CreateUserInformation(*types.User) error
	GetUserById(id string) (*types.User, error)
	// This function will be use to get users created in the last 5 minutes
	// to be sent to rabbitMQ
	GetLatestUserInformations() ([]*types.User, error)
	GetUsersInformations() ([]*types.User, error)
}

type MySQLStore struct {
	db *sql.DB
}

func NewMySQLStore() (*MySQLStore, error) {
	db, err := sql.Open("mysql", dsn())
	if err != nil {
		return nil, err
	}

	return &MySQLStore{
		db: db,
	}, nil
}

func (s *MySQLStore) Init() error {
	query := `CREATE TABLE IF NOT EXISTS user (
		id varchar(50) primary key,
		firstName varchar(50),
		lastName varchar(50),
		created_at timestamp
	)`
	_, err := s.db.Exec(query)
	if err != nil {
		panic(err)
	}

	log.Println("Database Initialized")
	return nil
}

func (s *MySQLStore) CreateUserInformation(user *types.User) error {

	userId := uuid.Must(uuid.NewRandom()).String()
	query := "INSERT INTO user (id, firstName, lastName, created_at) VALUES(?, ?, ?, ?)"
	_, err := s.db.Query(query, userId, user.FirstName, user.LastName, user.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (s *MySQLStore) GetUserById(id string) (*types.User, error) {
	rows, err := s.db.Query("SELECT * FROM user WHERE id = ?", id)
	if err != nil {
		return nil, err
	}

	user := new(types.User)
	for rows.Next() {
		user, err = scanIntoUser(rows)
		if err != nil {
			return nil, err
		}
	}

	return user, nil
}

func (s *MySQLStore) GetUsersInformations() ([]*types.User, error) {
	rows, err := s.db.Query("SELECT * FROM user")
	if err != nil {
		return nil, err
	}

	users := []*types.User{}
	for rows.Next() {
		user, err := scanIntoUser(rows)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (s *MySQLStore) GetLatestUserInformations() ([]*types.User, error) {
	// TODO verificar essa query
	rows, err := s.db.Query("SELECT * FROM user WHERE created_at >= UTC_TIMESTAMP() - INTERVAL 5 MINUTE")
	if err != nil {
		return nil, err
	}

	var users []*types.User

	for rows.Next() {
		user, err := scanIntoUser(rows)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}
