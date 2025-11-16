---
name: Microservices Architect
description: Generate Docker-based microservices with essential patterns. Creates API Gateway, services, databases, message queues, and monitoring - focused on simplicity and production readiness without unnecessary complexity.
---

# Microservices Architect

## Overview

This skill generates clean, Docker-based microservices architectures with essential patterns only. No over-engineering - just production-ready microservices with API gateway, service discovery, databases, messaging, and basic monitoring.

## Quick Start

### Generate Microservices Architecture
```bash
# Basic microservices with Docker Compose
python scripts/generate_microservices.py --name "EcommerceApp" --services user,product,order

# With API Gateway and Message Queue
python scripts/generate_microservices.py --name "MyApp" --services auth,catalog,payment --gateway kong --queue rabbitmq

# With specific databases
python scripts/generate_microservices.py --name "MyApp" --services user,post --databases postgres,mongodb
```

## Architecture Patterns Included

### Essential Components Only
- **API Gateway** - Single entry point, routing, rate limiting
- **Microservices** - Independent services with REST APIs
- **Databases** - One database per service pattern
- **Message Queue** - Async communication between services
- **Service Discovery** - Simple DNS-based discovery
- **Load Balancing** - Docker/Nginx load balancing
- **Monitoring** - Basic health checks and logging

### What's NOT Included (Keeping it Simple)
- ❌ Service Mesh (Istio, Linkerd) - too complex for most
- ❌ Kubernetes - Docker Compose is enough to start
- ❌ Complex orchestration - KISS principle
- ❌ CQRS/Event Sourcing - unless you really need it
- ❌ Distributed tracing - start with basic logging

## Generated Architecture

```
┌─────────────────┐
│   API Gateway   │ (Kong/Nginx)
│   Port: 8000    │
└────────┬────────┘
         │
    ┌────┴────┬──────────┬──────────┐
    ▼         ▼          ▼          ▼
┌────────┐ ┌────────┐ ┌────────┐ ┌────────┐
│ Auth   │ │Product │ │ Order  │ │Payment │
│Service │ │Service │ │Service │ │Service │
│:3001   │ │:3002   │ │:3003   │ │:3004   │
└───┬────┘ └───┬────┘ └───┬────┘ └───┬────┘
    │          │          │          │
    ▼          ▼          ▼          ▼
┌────────┐ ┌────────┐ ┌────────┐ ┌────────┐
│AuthDB  │ │ProdDB  │ │OrderDB │ │PayDB   │
│Postgres│ │MongoDB │ │Postgres│ │Postgres│
└────────┘ └────────┘ └────────┘ └────────┘
                │
         ┌──────┴──────┐
         ▼             ▼
    ┌─────────┐  ┌───────────┐
    │RabbitMQ │  │   Redis    │
    │ Queue   │  │   Cache    │
    └─────────┘  └───────────┘
```

## Service Template Structure

Each microservice includes:
```
service-name/
├── src/
│   ├── index.js         # Entry point
│   ├── routes/          # API routes
│   ├── controllers/     # Business logic
│   ├── models/          # Data models
│   ├── middleware/      # Auth, validation
│   └── utils/           # Helpers
├── tests/               # Unit & integration tests
├── Dockerfile           # Container definition
├── package.json         # Dependencies
├── .env.example         # Environment template
└── README.md           # Service documentation
```

## Docker Compose Configuration

### Complete Stack in One File
```yaml
version: '3.8'

services:
  # API Gateway
  gateway:
    image: kong:3.4-alpine
    environment:
      KONG_DATABASE: "off"
      KONG_DECLARATIVE_CONFIG: /kong/kong.yml
      KONG_PROXY_ACCESS_LOG: /dev/stdout
      KONG_ADMIN_ACCESS_LOG: /dev/stdout
    ports:
      - "8000:8000"  # API Gateway port
      - "8001:8001"  # Admin API
    volumes:
      - ./gateway/kong.yml:/kong/kong.yml
    networks:
      - microservices
    depends_on:
      - auth-service
      - product-service
      - order-service

  # Auth Service
  auth-service:
    build: ./services/auth
    environment:
      NODE_ENV: production
      PORT: 3000
      DB_HOST: auth-db
      DB_PORT: 5432
      DB_NAME: auth
      JWT_SECRET: ${JWT_SECRET}
      REDIS_URL: redis://cache:6379
    ports:
      - "3001:3000"
    depends_on:
      - auth-db
      - cache
    networks:
      - microservices
    restart: unless-stopped

  # Product Service
  product-service:
    build: ./services/product
    environment:
      NODE_ENV: production
      PORT: 3000
      MONGO_URL: mongodb://product-db:27017/products
      REDIS_URL: redis://cache:6379
      RABBITMQ_URL: amqp://rabbitmq:5672
    ports:
      - "3002:3000"
    depends_on:
      - product-db
      - cache
      - rabbitmq
    networks:
      - microservices
    restart: unless-stopped

  # Order Service
  order-service:
    build: ./services/order
    environment:
      NODE_ENV: production
      PORT: 3000
      DB_HOST: order-db
      DB_PORT: 5432
      DB_NAME: orders
      RABBITMQ_URL: amqp://rabbitmq:5672
    ports:
      - "3003:3000"
    depends_on:
      - order-db
      - rabbitmq
    networks:
      - microservices
    restart: unless-stopped

  # Databases
  auth-db:
    image: postgres:15-alpine
    environment:
      POSTGRES_DB: auth
      POSTGRES_USER: auth_user
      POSTGRES_PASSWORD: ${AUTH_DB_PASSWORD}
    volumes:
      - auth-db-data:/var/lib/postgresql/data
    networks:
      - microservices

  product-db:
    image: mongo:7
    environment:
      MONGO_INITDB_DATABASE: products
    volumes:
      - product-db-data:/data/db
    networks:
      - microservices

  order-db:
    image: postgres:15-alpine
    environment:
      POSTGRES_DB: orders
      POSTGRES_USER: order_user
      POSTGRES_PASSWORD: ${ORDER_DB_PASSWORD}
    volumes:
      - order-db-data:/var/lib/postgresql/data
    networks:
      - microservices

  # Message Queue
  rabbitmq:
    image: rabbitmq:3-management-alpine
    ports:
      - "5672:5672"
      - "15672:15672"  # Management UI
    environment:
      RABBITMQ_DEFAULT_USER: admin
      RABBITMQ_DEFAULT_PASS: ${RABBITMQ_PASSWORD}
    volumes:
      - rabbitmq-data:/var/lib/rabbitmq
    networks:
      - microservices

  # Cache
  cache:
    image: redis:7-alpine
    ports:
      - "6379:6379"
    volumes:
      - cache-data:/data
    networks:
      - microservices

networks:
  microservices:
    driver: bridge

volumes:
  auth-db-data:
  product-db-data:
  order-db-data:
  rabbitmq-data:
  cache-data:
```

## Service Implementation (Node.js/Express)

### Base Service Template
```javascript
// src/index.js
const express = require('express');
const cors = require('cors');
const helmet = require('helmet');
const morgan = require('morgan');
const { errorHandler } = require('./middleware/errorHandler');
const routes = require('./routes');
const { connectDB } = require('./config/database');
const { connectMessageQueue } = require('./config/messageQueue');

const app = express();
const PORT = process.env.PORT || 3000;

// Middleware
app.use(helmet());
app.use(cors());
app.use(express.json());
app.use(morgan('combined'));

// Health check
app.get('/health', (req, res) => {
  res.json({ 
    status: 'healthy',
    service: process.env.SERVICE_NAME,
    timestamp: new Date().toISOString()
  });
});

// API routes
app.use('/api/v1', routes);

// Error handling
app.use(errorHandler);

// Start server
async function startServer() {
  try {
    // Connect to database
    await connectDB();
    
    // Connect to message queue
    await connectMessageQueue();
    
    app.listen(PORT, () => {
      console.log(`Service running on port ${PORT}`);
    });
  } catch (error) {
    console.error('Failed to start service:', error);
    process.exit(1);
  }
}

startServer();

// Graceful shutdown
process.on('SIGTERM', () => {
  console.log('SIGTERM received, closing server gracefully');
  process.exit(0);
});
```

### Service Routes
```javascript
// src/routes/index.js
const router = require('express').Router();
const { authenticate } = require('../middleware/auth');
const { validate } = require('../middleware/validation');
const controller = require('../controllers');

// Public routes
router.post('/register', validate('register'), controller.register);
router.post('/login', validate('login'), controller.login);

// Protected routes
router.use(authenticate);
router.get('/profile', controller.getProfile);
router.put('/profile', validate('updateProfile'), controller.updateProfile);

module.exports = router;
```

### Database Connection
```javascript
// src/config/database.js
const { Sequelize } = require('sequelize');

const sequelize = new Sequelize({
  host: process.env.DB_HOST,
  port: process.env.DB_PORT,
  database: process.env.DB_NAME,
  username: process.env.DB_USER,
  password: process.env.DB_PASSWORD,
  dialect: 'postgres',
  logging: false,
  pool: {
    max: 10,
    min: 0,
    acquire: 30000,
    idle: 10000
  }
});

async function connectDB() {
  try {
    await sequelize.authenticate();
    console.log('Database connected');
    
    // Sync models
    await sequelize.sync({ alter: true });
  } catch (error) {
    console.error('Database connection failed:', error);
    throw error;
  }
}

module.exports = { sequelize, connectDB };
```

### Message Queue Integration
```javascript
// src/config/messageQueue.js
const amqp = require('amqplib');

let channel = null;

async function connectMessageQueue() {
  try {
    const connection = await amqp.connect(process.env.RABBITMQ_URL);
    channel = await connection.createChannel();
    
    console.log('Message queue connected');
    
    // Set up queues
    await channel.assertQueue('order.created');
    await channel.assertQueue('payment.processed');
    
    // Start consuming messages
    startConsumers();
  } catch (error) {
    console.error('Message queue connection failed:', error);
    throw error;
  }
}

function startConsumers() {
  // Listen for order created events
  channel.consume('order.created', async (msg) => {
    if (msg) {
      const order = JSON.parse(msg.content.toString());
      console.log('Order created:', order);
      
      // Process order
      await processOrder(order);
      
      // Acknowledge message
      channel.ack(msg);
    }
  });
}

async function publishMessage(queue, message) {
  if (!channel) {
    throw new Error('Message queue not connected');
  }
  
  channel.sendToQueue(queue, Buffer.from(JSON.stringify(message)));
}

module.exports = { connectMessageQueue, publishMessage };
```

## API Gateway Configuration (Kong)

```yaml
# gateway/kong.yml
_format_version: "3.0"

services:
  - name: auth-service
    url: http://auth-service:3000
    routes:
      - name: auth-routes
        paths:
          - /api/auth
        strip_path: true
        methods:
          - GET
          - POST
          - PUT
          - DELETE

  - name: product-service
    url: http://product-service:3000
    routes:
      - name: product-routes
        paths:
          - /api/products
        strip_path: true
        methods:
          - GET
          - POST
          - PUT
          - DELETE

  - name: order-service
    url: http://order-service:3000
    routes:
      - name: order-routes
        paths:
          - /api/orders
        strip_path: true
        methods:
          - GET
          - POST
          - PUT
          - DELETE

plugins:
  - name: rate-limiting
    config:
      minute: 100
      hour: 10000
      policy: local

  - name: cors
    config:
      origins:
        - "*"
      methods:
        - GET
        - POST
        - PUT
        - DELETE
      headers:
        - Accept
        - Authorization
        - Content-Type

  - name: jwt
    config:
      secret_is_base64: false
      key_claim_name: kid
```

## Inter-Service Communication

### REST Communication
```javascript
// src/services/apiClient.js
const axios = require('axios');

class ApiClient {
  constructor(baseURL) {
    this.client = axios.create({
      baseURL,
      timeout: 5000,
      headers: {
        'Content-Type': 'application/json'
      }
    });
    
    // Add request interceptor for auth
    this.client.interceptors.request.use(
      config => {
        // Add auth token if available
        const token = getServiceToken();
        if (token) {
          config.headers.Authorization = `Bearer ${token}`;
        }
        return config;
      },
      error => Promise.reject(error)
    );
    
    // Add response interceptor for error handling
    this.client.interceptors.response.use(
      response => response.data,
      error => {
        console.error('API call failed:', error.message);
        throw error;
      }
    );
  }
  
  async get(url, params) {
    return this.client.get(url, { params });
  }
  
  async post(url, data) {
    return this.client.post(url, data);
  }
  
  async put(url, data) {
    return this.client.put(url, data);
  }
  
  async delete(url) {
    return this.client.delete(url);
  }
}

// Service clients
const userService = new ApiClient('http://auth-service:3000');
const productService = new ApiClient('http://product-service:3000');
const orderService = new ApiClient('http://order-service:3000');

module.exports = {
  userService,
  productService,
  orderService
};
```

### Event-Driven Communication
```javascript
// src/events/publisher.js
const { publishMessage } = require('../config/messageQueue');

class EventPublisher {
  async publishOrderCreated(order) {
    await publishMessage('order.created', {
      eventType: 'ORDER_CREATED',
      timestamp: new Date().toISOString(),
      data: order
    });
  }
  
  async publishPaymentProcessed(payment) {
    await publishMessage('payment.processed', {
      eventType: 'PAYMENT_PROCESSED',
      timestamp: new Date().toISOString(),
      data: payment
    });
  }
  
  async publishUserRegistered(user) {
    await publishMessage('user.registered', {
      eventType: 'USER_REGISTERED',
      timestamp: new Date().toISOString(),
      data: {
        id: user.id,
        email: user.email,
        name: user.name
      }
    });
  }
}

module.exports = new EventPublisher();
```

## Database Models

### PostgreSQL Model (Sequelize)
```javascript
// src/models/User.js
const { DataTypes } = require('sequelize');
const { sequelize } = require('../config/database');
const bcrypt = require('bcryptjs');

const User = sequelize.define('User', {
  id: {
    type: DataTypes.UUID,
    defaultValue: DataTypes.UUIDV4,
    primaryKey: true
  },
  email: {
    type: DataTypes.STRING,
    unique: true,
    allowNull: false,
    validate: {
      isEmail: true
    }
  },
  password: {
    type: DataTypes.STRING,
    allowNull: false
  },
  name: {
    type: DataTypes.STRING,
    allowNull: false
  },
  role: {
    type: DataTypes.ENUM('user', 'admin'),
    defaultValue: 'user'
  },
  isActive: {
    type: DataTypes.BOOLEAN,
    defaultValue: true
  }
}, {
  timestamps: true,
  hooks: {
    beforeCreate: async (user) => {
      user.password = await bcrypt.hash(user.password, 10);
    }
  }
});

User.prototype.validatePassword = async function(password) {
  return bcrypt.compare(password, this.password);
};

User.prototype.toJSON = function() {
  const values = Object.assign({}, this.get());
  delete values.password;
  return values;
};

module.exports = User;
```

### MongoDB Model (Mongoose)
```javascript
// src/models/Product.js
const mongoose = require('mongoose');

const productSchema = new mongoose.Schema({
  name: {
    type: String,
    required: true,
    trim: true
  },
  description: {
    type: String,
    required: true
  },
  price: {
    type: Number,
    required: true,
    min: 0
  },
  category: {
    type: String,
    required: true,
    enum: ['electronics', 'clothing', 'food', 'books', 'other']
  },
  stock: {
    type: Number,
    default: 0,
    min: 0
  },
  images: [{
    url: String,
    alt: String
  }],
  ratings: {
    average: {
      type: Number,
      default: 0,
      min: 0,
      max: 5
    },
    count: {
      type: Number,
      default: 0
    }
  }
}, {
  timestamps: true
});

// Indexes
productSchema.index({ name: 'text', description: 'text' });
productSchema.index({ category: 1, price: 1 });

// Methods
productSchema.methods.decreaseStock = function(quantity) {
  if (this.stock < quantity) {
    throw new Error('Insufficient stock');
  }
  this.stock -= quantity;
  return this.save();
};

module.exports = mongoose.model('Product', productSchema);
```

## Authentication & Authorization

### JWT Authentication
```javascript
// src/middleware/auth.js
const jwt = require('jsonwebtoken');

function authenticate(req, res, next) {
  try {
    const token = req.headers.authorization?.replace('Bearer ', '');
    
    if (!token) {
      return res.status(401).json({ error: 'Authentication required' });
    }
    
    const decoded = jwt.verify(token, process.env.JWT_SECRET);
    req.user = decoded;
    next();
  } catch (error) {
    return res.status(401).json({ error: 'Invalid token' });
  }
}

function authorize(...roles) {
  return (req, res, next) => {
    if (!roles.includes(req.user.role)) {
      return res.status(403).json({ error: 'Insufficient permissions' });
    }
    next();
  };
}

module.exports = { authenticate, authorize };
```

## Error Handling

```javascript
// src/middleware/errorHandler.js
class AppError extends Error {
  constructor(message, statusCode) {
    super(message);
    this.statusCode = statusCode;
    this.isOperational = true;
  }
}

function errorHandler(err, req, res, next) {
  let { statusCode, message } = err;
  
  if (!err.isOperational) {
    statusCode = 500;
    message = 'Something went wrong';
    console.error('ERROR:', err);
  }
  
  res.status(statusCode).json({
    error: message,
    ...(process.env.NODE_ENV === 'development' && { stack: err.stack })
  });
}

module.exports = { AppError, errorHandler };
```

## Caching with Redis

```javascript
// src/services/cache.js
const redis = require('redis');

const client = redis.createClient({
  url: process.env.REDIS_URL
});

client.on('error', (err) => console.error('Redis error:', err));
client.connect();

class CacheService {
  async get(key) {
    try {
      const value = await client.get(key);
      return value ? JSON.parse(value) : null;
    } catch (error) {
      console.error('Cache get error:', error);
      return null;
    }
  }
  
  async set(key, value, ttl = 3600) {
    try {
      await client.setEx(key, ttl, JSON.stringify(value));
    } catch (error) {
      console.error('Cache set error:', error);
    }
  }
  
  async invalidate(pattern) {
    try {
      const keys = await client.keys(pattern);
      if (keys.length > 0) {
        await client.del(keys);
      }
    } catch (error) {
      console.error('Cache invalidate error:', error);
    }
  }
}

module.exports = new CacheService();
```

## Health Checks & Monitoring

```javascript
// src/utils/healthCheck.js
const { sequelize } = require('../config/database');
const cache = require('../services/cache');

async function checkHealth() {
  const checks = {
    service: 'healthy',
    database: 'unknown',
    cache: 'unknown',
    messageQueue: 'unknown'
  };
  
  // Check database
  try {
    await sequelize.authenticate();
    checks.database = 'healthy';
  } catch (error) {
    checks.database = 'unhealthy';
  }
  
  // Check cache
  try {
    await cache.set('health', 'check', 10);
    await cache.get('health');
    checks.cache = 'healthy';
  } catch (error) {
    checks.cache = 'unhealthy';
  }
  
  // Overall status
  const allHealthy = Object.values(checks).every(status => status === 'healthy');
  
  return {
    status: allHealthy ? 'healthy' : 'degraded',
    checks,
    timestamp: new Date().toISOString()
  };
}

module.exports = { checkHealth };
```

## Testing

### Unit Tests
```javascript
// tests/unit/user.test.js
const { expect } = require('chai');
const sinon = require('sinon');
const UserService = require('../../src/services/userService');

describe('UserService', () => {
  describe('createUser', () => {
    it('should create a new user', async () => {
      const userData = {
        email: 'test@example.com',
        password: 'password123',
        name: 'Test User'
      };
      
      const user = await UserService.createUser(userData);
      
      expect(user).to.have.property('id');
      expect(user.email).to.equal(userData.email);
      expect(user).to.not.have.property('password');
    });
    
    it('should throw error for duplicate email', async () => {
      const userData = {
        email: 'existing@example.com',
        password: 'password123',
        name: 'Test User'
      };
      
      await UserService.createUser(userData);
      
      try {
        await UserService.createUser(userData);
        expect.fail('Should have thrown error');
      } catch (error) {
        expect(error.message).to.include('already exists');
      }
    });
  });
});
```

### Integration Tests
```javascript
// tests/integration/api.test.js
const request = require('supertest');
const app = require('../../src/index');

describe('API Integration Tests', () => {
  let authToken;
  
  describe('POST /api/v1/register', () => {
    it('should register a new user', async () => {
      const response = await request(app)
        .post('/api/v1/register')
        .send({
          email: 'test@example.com',
          password: 'password123',
          name: 'Test User'
        });
      
      expect(response.status).toBe(201);
      expect(response.body).toHaveProperty('token');
      expect(response.body.user.email).toBe('test@example.com');
      
      authToken = response.body.token;
    });
  });
  
  describe('GET /api/v1/profile', () => {
    it('should get user profile with valid token', async () => {
      const response = await request(app)
        .get('/api/v1/profile')
        .set('Authorization', `Bearer ${authToken}`);
      
      expect(response.status).toBe(200);
      expect(response.body).toHaveProperty('email');
    });
    
    it('should return 401 without token', async () => {
      const response = await request(app)
        .get('/api/v1/profile');
      
      expect(response.status).toBe(401);
    });
  });
});
```

## Deployment Scripts

### Docker Build Script
```bash
#!/bin/bash
# scripts/build.sh

echo "Building microservices..."

# Build each service
for service in services/*; do
  if [ -d "$service" ]; then
    service_name=$(basename $service)
    echo "Building $service_name..."
    docker build -t $service_name:latest $service
  fi
done

echo "Build complete!"
```

### Start Script
```bash
#!/bin/bash
# scripts/start.sh

# Load environment variables
export $(cat .env | grep -v '^#' | xargs)

# Start services
docker-compose up -d

# Wait for services to be healthy
echo "Waiting for services to start..."
sleep 10

# Check health
for service in $(docker-compose ps --services); do
  echo "Checking $service..."
  docker-compose exec $service curl -f http://localhost:3000/health || echo "$service is not healthy"
done

echo "All services started!"
```

## Environment Configuration

```env
# .env
NODE_ENV=production

# JWT
JWT_SECRET=your-super-secret-jwt-key

# Database Passwords
AUTH_DB_PASSWORD=auth_password
ORDER_DB_PASSWORD=order_password
PAYMENT_DB_PASSWORD=payment_password

# RabbitMQ
RABBITMQ_PASSWORD=rabbitmq_password

# Redis
REDIS_PASSWORD=redis_password

# Service Ports
AUTH_SERVICE_PORT=3001
PRODUCT_SERVICE_PORT=3002
ORDER_SERVICE_PORT=3003
PAYMENT_SERVICE_PORT=3004

# API Gateway
GATEWAY_PORT=8000
```

## Makefile for Common Tasks

```makefile
# Makefile
.PHONY: help build start stop logs clean test

help:
	@echo "Available commands:"
	@echo "  make build    - Build all services"
	@echo "  make start    - Start all services"
	@echo "  make stop     - Stop all services"
	@echo "  make logs     - View logs"
	@echo "  make clean    - Clean up"
	@echo "  make test     - Run tests"

build:
	docker-compose build

start:
	docker-compose up -d
	@echo "Services running at http://localhost:8000"

stop:
	docker-compose down

logs:
	docker-compose logs -f

clean:
	docker-compose down -v
	docker system prune -f

test:
	docker-compose exec auth-service npm test
	docker-compose exec product-service npm test
	docker-compose exec order-service npm test

restart: stop start

status:
	docker-compose ps
```

## Best Practices

1. **One Database Per Service** - Each service owns its data
2. **API Gateway Pattern** - Single entry point for clients
3. **Service Discovery** - Services find each other via DNS
4. **Circuit Breaker** - Prevent cascading failures
5. **Health Checks** - Monitor service health
6. **Centralized Logging** - Aggregate logs from all services
7. **Message Queue** - Asynchronous communication
8. **Caching** - Redis for performance
9. **Rate Limiting** - Prevent abuse
10. **Graceful Shutdown** - Handle SIGTERM properly