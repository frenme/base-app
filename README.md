Template app

1) docker network create app_shared_network
2) mkcert in /kong/certs for https://localhost
3) docker compose -f app/docker-compose.dev.yml build
4) chmod +x run.sh stop.sh
5) cd app/order-service, swag init
6) ./run.sh


What else can be done:
- Dev and Prod Separation
-- Deploy two separate databases (Dev and Prod) on the same DB instances
-- Try to unify master and replicas into a single entity so that HAProxy can route requests automatically
-- Work on subdomain configuration in Kong
-- Consider which topics to publish to on the dev environment to avoid collisions
- Automatically restart services on failure
- Testing
