FRONTEND

frontend/vue/.env example

```
VITE_DOMAIN=localhost:5173
VITE_API=http://localhost:8000/api/v2
#VITE_API=http://staging.game.k8s.atpstealer.com/api/v2
VITE_MINIO_URL=https://minio.game.k8s.atpstealer.com/game-staging/assets
```

BACKEND

deploymnet/docker/api.env example

```
GOOGLE_API={"type":"service_account","project_id":"XXXXXXXX","private_key_id":"XXXXXXX","private_key":"-----BEGIN PRIVATE KEY-----\nXXXXXXXXX\n-----END PRIVATE KEY-----\n","client_email":"XXXXXXXX","client_id": "XXXXXXXXX","auth_uri": "https://accounts.google.com/o/oauth2/auth","token_uri":"https://oauth2.googleapis.com/token","auth_provider_x509_cert_url":"https://www.googleapis.com/oauth2/v1/certs","client_x509_cert_url":"XXXXXXXXXXXX","universe_domain":"googleapis.com"}
GOOGLE_SHEET_ID=18DblwMx9a-YnLh7hfqwSC5a3Sc-_bRybMWMWp-isgO0
INIT=true
MONGO_DATABASE=game-staging
MONGO_HOST=mongodb
MONGO_PASSWORD=MONGO_ROOT_PASSWORD
MONGO_PORT=27017
MONGO_USER=MONGO_ROOT_USERNAME
APP_PORT=8001
```

RUN

```
cd /deployment/docker-compose
docker compose up --build -d
```