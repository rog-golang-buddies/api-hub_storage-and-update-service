package dto

// Page represents consistent result from the repository with page data
type Page[T any] struct {
	Data    []T
	Page    int
	PerPage int
	Total   int
}
