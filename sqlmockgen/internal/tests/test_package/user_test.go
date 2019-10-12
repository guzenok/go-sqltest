package users_test

//go:generate sqlmockgen --source=user_test.go --destination=mock_test.go --package=users_test

type User struct {
	Name string
}

type Finder interface {
	FindUser(name string) User
	Add(u User)
}
