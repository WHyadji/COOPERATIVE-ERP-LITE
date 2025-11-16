# API Builder Skill

A comprehensive API generation skill that creates production-ready REST and GraphQL APIs with authentication, validation, error handling, and complete documentation.

## Features

### 🚀 Core Capabilities
- **REST API Generation**: Express.js, FastAPI, Django REST, Flask, NestJS, Spring Boot
- **GraphQL API Generation**: Apollo Server, GraphQL Yoga, Strawberry, Graphene
- **Authentication**: JWT, OAuth2, API Keys
- **Validation**: Schema validation, input sanitization
- **Error Handling**: Centralized error handling with proper status codes
- **Rate Limiting**: Redis-based rate limiting with configurable limits
- **Caching**: Redis caching with TTL and invalidation
- **Documentation**: Auto-generated OpenAPI/Swagger or GraphQL schema
- **Testing**: Unit and integration tests included
- **Monitoring**: Structured logging and metrics collection
- **Security**: CORS, helmet, input sanitization, SQL injection prevention
- **Database**: ORM setup (Sequelize, SQLAlchemy, Prisma)
- **Deployment**: Docker, Kubernetes, CI/CD pipelines

## Quick Start

### Generate REST API

```bash
# Express.js API
python scripts/generate_api.py --type rest --framework express --name "MyAPI"

# FastAPI API
python scripts/generate_api.py --type rest --framework fastapi --name "MyAPI"

# Django REST API
python scripts/generate_api.py --type rest --framework django --name "MyAPI"
```

### Generate GraphQL API

```bash
# Apollo GraphQL
python scripts/generate_api.py --type graphql --framework apollo --name "MyGraphQL"

# Strawberry GraphQL (Python)
python scripts/generate_api.py --type graphql --framework strawberry --name "MyGraphQL"
```

## Generated Project Structure

### REST API Structure
```
my-api/
├── src/                    # Source code
│   ├── controllers/        # Request handlers
│   ├── models/            # Database models
│   ├── routes/            # Route definitions
│   ├── middleware/        # Custom middleware
│   ├── services/          # Business logic
│   ├── validators/        # Input validation
│   ├── utils/            # Utility functions
│   └── config/           # Configuration files
├── tests/                 # Test files
│   ├── unit/             # Unit tests
│   └── integration/      # Integration tests
├── docs/                  # Documentation
├── logs/                  # Application logs
├── .env.example          # Environment variables template
├── Dockerfile            # Docker configuration
├── docker-compose.yml    # Docker Compose setup
└── README.md            # Project documentation
```

### GraphQL API Structure
```
my-graphql/
├── src/
│   ├── schema/           # GraphQL schema definitions
│   ├── resolvers/        # GraphQL resolvers
│   ├── dataSources/      # Data source classes
│   ├── models/           # Database models
│   ├── directives/       # Custom GraphQL directives
│   ├── middleware/       # GraphQL middleware
│   └── utils/           # Utility functions
├── tests/               # Test files
├── .env.example        # Environment variables
├── Dockerfile          # Docker configuration
└── README.md          # Documentation
```

## API Patterns Included

### Authentication & Authorization
- **JWT**: Access and refresh tokens with configurable expiry
- **OAuth2**: Google, GitHub, Facebook integration ready
- **API Keys**: Key-based authentication for service-to-service
- **Role-Based Access Control**: Admin, User, Moderator roles
- **Permission System**: Fine-grained permissions

### Request Validation
- **Schema Validation**: Automatic request/response validation
- **Type Checking**: Strong typing for all inputs
- **Custom Validators**: Email, phone, password strength
- **Sanitization**: XSS and SQL injection prevention

### Error Handling
- **Standardized Errors**: Consistent error format across API
- **HTTP Status Codes**: Proper status codes for all responses
- **Error Logging**: Structured error logging
- **User-Friendly Messages**: Clear error messages

### Performance Features
- **Caching**: Redis-based response caching
- **Rate Limiting**: Per-user and per-IP limiting
- **Pagination**: Cursor and offset-based pagination
- **Query Optimization**: N+1 query prevention
- **Response Compression**: Gzip compression

### Security Features
- **CORS**: Configurable CORS policies
- **Helmet**: Security headers
- **Input Sanitization**: HTML and SQL sanitization
- **Password Hashing**: Bcrypt with salt rounds
- **HTTPS Ready**: SSL/TLS configuration included

## Configuration Options

### REST API Frameworks

#### Express.js (Node.js)
- **ORM**: Sequelize with PostgreSQL/MySQL
- **Validation**: express-validator
- **Documentation**: Swagger/OpenAPI
- **Testing**: Jest + Supertest

#### FastAPI (Python)
- **ORM**: SQLAlchemy with PostgreSQL/MySQL
- **Validation**: Pydantic
- **Documentation**: Auto-generated OpenAPI
- **Testing**: Pytest

#### Django REST Framework
- **ORM**: Django ORM
- **Validation**: Django serializers
- **Documentation**: DRF spectacular
- **Testing**: Django test framework

### GraphQL Frameworks

#### Apollo Server (Node.js)
- **Database**: MongoDB with Mongoose / PostgreSQL with Prisma
- **Subscriptions**: WebSocket support
- **Federation**: Apollo Federation ready
- **Testing**: Jest with apollo-server-testing

#### Strawberry (Python)
- **Database**: SQLAlchemy
- **Type Safety**: Python type hints
- **Async Support**: Full async/await
- **Testing**: Pytest

## Environment Variables

The generated API includes comprehensive environment configuration:

```env
# Application
NODE_ENV=development
PORT=3000
API_VERSION=v1

# Database
DB_HOST=localhost
DB_PORT=5432
DB_NAME=api_db
DB_USER=postgres
DB_PASSWORD=postgres

# Authentication
JWT_SECRET=your-secret-key
JWT_EXPIRY=24h
REFRESH_TOKEN_EXPIRY=7d

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379

# Rate Limiting
RATE_LIMIT_WINDOW=15
RATE_LIMIT_MAX_REQUESTS=100

# External Services
SENDGRID_API_KEY=
AWS_ACCESS_KEY_ID=
GOOGLE_OAUTH_CLIENT_ID=
```

## Deployment

### Docker
Every generated API includes:
- Multi-stage Dockerfile for optimal image size
- Docker Compose for local development
- Health checks and graceful shutdown
- Non-root user configuration

### Kubernetes
- Deployment manifests
- Service configuration
- ConfigMaps and Secrets
- Horizontal Pod Autoscaling
- Health and readiness probes

### CI/CD
- GitHub Actions workflow
- GitLab CI configuration
- Testing pipeline
- Security scanning
- Automated deployment

## Testing

### Unit Tests
- Model tests
- Service tests
- Utility function tests
- Validation tests

### Integration Tests
- API endpoint tests
- Authentication flow tests
- Database integration tests
- Cache integration tests

### Load Testing
- K6 scripts included
- Stress test configurations
- Performance benchmarks

## Documentation

Each generated API includes:
- **README.md**: Comprehensive project documentation
- **API Documentation**: Swagger UI or GraphQL Playground
- **Postman Collection**: Ready-to-import collection
- **Code Comments**: Inline documentation
- **Architecture Diagram**: System design documentation

## Monitoring & Logging

### Logging
- Structured JSON logging
- Log levels (debug, info, warn, error)
- Request/response logging
- Error tracking
- Log aggregation ready

### Metrics
- Prometheus metrics endpoint
- Response time tracking
- Error rate monitoring
- Request count metrics
- Custom business metrics

## Best Practices Implemented

1. **Clean Architecture**: Separation of concerns
2. **SOLID Principles**: Single responsibility, dependency injection
3. **12-Factor App**: Environment-based configuration
4. **RESTful Design**: Proper HTTP methods and status codes
5. **GraphQL Best Practices**: Schema-first design
6. **Security First**: Input validation, authentication, authorization
7. **Performance**: Caching, pagination, query optimization
8. **Testing**: Comprehensive test coverage
9. **Documentation**: Auto-generated and maintained
10. **DevOps Ready**: Containerized and CI/CD enabled

## Customization

After generation, you can customize:
- Add custom endpoints/resolvers
- Modify validation rules
- Extend authentication strategies
- Add business logic
- Configure additional middleware
- Integrate external services

## Support Features

- **Hot Reloading**: Development server with auto-restart
- **Debug Mode**: Enhanced error messages in development
- **Database Migrations**: Version control for schema changes
- **Seed Data**: Sample data for development
- **API Versioning**: Support for multiple API versions
- **Internationalization**: i18n support structure

## Generated Files Overview

| File | Purpose |
|------|---------|
| `server.js` / `main.py` | Application entry point |
| `app.js` / `app.py` | Application configuration |
| `routes/` / `endpoints/` | API endpoint definitions |
| `models/` | Database models |
| `middleware/` | Request/response processing |
| `services/` | Business logic layer |
| `utils/` | Helper functions |
| `config/` | Configuration management |
| `tests/` | Test suites |
| `Dockerfile` | Container configuration |
| `docker-compose.yml` | Local development setup |
| `.env.example` | Environment variables template |
| `.github/workflows/` | CI/CD pipelines |

## Requirements

- Python 3.8+ (for running the generator)
- Node.js 16+ (for Node.js frameworks)
- Docker (optional, for containerization)
- PostgreSQL/MySQL/MongoDB (database)
- Redis (optional, for caching/rate limiting)

## License

This skill is provided as-is for use with Claude AI. The generated code is yours to use and modify as needed.