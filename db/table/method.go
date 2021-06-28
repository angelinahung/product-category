package table

type Table interface {
	IsBadRequest() bool
	IsRquired() bool
}
