# Warehouse Management Microservice

## 1. Kien truc

- `domain`: entity, DTO, event contract.
- `repository`: interface cho PostgreSQL, Redis, Elasticsearch, Kafka.
- `usecase`: business rules, context timeout, cache-aside va publish event.
- `delivery/http`: Gin REST, validation, Swagger annotations.
- `delivery/grpc`: `CheckStock`.
- `config`: environment-driven dependency wiring.

Luon di theo luong `delivery -> usecase -> repository interface -> adapter`. Handler khong truy cap GORM, Redis, Elasticsearch hoac Kafka truc tiep.

## 2. Cai dat thu cong

Yeu cau:

- Go 1.26+
- PostgreSQL 15+
- Redis 7+
- Elasticsearch 8+
- Kafka 3.7+
- Docker Engine 24+ (neu dung Swarm)
- `protoc` va plugin Go neu muon sinh protobuf client/server chuan

Tao file moi tu mau:

```bash
cp .env.example .env
```

Dat gia tri that cho `DB_PASSWORD`, `REDIS_PASSWORD`, `ES_PASSWORD`. Khong commit `.env` va khong ghi secret vao source.

Chay cac dependency bang Docker rieng, sau do dat host trong `.env`:

```bash
docker run -d --name postgres -e POSTGRES_USER="$DB_USER" \
  -e POSTGRES_PASSWORD="$DB_PASSWORD" -e POSTGRES_DB="$DB_NAME" \
  -p 5432:5432 postgres:15-alpine

docker run -d --name redis -p 6379:6379 redis:7-alpine \
  redis-server --requirepass "$REDIS_PASSWORD"
```

Elasticsearch va Kafka can cau hinh authentication/TLS phu hop voi moi truong. Cap nhat `ES_ADDRESS`, `ES_USERNAME`, `ES_PASSWORD`, `KAFKA_BROKERS`, `KAFKA_TLS`.

Cai dependency va build:

```bash
go mod download
go mod verify
go test ./...
go vet ./...
go build -trimpath -ldflags='-s -w' -o bin/warehouse .
```

## 3. Chay local

```bash
go run .
```

- REST va Swagger: `http://localhost:8080`
- gRPC: `localhost:9090`
- Swagger UI: `http://localhost:8080/swagger/index.html`
- Liveness: `http://localhost:8080/health/live`
- Readiness: `http://localhost:8080/health/ready`

App dung timeout cho DB, Redis, Elasticsearch, Kafka va gRPC business call. Graceful shutdown duoc kich hoat boi `SIGINT`/`SIGTERM`.

## 4. Docker build va Docker Compose local

```bash
docker build -t warehouse-management:local .
docker compose --env-file .env config
```

`docker-compose.yml` hien tai la stack cho Docker Swarm, khong phai bo local standalone. Khong dung `docker compose up -d` voi file nay; lenh do se thieu cac bien bat buoc nhu `BACKEND_IMAGE` va cac password. `docker compose config` chi dung de kiem tra interpolation.

De chay local nhanh, dung Docker rieng cho tung dependency theo phan 2, hoac deploy Swarm theo phan 5. Swarm khong ho tro `depends_on.condition`; backend can duoc theo doi qua health endpoint va service logs sau deploy.

## 5. Docker Swarm + Traefik

Khoi tao Swarm tren manager:

```bash
docker swarm init
```

Build va push image len registry ma cac node truy cap duoc:

```bash
docker build -t registry.example.com/warehouse-management:1.0.0 .
docker push registry.example.com/warehouse-management:1.0.0
export BACKEND_IMAGE=registry.example.com/warehouse-management:1.0.0
```

Nap secrets vao shell hoac dung secret manager cua CI/CD:

```bash
export DB_USER=warehouse
export DB_PASSWORD='replace-with-strong-secret'
export DB_NAME=warehouse
export REDIS_PASSWORD='replace-with-strong-secret'
export ES_USERNAME=elastic
export ES_PASSWORD='replace-with-strong-secret'
export ES_ADDRESS=https://elasticsearch:9200
export REDIS_ADDR=redis-cache:6379
export KAFKA_BROKERS=kafka:9092
export KAFKA_TLS=true
export APP_ENV=production
export GIN_MODE=release
```

Kiem tra interpolation truoc khi deploy:

```bash
docker compose --env-file .env config
```

Deploy:

```bash
docker stack deploy --with-registry-auth \
  --compose-file docker-compose.yml warehouse

docker stack services warehouse
docker stack ps warehouse
```

Traefik route:

- HTTP: port `80` -> backend port `8080`
- gRPC TCP: port `9090` -> backend port `9090`

Cap nhat image:

```bash
docker service update --image registry.example.com/warehouse-management:1.0.1 \
  warehouse_backend
docker stack rm warehouse
```

Khong public PostgreSQL, Redis, Elasticsearch hoac Kafka ra Internet. Trong stack, cac service chi nam tren overlay network `warehouse`.

## 6. Kiem tra API

Tao product:

```bash
curl -i -X POST http://localhost:8080/products \
  -H 'Content-Type: application/json' \
  -d '{"code":"P001","name":"Laptop","price":15000000,"stock":20}'
```

Doc product va kiem tra cache-aside:

```bash
curl -i http://localhost:8080/products/1
redis-cli -h "${REDIS_ADDR%:*}" -p "${REDIS_ADDR##*:}" -a "$REDIS_PASSWORD" GET product:1
```

Tim kiem Elasticsearch:

```bash
curl -G http://localhost:8080/search --data-urlencode 'q=laptop'
```

Tao order, sau do kiem tra Kafka `warehouse.order.created`:

```bash
curl -i -X POST http://localhost:8080/orders \
  -H 'Content-Type: application/json' \
  -d '{"customer_id":1,"items":[{"product_id":1,"quantity":2}]}'
```

Kiem tra loi input:

```bash
curl -i -X POST http://localhost:8080/products \
  -H 'Content-Type: application/json' -d '{"price":-1}'
```

Ky vong `400`, khong co stack trace hoac loi DB noi bo trong response.

## 7. gRPC CheckStock

Contract nam tai `proto/warehouse.proto`. Neu da cai protoc:

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
protoc --go_out=. --go-grpc_out=. proto/warehouse.proto
```

Server hien tai dang dung service descriptor va JSON codec de chay duoc khi workspace chua co `protoc`; client phai dung codec tuong ung. Khi production can protobuf wire format, generate code tu proto va thay delivery adapter bang generated service.

## 8. Test, security va monitoring

Lenh CI toi thieu:

```bash
go test ./...
go test -race ./...
go vet ./...
go build -trimpath -ldflags='-s -w' ./...
docker build --no-cache -t warehouse-management:ci .
```

Kiem tra:

- `GET /health/live`: process con song.
- `GET /health/ready`: endpoint readiness cho load balancer.
- `docker service ps warehouse_backend`: replicas va restart.
- `docker service logs -f warehouse_backend`: loi runtime, khong log password/token.
- `docker service logs -f warehouse_traefik`: routing va HTTP status.
- `docker stats`: CPU/RAM.
- PostgreSQL slow query log, Redis memory/eviction, Elasticsearch cluster health, Kafka consumer lag can duoc dua vao Prometheus/Grafana o tang van hanh.

Security checklist:

- Secret chi tu environment/secret manager.
- TLS cho client-to-edge, Elasticsearch, Redis va Kafka trong production.
- Khong dung `KEYS *` hoac `FLUSHALL` tren Redis production.
- Khong cho phep public management ports.
- Rotation credentials dinh ky.
- Alert khi HTTP 5xx, gRPC errors, consumer lag, DB connection exhaustion, Redis evictions va ES cluster khong green.
