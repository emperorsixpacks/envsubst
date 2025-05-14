# envsubst

A lightweight Go package for substituting environment variables in configuration strings, similar to Docker Compose's `${VAR}` syntax.

## Features

- ✅ Replaces `${VAR}` with the value of the corresponding environment variable
- ❌ Does **not yet** support default values (`${VAR:-default}`) or required checks (`${VAR:?error}`) — *coming soon!*

## Installation

```bash
go get github.com/emperorsixpacks/envsubst
````

## Usage

```go
import (
	"fmt"
	"os"

	"github.com/emperorsixpacks/envsubst"
)

type Config struct {
	Name    string   `yaml:"name"`
	Version string   `yaml:"version"`
	Port    []int    `yaml:"port"`
	Network string   `yaml:"network"`
}

func main() {
	os.Setenv("CLIENT_NAME", "lighthouse")
	yml := []byte(`
name: ${CLIENT_NAME}
version: v1.2.3
port:
  - 30303
  - 8545
network: testnet
`)

	var cfg Config
	if err := envsubst.Unmarshal(yml, &cfg); err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", cfg)
}
```

## License

MIT

```


