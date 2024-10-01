<br />

<h1 align="center">
<img width=32 style="transform: translateY(6px)" src="https://raw.githubusercontent.com/surrealdb/icons/main/surreal.svg" />
&nbsp; Surgo &nbsp;
<img width=32 style="transform: translateY(6px)" src="https://raw.githubusercontent.com/surrealdb/icons/main/golang.svg" />
</h1>
<p align=center>QOL features and sqlx-like mappings for <code><a href="https://github.com/surrealdb/surrealdb">github.com/surrealdb/surrealdb</a></code></p>

<br />

<h2 align=center>New Features</h2>
<p align="center">Simplified database connection</p>
<p align="center">Ability to directly scan the result into a struct using sqlx-like syntax</p>
<p align="center">A consistent <code>Result</code> type instead of using <code>interface{}</code></p>
<p align="center">Consistent error handling</p>
<p align="center">Up-to-date Documentation</p>
<p align="center">Support for struct tags</p>
<p align="center">Support for <code>time.Duration</code> and <code>time.Time</code></p>
<p align="center">Support for <code>context.Context</code></p>
<p align="center">Supports tracing</p>
<br>
<p align="center"><b>Live Notifications are not yet supported</b></p>

## Installation
```bash
go get github.com/NoBypass/surgo
```
> Make sure that your Go project runs on version 1.23 or later!

## Documentation

### Connecting to the Database

```go
db, err := surgo.Connect("ws://localhost:8000", &surgo.Credentials{
    // any of these fields can be omitted if they are not needed
    Namespace: "test",
    Database:  "default",
    Username: "admin",
    Password: "1234", 
	Scope:    "myScope",
})
```
**Important:**
Namespace, Database and Scope are optional, if not provided the signin will happen on the Root, Namespace or Database level respectively.

There are a couple of available options which you can pass to the `Connect` function:
- `WithDefaultTimeout`: The default timeout value is 10 seconds. You can use this option to change it or pass contexts to queries to individually specify a timeout.
- `WithLogger`: Use a custom logger/tracer. More about this in the [Tracing](#tracing) section.
- `WithDisableLogging`: Disable logging.
- `WithFallbackTag`: Use a fallback tag for struct tags. More about this in the [Fallback Tag](#fallback-tag) section.

### Querying the Database

Example:

```go
result := db.Query("SELECT * FROM ONLY $john", map[string]any{
	"john": "users:john",
})

// get the first result. If you only sent one query using this function makes the most sense.
resp, err := result.First()

// if you sent multiple queries you can have the following options:
resp, err := result.Last()
resp, err := result.At(0)

// or you can iterate over the results
for resp, err := range result.Iter() {
    // do something with res
}
```

#### Unmarshal

If you want to scan the result from such a query into a struct, you can use the `Marshaler.Unmarshal` function:

```go
resp, err := result.First()
if err != nil {
    // handle error
}

var john User
if err := db.Marshaler.Unmarshal(&john, resp); err != nil {
    // handle error
}
```

#### Scan Directly

If you want to directly scan the result into a struct you can use the `Scan` function:
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
    // this field will be mapped to the default "Sleep" key and parsed to SurrealDB's duration format
	Sleep time.Duration
    // this field will be mapped to the "birthday" key and parsed to SurrealDB's datetime format
	Birthday time.Time `db:"birthday"`
}

err := db.Query("CREATE $john CONTENT $data", map[string]any{
    "john": "users:john",
	"data": User{
        Name: "",
        Age:  42,
        Sleep: time.Hour * 8 + time.Minute * 30, 
		Birthday: time.Date(1980, 1, 1, 0, 0, 0, 0, time.UTC),
    },
}).Error
if err != nil {
    // handle error
}
```

The query above will result in the following entry in SurrealDB:

```
{
    id: 'someRandomString', 
    Sleep: '8h30m',
    birthday: '1980-01-01T00:00:00Z'
}
```

Unmarshal and Scan functions will automatically convert the SurrealDB formats back to the Go types.

#### Fallback Tag
If you don't like using the `db` tag, or your struct already uses it for something else, you can use the `fallback` tag.
For example if most of your structs use the `json` tag, you can set the fallback tag to `json`. This way for the fields
which don't have a `db` tag, the library will look for a `json` tag. Here is an example:

### Tracing & Context
You can use the `WithLogger` option to pass a custom logger/tracer to the `Connect` function. The logger/tracer must 
implement the `surgo.Logger` interface. Here is an example of a simple logger:

```go
type Logger struct{}

func (l *Logger) Error(err error) {
    log.Println(err)
}

func (l *Logger) Trace(ctx context.Context, t TraceType, data any) {
    l.Printf("trace: %v | %v | %v\n", t, ctx.Value("traceID"), data)
}
```

So the value from the context works you will have to pass your own context to the query you called. Here is an example:

```go
ctx := context.WithValue(context.Background(), "traceID", "yourTraceID")
result := db.WithContext(ctx).Query("SELECT * FROM ONLY $john", map[string]any{
    "john": "users:john",
})
```

Of course, you can also use a cancelable context or timeout context.
