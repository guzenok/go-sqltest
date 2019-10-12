//go:generate sqlmockgen -destination ../source_mock.go -source=source.go
//go:generate sqlmockgen -package source -destination source_mock.go -source=source.go
package source

type X struct{}

type S interface {
	F(X)
}
