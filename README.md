# vancone-web-common-go

## Quick Start

Step 1: Add component
```bash
go get github.com/vancone/vancone-web-common-go
```

Step 2: Initialize config

```go
package main

import (
	commonconf "github.com/vancone/vancone-web-common-go/pkg/config"
)

func main() {
	viperConfig := commonconf.ReadConfig()
	// ...
}
```



## Config

Common Config Map YAML

```yaml
app:
  name:
database:
  url:
  username:
  password:
mail:
  protocol:
  host:
  port:
  username:
  password:
redis:
  url:
  password:
server:
  port:
```





## Response Model

Response & ResponsePage



## Utils

