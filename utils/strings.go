package utils

import (
	"strings"
)

// Builder Builder
type Builder struct {
	strings.Builder
}

// NewBuilder combine string as Builder
func NewBuilder(ss ...string) *Builder {
	b := &Builder{}
	b.WriteStrings(ss...)
	return b
}

// WriteStrings WriteStrings
func (b *Builder) WriteStrings(ss ...string) (int, error) {
	iLen := 0
	for _, s := range ss {
		i, _ := b.WriteString(s)
		iLen += i
	}
	return iLen, nil
}
