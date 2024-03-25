# SurGo
A simple sqlx-like library for using SurrealDB in Go.

## Installation

Add the library to your go project using the following command:
```bash
go get github.com/NoBypass/surgo
```
**Make sure that your Go project runs on version 1.22 or later!**

## Functions

### Connect

Connect to a database and get a DB object and an error.

Example usage:
```go
db, err := surgo.Connect("127.0.0.1:8000")
```

### MustConnect

Works the same as [Connect](#Connect), but panics if an error occurs.

Example usage:
```go
db := surgo.MustConnect("127.0.0.1:8000")
```

### (Configuration)

To configure the database you can use the following functions:
- `User` - Set the username for the database.
- `Password` - Set the password for the connection.
- `Database` - Set the database to use.
- `Namespace` - Set the namespace to use.

Example usage with all functions:
```go
db, err := surgo.Connect("127.0.0.1:8000",
    surgo.User("user"),
    surgo.Password("password"),
    surgo.Database("database"),
    surgo.Namespace("namespace"))
```

## Types

### Result
The result of a query. It is returned by the [Exec](#Exec) and
[MustExec](#MustExec) functions.

```go 
type Result struct {
	Data     any
	Error    error
	Duration time.Duration
}
```

Data will contain either a `map[string]any` or a `[]map[string]any` depending
on the query. The other fields are self-explanatory.

### DB
DB is the main object used to interact with the database.

#### Close
Close the connection to the database.

#### Scan
Scan the data from the result of the query to a struct or slice of
structs. The first argument is the object to scan the result to, the
second argument is the query, and the third argument is the parameters.

If the given query string contains multiple queries, the last one will
be scanned to the given object. The values from the result are either
assigned to the fields of the given struct as lowercase names, or to their
`db` tags if they are present. If a struct or teh result contains fields
that are not present in the other, they will be ignored.

You can input a single struct or a slice of structs as the first argument.
What matters is that it has to be a pointer to either. Parameters work the
same way as in the [Exec](#Exec) function.

Example usage:
```go
type User struct {
    CreatedAt int `db:"created_at"`
    Name      string
}

var user User
err := db.Scan(&user, "SELECT * FROM users:$1 WHERE", 1)
```

#### Exec
Execute a query and return the result. The parameters can be just
normal values, a map, or a struct. If a map or a struct is used, the
keys or the fields of the map or struct will be used as the names of
the parameters. If the `db` tag is present for some fields it will be
used as the names of the parameters instead.

The parameters can be represented with `$myvar` in the query string,
where `myvar` is the name of the parameter. This works when either
structs or maps are used but if normal values are used, you will have 
to use the `$1`, `$2`, etc. syntax.

Example usage:
```go
result, err := db.Exec("INSERT INTO users (name, age) VALUES ($name, $age)", map[string]any{
    "name": "John",
    "age":  25,
})

result, err := db.Exec("INSERT INTO users (name, age) VALUES ($1, $2)", "John", 25)

type User struct {
    Name string
    Age  int
}
result, err := db.Exec("INSERT INTO users (name, age) VALUES ($name, $age)", User{
	Name: "John",
    Age:  25,
})
```

#### MustExec

Works the same as [Exec](#Exec), but panics if an error occurs.

Example usage:
```go
result := db.MustExec("DEFINE TABLE users")
```

## Contributing

Just make a pull request, and it will be reviewed as soon as possible.

## To-Do

- [ ] Optimize parameter parsing
- [x] Automatically use string syntax for string IDs
- [ ] Document IDs and Ranged IDs
