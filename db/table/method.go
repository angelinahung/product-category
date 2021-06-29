package table

// Table methods
type Table interface {
	IsBadRequest() bool
	IsRequired() bool
}
