# Frontend
## Installing
```bash
npm install
```

## Getting Started
The recommended approach is to use `devenv up` [from the root directory](../README.md) 

First, run the development server:

```bash
npm run dev
```

Open [http://localhost:3000](http://localhost:3000) with your browser to see the result.

## Development
Development of the backend is reflected in the [OpenAPI file](./../backend/api.yaml)

When updates are made, you must generate a new client file to get all of the typesafety features

```bash
npx openapi-typescript ../backend/api.yaml --output src/lib/contacts.ts
```
