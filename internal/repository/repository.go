package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

// Repository - repo
type Repository struct {
	Producter
	Authorizer
}

// NewRepository - constructor
func NewRepository(db *sqlx.DB) *Repository {
	createProductTable(db)
	insertSeedData(db) // remove if seed-data is not needed
	createUsersTable(db)

	return &Repository{
		Producter:  NewProduct(db),
		Authorizer: NewAuth(db),
	}
}

// createProductTable - creating memory table
func createProductTable(db *sqlx.DB) {
	query := `CREATE TABLE IF NOT EXISTS product (
		"id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,			
		"name" TEXT NOT NULL,
		"price" INTEGER DEFAULT 0
	);`

	statement, err := db.Prepare(query)
	defer statement.Close() //nolint
	if err != nil {
		logrus.Panic("create product table error", err.Error())
	}
	_, err = statement.Exec()
	if err != nil {
		logrus.Panic("create product table error", err.Error())
	}
	logrus.Info("product table created 🗂")
}

// insertSeedData - insert seeds
func insertSeedData(db *sqlx.DB) {
	query := `INSERT INTO product (name, price)
		VALUES ('Soflyy T-Shirt', 42),
		('Soflyy T-Shirt (blue)', 43),
		('Soflyy T-Shirt (red)', 44),
		('Soflyy Mug', 77),
		('Soflyy Mug 1', 777),
		('Soflyy Mug 2', 7774);`

	statement, err := db.Prepare(query)
	defer statement.Close() //nolint
	if err != nil {
		logrus.Panic("insert product table error", err.Error())
	}
	_, err = statement.Exec()
	if err != nil {
		logrus.Panic("create product table error", err.Error())
	}
	logrus.Info("product table seed data added 🗂")
}

// createUsersTable - creating memory table
func createUsersTable(db *sqlx.DB) {
	query := `CREATE TABLE IF NOT EXISTS users (
		"id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,			
		"name" TEXT NOT NULL,
		"email" TEXT NOT NULL,
		"password" TEXT NOT NULL,
		"registered_at" DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL
	);`

	statement, err := db.Prepare(query)
	defer statement.Close() //nolint
	if err != nil {
		logrus.Panic("create users table error", err.Error())
	}
	_, err = statement.Exec()
	if err != nil {
		logrus.Panic("create users table error", err.Error())
	}
	logrus.Info("product users created 🗂")
}
