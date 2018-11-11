package utils

type Builder struct {
	arr []rune
	loc int
}

func (b *Builder) Grow(n int) {
	b.arr = make([]rune, 0, n)
	b.loc = 0
}

func (b *Builder) WriteRune(r rune) {
	b.arr = append(b.arr, r)
	b.loc++
}

func (b *Builder) WriteString(s string) {
	for _, c := range s {
		b.WriteRune(c)
	}
}

func (b *Builder) String() string {
	return string(b.arr)
}
