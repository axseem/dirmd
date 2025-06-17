<div align="center">
    <img alt="dirmd logo" src="./dirmd.svg">
    <p><i>Bundle directories into a single Markdown file</i></p>
</div>

# What is `dirmd`?

`dirmd` is a CLI tool that bundles multiple files into a single, well-formatted Markdown file. It recursively scans a directory, omitting files specified in your `.gitignore` or custom ignore rules.

# Why?

I find that using LLMs in a browser, rather than an IDE, introduces a healthy amount of friction. This simple step helps prevent the overuse of code generation, leading to a workflow that feels more productive and enjoyable.

This approach is especially useful with web-based tools like Google AI Studio, which allows you to use state-of-the-art Gemini models for free.

# Showcase

![demo](./demo.gif)

## Output Example 

- `hello/`
  - `go.mod`
  - `main.go`
  - `name/`
    - `name.go`


`go.mod`
```go
module hello

go 1.24.3
```

`main.go`
```go
package main

import (
	"fmt"
	"hello/name"
)

func main() {
	fmt.Printf("hello %s!", name.SystemUser())
}
```

`name/name.go`
```go
package name

import (
	"os/user"
)

func SystemUser() string {
	systemUser, _ := user.Current()
	return systemUser.Name
}
```

# Installation

```sh
# Nix
nix profile install github:axseem/dirmd

# Go
go install github.com/axseem/dirmd@latest
```

You can also run `dirmd` directly from the GitHub repository without a permanent installation:

```bash
nix run github:axseem/dirmd
```