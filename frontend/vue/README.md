# Simple way to start application localy is:
1) Adding `127.0.0.1 local.staging.game.kube.atpstealer.com`
to hosts in OS.
2) run `npm install`
3) run `npm run dev`
4) You have app on http://local.staging.game.kube.atpstealer.com:5173/

If you try to use HTTPS you will receive certificate misconfiguration. <br> 
localhost:3000 also impossible, because browser wants to send cookie to staging.game.kube.atpstealer.com/api

# Swagger API generate

back

```aiignore
swag init --parseDependency --parseInternal
```

front
```aiignore
npx @hey-api/openapi-ts --useOptions -i ../../backend/go/docs/swagger.json -o ./src/api -c @hey-api/client-fetch
```

Replace client in services.gen.ts

```aiignore
export const client = createClient(createConfig({
  baseUrl: `${import.meta.env.VITE_API}`,
  credentials: 'include'
}))
```
