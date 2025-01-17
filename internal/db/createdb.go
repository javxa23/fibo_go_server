package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func createSchemaAndTables(db *sql.DB) {
	ctx := context.Background()
	var schemaExists bool
	err := db.QueryRowContext(ctx, `SELECT EXISTS (SELECT 1 FROM pg_catalog.pg_namespace WHERE nspname = 'blog')`).Scan(&schemaExists)
	if err != nil {
		log.Fatalf("Error checking schema existence: %v", err)
	}

	if !schemaExists {
		createSchema := `CREATE SCHEMA IF NOT EXISTS blog;`
		if _, err := db.ExecContext(ctx, createSchema); err != nil {
			log.Fatalf("Error creating schema: %v", err)
		}
		fmt.Println("Schema 'blog' created.")
	}
	tables := []struct {
		name   string
		create string
	}{
		{
			name: "users",
			create: `
				CREATE TABLE IF NOT EXISTS blog.users (
					id SERIAL PRIMARY KEY,
					username VARCHAR(50) NOT NULL UNIQUE,
					email VARCHAR(100) NOT NULL UNIQUE,
					password_hash TEXT NOT NULL,
					created_at TIMESTAMP DEFAULT NOW()
				);`,
		},
		{
			name: "posts",
			create: `
				CREATE TABLE IF NOT EXISTS blog.posts (
					id SERIAL PRIMARY KEY,
					user_id INT REFERENCES blog.users(id) ON DELETE CASCADE,
					title VARCHAR(200) NOT NULL,
					content TEXT NOT NULL,
					created_at TIMESTAMP DEFAULT NOW()
				);`,
		},
		{
			name: "comments",
			create: `
				CREATE TABLE IF NOT EXISTS blog.comments (
					id SERIAL PRIMARY KEY,
					post_id INT REFERENCES blog.posts(id) ON DELETE CASCADE,
					user_id INT REFERENCES blog.users(id) ON DELETE CASCADE,
					comment_text TEXT NOT NULL,
					created_at TIMESTAMP DEFAULT NOW()
				);`,
		},
		{
			name: "ratings",
			create: `
				CREATE TABLE IF NOT EXISTS blog.ratings (
					id SERIAL PRIMARY KEY,
					post_id INT REFERENCES blog.posts(id) ON DELETE CASCADE,
					user_id INT REFERENCES blog.users(id) ON DELETE CASCADE,
					rating INT CHECK (rating BETWEEN 1 AND 5),
					created_at TIMESTAMP DEFAULT NOW()
				);`,
		},
	}
	for _, table := range tables {
		var tableExists bool
		err := db.QueryRowContext(ctx, `SELECT EXISTS (SELECT 1 FROM pg_catalog.pg_tables WHERE schemaname = 'blog' AND tablename = $1)`, table.name).Scan(&tableExists)
		if err != nil {
			log.Fatalf("Error checking table existence: %v", err)
		}

		if !tableExists {
			if _, err := db.ExecContext(ctx, table.create); err != nil {
				log.Fatalf("Error creating table %s: %v", table.name, err)
			}
			fmt.Printf("Table '%s' created.\n", table.name)
		} else {
			fmt.Printf("Table '%s' already exists.\n", table.name)
		}
	}

	fmt.Println("Schema and tables check complete.")
}
