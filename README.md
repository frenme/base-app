Template golang app

1) docker network create app_shared_network
2) mkcert in /kong/certs for https://localhost
3) docker compose -f app/docker-compose.${mode_env}.yml build
4) chmod +x run.sh stop.sh
5) cd app/order-service, swag init
6) ./run.sh