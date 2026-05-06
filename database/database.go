package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"waweb-v2/config"
	"waweb-v2/models"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func Init() {
	cfg := config.LoadConfig()

	var err error
	DB, err = sql.Open("postgres", cfg.DSN())
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	log.Println("Connected to PostgreSQL database")

	if err = migrate(); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}
}

func migrate() error {
	_, err := DB.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id          VARCHAR(36) PRIMARY KEY,
			email       VARCHAR(255) UNIQUE NOT NULL,
			password    VARCHAR(255) NOT NULL,
			name        VARCHAR(255) NOT NULL DEFAULT '',
			created_at  TIMESTAMP NOT NULL DEFAULT NOW()
		);

		CREATE TABLE IF NOT EXISTS workspaces (
			id           VARCHAR(36) PRIMARY KEY,
			user_id      VARCHAR(36) NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			name         VARCHAR(255) NOT NULL,
			template_id  VARCHAR(255) NOT NULL DEFAULT '',
			subdomain    VARCHAR(255) NOT NULL DEFAULT '',
			values       JSONB NOT NULL DEFAULT '{}',
			is_published BOOLEAN NOT NULL DEFAULT FALSE,
			created_at   TIMESTAMP NOT NULL DEFAULT NOW(),
			updated_at   TIMESTAMP NOT NULL DEFAULT NOW()
		);

		CREATE INDEX IF NOT EXISTS idx_workspaces_user_id ON workspaces(user_id);
		CREATE UNIQUE INDEX IF NOT EXISTS idx_workspaces_subdomain ON workspaces(subdomain) WHERE subdomain != '';
	`)
	if err != nil {
		return fmt.Errorf("migrate error: %w", err)
	}
	log.Println("Database migrations applied")
	return nil
}

// ── User functions ──────────────────────────────────────────────────────────

func GetUserByEmail(email string) (*models.User, bool) {
	user := &models.User{}
	err := DB.QueryRow(
		`SELECT id, email, password, name, created_at FROM users WHERE email = $1`, email,
	).Scan(&user.ID, &user.Email, &user.Password, &user.Name, &user.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, false
	}
	if err != nil {
		log.Printf("GetUserByEmail error: %v", err)
		return nil, false
	}
	return user, true
}

func CreateUser(user *models.User) error {
	_, err := DB.Exec(
		`INSERT INTO users (id, email, password, name, created_at) VALUES ($1, $2, $3, $4, $5)`,
		user.ID, user.Email, user.Password, user.Name, user.CreatedAt,
	)
	return err
}

// ── Workspace functions ─────────────────────────────────────────────────────

func GetWorkspacesByUserID(userID string) []*models.Workspace {
	rows, err := DB.Query(
		`SELECT id, user_id, name, template_id, subdomain, values, is_published, created_at, updated_at
		 FROM workspaces WHERE user_id = $1 ORDER BY created_at DESC`, userID,
	)
	if err != nil {
		log.Printf("GetWorkspacesByUserID error: %v", err)
		return nil
	}
	defer rows.Close()

	var result []*models.Workspace
	for rows.Next() {
		ws, err := scanWorkspace(rows)
		if err != nil {
			log.Printf("scan workspace error: %v", err)
			continue
		}
		result = append(result, ws)
	}
	return result
}

func GetWorkspaceByID(id string) (*models.Workspace, bool) {
	row := DB.QueryRow(
		`SELECT id, user_id, name, template_id, subdomain, values, is_published, created_at, updated_at
		 FROM workspaces WHERE id = $1`, id,
	)
	ws, err := scanWorkspaceRow(row)
	if err == sql.ErrNoRows {
		return nil, false
	}
	if err != nil {
		log.Printf("GetWorkspaceByID error: %v", err)
		return nil, false
	}
	return ws, true
}

func CreateWorkspace(ws *models.Workspace) error {
	valuesJSON, err := json.Marshal(ws.Values)
	if err != nil {
		return err
	}
	_, err = DB.Exec(
		`INSERT INTO workspaces (id, user_id, name, template_id, subdomain, values, is_published, created_at, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
		ws.ID, ws.UserID, ws.Name, ws.TemplateID, ws.Subdomain,
		string(valuesJSON), ws.IsPublished, ws.CreatedAt, ws.UpdatedAt,
	)
	return err
}

func UpdateWorkspace(ws *models.Workspace) error {
	valuesJSON, err := json.Marshal(ws.Values)
	if err != nil {
		return err
	}
	_, err = DB.Exec(
		`UPDATE workspaces
		 SET name=$1, template_id=$2, subdomain=$3, values=$4, is_published=$5, updated_at=$6
		 WHERE id=$7`,
		ws.Name, ws.TemplateID, ws.Subdomain,
		string(valuesJSON), ws.IsPublished, ws.UpdatedAt, ws.ID,
	)
	return err
}

func DeleteWorkspace(id string) error {
	_, err := DB.Exec(`DELETE FROM workspaces WHERE id = $1`, id)
	return err
}

// ── helpers ──────────────────────────────────────────────────────────────────

type rowScanner interface {
	Scan(dest ...interface{}) error
}

func scanWorkspace(scanner rowScanner) (*models.Workspace, error) {
	ws := &models.Workspace{}
	var valuesJSON string
	err := scanner.Scan(
		&ws.ID, &ws.UserID, &ws.Name, &ws.TemplateID,
		&ws.Subdomain, &valuesJSON, &ws.IsPublished,
		&ws.CreatedAt, &ws.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal([]byte(valuesJSON), &ws.Values); err != nil {
		ws.Values = make(map[string]interface{})
	}
	return ws, nil
}

func scanWorkspaceRow(row *sql.Row) (*models.Workspace, error) {
	ws := &models.Workspace{}
	var valuesJSON string
	err := row.Scan(
		&ws.ID, &ws.UserID, &ws.Name, &ws.TemplateID,
		&ws.Subdomain, &valuesJSON, &ws.IsPublished,
		&ws.CreatedAt, &ws.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal([]byte(valuesJSON), &ws.Values); err != nil {
		ws.Values = make(map[string]interface{})
	}
	return ws, nil
}
