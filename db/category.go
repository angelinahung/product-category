package db

// Category definition.
type Category struct {
	ID          int64  `json:"ID"`
	Name        string `json:"name"`
	IsInvisible int16  `json:"is_invisible"`
}
