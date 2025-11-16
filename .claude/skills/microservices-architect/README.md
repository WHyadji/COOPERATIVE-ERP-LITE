# Microservices Architect Skill

A focused skill for generating Docker-based microservices architectures with essential patterns only. No over-engineering - just clean, production-ready microservices.

## Philosophy

This skill follows the **KISS principle** - Keep It Simple, Stupid. It generates only what you need:
- ✅ Docker Compose (not Kubernetes until you need it)
- ✅ Simple API Gateway (not service mesh)
- ✅ Basic message queue (not complex event sourcing)
- ✅ One database per service (simple pattern)
- ✅ REST APIs (not GraphQL unless needed)
- ✅ JWT auth (straightforward security)

## Features

### Essential Components
- **API Gateway** - Single entry point (Nginx or Kong)
- **Microservices** - Independent Node.js services
- **Databases** - PostgreSQL/MongoDB per service
- **Message Queue** - RabbitMQ for async communication
- **Cache** - Redis for performance
- **Load Balancing** - Built-in with Docker
- **Health Checks** - Service monitoring
- **Authentication** - JWT-based auth

### What's NOT Included (On Purpose)
- ❌ Kubernetes - Docker Compose is enough to start
- ❌ Service Mesh (Istio/Linkerd) - Too complex for most
- ❌ CQRS/Event Sourcing - Unless you really need it
- ❌ GraphQL Federation - REST is simpler
- ❌ Distributed Tracing - Start with basic logs
- ❌ Complex CI/CD - Keep deployment simple

## Quick Start

### Generate Basic Microservices
```bash
# Simple e-commerce microservices
python scripts/generate_microservices.py \
  --name "EcommerceApp" \
  --services user,product,order

# With specific gateway and queue
python scripts/generate_microservices.py \
  --name "MyApp" \
  --services auth,catalog,payment \
  --gateway kong \
  --queue rabbitmq
```

### Start Everything with One Command
```bash
cd ecommerce-app
docker-compose up -d

# That's it! Access at http://localhost:8000
```

## Generated Architecture

```
Your App
│
├── services/               # Microservices
│   ├── user/              # User service
│   ├── product/           # Product service
│   └── order/             # Order service
│
├── gateway/               # API Gateway config
│   └── nginx.conf         # or kong.yml
│
├── docker-compose.yml     # Everything runs here
├── .env.example          # Environment variables
└── Makefile              # Helpful commands
```

## Each Service Includes

```
service-name/
├── src/
│   ├── index.js          # Service entry point
│   ├── routes/           # API endpoints
│   ├── controllers/      # Business logic
│   ├── models/           # Database models
│   ├── middleware/       # Auth, validation
│   ├── config/           # Configurations
│   └── utils/            # Helpers
├── Dockerfile            # Container setup
├── package.json          # Dependencies
└── .env.example         # Service config
```

## Simple Docker Compose Setup

Everything runs with one file:
```yaml
version: '3.8'

services:
  gateway:         # API Gateway (port 8000)
  user-service:    # User service (port 3001)
  product-service: # Product service (port 3002)
  order-service:   # Order service (port 3003)
  user-db:         # PostgreSQL for users
  product-db:      # MongoDB for products
  order-db:        # PostgreSQL for orders
  rabbitmq:        # Message queue
  redis:           # Cache
```

## API Endpoints (Auto-Generated)

Each service gets these endpoints:
- `GET /api/{service}/` - Service info
- `POST /api/{service}/register` - User registration
- `POST /api/{service}/login` - User login
- `GET /api/{service}/data` - Get data
- `POST /api/{service}/data` - Create data
- `PUT /api/{service}/data/:id` - Update data
- `DELETE /api/{service}/data/:id` - Delete data

Access via gateway: `http://localhost:8000/api/user/login`

## Service Communication

### REST (Simple)
```javascript
// Call another service
const response = await axios.get('http://product-service:3000/api/v1/products');
```

### Message Queue (Async)
```javascript
// Publish event
await publishMessage('order.created', { orderId: 123 });

// Subscribe to events
channel.consume('order.created', (msg) => {
  // Handle event
});
```

## Development Commands

```bash
make build    # Build all services
make start    # Start everything
make stop     # Stop everything
make logs     # View logs
make status   # Check status
make test     # Run tests
make clean    # Clean up
```

## Configuration Options

### Services
Any combination of services:
```bash
--services user,product,order,payment,shipping,inventory
```

### API Gateway
- `nginx` - Simple, fast, reliable (default)
- `kong` - More features, API management

### Message Queue
- `rabbitmq` - Full-featured message broker (default)
- `redis` - Simpler pub/sub

### Databases
- `postgres` - Relational data (default)
- `mongodb` - Document store
- `mysql` - Alternative to postgres

## Best Practices Included

1. **One Database Per Service** - Microservices 101
2. **API Gateway Pattern** - Single entry point
3. **Health Checks** - Know when services are down
4. **Environment Variables** - Easy configuration
5. **Docker Networks** - Services communicate securely
6. **Graceful Shutdown** - Handle SIGTERM properly
7. **Error Handling** - Consistent error responses
8. **Logging** - Structured logs with Winston
9. **Authentication** - JWT tokens
10. **Message Queue** - Async communication

## When to Use This

✅ **Perfect for:**
- Starting a new microservices project
- Learning microservices patterns
- Proof of concepts
- Small to medium projects
- Teams new to microservices

❌ **Not for:**
- Existing Kubernetes deployments
- Need for service mesh features
- Complex event sourcing requirements
- Massive scale (1000+ requests/sec)

## Scaling Later

When you need to scale:
1. **First**: Scale with Docker Compose (multiple replicas)
2. **Then**: Add a load balancer (HAProxy/Traefik)
3. **Later**: Move to Kubernetes if needed
4. **Finally**: Add service mesh if complexity demands it

## Example: E-commerce Microservices

```bash
python scripts/generate_microservices.py \
  --name "ShopAPI" \
  --services user,product,cart,order,payment \
  --gateway kong \
  --queue rabbitmq

cd shop-api
docker-compose up -d

# Your microservices are running!
# API Gateway: http://localhost:8000
# RabbitMQ UI: http://localhost:15672
```

## Tips

1. **Start Small** - 3-5 services maximum initially
2. **Use Docker Compose** - Don't jump to Kubernetes
3. **REST First** - GraphQL adds complexity
4. **Simple Auth** - JWT is enough for most cases
5. **Basic Queue** - RabbitMQ handles most needs
6. **Logs Over Tracing** - Distributed tracing can wait
7. **Monolith First** - Consider starting with a modular monolith

## Philosophy: YAGNI

**You Aren't Gonna Need It** - This skill embraces YAGNI:
- No Kubernetes orchestration
- No service mesh
- No CQRS/Event Sourcing
- No distributed tracing
- No complex CI/CD pipelines

Add these only when you actually need them!

## License

This skill is provided for use with Claude AI. Keep it simple!