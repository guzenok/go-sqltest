//go:generate sqlmockgen -destination subdir/internal/pkg/reflect_output/mock.go github.com/guzenok/go-sqltest/sqlmockgen/internal/tests/internal_pkg/subdir/internal/pkg Intf
//go:generate sqlmockgen -source subdir/internal/pkg/input.go -destination subdir/internal/pkg/source_output/mock.go
package test
