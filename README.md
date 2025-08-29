Template app

1) docker network create app_shared_network
2) setup mkcert in /kong/certs for https://localhost
brew install mkcert
mkcert -install
cd kong/certs
mkcert localhost 127.0.0.1 ::1
3) docker compose -f app/docker-compose.dev.yml build
4) cd app/{service_name}, swag init --parseDependency --parseInternal
5) chmod +x run.sh stop.sh
6) ./run.sh
7) Check
https://localhost:8443/api/user-service/swagger/index.html
http://localhost:8000/api/user-service/swagger/index.html


What else can be done:
- **Dev and Prod Separation**
    - Deploy two separate databases (Dev and Prod) on the same DB instances
    - Work on subdomain configuration in Kong
- **Automatically restart services on failure (k8s)**
- **Testing**
