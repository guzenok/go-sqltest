//go:generate sqlmockgen -package empty_interface -destination mock.go -source input.go
package empty_interface

type Empty interface{}
