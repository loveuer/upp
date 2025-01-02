# UPP - your app

### Usage

> usage present
```go
app := upp.New()

app.With(db, es, api)

app.Run(ctx)
```

> simple example
```go
type Record struct {
	Id        uint64 `json:"id" gorm:"primaryKey;column:id"`
	CreatedAt int64  `json:"created_at" gorm:"column:created_at;autoCreateTime:milli"`
	Name      string `json:"name" gorm:"column:name"`
}

func main() {
	app := upp.New()

	app.With(upp.InitDB("sqlite://data.db", &Record{}))
	app.With(upp.InitApi(api.New()))

	app.GET("/hello/:name", func(c *api.Ctx) error {
		name := c.Param("name")
		c.UseLogger().Info("[hello] got name = %s", name)
		record := &Record{Name: name}
		err := c.UseDB().Create(record).Error
		return c.JSON(map[string]any{"record": record, "err": err})
	})

	app.RunSignal()
}

```
> run with flags
```sh
go run . --debug --listen.http '0.0.0.0:8080'
```