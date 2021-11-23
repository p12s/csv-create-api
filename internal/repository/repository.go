package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type Repository struct {
	Producter
}

func NewRepository(db *sqlx.DB) *Repository {
	createProductTable(db)
	insertSeedData(db) // remove if seed data is not needed

	return &Repository{
		Producter: NewProduct(db),
	}
}

func createProductTable(db *sqlx.DB) {
	query := `CREATE TABLE IF NOT EXISTS product (
		"id" integer NOT NULL PRIMARY KEY AUTOINCREMENT,			
		"name" TEXT NOT NULL,
		"price" INTEGER DEFAULT 0
	  );`

	statement, err := db.Prepare(query)
	if err != nil {
		logrus.Panic("create product table error", err.Error())
	}
	statement.Exec()
	logrus.Info("product table created 🗂")
}

func insertSeedData(db *sqlx.DB) {
	query := `INSERT INTO product (name, price)
		VALUES ('Soflyy T-Shirt', 42),
		('Soflyy T-Shirt (blue)', 43),
		('Soflyy T-Shirt (red)', 44),
		('Soflyy Mug', 77),
		('Soflyy Mug 1', 777),
		('Soflyy Mug 2', 7774);`

	statement, err := db.Prepare(query)
	if err != nil {
		logrus.Panic("insert product table error", err.Error())
	}
	statement.Exec()
	logrus.Info("product table seed data added 🗂")
}
