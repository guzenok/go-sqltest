package vendor_dep

//go:generate sqlmockgen -package vendor_dep -destination mock.go github.com/guzenok/go-sqltest/sqlmockgen/internal/tests/vendor_dep VendorsDep
//go:generate sqlmockgen -destination source_mock_package/mock.go -source=vendor_dep.go
