# Cheatsheet

## Local

```bash
cp .env.example .env
# dien secret vao .env
go mod download
go mod verify
go test ./...
go test -race ./...
go vet ./...
go run .
```

## Build

```bash
go build -trimpath -ldflags='-s -w' -o bin/warehouse .
docker build -t warehouse-management:local .
```

## Swarm

```bash
docker swarm init
export BACKEND_IMAGE=registry.example.com/warehouse-management:1.0.0
docker compose --env-file .env config
docker stack deploy --with-registry-auth -c docker-compose.yml warehouse
docker stack services warehouse
docker stack ps warehouse
docker service logs -f warehouse_backend
docker stack rm warehouse
```

## Health and API

```bash
curl -i http://localhost/health/live
curl -i http://localhost/health/ready
curl -i http://localhost/swagger/index.html
curl -i -X POST http://localhost/products -H 'Content-Type: application/json' \
  -d '{"code":"P001","name":"Laptop","price":15000000,"stock":20}'
curl -i http://localhost/products/1
curl -G http://localhost/search --data-urlencode 'q=laptop'
curl -i -X POST http://localhost/orders -H 'Content-Type: application/json' \
  -d '{"customer_id":1,"items":[{"product_id":1,"quantity":2}]}'
```

## Runtime inspection

```bash
docker service ps warehouse_backend
docker service logs -f warehouse_backend
docker service logs -f warehouse_traefik
docker stats
nc -zv localhost 80
nc -zv localhost 9090
```

## Dependency checks

```bash
curl -u "$ES_USERNAME:$ES_PASSWORD" "$ES_ADDRESS/_cluster/health"
redis-cli -h "${REDIS_ADDR%:*}" -p "${REDIS_ADDR##*:}" -a "$REDIS_PASSWORD" PING
# Kafka: use kafka-consumer-groups.sh to inspect warehouse-service lag
```

## Expected behavior

- Product GET: Redis hit first; PostgreSQL fallback then cache fill.
- Product create/update: PostgreSQL write, Redis cache update, Elasticsearch index.
- Product delete: PostgreSQL delete, Redis and Elasticsearch cleanup.
- Order create: PostgreSQL write followed by `warehouse.order.created` publish.
- Invalid HTTP/gRPC input: safe `400`/`InvalidArgument`, no internal error details.
