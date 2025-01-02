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
var config = struct {
	ES string
}{}

func init() {
	flag.StringVar(&config.ES, "es", "http://es.dev:9200", "")

	flag.Parse()

	tool.TablePrinter(config)
}

type Record struct {
	Id        uint64 `json:"id" gorm:"primaryKey;column:id"`
	CreatedAt int64  `json:"created_at" gorm:"column:created_at;autoCreateTime:milli"`
	Name      string `json:"name" gorm:"column:name"`
}

func main() {
	app := upp.New(upp.Config{Debug: false})

	app.With(upp.InitDB("sqlite://data.db", &Record{}))
	app.With(upp.InitApi(api.New()))

	app.With(upp.InitFn(func(u interfaces.Upp) {
		u.UseLogger().Debug("[init] create init record")
		u.UseDB().Create(&Record{Name: "init"})
	}))

	app.GET("/hello/:name", func(c *api.Ctx) error {
		name := c.Param("name")
		c.UseLogger().Debug("[hello] got name = %s", name)
		record := &Record{Name: name}
		err := c.UseDB().Create(record).Error
		return c.JSON(map[string]any{"record": record, "err": err})
	})

	app.RunSignal()
}


```
> run with env
```sh
DEBUG=true go run . 

DEBUG=true LISTEN_HTTP=0.0.0.0:8080 go run .
```