# Backend
## About
---
Golang server using the echo framework.

## Generating the Server
---
Development of endpoints must happen first in the [OpenAPI file](./api.yaml)
Once the file is updated, you must generate a new `server.gen.go` using the following command:
```bash
cd api/
go generate .
```

## Development
---
`repo/` - Contains all interactions with the database
`api/` - Contains all API routes
`service/` - Services needed to perform more complex transactions, especially if they span multiple `repos`
