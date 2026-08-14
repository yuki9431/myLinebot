logger
====
## Description
Write log in custom output

## Requirement
- Go 1.10 or later


## Install
```bash:#
go get github.com/yuki9431/logger
```

## Configuration
```go:main.go
import (
	"github.com/yuki9431/logger"
)

func main() {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	l := logger.New(file)
	l.Write("Hello World")
	...
}
```

## Contribution
1. Fork ([https://github.com/yuki9431/logger](https://github.com/yuki9431/logger))
2. Create a feature branch
3. Commit your changes
4. Rebase your local changes against the master branch
5. Create new Pull Request


## Author
[Dillen H. Tomida](https://github.com/yuki9431)
