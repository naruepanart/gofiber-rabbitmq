# gofiber-rabbitmq
 
go run main.go

cd worker && go run main.go

```
version: '3.6'

services:
  rabbitmq:
    image: rabbitmq:3-management-alpine
    container_name: rabbitmq-management-alpine
    restart: on-failure
    environment:
    - RABBITMQ_DEFAULT_USER=rabbitmq
    - RABBITMQ_DEFAULT_PASS=mypassword
    ports:
      - 5672:5672
      - 15672:15672
    volumes:
      - /var/lib/rabbitmq:/data/db/rabbitmq
    networks:
      -  star-network
      
networks:
  star-network:
    external: true
```

## UI

http://localhost:15672