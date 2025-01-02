package database

import (
	"fmt"
	"log"
	ds "orderFood-server/internal/dataSchema"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/joho/godotenv/autoload"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Service represents a service that interacts with a database.
type Service interface {
	AddUser(username, password string) error
	DelUser(id uint64) error
	GetUserById(id string) (ds.User, error)

	// Close terminates the database connection.
	// It returns an error if the connection cannot be closed.
	Close() error
}

type service struct {
	db *gorm.DB
}

var (
	database   = os.Getenv("BLUEPRINT_DB_DATABASE")
	password   = os.Getenv("BLUEPRINT_DB_PASSWORD")
	username   = os.Getenv("BLUEPRINT_DB_USERNAME")
	port       = os.Getenv("BLUEPRINT_DB_PORT")
	host       = os.Getenv("BLUEPRINT_DB_HOST")
	schema     = os.Getenv("BLUEPRINT_DB_SCHEMA")
	dbInstance *service
)

func New() Service {
	// Reuse Connection
	if dbInstance != nil {
		return dbInstance
	}
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable&search_path=%s", username, password, host, port, database, schema)
	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	// add table
	err = db.AutoMigrate(&ds.User{})
	if err != nil {
		log.Fatal(err)
	}

	dbInstance = &service{
		db: db,
	}
	return dbInstance
}

func (s *service) AddUser(username, password string) error {
	user := ds.User{
		Username: username,
		Password: password,
	}
	result := s.db.Create(&user)
	if result.Error != nil {
		log.Printf("Error adding user: %v", result.Error)
		return result.Error
	}
	log.Printf("User added successfully, id: %d", user.ID)

	return nil
}

func (s *service) DelUser(id uint64) error {
	result := s.db.Model(&ds.User{}).Where("id = ?", id).Update("deleted_at", time.Now())
	if result.Error != nil {
		log.Printf("Error deleting user: %v", result.Error)
		return result.Error
	}
	log.Printf("User deleted successfully, id: %d", id)
	return nil
}

func (s *service) GetUserById(id string) (ds.User, error) {
	user := ds.User{}
	result := s.db.First(&user, id)
	if result.Error != nil {
		log.Printf("Error getting users: %v", result.Error)
		return user, result.Error
	}
	return user, nil
}

// Close closes the database connection.
// It logs a message indicating the disconnection from the specific database.
// If the connection is successfully closed, it returns nil.
// If an error occurs while closing the connection, it returns the error.
func (s *service) Close() error {
	log.Printf("Disconnected from database: %s", database)
	sqlDB, err := s.db.DB() // 获取原始的 sql.DB 对象
	if err != nil {
		log.Printf("failed to get sql.DB: %v", err)
		return err
	}
	return sqlDB.Close()
}
