#!/usr/bin/env python3
"""
Microservices Generator
Generates Docker-based microservices architecture with essential patterns
"""

import os
import sys
import json
import yaml
import argparse
from pathlib import Path
from typing import Dict, List, Optional

class MicroservicesGenerator:
    def __init__(self, name: str, services: List[str], options: Dict = None):
        self.name = name
        self.services = services
        self.options = options or {}
        self.output_dir = Path(self.options.get('output', f"./{name.lower().replace(' ', '-')}"))
        
        # Configuration
        self.gateway = self.options.get('gateway', 'nginx')  # nginx or kong
        self.queue = self.options.get('queue', 'rabbitmq')   # rabbitmq or redis
        self.databases = self.options.get('databases', ['postgres'])  # postgres, mongodb, mysql
        
    def generate(self):
        """Main generation process"""
        print(f"🚀 Generating microservices architecture: {self.name}")
        print(f"   Services: {', '.join(self.services)}")
        print(f"   Gateway: {self.gateway}")
        print(f"   Queue: {self.queue}")
        
        # Create root directory
        self.output_dir.mkdir(parents=True, exist_ok=True)
        
        # Generate structure
        self._generate_services()
        self._generate_gateway()
        self._generate_docker_compose()
        self._generate_common_files()
        
        print(f"\n✅ Microservices architecture generated in {self.output_dir}")
        print(f"\n📖 Next steps:")
        print(f"   1. cd {self.output_dir}")
        print(f"   2. cp .env.example .env")
        print(f"   3. docker-compose build")
        print(f"   4. docker-compose up -d")
        print(f"   5. Access API Gateway at http://localhost:8000")
    
    def _create_file(self, path: str, content: str):
        """Create a file with content"""
        file_path = self.output_dir / path
        file_path.parent.mkdir(parents=True, exist_ok=True)
        file_path.write_text(content)
        print(f"    ✓ Created {path}")
    
    def _generate_services(self):
        """Generate each microservice"""
        for service_name in self.services:
            print(f"\n  📦 Creating {service_name} service...")
            service_dir = f"services/{service_name}"
            
            # Create service structure
            dirs = [
                f"{service_dir}/src/routes",
                f"{service_dir}/src/controllers",
                f"{service_dir}/src/models",
                f"{service_dir}/src/middleware",
                f"{service_dir}/src/services",
                f"{service_dir}/src/utils",
                f"{service_dir}/src/config",
                f"{service_dir}/tests/unit",
                f"{service_dir}/tests/integration"
            ]
            
            for dir_path in dirs:
                (self.output_dir / dir_path).mkdir(parents=True, exist_ok=True)
            
            # Generate service files
            self._create_file(f"{service_dir}/package.json", self._get_package_json(service_name))
            self._create_file(f"{service_dir}/Dockerfile", self._get_dockerfile())
            self._create_file(f"{service_dir}/src/index.js", self._get_service_index(service_name))
            self._create_file(f"{service_dir}/src/routes/index.js", self._get_routes())
            self._create_file(f"{service_dir}/src/controllers/index.js", self._get_controller(service_name))
            self._create_file(f"{service_dir}/src/config/database.js", self._get_database_config())
            self._create_file(f"{service_dir}/src/config/messageQueue.js", self._get_message_queue_config())
            self._create_file(f"{service_dir}/src/middleware/auth.js", self._get_auth_middleware())
            self._create_file(f"{service_dir}/src/middleware/errorHandler.js", self._get_error_handler())
            self._create_file(f"{service_dir}/src/utils/logger.js", self._get_logger())
            self._create_file(f"{service_dir}/.env.example", self._get_service_env(service_name))
            self._create_file(f"{service_dir}/README.md", self._get_service_readme(service_name))
    
    def _generate_gateway(self):
        """Generate API Gateway configuration"""
        print(f"\n  🌐 Creating API Gateway ({self.gateway})...")
        
        if self.gateway == 'kong':
            self._create_file("gateway/kong.yml", self._get_kong_config())
        else:
            self._create_file("gateway/nginx.conf", self._get_nginx_config())
    
    def _generate_docker_compose(self):
        """Generate Docker Compose configuration"""
        print(f"\n  🐳 Creating Docker Compose configuration...")
        self._create_file("docker-compose.yml", self._get_docker_compose())
        self._create_file("docker-compose.dev.yml", self._get_docker_compose_dev())
    
    def _generate_common_files(self):
        """Generate common project files"""
        print(f"\n  📄 Creating common files...")
        self._create_file("README.md", self._get_readme())
        self._create_file(".env.example", self._get_env_example())
        self._create_file(".gitignore", self._get_gitignore())
        self._create_file("Makefile", self._get_makefile())
        self._create_file("scripts/build.sh", self._get_build_script())
        self._create_file("scripts/start.sh", self._get_start_script())
        self._create_file("scripts/test.sh", self._get_test_script())
    
    # Template methods
    def _get_package_json(self, service_name: str) -> str:
        return json.dumps({
            "name": f"{service_name}-service",
            "version": "1.0.0",
            "description": f"{service_name.title()} microservice",
            "main": "src/index.js",
            "scripts": {
                "start": "node src/index.js",
                "dev": "nodemon src/index.js",
                "test": "jest",
                "test:watch": "jest --watch"
            },
            "dependencies": {
                "express": "^4.18.2",
                "cors": "^2.8.5",
                "helmet": "^7.0.0",
                "morgan": "^1.10.0",
                "dotenv": "^16.0.3",
                "axios": "^1.5.0",
                "jsonwebtoken": "^9.0.2",
                "bcryptjs": "^2.4.3",
                "sequelize": "^6.33.0",
                "pg": "^8.11.3",
                "pg-hstore": "^2.3.4",
                "redis": "^4.6.10",
                "amqplib": "^0.10.3",
                "express-validator": "^7.0.1",
                "uuid": "^9.0.1"
            },
            "devDependencies": {
                "nodemon": "^3.0.1",
                "jest": "^29.7.0",
                "supertest": "^6.3.3"
            }
        }, null, 2)
    
    def _get_dockerfile(self) -> str:
        return """FROM node:18-alpine

WORKDIR /app

# Install dependencies
COPY package*.json ./
RUN npm ci --only=production

# Copy application
COPY . .

# Create non-root user
RUN addgroup -g 1001 -S nodejs && \\
    adduser -S nodejs -u 1001
USER nodejs

EXPOSE 3000

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \\
    CMD node -e "require('http').get('http://localhost:3000/health', (r) => {process.exit(r.statusCode === 200 ? 0 : 1)})"

CMD ["node", "src/index.js"]"""
    
    def _get_service_index(self, service_name: str) -> str:
        return f"""require('dotenv').config();
const express = require('express');
const cors = require('cors');
const helmet = require('helmet');
const morgan = require('morgan');

const routes = require('./routes');
const {{ errorHandler }} = require('./middleware/errorHandler');
const {{ connectDB }} = require('./config/database');
const {{ connectMessageQueue }} = require('./config/messageQueue');
const logger = require('./utils/logger');

const app = express();
const PORT = process.env.PORT || 3000;
const SERVICE_NAME = '{service_name}';

// Middleware
app.use(helmet());
app.use(cors());
app.use(express.json());
app.use(morgan('combined'));

// Health check
app.get('/health', (req, res) => {{
  res.json({{
    status: 'healthy',
    service: SERVICE_NAME,
    timestamp: new Date().toISOString()
  }});
}});

// API routes
app.use('/api/v1', routes);

// Error handling
app.use(errorHandler);

// Start server
async function startServer() {{
  try {{
    // Connect to database
    await connectDB();
    logger.info('Database connected');
    
    // Connect to message queue
    await connectMessageQueue();
    logger.info('Message queue connected');
    
    app.listen(PORT, () => {{
      logger.info(`${{SERVICE_NAME}} service running on port ${{PORT}}`);
    }});
  }} catch (error) {{
    logger.error('Failed to start service:', error);
    process.exit(1);
  }}
}}

startServer();

// Graceful shutdown
process.on('SIGTERM', () => {{
  logger.info('SIGTERM received, shutting down gracefully');
  process.exit(0);
}});"""
    
    def _get_routes(self) -> str:
        return """const router = require('express').Router();
const { authenticate } = require('../middleware/auth');
const controller = require('../controllers');

// Public routes
router.get('/', controller.getInfo);
router.post('/register', controller.register);
router.post('/login', controller.login);

// Protected routes
router.use(authenticate);
router.get('/data', controller.getData);
router.post('/data', controller.createData);
router.put('/data/:id', controller.updateData);
router.delete('/data/:id', controller.deleteData);

module.exports = router;"""
    
    def _get_controller(self, service_name: str) -> str:
        return f"""const logger = require('../utils/logger');

class Controller {{
  async getInfo(req, res) {{
    res.json({{
      service: '{service_name}',
      version: '1.0.0',
      status: 'running'
    }});
  }}
  
  async register(req, res, next) {{
    try {{
      const {{ email, password, name }} = req.body;
      
      // TODO: Implement registration logic
      
      res.status(201).json({{
        message: 'User registered successfully',
        user: {{ email, name }}
      }});
    }} catch (error) {{
      next(error);
    }}
  }}
  
  async login(req, res, next) {{
    try {{
      const {{ email, password }} = req.body;
      
      // TODO: Implement login logic
      
      res.json({{
        token: 'jwt-token-here',
        user: {{ email }}
      }});
    }} catch (error) {{
      next(error);
    }}
  }}
  
  async getData(req, res, next) {{
    try {{
      // TODO: Implement data retrieval
      
      res.json({{
        data: [],
        total: 0
      }});
    }} catch (error) {{
      next(error);
    }}
  }}
  
  async createData(req, res, next) {{
    try {{
      const data = req.body;
      
      // TODO: Implement data creation
      
      res.status(201).json({{
        message: 'Data created successfully',
        data
      }});
    }} catch (error) {{
      next(error);
    }}
  }}
  
  async updateData(req, res, next) {{
    try {{
      const {{ id }} = req.params;
      const updates = req.body;
      
      // TODO: Implement data update
      
      res.json({{
        message: 'Data updated successfully',
        data: {{ id, ...updates }}
      }});
    }} catch (error) {{
      next(error);
    }}
  }}
  
  async deleteData(req, res, next) {{
    try {{
      const {{ id }} = req.params;
      
      // TODO: Implement data deletion
      
      res.json({{
        message: 'Data deleted successfully',
        id
      }});
    }} catch (error) {{
      next(error);
    }}
  }}
}}

module.exports = new Controller();"""
    
    def _get_docker_compose(self) -> str:
        services = {}
        
        # Add gateway
        if self.gateway == 'kong':
            services['gateway'] = {
                'image': 'kong:3.4-alpine',
                'environment': {
                    'KONG_DATABASE': 'off',
                    'KONG_DECLARATIVE_CONFIG': '/kong/kong.yml',
                    'KONG_PROXY_ACCESS_LOG': '/dev/stdout'
                },
                'ports': ['8000:8000', '8001:8001'],
                'volumes': ['./gateway/kong.yml:/kong/kong.yml'],
                'networks': ['microservices'],
                'depends_on': [f"{s}-service" for s in self.services]
            }
        else:
            services['gateway'] = {
                'image': 'nginx:alpine',
                'ports': ['8000:80'],
                'volumes': ['./gateway/nginx.conf:/etc/nginx/nginx.conf:ro'],
                'networks': ['microservices'],
                'depends_on': [f"{s}-service" for s in self.services]
            }
        
        # Add services
        port = 3001
        for service in self.services:
            services[f"{service}-service"] = {
                'build': f'./services/{service}',
                'environment': {
                    'NODE_ENV': 'production',
                    'PORT': '3000',
                    'DB_HOST': f'{service}-db',
                    'REDIS_URL': 'redis://cache:6379',
                    'RABBITMQ_URL': 'amqp://rabbitmq:5672'
                },
                'ports': [f'{port}:3000'],
                'depends_on': [f'{service}-db', 'cache', 'rabbitmq'],
                'networks': ['microservices'],
                'restart': 'unless-stopped'
            }
            
            # Add database for service
            services[f"{service}-db"] = {
                'image': 'postgres:15-alpine',
                'environment': {
                    'POSTGRES_DB': service,
                    'POSTGRES_USER': f'{service}_user',
                    'POSTGRES_PASSWORD': '${' + f'{service.upper()}_DB_PASSWORD' + '}'
                },
                'volumes': [f'{service}-db-data:/var/lib/postgresql/data'],
                'networks': ['microservices']
            }
            port += 1
        
        # Add message queue
        if self.queue == 'rabbitmq':
            services['rabbitmq'] = {
                'image': 'rabbitmq:3-management-alpine',
                'ports': ['5672:5672', '15672:15672'],
                'environment': {
                    'RABBITMQ_DEFAULT_USER': 'admin',
                    'RABBITMQ_DEFAULT_PASS': '${RABBITMQ_PASSWORD}'
                },
                'volumes': ['rabbitmq-data:/var/lib/rabbitmq'],
                'networks': ['microservices']
            }
        
        # Add cache
        services['cache'] = {
            'image': 'redis:7-alpine',
            'ports': ['6379:6379'],
            'volumes': ['cache-data:/data'],
            'networks': ['microservices']
        }
        
        # Create docker-compose structure
        compose = {
            'version': '3.8',
            'services': services,
            'networks': {
                'microservices': {
                    'driver': 'bridge'
                }
            },
            'volumes': {}
        }
        
        # Add volumes
        for service in self.services:
            compose['volumes'][f'{service}-db-data'] = None
        compose['volumes']['rabbitmq-data'] = None
        compose['volumes']['cache-data'] = None
        
        return yaml.dump(compose, default_flow_style=False, sort_keys=False)
    
    def _get_kong_config(self) -> str:
        services = []
        
        for service in self.services:
            services.append({
                'name': f'{service}-service',
                'url': f'http://{service}-service:3000',
                'routes': [{
                    'name': f'{service}-routes',
                    'paths': [f'/api/{service}'],
                    'strip_path': True,
                    'methods': ['GET', 'POST', 'PUT', 'DELETE']
                }]
            })
        
        config = {
            '_format_version': '3.0',
            'services': services,
            'plugins': [
                {
                    'name': 'rate-limiting',
                    'config': {
                        'minute': 100,
                        'hour': 10000,
                        'policy': 'local'
                    }
                },
                {
                    'name': 'cors',
                    'config': {
                        'origins': ['*'],
                        'methods': ['GET', 'POST', 'PUT', 'DELETE'],
                        'headers': ['Accept', 'Authorization', 'Content-Type']
                    }
                }
            ]
        }
        
        return yaml.dump(config, default_flow_style=False, sort_keys=False)
    
    def _get_nginx_config(self) -> str:
        upstreams = ""
        locations = ""
        
        port = 3001
        for service in self.services:
            upstreams += f"""
    upstream {service}_service {{
        server {service}-service:3000;
    }}
"""
            locations += f"""
        location /api/{service} {{
            rewrite ^/api/{service}/(.*) /api/v1/$1 break;
            proxy_pass http://{service}_service;
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
        }}
"""
            port += 1
        
        return f"""events {{
    worker_connections 1024;
}}

http {{
    {upstreams}
    
    server {{
        listen 80;
        
        # Health check
        location /health {{
            access_log off;
            return 200 "healthy";
        }}
        
        # Service routes
        {locations}
    }}
}}"""
    
    def _get_readme(self) -> str:
        services_list = "\n".join([f"- **{s.title()} Service** - Port 300{i+1}" 
                                   for i, s in enumerate(self.services)])
        
        return f"""# {self.name} Microservices

Docker-based microservices architecture with essential patterns.

## Architecture

- **API Gateway** ({self.gateway}) - Port 8000
- **Message Queue** ({self.queue})
- **Cache** (Redis)
- **Databases** (PostgreSQL for each service)

## Services

{services_list}

## Quick Start

1. **Setup environment**
   ```bash
   cp .env.example .env
   # Edit .env with your passwords
   ```

2. **Build services**
   ```bash
   docker-compose build
   # or
   make build
   ```

3. **Start services**
   ```bash
   docker-compose up -d
   # or
   make start
   ```

4. **Access services**
   - API Gateway: http://localhost:8000
   - RabbitMQ Management: http://localhost:15672 (admin/admin)

## API Endpoints

Each service exposes:
- `GET /api/{{service}}/` - Service info
- `POST /api/{{service}}/register` - User registration
- `POST /api/{{service}}/login` - User login
- `GET /api/{{service}}/data` - Get data (authenticated)
- `POST /api/{{service}}/data` - Create data (authenticated)
- `PUT /api/{{service}}/data/:id` - Update data (authenticated)
- `DELETE /api/{{service}}/data/:id` - Delete data (authenticated)

## Development

### Run in development mode
```bash
docker-compose -f docker-compose.yml -f docker-compose.dev.yml up
```

### View logs
```bash
docker-compose logs -f [service-name]
# or
make logs
```

### Run tests
```bash
make test
```

### Stop services
```bash
docker-compose down
# or
make stop
```

## Monitoring

### Health checks
```bash
curl http://localhost:8000/health
```

### Service status
```bash
docker-compose ps
# or
make status
```

## Useful Commands

```bash
make help     # Show available commands
make build    # Build all services
make start    # Start all services
make stop     # Stop all services
make restart  # Restart all services
make logs     # View logs
make test     # Run tests
make clean    # Clean up
```"""
    
    def _get_env_example(self) -> str:
        db_passwords = "\n".join([f"{s.upper()}_DB_PASSWORD={s}_password" 
                                  for s in self.services])
        
        return f"""# Environment Configuration
NODE_ENV=development

# JWT Secret
JWT_SECRET=your-super-secret-jwt-key-change-this

# Database Passwords
{db_passwords}

# RabbitMQ
RABBITMQ_PASSWORD=rabbitmq_password

# Redis
REDIS_PASSWORD=redis_password"""
    
    def _get_makefile(self) -> str:
        return """.PHONY: help build start stop restart logs status clean test

help:
	@echo "Available commands:"
	@echo "  make build    - Build all services"
	@echo "  make start    - Start all services"
	@echo "  make stop     - Stop all services"
	@echo "  make restart  - Restart all services"
	@echo "  make logs     - View logs"
	@echo "  make status   - Show service status"
	@echo "  make clean    - Clean up everything"
	@echo "  make test     - Run tests"

build:
	docker-compose build

start:
	docker-compose up -d
	@echo "Services running at http://localhost:8000"

stop:
	docker-compose down

restart: stop start

logs:
	docker-compose logs -f

status:
	docker-compose ps

clean:
	docker-compose down -v
	docker system prune -f

test:
	./scripts/test.sh"""
    
    def _get_database_config(self) -> str:
        return """const { Sequelize } = require('sequelize');
const logger = require('../utils/logger');

const sequelize = new Sequelize({
  host: process.env.DB_HOST || 'localhost',
  port: process.env.DB_PORT || 5432,
  database: process.env.DB_NAME || 'service_db',
  username: process.env.DB_USER || 'postgres',
  password: process.env.DB_PASSWORD || 'postgres',
  dialect: 'postgres',
  logging: process.env.NODE_ENV === 'development' ? logger.debug : false,
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
    await sequelize.sync({ alter: true });
    logger.info('Database connected and synced');
  } catch (error) {
    logger.error('Database connection failed:', error);
    throw error;
  }
}

module.exports = { sequelize, connectDB };"""
    
    def _get_message_queue_config(self) -> str:
        return """const amqp = require('amqplib');
const logger = require('../utils/logger');

let connection = null;
let channel = null;

async function connectMessageQueue() {
  try {
    connection = await amqp.connect(process.env.RABBITMQ_URL || 'amqp://localhost');
    channel = await connection.createChannel();
    
    logger.info('Message queue connected');
    
    // Setup queues
    await setupQueues();
    
    // Start consumers
    startConsumers();
  } catch (error) {
    logger.error('Message queue connection failed:', error);
    // Retry connection after 5 seconds
    setTimeout(connectMessageQueue, 5000);
  }
}

async function setupQueues() {
  // Declare queues
  await channel.assertQueue('events');
  await channel.assertQueue('notifications');
}

function startConsumers() {
  // Consume messages
  channel.consume('events', (msg) => {
    if (msg) {
      const event = JSON.parse(msg.content.toString());
      logger.info('Event received:', event);
      
      // Process event
      // ...
      
      channel.ack(msg);
    }
  });
}

async function publishMessage(queue, message) {
  if (!channel) {
    throw new Error('Message queue not connected');
  }
  
  channel.sendToQueue(queue, Buffer.from(JSON.stringify(message)));
  logger.info(`Message published to ${queue}:`, message);
}

module.exports = { connectMessageQueue, publishMessage };"""
    
    def _get_auth_middleware(self) -> str:
        return """const jwt = require('jsonwebtoken');

function authenticate(req, res, next) {
  try {
    const token = req.headers.authorization?.replace('Bearer ', '');
    
    if (!token) {
      return res.status(401).json({ error: 'Authentication required' });
    }
    
    const decoded = jwt.verify(token, process.env.JWT_SECRET || 'secret');
    req.user = decoded;
    next();
  } catch (error) {
    return res.status(401).json({ error: 'Invalid token' });
  }
}

function authorize(...roles) {
  return (req, res, next) => {
    if (!req.user) {
      return res.status(401).json({ error: 'Authentication required' });
    }
    
    if (roles.length && !roles.includes(req.user.role)) {
      return res.status(403).json({ error: 'Insufficient permissions' });
    }
    
    next();
  };
}

module.exports = { authenticate, authorize };"""
    
    def _get_error_handler(self) -> str:
        return """const logger = require('../utils/logger');

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
    logger.error('ERROR:', err);
  }
  
  res.status(statusCode || 500).json({
    error: message,
    ...(process.env.NODE_ENV === 'development' && { stack: err.stack })
  });
}

module.exports = { AppError, errorHandler };"""
    
    def _get_logger(self) -> str:
        return """const winston = require('winston');

const logger = winston.createLogger({
  level: process.env.LOG_LEVEL || 'info',
  format: winston.format.combine(
    winston.format.timestamp(),
    winston.format.json()
  ),
  transports: [
    new winston.transports.Console({
      format: winston.format.combine(
        winston.format.colorize(),
        winston.format.simple()
      )
    })
  ]
});

module.exports = logger;"""
    
    def _get_service_env(self, service_name: str) -> str:
        return f"""# {service_name.title()} Service Configuration
NODE_ENV=development
PORT=3000

# Database
DB_HOST=localhost
DB_PORT=5432
DB_NAME={service_name}
DB_USER={service_name}_user
DB_PASSWORD={service_name}_password

# JWT
JWT_SECRET=your-jwt-secret

# Redis
REDIS_URL=redis://localhost:6379

# RabbitMQ
RABBITMQ_URL=amqp://localhost:5672

# Logging
LOG_LEVEL=info"""
    
    def _get_service_readme(self, service_name: str) -> str:
        return f"""# {service_name.title()} Service

Microservice for {service_name} functionality.

## Development

### Install dependencies
```bash
npm install
```

### Run locally
```bash
npm run dev
```

### Run tests
```bash
npm test
```

## API Endpoints

- `GET /api/v1/` - Service info
- `POST /api/v1/register` - Register user
- `POST /api/v1/login` - Login user
- `GET /api/v1/data` - Get data
- `POST /api/v1/data` - Create data
- `PUT /api/v1/data/:id` - Update data
- `DELETE /api/v1/data/:id` - Delete data

## Environment Variables

See `.env.example` for required environment variables."""
    
    def _get_docker_compose_dev(self) -> str:
        return """version: '3.8'

# Development overrides
services:
  gateway:
    volumes:
      - ./gateway:/gateway:ro
    environment:
      - DEBUG=true

  cache:
    ports:
      - "6379:6379"

  rabbitmq:
    environment:
      - RABBITMQ_DEFAULT_USER=admin
      - RABBITMQ_DEFAULT_PASS=admin"""
    
    def _get_gitignore(self) -> str:
        return """# Dependencies
node_modules/
package-lock.json

# Environment
.env
.env.local

# Logs
logs/
*.log

# IDE
.vscode/
.idea/
*.swp
.DS_Store

# Testing
coverage/
.nyc_output/

# Build
dist/
build/

# Docker
.docker/"""
    
    def _get_build_script(self) -> str:
        return """#!/bin/bash
set -e

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
"""
    
    def _get_start_script(self) -> str:
        return """#!/bin/bash
set -e

# Load environment variables
if [ -f .env ]; then
  export $(cat .env | grep -v '^#' | xargs)
fi

# Start services
echo "Starting services..."
docker-compose up -d

# Wait for services
echo "Waiting for services to be healthy..."
sleep 10

# Check health
echo "Checking service health..."
for service in $(docker-compose ps --services); do
  echo -n "Checking $service... "
  if docker-compose exec -T $service wget -q --spider http://localhost:3000/health 2>/dev/null; then
    echo "✓"
  else
    echo "✗"
  fi
done

echo ""
echo "Services available at:"
echo "  - API Gateway: http://localhost:8000"
echo "  - RabbitMQ Management: http://localhost:15672"
"""
    
    def _get_test_script(self) -> str:
        return """#!/bin/bash
set -e

echo "Running tests..."

# Run tests for each service
for service in services/*; do
  if [ -d "$service" ]; then
    service_name=$(basename $service)
    echo "Testing $service_name..."
    docker-compose exec -T $service_name-service npm test
  fi
done

echo "All tests complete!"
"""

def main():
    parser = argparse.ArgumentParser(description="Generate microservices architecture")
    parser.add_argument(
        "--name",
        required=True,
        help="Project name"
    )
    parser.add_argument(
        "--services",
        required=True,
        help="Comma-separated list of services (e.g., user,product,order)"
    )
    parser.add_argument(
        "--gateway",
        choices=["nginx", "kong"],
        default="nginx",
        help="API Gateway to use"
    )
    parser.add_argument(
        "--queue",
        choices=["rabbitmq", "redis"],
        default="rabbitmq",
        help="Message queue to use"
    )
    parser.add_argument(
        "--databases",
        default="postgres",
        help="Comma-separated list of databases (postgres,mongodb,mysql)"
    )
    parser.add_argument(
        "--output",
        help="Output directory"
    )
    
    args = parser.parse_args()
    
    # Parse services
    services = [s.strip() for s in args.services.split(',')]
    
    # Parse databases
    databases = [d.strip() for d in args.databases.split(',')]
    
    options = {
        "gateway": args.gateway,
        "queue": args.queue,
        "databases": databases,
        "output": args.output
    }
    
    generator = MicroservicesGenerator(
        name=args.name,
        services=services,
        options=options
    )
    
    try:
        generator.generate()
    except Exception as e:
        print(f"❌ Error: {e}", file=sys.stderr)
        sys.exit(1)

if __name__ == "__main__":
    main()