# gonfig

All-in-one solution for managing configuration from environment values and flags.

## Usage

Running the following using the following command would yield `localhost:3000?key=secret`

```bash
go run ./main.go -api-key=secret
```

```go
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/axatol/gonfig"
)

type config struct {
	Host   string `env:"HOST"`
	Port   int    `env:"PORT" default:"3000"`
	ApiKey string `flag:"api-key"`
}

func main() {
	os.Setenv("HOST", "localhost")
	cfg := config{}

	// load config using defaults
	_ = gonfig.Load(&cfg)

	// OR manually run each step

	// initialise configurator
	fs, _ := gonfig.NewConfig(&cfg)

	// load env vars first
	_ = fs.ReadEnv()

	// register cli flags
	_ = fs.BindFlags(flag.CommandLine)

	// process args as usual
	flag.Parse()

	// ensure required fields are set
	_ = fs.Validate()

	fmt.Printf(
		"%s:%d?key=%s\n",
		cfg.Host,
		cfg.Port,
		cfg.ApiKey,
	)
}
```
