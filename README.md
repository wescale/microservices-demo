
# Microservices demo

This is a sample application that demonstrates a microservice architecture.

## Architecture

![Architecture diagram](./microservice-demo.svg)

The application is composed of the following microservices:
- **[article-service](./article-service)**: Golang application that manages the articles.
  - Connected to a MongoDB database.
    - GET /healthz
    - GET /article/
    - POST /article/
    - DELETE /article/:articleId/
- **[cart-service](./cart-service)**: Golang application that manages the shopping cart.
  - Connected to a Redis database.
    - GET /healthz
    - GET /cart/:cartId/
    - PUT /cart/:cartId/
    - DELETE /cart/:cartId/
- **[User frontend](./front-user)**: Vuejs application that serves as the frontend.
  - GET /
  - GET /shop
  - GET /cart
- **[Admin frontend](./front-admin)**: Vuejs application that serves as the admin 
frontend.
  - GET /
  - GET /articles
  - GET /about

## Run the application

### `docker-compose`

A docker-compose file is provided at the root of the repository to run the application.

```sh
docker-compose up -d
```
### Publish

To publish all the container images to the GCP training registry, just run the following commands:
```sh
gcloud auth configure-docker europe-west1-docker.pkg.dev
docker-compose -f docker-compose-push.yml build
docker-compose -f docker-compose-push.yml push
```

### Kubernetes

> In a near future :)
