package utils

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
	err := db.QueryRowContext(ctx, `SELECT EXISTS (SELECT 1 FROM pg_catalog.pg_namespace WHERE nspname = 'fiboblog')`).Scan(&schemaExists)
	if err != nil {
		log.Fatalf("Error checking schema existence: %v", err)
	}

	if !schemaExists {
		createSchema := `CREATE SCHEMA IF NOT EXISTS fiboblog;`
		if _, err := db.ExecContext(ctx, createSchema); err != nil {
			log.Fatalf("Error creating schema: %v", err)
		}
		fmt.Println("Schema 'fiboblog' created.")
	}

	tables := []struct {
		name   string
		create string
	}{
		{
			name: "users",
			create: `
				CREATE TABLE IF NOT EXISTS fiboblog.users (
				   	id SERIAL PRIMARY KEY,
				    username VARCHAR(50) UNIQUE NOT NULL,
				    email VARCHAR(100) UNIQUE NOT NULL,
				    password VARCHAR(255) NOT NULL,
				    points INT DEFAULT 0,
				    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
				);`,
		},
		{
			name: "categories",
			create: `
				CREATE TABLE IF NOT EXISTS fiboblog.categories (
					id SERIAL PRIMARY KEY,
					name VARCHAR(50) UNIQUE NOT NULL,
					description TEXT,
					created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
				);
				INSERT INTO fiboblog.categories (name, description)
					VALUES
					  ('Technology', 'Articles related to the latest advancements in technology.'),
					  ('Health', 'Posts discussing wellness, fitness, and healthy living.'),
					  ('Business', 'Business strategies, tips, and industry news.'),
					  ('Lifestyle', 'Content on lifestyle, fashion, and everyday living.'),
					  ('Education', 'Resources and tips for learning and personal growth.');`,
		},
		{
			name: "posts",
			create: `
				CREATE TABLE IF NOT EXISTS fiboblog.posts (
					id SERIAL PRIMARY KEY,
					user_id INT NOT NULL,
					category_id INT,
					title VARCHAR(255) NOT NULL,
					content TEXT NOT NULL,
					is_approved BOOLEAN DEFAULT FALSE,
					view_count INT DEFAULT 0,
					created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
					updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
					CONSTRAINT fk_user FOREIGN KEY (user_id) REFERENCES fiboblog.users(id),
					CONSTRAINT fk_category FOREIGN KEY (category_id) REFERENCES fiboblog.categories(id)
				);`,
		},
		{
			name: "comments",
			create: `
				CREATE TABLE IF NOT EXISTS fiboblog.comments (
					id SERIAL PRIMARY KEY,
					post_id INT NOT NULL,
					user_id INT,
					parent_comment_id INT DEFAULT NULL,
					content TEXT NOT NULL,
					created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
					CONSTRAINT fk_post_comment FOREIGN KEY (post_id) REFERENCES fiboblog.posts(id),
					CONSTRAINT fk_user_comment FOREIGN KEY (user_id) REFERENCES fiboblog.users(id),
					CONSTRAINT fk_parent_comment FOREIGN KEY (parent_comment_id) REFERENCES fiboblog.comments(id)
				);`,
		},
		{
			name: "likes",
			create: `
				CREATE TABLE IF NOT EXISTS fiboblog.likes (
					id SERIAL PRIMARY KEY,
					post_id INT NOT NULL,
					user_id INT,
					created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
					CONSTRAINT fk_post_like FOREIGN KEY (post_id) REFERENCES fiboblog.posts(id),
					CONSTRAINT fk_user_like FOREIGN KEY (user_id) REFERENCES fiboblog.users(id)
				);`,
		},
		// {
		// 	name: "salaries",
		// 	create: `
		// 		CREATE TABLE IF NOT EXISTS blog.salaries (
		// 			id SERIAL PRIMARY KEY,
		// 			user_id INT NOT NULL,
		// 			month_year DATE NOT NULL,
		// 			reputation_points INT NOT NULL,
		// 			salary_amount NUMERIC(10, 2) NOT NULL,
		// 			created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
		// 			CONSTRAINT fk_user_salary FOREIGN KEY (user_id) REFERENCES blog.users(id)
		// 		);`,
		// },
	}

	for _, table := range tables {
		var tableExists bool
		err := db.QueryRowContext(ctx, `SELECT EXISTS (SELECT 1 FROM pg_catalog.pg_tables WHERE schemaname = 'fiboblog' AND tablename = $1)`, table.name).Scan(&tableExists)
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
