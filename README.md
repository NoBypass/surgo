# SurGo
A simple sqlx-like library for using SurrealDB in Go.

**Table of Contents**
* [Installation](#installation)
* [To-Do](#to-do)
* [Connecting to a database](#connecting-to-a-database)
  * [Configure the connection](#configure-the-connection)
* [Querying](#querying)
  * [Scanning](#scanning)
  * [Exec](#exec)
  * [MustExec](#mustexec)
  * [Result](#result)
  * [IDs (Records)](#ids-records)
    * [Normal/Array IDs](#normalarray-ids)
    * [Ranged IDs](#ranged-ids)
    * [Using multiple IDs](#using-multiple-ids)
  * [Datetimes and Durations](#datetimes-and-durations)

## Installation
Add the library to your go project using the following command:
```bash
go get github.com/NoBypass/surgo
```
**Make sure that your Go project runs on version 1.22 or later!**

## To-Do
- Fix not parsing back from string to `time.Duration` in scanning
- Allow pointers to structs for parameters
- Allow scanning to a nil pointer
- Improve/Update Docs
- Use SurrealDB variables in ranged IDs
- Improve error messages/errors in general

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

#### Behavior
The `dest` input must have the same structure as the expected return value 
from the query. This means you will have to provide a slice even if you used
something like `LIMIT 1` and expect only one result. If you want to scan to a
single struct, you have to use the `ONLY` keyword in the query (or index the
returned array).

Since the scanning is that strict, you also have to play attention in other
cases like for example when mapping to a single variable. As an example, if
you use count in SurrealDB, the result will be:
```SQL
SELECT count() FROM test GROUP ALL
/* returns:
[
  {
    "count": 2
  }
]
*/
```

What you could do here is to wrap the query in a `RETURN` statement and then
it will work just fine:
```SQL
RETURN (SELECT count() FROM test GROUP ALL)[0].count
-- returns: 2
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

### Datetimes and Durations
You can use the `time.Time` type for datetime values. The library will automatically convert them to the correct format
for SurrealDB. Some goes for the `time.Duration` type. It will be converted to the highest possible unit (e.g. seconds)
that keeps the precision. Both of these types can be used as parameters in query functions. It does not matter if they are used as normal values,
in maps, or in structs. When scanning if the destination field is of type `time.Time` and the source is a string, the library will try to parse
the string to a `time.Time` object.
