# Simple way to start application localy is:
1) Adding `127.0.0.1 local.staging.game.k8s.atpstealer.com`
to hosts in OS.
2) run `npm install`
3) run `npm run dev`
4) You have app on http://local.staging.game.k8s.atpstealer.com:5173/

If you try to use HTTPS you will receive certificate misconfiguration. <br> 
localhost:3000 also impossible, because browser wants to send cookie to api.staging.game.k8s.atpstealer.com

Prod with SSL and HSTS here https://game.k8s.atpstealer.com/

