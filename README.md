sqltest [![Build Status][travis-ci-badge]][travis-ci]
======

SqlTest is a test code generator for the [Go programming language][golang].
It runs your tests on real db, records sql-traffic into [sqlmock][sqlmock]
and makes your tests work without real db.


Installation
------------

Once you have [installed Go][golang-install], run these commands
to install the `sqlmockgen` tool:

    go get github.com/guzenok/go-sqltest
    go install github.com/guzenok/go-sqltest/sqlmockgen


Running sqlmockgen
---------------

The `sqlmockgen` command is used to generate offline tests code according to you [special named](#Test-function-naming-agreement) test functions.
It takes 1 argument:
	
 * absolute or relative import path of package to generate offline tests for;
	
and supports the following flags:

 *  `-out`: a file to which to write the resulting source code;

 *  `-db`: a connection string to real db;

 *  `-copyright`: copyright file used to add copyright header to the resulting source code;

Example for installed:

```sh
sqlmockgen -out=sql_test.go -db=postgresql://postgres:postgres@localhost:5432/test?sslmode=disable .
```

Example for gotten:

```sh
go run github.com/guzenok/go-sqltest/sqlmockgen -out=sql_test.go -db=postgresql://postgres:postgres@localhost:5432/test?sslmode=disable .
```
	
Example for go generate:

```go
//go:generate go run github.com/guzenok/go-sqltest/sqlmockgen -out=sql_test.go -db=postgresql://postgres:postgres@localhost:5432/test?sslmode=disable .
```

For an example of the `sqlmockgen` using, see the [sample/](./sample) directory.


Test function naming agreement
--------------

Your function for test db initialization (migrations, data fixtures, etc.) should be:

```go
func initTestDb(dbUrl string) (*sql.DB, error) {
  // open connection and prepare data
}

```

Your test functions should be like:

```go
func test<TESTNAME>(*testing.T, *sql.DB) {
  // test your code with db connection
}
```
(with different \<TESTNAME\>)

Then sqlmockgen will generate go tests:

```go
func Test<TESTNAME>(*testing.T) {
  db, err := test<TESTNAME>SqlMock()
  if err != nil {
    t.Fatal(err)
  }
  test<TESTNAME>(t, db)
}

func test<TESTNAME>SqlMock() (*sql.DB, error) {
    // generated code here:
    // github.com/DATA-DOG/go-sqlmock initialization.
}
```


[sqlmock]:         https://github.com/DATA-DOG/go-sqlmock
[golang]:          http://golang.org/
[golang-install]:  http://golang.org/doc/install.html#releases
[travis-ci-badge]: https://travis-ci.org/guzenok/go-sqltest.svg?branch=master
[travis-ci]:       https://travis-ci.org/guzenok/go-sqltest
