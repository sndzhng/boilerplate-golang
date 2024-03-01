# Development Setup

#### Required files
`./env/$ENVIRONMENT`
`./cloud-storage-credential.json`

#### Start database:
```bash
docker compose up
```

#### Start service:
```bash
go run cmd/main.go $ENVIRONMENT
```
or with air live reloading
```base
air
```

#### Generate mocks (reflect mode):
```bash
go generate ./...
```

#### Run tests:
```bash
go test ./...
```
