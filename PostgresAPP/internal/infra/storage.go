package infra

import (
	"database/sql"
	"log"

	"github.com/BernardoDenkvitts/PostgresAPP/internal/types"
	"github.com/BernardoDenkvitts/PostgresAPP/internal/utils"
	_ "github.com/lib/pq"
)

type Storage interface {
	Init() error
	CreateUserInformation(*types.User) error
	GetUserById(id string) (*types.User, error)
	// This function will be use to get users created in the last 5 minutes
	// to be sent to rabbitMQ
	GetLatestUserInformations() ([]*types.User, error)
	GetUsersInformations() ([]*types.User, error)
}

type PostgresStore struct {
	db *sql.DB
}

func NewPostgresStore() (*PostgresStore, error) {
	db, err := sql.Open("postgres", "user=postgres dbname=postgres password=root sslmode=disable")
	utils.FailOnError(err, "Failed to connect to database")

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresStore{
		db: db,
	}, nil
}

func (s *PostgresStore) Init() error {
	query := `CREATE TABLE IF NOT EXISTS userinfo (
		id VARCHAR(50) PRIMARY KEY,
		firstName VARCHAR(50), 
		lastName VARCHAR(50),
		created_at TIMESTAMP
	);`

	_, err := s.db.Exec(query)
	utils.FailOnError(err, "Failed to create table")
	log.Println("Table created")

	return nil
}

func (s *PostgresStore) CreateUserInformation(user *types.User) error {
	query := "INSERT INTO userinfo (id, firstName, lastName, created_at) VALUES ($1, $2, $3, $4)"
	_, err := s.db.Query(query, user.Id, user.FirstName, user.LastName, user.CreatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (s *PostgresStore) GetUserById(id string) (*types.User, error) {
	query := "SELECT * FROM userinfo WHERE userinfo.id = $1;"
	row, err := s.db.Query(query, id)
	if err != nil {
		return nil, err
	}

	user := new(types.User)
	for row.Next() {
		user, err = scanIntoUser(row)
		if err != nil {
			return nil, err
		}
	}

	return user, nil
}

func (s *PostgresStore) GetLatestUserInformations() ([]*types.User, error) {
	query := "SELECT * FROM userinfo WHERE userinfo.created_at >= NOW() at time zone 'utc' - INTERVAL '30 second'"
	rows, err := s.db.Query(query)
	if err != nil {
		return nil, err
	}

	latestUsers := []*types.User{}
	for rows.Next() {
		user, err := scanIntoUser(rows)
		if err != nil {
			return nil, err
		}
		latestUsers = append(latestUsers, user)
	}

	return latestUsers, nil
}

func (s *PostgresStore) GetUsersInformations() ([]*types.User, error) {
	query := "SELECT * FROM userinfo"
	rows, err := s.db.Query(query)
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
