# Go engine binaries repository

This repo contains the Prisma query engine as .go files. It uses build constraints to only include the required
engine in the final go binary.

## Fetch CLI

To fetch new binaries, run

```shell script
go run ./fetch/
```

To adapt the version, modify the `EngineVersion` variable in `fetch/binaries.go`.
