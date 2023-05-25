package database

import (
    "gorm.io/gorm"
    "gorm.io/driver/postgres"
    "fmt"
)

func NewGormDatabase() (*gorm.DB, error) {
    dsn := createPostgresDSN("todo-sample-database", "test", "test", "todo-sample")
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        return nil, err
    }

    if err := Migrate(db); err != nil {
        return nil, err
    }
    
    return db, nil
}

func Migrate(db *gorm.DB) error {
    return db.AutoMigrate(Todo{})
}

func createPostgresDSN(host string, username string, password string, database string) string {
    return fmt.Sprintf("postgres://%s:%s@%s:5432/%s", username, password, host, database)
}
