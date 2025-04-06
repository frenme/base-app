Template app

1) docker network create app_shared_network
2) mkcert in /kong/certs for https://localhost
3) docker compose -f app/docker-compose.dev.yml build
4) chmod +x run.sh stop.sh
5) cd app/order-service, swag init
6) ./run.sh


ЧТО ЕЩЕ МОЖНО СДЕЛАТЬ
- дев и прод разделение 
  - поднять две разные БД в тех же инстансах БД
  - попробовать унифицировать мастер и реплики в единую сущность чтоб haproxy сам понимал куда направить запрос
  - поработать над субдоменами в Kong
  - подумать в какие топики писать на дев стенде, иначе может быть пересечение
- поднимать сервисы автоматически при ошибках
- тестирование
