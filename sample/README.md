example
======


Codegeneration
------------

Database is needed

```sh
make db-start
go generate ./...
make db-stop
```


Test
---------------

No database needed:

```sh
go test ./...
```

