example
======


Codegeneration
------------

Database is needed:

```sh
make db-start
go generate ./...
make db-stop
```


Test
---------------

Database is not needed, [go-sqlmock](https://github.com/DATA-DOG/go-sqlmock) instead:

```sh
go test ./...
```

