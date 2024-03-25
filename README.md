# SurGo
A simple sqlx-like library for using SurrealDB in Go.

## Table of Contents
* [Table of Contents](#table-of-contents)
* [Installation](#installation)
* [Connecting to a database](#connecting-to-a-database)
  * [Configure the connection](#configure-the-connection)
  * [Query agent](#query-agent)
* [Querying](#querying)
  * [Scanning](#scanning)
  * [Exec](#exec)
  * [MustExec](#mustexec)
  * [Result](#result)
  * [IDs (Records)](#ids-records)
    * [Normal/Array IDs](#normalarray-ids)
    * [Ranged IDs](#ranged-ids)
    * [Using multiple IDs](#using-multiple-ids)
* [Contributing](#contributing)
* [To-Do](#to-do)

## Installation
Add the library to your go project using the following command:
```bash
go get github.com/NoBypass/surgo
```
**Make sure that your Go project runs on version 1.22 or later!**

## Connecting to a database
Connect to a database and get a DB object and an error.

Example usage:
```go
db, err := surgo.Connect("127.0.0.1:8000")
```

`MustConnect` works the same as [Connect](#Connect), but panics if an error occurs.

Example usage:
```go
db := surgo.MustConnect("127.0.0.1:8000")
```

### Configure the connection
To configure the database you can use the following functions:
- `User` - Set the username for the database.
- `Password` - Set the password for the connection.
- `Database` - Set the database to use.
- `Namespace` - Set the namespace to use.
- `CustomAgent` - Set a custom query agent to use.

Example usage with all functions:
```go
db, err := surgo.Connect("127.0.0.1:8000",
    surgo.User("user"),
    surgo.Password("password"),
    surgo.Database("database"),
    surgo.Namespace("namespace"))
```

### Query agent
The idea of the query agent is to allow the user to define their own way to query tha database. Query agent is an
interface that has to be implemented in order to use it. The interface is as follows:
```go
type QueryAgent interface {
    Query(query string, params map[string]any) (map[string]any, error)
    Close() error
}
```
Setting a custom QueryAgent is useful when you want to use a different way to query the database or implement for
example a custom caching or tracing system.

## Querying

### Scanning
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

### Exec
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

### MustExec
Works the same as [Exec](#Exec), but panics if an error occurs.

Example usage:
```go
result := db.MustExec("DEFINE TABLE users")
```

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

### IDs (Records)
IDs are used to reference the data in the database. They are a big part of SurrealDB. There are multiple ways these
can be used. You can have just a normal ID and an array ID which contains multiple values that make up the ID. With
the latter, you can also use [ranged IDs](https://surrealdb.com/docs/surrealdb/surrealql/datamodel/ids/#record-ranges) 
while querying the database.

#### Normal/Array IDs
Normal IDs can simply be used like this:
```go
db.Exec("SELECT * FROM foo:$", surgo.ID{"myid"})      // single value
                                                      // will be parsed to: SELECT * FROM foo:`myid`

db.Exec("SELECT * FROM foo:$", surgo.ID{"myid", 123}) // multiple values
                                                      // will be parsed to: SELECT * FROM foo:['myid', 123]
```

Notice that the ID was represented by a single `$` in the query string. This is done to separate the IDs from the other
parameters. You can still mix them with other parameters.

#### Ranged IDs
Ranged IDs can be used like this:
```go
db.Exec("SELECT * FROM foo:$", surgo.Range{surgo.ID{"myid", 123}, surgo.ID{"myid", 456}})
// will be parsed to: SELECT * FROM foo:['myid', 123]..['myid', 456]
``` 

#### Using multiple IDs
You can use multiple IDs in a single query, but you have to keep the same order in the parameters as in the query string.
While you can still mix them with other parameters, you have to play attention to their order when mixing them with
normal parameters (like `$1`, `$2` etc.) as their index might get mixed up.
```go
db.Exec("RELATE foo:$->edge->bar:$", surgo.ID{123}, surgo.ID{456}
// will be parsed to: RELATE foo:123->edge->bar:456
```

## Contributing
Just make a pull request, and it will be reviewed as soon as possible.

## To-Do
- [ ] Optimize parameter parsing
- [x] Automatically use string syntax for string IDs
- [x] Document IDs and Ranged IDs
