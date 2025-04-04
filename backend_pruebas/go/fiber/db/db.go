package db

import (
	stdsql "database/sql" // Estándar (para sql.Open)
	"log"

	entSql "entgo.io/ent/dialect/sql" // Para el driver de Ent

	_ "github.com/jackc/pgx/v5/stdlib" // Soporte para usar pgx con database/sql
)

// Devuelve un driver para Ent
func GetEntDriver(databaseURL string) *entSql.Driver {
	// Convertimos pgx en un sql.DB con el stdlib de pgx
	db, err := stdsql.Open("pgx", databaseURL)
	if err != nil {
		log.Fatalf("No se pudo abrir la conexión: %v", err)
	}

	// Probamos la conexión
	if err := db.Ping(); err != nil {
		log.Fatalf("No se pudo hacer ping a la DB: %v", err)
	}

	// Creamos el driver que Ent necesita
	driver := entSql.OpenDB("postgres", db)
	return driver
}
