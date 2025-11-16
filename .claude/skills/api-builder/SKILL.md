---
name: API Builder
description: Generate REST/GraphQL APIs with consistent patterns, error handling, authentication, validation, and documentation. Supports Express, FastAPI, Django, and GraphQL with best practices built-in.
---

# API Builder

## Overview

This skill provides comprehensive API generation with production-ready patterns for REST and GraphQL APIs. It includes authentication, validation, error handling, rate limiting, caching, monitoring, and automatic documentation generation.

## Dependencies

- Python 3.8+ (for running the generator)
- Node.js 16+ (for Node.js frameworks)
- PostgreSQL/MySQL (for database)
- Redis (optional, for caching)

## Quick Start

### Generate Complete REST API
```bash
# Express.js API
python scripts/generate_api.py --type rest --framework express --name "UserAPI"

# FastAPI API
python scripts/generate_api.py --type rest --framework fastapi --name "UserAPI"

# Django REST API
python scripts/generate_api.py --type rest --framework django --name "UserAPI"
```

### Generate GraphQL API
```bash
# Apollo GraphQL
python scripts/generate_api.py --type graphql --framework apollo --name "UserAPI"

# GraphQL with Strawberry (Python)
python scripts/generate_api.py --type graphql --framework strawberry --name "UserAPI"
```

## Supported Frameworks

### REST APIs
- **Express.js** (Node.js)
- **FastAPI** (Python)
- **Django REST Framework** (Python)
- **Flask-RESTful** (Python)
- **NestJS** (TypeScript)
- **Spring Boot** (Java)

### GraphQL APIs
- **Apollo Server** (Node.js)
- **GraphQL Yoga** (Node.js)
- **Strawberry** (Python)
- **Graphene** (Python)

## API Patterns

### RESTful Design Principles
```yaml
patterns:
  resource_naming:
    - Use nouns, not verbs
    - Plural for collections (/users)
    - Singular for specific resources (/users/{id})
  
  http_methods:
    GET: Read operations
    POST: Create new resources
    PUT: Full update (replace)
    PATCH: Partial update
    DELETE: Remove resources
  
  status_codes:
    200: OK - Successful GET/PUT
    201: Created - Successful POST
    204: No Content - Successful DELETE
    400: Bad Request - Invalid input
    401: Unauthorized - Missing/invalid auth
    403: Forbidden - No permission
    404: Not Found - Resource doesn't exist
    409: Conflict - Duplicate resource
    422: Unprocessable Entity - Validation failed
    429: Too Many Requests - Rate limited
    500: Internal Server Error
```

## Authentication & Authorization

### JWT Authentication
```javascript
// Express.js Example
const jwt = require('jsonwebtoken');
const bcrypt = require('bcryptjs');

class AuthService {
  async generateToken(user) {
    const payload = {
      id: user.id,
      email: user.email,
      roles: user.roles
    };
    
    return jwt.sign(payload, process.env.JWT_SECRET, {
      expiresIn: '24h',
      issuer: 'api.example.com',
      audience: 'app.example.com'
    });
  }
  
  async verifyToken(token) {
    try {
      return jwt.verify(token, process.env.JWT_SECRET);
    } catch (error) {
      throw new UnauthorizedError('Invalid token');
    }
  }
  
  async hashPassword(password) {
    return bcrypt.hash(password, 12);
  }
  
  async comparePassword(password, hash) {
    return bcrypt.compare(password, hash);
  }
}
```

### OAuth2 Integration
```python
# FastAPI Example
from fastapi import Depends, HTTPException, status
from fastapi.security import OAuth2PasswordBearer
from jose import JWTError, jwt
from datetime import datetime, timedelta

oauth2_scheme = OAuth2PasswordBearer(tokenUrl="auth/token")

class OAuth2Service:
    SECRET_KEY = os.getenv("SECRET_KEY")
    ALGORITHM = "HS256"
    ACCESS_TOKEN_EXPIRE_MINUTES = 30
    
    def create_access_token(self, data: dict):
        to_encode = data.copy()
        expire = datetime.utcnow() + timedelta(minutes=self.ACCESS_TOKEN_EXPIRE_MINUTES)
        to_encode.update({"exp": expire})
        encoded_jwt = jwt.encode(to_encode, self.SECRET_KEY, algorithm=self.ALGORITHM)
        return encoded_jwt
    
    async def get_current_user(self, token: str = Depends(oauth2_scheme)):
        credentials_exception = HTTPException(
            status_code=status.HTTP_401_UNAUTHORIZED,
            detail="Could not validate credentials",
            headers={"WWW-Authenticate": "Bearer"},
        )
        
        try:
            payload = jwt.decode(token, self.SECRET_KEY, algorithms=[self.ALGORITHM])
            username: str = payload.get("sub")
            if username is None:
                raise credentials_exception
        except JWTError:
            raise credentials_exception
        
        user = await self.get_user(username=username)
        if user is None:
            raise credentials_exception
        return user
```

### API Key Authentication
```typescript
// NestJS Example
import { Injectable, CanActivate, ExecutionContext } from '@nestjs/common';
import { Observable } from 'rxjs';

@Injectable()
export class ApiKeyGuard implements CanActivate {
  canActivate(
    context: ExecutionContext,
  ): boolean | Promise<boolean> | Observable<boolean> {
    const request = context.switchToHttp().getRequest();
    const apiKey = request.headers['x-api-key'];
    
    if (!apiKey) {
      return false;
    }
    
    return this.validateApiKey(apiKey);
  }
  
  private async validateApiKey(apiKey: string): Promise<boolean> {
    // Check API key against database
    const key = await ApiKey.findOne({ key: apiKey, active: true });
    
    if (key) {
      // Update last used timestamp
      await key.update({ lastUsed: new Date() });
      return true;
    }
    
    return false;
  }
}
```

## Request Validation

### Schema Validation
```python
# FastAPI with Pydantic
from pydantic import BaseModel, Field, validator
from typing import Optional, List
from datetime import datetime
from email_validator import validate_email

class UserCreate(BaseModel):
    email: str = Field(..., description="User's email address")
    password: str = Field(..., min_length=8, description="Password (min 8 chars)")
    username: str = Field(..., min_length=3, max_length=50)
    full_name: Optional[str] = Field(None, max_length=100)
    age: Optional[int] = Field(None, ge=13, le=120)
    roles: List[str] = Field(default_factory=list)
    
    @validator('email')
    def email_valid(cls, v):
        try:
            validate_email(v)
        except:
            raise ValueError('Invalid email address')
        return v.lower()
    
    @validator('password')
    def password_strength(cls, v):
        if not any(char.isdigit() for char in v):
            raise ValueError('Password must contain at least one digit')
        if not any(char.isupper() for char in v):
            raise ValueError('Password must contain at least one uppercase letter')
        if not any(char.islower() for char in v):
            raise ValueError('Password must contain at least one lowercase letter')
        return v
    
    @validator('roles')
    def validate_roles(cls, v):
        valid_roles = {'admin', 'user', 'moderator', 'viewer'}
        for role in v:
            if role not in valid_roles:
                raise ValueError(f'Invalid role: {role}')
        return v
```

### Express.js Validation
```javascript
const { body, validationResult } = require('express-validator');

const userValidationRules = () => {
  return [
    body('email')
      .isEmail().withMessage('Must be a valid email')
      .normalizeEmail(),
    body('password')
      .isLength({ min: 8 }).withMessage('Password must be at least 8 characters')
      .matches(/^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[@$!%*?&])[A-Za-z\d@$!%*?&]/)
      .withMessage('Password must contain uppercase, lowercase, number and special character'),
    body('username')
      .isLength({ min: 3, max: 50 }).withMessage('Username must be 3-50 characters')
      .isAlphanumeric().withMessage('Username must be alphanumeric'),
    body('age')
      .optional()
      .isInt({ min: 13, max: 120 }).withMessage('Age must be between 13 and 120'),
    body('phoneNumber')
      .optional()
      .isMobilePhone().withMessage('Must be a valid phone number')
  ];
};

const validate = (req, res, next) => {
  const errors = validationResult(req);
  if (!errors.isEmpty()) {
    return res.status(422).json({ 
      error: 'Validation failed',
      details: errors.array() 
    });
  }
  next();
};
```

## Error Handling

### Centralized Error Handler
```typescript
// Express/TypeScript Error Handler
export class AppError extends Error {
  statusCode: number;
  isOperational: boolean;
  
  constructor(message: string, statusCode: number) {
    super(message);
    this.statusCode = statusCode;
    this.isOperational = true;
    
    Error.captureStackTrace(this, this.constructor);
  }
}

export class ValidationError extends AppError {
  constructor(message: string) {
    super(message, 422);
  }
}

export class AuthenticationError extends AppError {
  constructor(message: string = 'Authentication failed') {
    super(message, 401);
  }
}

export class AuthorizationError extends AppError {
  constructor(message: string = 'Permission denied') {
    super(message, 403);
  }
}

export class NotFoundError extends AppError {
  constructor(resource: string) {
    super(`${resource} not found`, 404);
  }
}

export class ConflictError extends AppError {
  constructor(message: string) {
    super(message, 409);
  }
}

// Global error handler middleware
export const errorHandler = (err: Error, req: Request, res: Response, next: NextFunction) => {
  if (!(err instanceof AppError)) {
    // Log unexpected errors
    logger.error('Unexpected error:', err);
    
    return res.status(500).json({
      error: 'Internal server error',
      message: process.env.NODE_ENV === 'development' ? err.message : 'Something went wrong'
    });
  }
  
  // Operational errors
  res.status(err.statusCode).json({
    error: err.name,
    message: err.message,
    ...(process.env.NODE_ENV === 'development' && { stack: err.stack })
  });
};
```

### FastAPI Error Handling
```python
from fastapi import HTTPException, Request, status
from fastapi.responses import JSONResponse
from fastapi.exceptions import RequestValidationError
from starlette.exceptions import HTTPException as StarletteHTTPException
import logging

logger = logging.getLogger(__name__)

class APIException(Exception):
    def __init__(self, status_code: int, detail: str, headers: dict = None):
        self.status_code = status_code
        self.detail = detail
        self.headers = headers

class ValidationException(APIException):
    def __init__(self, detail: str):
        super().__init__(status.HTTP_422_UNPROCESSABLE_ENTITY, detail)

class NotFoundException(APIException):
    def __init__(self, resource: str):
        super().__init__(status.HTTP_404_NOT_FOUND, f"{resource} not found")

class UnauthorizedException(APIException):
    def __init__(self, detail: str = "Unauthorized"):
        super().__init__(status.HTTP_401_UNAUTHORIZED, detail)

class ForbiddenException(APIException):
    def __init__(self, detail: str = "Forbidden"):
        super().__init__(status.HTTP_403_FORBIDDEN, detail)

async def api_exception_handler(request: Request, exc: APIException):
    return JSONResponse(
        status_code=exc.status_code,
        content={
            "error": exc.__class__.__name__,
            "detail": exc.detail,
            "path": request.url.path
        },
        headers=exc.headers
    )

async def validation_exception_handler(request: Request, exc: RequestValidationError):
    return JSONResponse(
        status_code=status.HTTP_422_UNPROCESSABLE_ENTITY,
        content={
            "error": "ValidationError",
            "detail": exc.errors(),
            "body": exc.body
        }
    )

async def generic_exception_handler(request: Request, exc: Exception):
    logger.error(f"Unhandled exception: {exc}", exc_info=True)
    return JSONResponse(
        status_code=status.HTTP_500_INTERNAL_SERVER_ERROR,
        content={
            "error": "InternalServerError",
            "detail": "An unexpected error occurred"
        }
    )
```

## Rate Limiting

### Redis-based Rate Limiting
```javascript
// Express.js with Redis
const redis = require('redis');
const { RateLimiterRedis } = require('rate-limiter-flexible');

const redisClient = redis.createClient({
  host: process.env.REDIS_HOST,
  port: process.env.REDIS_PORT,
});

const rateLimiter = new RateLimiterRedis({
  storeClient: redisClient,
  keyPrefix: 'rate_limit',
  points: 100, // Number of requests
  duration: 900, // Per 15 minutes
  blockDuration: 900, // Block for 15 minutes
});

const rateLimitMiddleware = async (req, res, next) => {
  try {
    const key = req.user ? req.user.id : req.ip;
    await rateLimiter.consume(key);
    
    // Add rate limit headers
    const rateLimitInfo = await rateLimiter.get(key);
    res.set({
      'X-RateLimit-Limit': rateLimiter.points,
      'X-RateLimit-Remaining': rateLimitInfo.remainingPoints,
      'X-RateLimit-Reset': rateLimitInfo.msBeforeNext
    });
    
    next();
  } catch (rejRes) {
    res.set({
      'Retry-After': Math.round(rejRes.msBeforeNext / 1000) || 1,
      'X-RateLimit-Limit': rateLimiter.points,
      'X-RateLimit-Remaining': rejRes.remainingPoints || 0,
      'X-RateLimit-Reset': rejRes.msBeforeNext
    });
    
    res.status(429).json({
      error: 'Too Many Requests',
      message: 'Rate limit exceeded',
      retryAfter: Math.round(rejRes.msBeforeNext / 1000)
    });
  }
};
```

### Python/FastAPI Rate Limiting
```python
from slowapi import Limiter, _rate_limit_exceeded_handler
from slowapi.util import get_remote_address
from slowapi.errors import RateLimitExceeded
from slowapi.middleware import SlowAPIMiddleware

limiter = Limiter(
    key_func=get_remote_address,
    default_limits=["100/15minute"],
    storage_uri="redis://localhost:6379"
)

app.state.limiter = limiter
app.add_exception_handler(RateLimitExceeded, _rate_limit_exceeded_handler)
app.add_middleware(SlowAPIMiddleware)

# Custom rate limits for specific endpoints
@app.post("/api/expensive-operation")
@limiter.limit("5/hour")
async def expensive_operation(request: Request):
    # Expensive operation
    pass

# User-based rate limiting
def get_user_id(request: Request):
    # Extract user ID from JWT token
    token = request.headers.get("Authorization")
    if token:
        payload = jwt.decode(token.replace("Bearer ", ""), SECRET_KEY)
        return payload.get("user_id")
    return get_remote_address(request)

@app.get("/api/user-data")
@limiter.limit("1000/hour", key_func=get_user_id)
async def get_user_data(request: Request):
    pass
```

## Caching Strategies

### Redis Caching
```python
import redis
import json
from functools import wraps
from typing import Optional
import hashlib

class CacheService:
    def __init__(self):
        self.redis_client = redis.Redis(
            host='localhost',
            port=6379,
            decode_responses=True
        )
    
    def cache_key(self, prefix: str, **kwargs) -> str:
        """Generate cache key from prefix and parameters"""
        key_data = json.dumps(kwargs, sort_keys=True)
        key_hash = hashlib.md5(key_data.encode()).hexdigest()
        return f"{prefix}:{key_hash}"
    
    def get(self, key: str) -> Optional[dict]:
        """Get cached data"""
        data = self.redis_client.get(key)
        if data:
            return json.loads(data)
        return None
    
    def set(self, key: str, value: dict, ttl: int = 3600):
        """Set cache with TTL"""
        self.redis_client.setex(
            key,
            ttl,
            json.dumps(value)
        )
    
    def invalidate(self, pattern: str):
        """Invalidate cache by pattern"""
        keys = self.redis_client.keys(pattern)
        if keys:
            self.redis_client.delete(*keys)
    
    def cache_decorator(self, prefix: str, ttl: int = 3600):
        """Decorator for caching function results"""
        def decorator(func):
            @wraps(func)
            async def wrapper(*args, **kwargs):
                # Generate cache key
                cache_key = self.cache_key(prefix, args=args, kwargs=kwargs)
                
                # Check cache
                cached = self.get(cache_key)
                if cached:
                    return cached
                
                # Execute function
                result = await func(*args, **kwargs)
                
                # Cache result
                self.set(cache_key, result, ttl)
                
                return result
            return wrapper
        return decorator

# Usage example
cache = CacheService()

@cache.cache_decorator("users", ttl=300)
async def get_user_by_id(user_id: int):
    # Expensive database query
    user = await db.users.find_one({"id": user_id})
    return user
```

## Database Integration

### ORM Setup (Sequelize)
```javascript
// models/index.js
const { Sequelize } = require('sequelize');

const sequelize = new Sequelize({
  dialect: 'postgres',
  host: process.env.DB_HOST,
  port: process.env.DB_PORT,
  database: process.env.DB_NAME,
  username: process.env.DB_USER,
  password: process.env.DB_PASSWORD,
  pool: {
    max: 20,
    min: 0,
    acquire: 30000,
    idle: 10000
  },
  logging: process.env.NODE_ENV === 'development' ? console.log : false
});

// models/user.js
module.exports = (sequelize, DataTypes) => {
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
    username: {
      type: DataTypes.STRING,
      unique: true,
      allowNull: false,
      validate: {
        len: [3, 50]
      }
    },
    password: {
      type: DataTypes.STRING,
      allowNull: false
    },
    isActive: {
      type: DataTypes.BOOLEAN,
      defaultValue: true
    },
    emailVerified: {
      type: DataTypes.BOOLEAN,
      defaultValue: false
    },
    lastLogin: {
      type: DataTypes.DATE
    }
  }, {
    timestamps: true,
    paranoid: true, // Soft deletes
    indexes: [
      {
        fields: ['email']
      },
      {
        fields: ['username']
      }
    ],
    hooks: {
      beforeCreate: async (user) => {
        user.password = await bcrypt.hash(user.password, 12);
      },
      beforeUpdate: async (user) => {
        if (user.changed('password')) {
          user.password = await bcrypt.hash(user.password, 12);
        }
      }
    }
  });
  
  User.prototype.validatePassword = async function(password) {
    return bcrypt.compare(password, this.password);
  };
  
  User.associate = (models) => {
    User.hasMany(models.Post);
    User.belongsToMany(models.Role, { through: 'UserRoles' });
  };
  
  return User;
};
```

### SQLAlchemy (Python)
```python
from sqlalchemy import create_engine, Column, String, Boolean, DateTime, ForeignKey
from sqlalchemy.ext.declarative import declarative_base
from sqlalchemy.orm import sessionmaker, relationship
from sqlalchemy.dialects.postgresql import UUID
import uuid
from datetime import datetime

Base = declarative_base()

class User(Base):
    __tablename__ = 'users'
    
    id = Column(UUID(as_uuid=True), primary_key=True, default=uuid.uuid4)
    email = Column(String(255), unique=True, nullable=False, index=True)
    username = Column(String(50), unique=True, nullable=False, index=True)
    password_hash = Column(String(255), nullable=False)
    is_active = Column(Boolean, default=True)
    email_verified = Column(Boolean, default=False)
    created_at = Column(DateTime, default=datetime.utcnow)
    updated_at = Column(DateTime, default=datetime.utcnow, onupdate=datetime.utcnow)
    last_login = Column(DateTime)
    
    # Relationships
    posts = relationship("Post", back_populates="author", cascade="all, delete-orphan")
    roles = relationship("Role", secondary="user_roles", back_populates="users")
    
    def set_password(self, password):
        self.password_hash = bcrypt.hashpw(password.encode('utf-8'), bcrypt.gensalt())
    
    def check_password(self, password):
        return bcrypt.checkpw(password.encode('utf-8'), self.password_hash)
    
    def to_dict(self):
        return {
            'id': str(self.id),
            'email': self.email,
            'username': self.username,
            'is_active': self.is_active,
            'email_verified': self.email_verified,
            'created_at': self.created_at.isoformat(),
            'roles': [role.name for role in self.roles]
        }

# Database connection
engine = create_engine(
    f"postgresql://{DB_USER}:{DB_PASSWORD}@{DB_HOST}:{DB_PORT}/{DB_NAME}",
    pool_size=20,
    max_overflow=40,
    pool_pre_ping=True,
    echo=False
)

SessionLocal = sessionmaker(autocommit=False, autoflush=False, bind=engine)

# Dependency for FastAPI
def get_db():
    db = SessionLocal()
    try:
        yield db
    finally:
        db.close()
```

## GraphQL Implementation

### Schema Definition
```graphql
# schema.graphql
type Query {
  user(id: ID!): User
  users(
    page: Int = 1
    limit: Int = 20
    search: String
    sortBy: UserSortField = CREATED_AT
    sortOrder: SortOrder = DESC
  ): UserConnection!
  me: User
}

type Mutation {
  createUser(input: CreateUserInput!): UserPayload!
  updateUser(id: ID!, input: UpdateUserInput!): UserPayload!
  deleteUser(id: ID!): DeletePayload!
  login(email: String!, password: String!): AuthPayload!
  refreshToken(token: String!): AuthPayload!
  logout: LogoutPayload!
}

type Subscription {
  userCreated: User!
  userUpdated(id: ID!): User!
}

type User {
  id: ID!
  email: String!
  username: String!
  fullName: String
  roles: [Role!]!
  posts(page: Int, limit: Int): PostConnection!
  createdAt: DateTime!
  updatedAt: DateTime!
}

type UserConnection {
  edges: [UserEdge!]!
  pageInfo: PageInfo!
  totalCount: Int!
}

type UserEdge {
  cursor: String!
  node: User!
}

type PageInfo {
  hasNextPage: Boolean!
  hasPreviousPage: Boolean!
  startCursor: String
  endCursor: String
}

input CreateUserInput {
  email: String!
  username: String!
  password: String!
  fullName: String
  roles: [String!]
}

input UpdateUserInput {
  email: String
  username: String
  fullName: String
  roles: [String!]
}

type UserPayload {
  user: User
  errors: [Error!]
}

type AuthPayload {
  token: String
  refreshToken: String
  user: User
  errors: [Error!]
}

type Error {
  field: String
  message: String!
}

enum UserSortField {
  CREATED_AT
  UPDATED_AT
  USERNAME
  EMAIL
}

enum SortOrder {
  ASC
  DESC
}

scalar DateTime
```

## Best Practices Summary

1. **Design First**: Define API contracts before implementation
2. **Version Your API**: Use URL versioning (/v1, /v2)
3. **Use Proper HTTP Methods**: GET for read, POST for create, etc.
4. **Implement Pagination**: Never return unlimited results
5. **Rate Limit**: Protect against abuse
6. **Cache Aggressively**: Use Redis/CDN for performance
7. **Validate Everything**: Never trust client input
8. **Handle Errors Gracefully**: Consistent error format
9. **Secure by Default**: HTTPS, authentication, authorization
10. **Document Everything**: OpenAPI/GraphQL schema
11. **Monitor & Log**: Track performance and errors
12. **Test Thoroughly**: Unit, integration, and load tests

## Quick Reference

### Common HTTP Status Codes
- **2xx Success**: 200 OK, 201 Created, 204 No Content
- **3xx Redirection**: 301 Moved, 304 Not Modified
- **4xx Client Error**: 400 Bad Request, 401 Unauthorized, 403 Forbidden, 404 Not Found, 422 Unprocessable
- **5xx Server Error**: 500 Internal Error, 502 Bad Gateway, 503 Service Unavailable

### RESTful URL Examples
```
GET    /api/v1/users           # List users
GET    /api/v1/users/123       # Get specific user
POST   /api/v1/users           # Create user
PUT    /api/v1/users/123       # Replace user
PATCH  /api/v1/users/123       # Update user
DELETE /api/v1/users/123       # Delete user
GET    /api/v1/users/123/posts # Get user's posts
```

### GraphQL Query Examples
```graphql
# Query
query GetUser($id: ID!) {
  user(id: $id) {
    id
    email
    posts {
      title
      content
    }
  }
}

# Mutation
mutation CreateUser($input: CreateUserInput!) {
  createUser(input: $input) {
    user {
      id
      email
    }
    errors {
      field
      message
    }
  }
}

# Subscription
subscription OnUserCreated {
  userCreated {
    id
    email
    createdAt
  }
}
```