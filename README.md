# template

This project is a project template for [go](https://golang.org/) and [gonew](https://golang.org/x/tools/cmd/gonew).

## prerequisites

- [go](https://golang.org/): latest version
- gonew: `go install golang.org/x/tools/cmd/gonew@latest`

## usage

```sh
gonew github.com/snowmerak/template/<template-name> <your-module-name>
```

or you can use `gg` tool to generate a project from a template.

```sh
go install github.com/snowmerak/template/ggx@latest
```

```sh
# generate go workspace with root module
ggx init

# clone a module from template
ggx clone
```

## templates

- [hello](./hello): a simple hello world program
- [mono](./mono): a monorepo template
- [buf-build](./bufbuild): a buf build template
- [postgres](./postgres): a sqlc postgresql client template
- [redis](./redis): a redis client template
- [httpserver](./httpserver): a http server template
- [httpclient](./httpclient): a http client template
- [executable](./executable): a executable template
- [cassandra](./cassandra): a cassandra client template
- [s3](./s3): a s3 client template
