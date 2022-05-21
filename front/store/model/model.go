package model

import "time"

// Item represents a single TODO item in the store.
type Item struct {
	Title     string
	Completed bool
}

// FilterState represents a viewing filter for TODO items in the store.
type FilterState int

const (
	// All is a FilterState which shows all items.
	All FilterState = iota

	// Active is a FilterState which shows only non-completed items.
	Active

	// Completed is a FilterState which shows only completed items.
	Completed
)

// FIXME Maybe models could be shared between frontend and backend, not sure we'll have to think about it?
type Section struct {
	ID          uint
	Title       string
	Description string
	Placement   int
	BoardId     uint
	Items       []interface{}
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Board struct {
	ID          uint
	Title       string
	Description string
	WorkspaceId uint
	Sections    []Section
	CreatedAt   string `mapstructure:"created_at"`
	UpdatedAt   string `mapstructure:"updated_at"`
}

type Workspace struct {
	ID          uint
	Title       string
	Description string
	Boards      []Board
	CreatedAt   string `mapstructure:"created_at"`
	UpdatedAt   string `mapstructure:"updated_at"`
}

type ListWorkspace struct {
	Workspaces []Workspace
	Count      int
}

type ListSections struct {
	Sections []Workspace
	Count    int
}
