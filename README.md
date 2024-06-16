# template

This project is a project template for [go](https://golang.org/) and [gonew](https://golang.org/x/tools/cmd/gonew).

## prerequisites

- [go](https://golang.org/): latest version
- gonew: `go install golang.org/x/tools/cmd/gonew@latest`

## usage

```sh
gonew github.com/snowmerak/template/<template-name> <your-module-name>
```

## templates

- [hello](./hello/README.md): a simple hello world program
- [mono](./mono/README.md): a monorepo template
- [buf-build](./bufbuild/README.md): a buf build template
- [postgres](./postgres/README.md): a sqlc postgresql client template
- [redis](./redis/README.md): a redis client template