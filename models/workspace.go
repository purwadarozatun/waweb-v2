package models

import "time"

type Workspace struct {
	ID          string                 `json:"id"`
	UserID      string                 `json:"user_id"`
	Name        string                 `json:"name"`
	TemplateID  string                 `json:"template_id"`
	Subdomain   string                 `json:"subdomain"`
	Values      map[string]interface{} `json:"values"`
	IsPublished bool                   `json:"is_published"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}

type Template struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Thumbnail   string                 `json:"thumbnail"`
	Spec        map[string]interface{} `json:"spec"`
}
