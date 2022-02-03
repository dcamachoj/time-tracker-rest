package dbx

import (
	"fmt"
	"strings"
)

type Builder struct {
	sb strings.Builder
}

// NewBuider constructor
func NewBuilder() *Builder {
	return &Builder{}
}

func (b *Builder) Add(args ...interface{}) *Builder {
	b.sb.WriteString(fmt.Sprint(args...))
	return b
}
func (b *Builder) Addf(format string, args ...interface{}) *Builder {
	b.sb.WriteString(fmt.Sprintf(format, args...))
	return b
}
func (b *Builder) Addln(args ...interface{}) *Builder {
	b.sb.WriteString(fmt.Sprintln(args...))
	return b
}
func (b *Builder) String() string {
	return b.sb.String()
}
