<br />

<h1 align="center">
<img width=32 style="transform: translateY(6px)" src="https://raw.githubusercontent.com/surrealdb/icons/main/surreal.svg" />
&nbsp; Surgo &nbsp;
<img width=32 style="transform: translateY(6px)" src="https://raw.githubusercontent.com/surrealdb/icons/main/golang.svg" />
</h1>
<p align=center>QOL features and sqlx-like mappings for <code><a href="https://github.com/surrealdb/surrealdb.go">github.com/surrealdb/surrealdb.go</a></code></p>

<br />

<h2 align=center>Features</h2>
<p align=center><b>Features over the original library:</b></p>
<br />
<p align="center">Simplified database connection</p>
<p align="center">Ability to directly scan the result into a struct using sqlx-like syntax</p>
<p align="center">A consistent <code>Result</code> type instead of using <code>interface{}</code></p>
<p align="center">Consistent error handling</p>
<p align="center">Up-to-date Documentation</p>
<p align="center">Support for struct tags</p>
<p align="center">Support for `time.Duration` and `time.Time`</p>

## Installation
```bash
go get github.com/NoBypass/surgo
```
> Make sure that your Go project runs on version 1.22 or later!

## Roadmap
- Create a mapping function for each one which exists in the original surrealdb.go library (e.g. `Use`, `Merge`, etc.) but with the scan functionality like in the `Scan` function.
- Support `LiveNotifications` with the `Result` struct.
- Eventually remove the extra layer of the surrealdb.go library and directly use the websocket code.
- Keep the library up-to-date with the original library.
- (Maybe) Make it so that the library can be used with the `database/sql` package.

## Documentation

### Connecting to the Database

```go
db, err := surgo.Connect("ws://localhost:8000", &surgo.Credentials{
    // any of these fields can be omitted if they are not needed
    Namespace: "test",
    Database:  "default",
    Username: "admin",
    Password: "1234", 
	Scope:    "public",
})
```

```go
// panics if connection could not be established
db := surgo.MustConnect("ws://localhost:8000", &surgo.Credentials{
    Username: "admin",
    Password: "1234",
})

defer db.Close()
```

### Querying the Database

```go
result, err := db.Query("SELECT * FROM ONLY $john", map[string]any{
	"john": "users:john",
})
```

The result of the query above will be of type `[]surgo.Result`. The important fields of the `Result` struct are:
- `Data` which is of type `any` containing the data of the query (with a simple query like above a `map[string]any`).
- `Error` which is of type `error` containing the error of the query (if there is one).

If you want to scan the result from such a query into a struct, you can use the `Result.Unmarshal` function:

```go
type User struct {
    Name string
}

var john User
result[0].Unmarshal(&john)
```

### Directly Scanning the Result

```go
var john User
err := db.Scan(&john, "SELECT * FROM ONLY $john", map[string]any{
    "john": "users:john",
})
```

### Struct Tags
Struct tags essentially work the same way as in the `json` package. A full example would look like this:

```go
type User struct {
	// this field will be omitted if it is empty, otherwise it will be mapped to the "name" key
    Name string `db:"name,omitempty"`
	// this field will be ignored
    Age  int    `db:"-"`
    // this field will be mapped to the "Address" key
	Address string
}

_, err := db.Query("CREATE $john CONTENT $data", map[string]any{
    "john": "users:john",
	"data": User{
        Name: "John Doe",
        Age:  42,
    },
})
```

If your `db` tags would often end up being the same as other tags (for example the `json` tag) you can specify a
fallback tag using the `SURGO_FALLBACK_TAG` environment variable. If this variable is set, the library will use the
specified tag as a fallback if no `db` tag is present.

### Time and Duration
The library supports `time.Time` and `time.Duration` out of the box. They are automatically converted to and from the
formats which SurrealDB uses. (`datetime` and `duration` respectively) Here is an example:

```go
type User struct {
    Name string
    CreatedAt time.Time
    OneMilePB time.Duration
}

_, = db.MustQuery("CREATE $john CONTENT $data", map[string]any{
    "john": "users:john",
    "data": User{
        Name:      "John Doe",
        CreatedAt: time.Now(),
		OneMilePB: time.Minute * 6,
    },
})

var john User
err := db.Scan(&john, "SELECT * FROM ONLY $john", map[string]any{
    "john": "users:john",
})
```
