#!/usr/bin/env python3
"""
Template Generator - Additional API Components
Generates validators, tests, and configuration files
"""

from typing import Dict

class TemplateGenerator:
    @staticmethod
    def get_express_user_controller() -> str:
        return """const { User } = require('../models');
const { NotFoundError, ForbiddenError } = require('../utils/errors');
const logger = require('../utils/logger');
const CacheService = require('../services/cacheService');

const cache = new CacheService();

class UserController {
  async getUsers(req, res, next) {
    try {
      const { page = 1, limit = 20, search, sortBy = 'createdAt', sortOrder = 'DESC' } = req.query;
      
      const cacheKey = `users:${page}:${limit}:${search}:${sortBy}:${sortOrder}`;
      const cached = await cache.get(cacheKey);
      
      if (cached) {
        return res.json(cached);
      }
      
      const offset = (page - 1) * limit;
      const where = search ? {
        [Op.or]: [
          { username: { [Op.iLike]: `%${search}%` } },
          { email: { [Op.iLike]: `%${search}%` } },
          { fullName: { [Op.iLike]: `%${search}%` } }
        ]
      } : {};
      
      const { rows, count } = await User.findAndCountAll({
        where,
        limit: parseInt(limit),
        offset,
        order: [[sortBy, sortOrder]],
        attributes: { exclude: ['password', 'refreshToken'] }
      });
      
      const result = {
        users: rows,
        total: count,
        page: parseInt(page),
        pages: Math.ceil(count / limit),
        hasNext: page < Math.ceil(count / limit),
        hasPrev: page > 1
      };
      
      await cache.set(cacheKey, result, 60); // Cache for 1 minute
      
      res.json(result);
    } catch (error) {
      next(error);
    }
  }
  
  async getUser(req, res, next) {
    try {
      const { id } = req.params;
      
      const user = await User.findByPk(id, {
        attributes: { exclude: ['password', 'refreshToken'] }
      });
      
      if (!user) {
        throw new NotFoundError('User');
      }
      
      res.json(user);
    } catch (error) {
      next(error);
    }
  }
  
  async updateUser(req, res, next) {
    try {
      const { id } = req.params;
      const updates = req.body;
      
      // Check permissions
      if (req.user.id !== id && req.user.role !== 'admin') {
        throw new ForbiddenError('You can only update your own profile');
      }
      
      const user = await User.findByPk(id);
      
      if (!user) {
        throw new NotFoundError('User');
      }
      
      // Prevent role escalation
      if (updates.role && req.user.role !== 'admin') {
        delete updates.role;
      }
      
      await user.update(updates);
      
      // Invalidate cache
      await cache.invalidate('users:*');
      
      logger.info(`User updated: ${user.email}`);
      
      res.json({
        message: 'User updated successfully',
        user: user.toJSON()
      });
    } catch (error) {
      next(error);
    }
  }
  
  async deleteUser(req, res, next) {
    try {
      const { id } = req.params;
      
      // Check permissions
      if (req.user.id !== id && req.user.role !== 'admin') {
        throw new ForbiddenError('You can only delete your own account');
      }
      
      const user = await User.findByPk(id);
      
      if (!user) {
        throw new NotFoundError('User');
      }
      
      await user.destroy(); // Soft delete if paranoid is enabled
      
      // Invalidate cache
      await cache.invalidate('users:*');
      
      logger.info(`User deleted: ${user.email}`);
      
      res.status(204).send();
    } catch (error) {
      next(error);
    }
  }
  
  async getCurrentUser(req, res) {
    res.json(req.user);
  }
  
  async updateProfile(req, res, next) {
    try {
      const updates = req.body;
      const user = await User.findByPk(req.user.id);
      
      // Remove fields that users can't update themselves
      delete updates.role;
      delete updates.emailVerified;
      delete updates.isActive;
      
      await user.update(updates);
      
      res.json({
        message: 'Profile updated successfully',
        user: user.toJSON()
      });
    } catch (error) {
      next(error);
    }
  }
}

module.exports = new UserController();
"""

    @staticmethod
    def get_express_validation_rules() -> str:
        return """const { body, param, query, validationResult } = require('express-validator');

// Validation middleware
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

// Auth validations
const registerValidation = [
  body('email')
    .isEmail().withMessage('Must be a valid email')
    .normalizeEmail(),
  body('username')
    .isLength({ min: 3, max: 50 }).withMessage('Username must be 3-50 characters')
    .isAlphanumeric().withMessage('Username must be alphanumeric'),
  body('password')
    .isLength({ min: 8 }).withMessage('Password must be at least 8 characters')
    .matches(/^(?=.*[a-z])(?=.*[A-Z])(?=.*\\d)(?=.*[@$!%*?&])[A-Za-z\\d@$!%*?&]/)
    .withMessage('Password must contain uppercase, lowercase, number and special character'),
  body('fullName')
    .optional()
    .isLength({ max: 100 }).withMessage('Full name must not exceed 100 characters'),
  validate
];

const loginValidation = [
  body('email')
    .isEmail().withMessage('Must be a valid email')
    .normalizeEmail(),
  body('password')
    .notEmpty().withMessage('Password is required'),
  validate
];

// User validations
const updateUserValidation = [
  param('id')
    .isUUID().withMessage('Invalid user ID'),
  body('email')
    .optional()
    .isEmail().withMessage('Must be a valid email')
    .normalizeEmail(),
  body('username')
    .optional()
    .isLength({ min: 3, max: 50 }).withMessage('Username must be 3-50 characters')
    .isAlphanumeric().withMessage('Username must be alphanumeric'),
  body('fullName')
    .optional()
    .isLength({ max: 100 }).withMessage('Full name must not exceed 100 characters'),
  body('role')
    .optional()
    .isIn(['admin', 'user', 'moderator']).withMessage('Invalid role'),
  validate
];

const getUsersValidation = [
  query('page')
    .optional()
    .isInt({ min: 1 }).withMessage('Page must be a positive integer'),
  query('limit')
    .optional()
    .isInt({ min: 1, max: 100 }).withMessage('Limit must be between 1 and 100'),
  query('sortBy')
    .optional()
    .isIn(['createdAt', 'updatedAt', 'username', 'email']).withMessage('Invalid sort field'),
  query('sortOrder')
    .optional()
    .isIn(['ASC', 'DESC']).withMessage('Sort order must be ASC or DESC'),
  validate
];

module.exports = {
  validate,
  registerValidation,
  loginValidation,
  updateUserValidation,
  getUsersValidation
};
"""

    @staticmethod
    def get_fastapi_user_endpoint() -> str:
        return """from typing import List, Optional
from fastapi import APIRouter, Depends, HTTPException, Query, status
from sqlalchemy.orm import Session
from sqlalchemy import or_

from app.db.session import get_db
from app.models.user import User
from app.schemas.user import UserCreate, UserUpdate, UserResponse, UserList
from app.core.security import get_current_user, get_current_admin_user
from app.services.cache_service import cache_service

router = APIRouter()

@router.get("/", response_model=UserList)
async def get_users(
    page: int = Query(1, ge=1),
    limit: int = Query(20, ge=1, le=100),
    search: Optional[str] = None,
    sort_by: str = Query("created_at", regex="^(created_at|updated_at|username|email)$"),
    sort_order: str = Query("desc", regex="^(asc|desc)$"),
    db: Session = Depends(get_db),
    current_user: User = Depends(get_current_user)
):
    \"\"\"Get paginated list of users\"\"\"
    
    # Check cache
    cache_key = f"users:{page}:{limit}:{search}:{sort_by}:{sort_order}"
    cached = await cache_service.get(cache_key)
    if cached:
        return cached
    
    # Build query
    query = db.query(User)
    
    if search:
        query = query.filter(
            or_(
                User.username.ilike(f"%{search}%"),
                User.email.ilike(f"%{search}%"),
                User.full_name.ilike(f"%{search}%")
            )
        )
    
    # Get total count
    total = query.count()
    
    # Apply sorting
    order_column = getattr(User, sort_by)
    if sort_order == "desc":
        order_column = order_column.desc()
    query = query.order_by(order_column)
    
    # Apply pagination
    offset = (page - 1) * limit
    users = query.offset(offset).limit(limit).all()
    
    result = {
        "users": users,
        "total": total,
        "page": page,
        "pages": (total + limit - 1) // limit,
        "has_next": page < (total + limit - 1) // limit,
        "has_prev": page > 1
    }
    
    # Cache result
    await cache_service.set(cache_key, result, ttl=60)
    
    return result

@router.get("/me", response_model=UserResponse)
async def get_current_user_profile(
    current_user: User = Depends(get_current_user)
):
    \"\"\"Get current user profile\"\"\"
    return current_user

@router.get("/{user_id}", response_model=UserResponse)
async def get_user(
    user_id: str,
    db: Session = Depends(get_db),
    current_user: User = Depends(get_current_user)
):
    \"\"\"Get user by ID\"\"\"
    user = db.query(User).filter(User.id == user_id).first()
    if not user:
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND,
            detail="User not found"
        )
    return user

@router.patch("/{user_id}", response_model=UserResponse)
async def update_user(
    user_id: str,
    user_update: UserUpdate,
    db: Session = Depends(get_db),
    current_user: User = Depends(get_current_user)
):
    \"\"\"Update user\"\"\"
    # Check permissions
    if current_user.id != user_id and current_user.role != "admin":
        raise HTTPException(
            status_code=status.HTTP_403_FORBIDDEN,
            detail="You can only update your own profile"
        )
    
    user = db.query(User).filter(User.id == user_id).first()
    if not user:
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND,
            detail="User not found"
        )
    
    # Prevent role escalation
    update_data = user_update.dict(exclude_unset=True)
    if "role" in update_data and current_user.role != "admin":
        del update_data["role"]
    
    for field, value in update_data.items():
        setattr(user, field, value)
    
    db.commit()
    db.refresh(user)
    
    # Invalidate cache
    await cache_service.invalidate("users:*")
    
    return user

@router.delete("/{user_id}", status_code=status.HTTP_204_NO_CONTENT)
async def delete_user(
    user_id: str,
    db: Session = Depends(get_db),
    current_user: User = Depends(get_current_admin_user)
):
    \"\"\"Delete user (admin only)\"\"\"
    user = db.query(User).filter(User.id == user_id).first()
    if not user:
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND,
            detail="User not found"
        )
    
    db.delete(user)
    db.commit()
    
    # Invalidate cache
    await cache_service.invalidate("users:*")
    
    return None
"""

    @staticmethod
    def get_graphql_schema() -> str:
        return """type Query {
  # User queries
  user(id: ID!): User
  users(
    page: Int = 1
    limit: Int = 20
    search: String
    sortBy: UserSortField = CREATED_AT
    sortOrder: SortOrder = DESC
  ): UserConnection!
  me: User
  
  # Search
  searchUsers(query: String!, limit: Int = 10): [User!]!
}

type Mutation {
  # Authentication
  register(input: RegisterInput!): AuthPayload!
  login(email: String!, password: String!): AuthPayload!
  refreshToken(token: String!): AuthPayload!
  logout: LogoutPayload!
  
  # User management
  updateUser(id: ID!, input: UpdateUserInput!): UserPayload!
  deleteUser(id: ID!): DeletePayload!
  changePassword(oldPassword: String!, newPassword: String!): UserPayload!
  
  # Admin operations
  createUser(input: CreateUserInput!): UserPayload! @auth(requires: ADMIN)
  activateUser(id: ID!): UserPayload! @auth(requires: ADMIN)
  deactivateUser(id: ID!): UserPayload! @auth(requires: ADMIN)
}

type Subscription {
  userCreated: User! @auth
  userUpdated(id: ID!): User! @auth
  userDeleted(id: ID!): ID! @auth(requires: ADMIN)
}

# Types
type User {
  id: ID!
  email: String!
  username: String!
  fullName: String
  role: Role!
  isActive: Boolean!
  emailVerified: Boolean!
  createdAt: DateTime!
  updatedAt: DateTime!
  lastLogin: DateTime
  posts(page: Int, limit: Int): PostConnection!
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
  page: Int!
  pages: Int!
}

type AuthPayload {
  user: User
  accessToken: String
  refreshToken: String
  expiresIn: Int
  errors: [Error!]
}

type UserPayload {
  user: User
  errors: [Error!]
}

type LogoutPayload {
  success: Boolean!
  message: String
}

type DeletePayload {
  success: Boolean!
  message: String
  deletedId: ID
}

type Error {
  field: String
  message: String!
  code: String
}

# Inputs
input RegisterInput {
  email: String!
  username: String!
  password: String!
  fullName: String
}

input CreateUserInput {
  email: String!
  username: String!
  password: String!
  fullName: String
  role: Role = USER
  isActive: Boolean = true
}

input UpdateUserInput {
  email: String
  username: String
  fullName: String
  role: Role
  isActive: Boolean
}

# Enums
enum Role {
  ADMIN
  USER
  MODERATOR
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

# Custom scalars
scalar DateTime
scalar Email
scalar URL

# Directives
directive @auth(requires: Role) on FIELD_DEFINITION
directive @rateLimit(max: Int!, window: Int!) on FIELD_DEFINITION
directive @deprecated(reason: String) on FIELD_DEFINITION | ENUM_VALUE
"""

    @staticmethod
    def get_api_tests() -> str:
        return """import pytest
from fastapi.testclient import TestClient
from sqlalchemy import create_engine
from sqlalchemy.orm import sessionmaker

from app.main import app
from app.db.base import Base
from app.db.session import get_db
from app.core.security import create_access_token

# Create test database
SQLALCHEMY_DATABASE_URL = "sqlite:///./test.db"
engine = create_engine(SQLALCHEMY_DATABASE_URL, connect_args={"check_same_thread": False})
TestingSessionLocal = sessionmaker(autocommit=False, autoflush=False, bind=engine)

Base.metadata.create_all(bind=engine)

def override_get_db():
    try:
        db = TestingSessionLocal()
        yield db
    finally:
        db.close()

app.dependency_overrides[get_db] = override_get_db

client = TestClient(app)

class TestAuth:
    def test_register_user(self):
        response = client.post(
            "/api/v1/auth/register",
            json={
                "email": "test@example.com",
                "username": "testuser",
                "password": "Test123!@#",
                "fullName": "Test User"
            }
        )
        assert response.status_code == 201
        data = response.json()
        assert data["user"]["email"] == "test@example.com"
        assert "accessToken" in data
        assert "refreshToken" in data

    def test_login_user(self):
        # First register
        client.post(
            "/api/v1/auth/register",
            json={
                "email": "login@example.com",
                "username": "loginuser",
                "password": "Test123!@#"
            }
        )
        
        # Then login
        response = client.post(
            "/api/v1/auth/login",
            json={
                "email": "login@example.com",
                "password": "Test123!@#"
            }
        )
        assert response.status_code == 200
        data = response.json()
        assert "accessToken" in data
        assert "refreshToken" in data

    def test_login_invalid_credentials(self):
        response = client.post(
            "/api/v1/auth/login",
            json={
                "email": "wrong@example.com",
                "password": "WrongPassword"
            }
        )
        assert response.status_code == 401

class TestUsers:
    @pytest.fixture
    def auth_headers(self):
        # Create a test user and get token
        response = client.post(
            "/api/v1/auth/register",
            json={
                "email": "auth@example.com",
                "username": "authuser",
                "password": "Test123!@#"
            }
        )
        token = response.json()["accessToken"]
        return {"Authorization": f"Bearer {token}"}

    def test_get_users(self, auth_headers):
        response = client.get("/api/v1/users", headers=auth_headers)
        assert response.status_code == 200
        data = response.json()
        assert "users" in data
        assert "total" in data
        assert "page" in data

    def test_get_current_user(self, auth_headers):
        response = client.get("/api/v1/users/me", headers=auth_headers)
        assert response.status_code == 200
        data = response.json()
        assert data["email"] == "auth@example.com"

    def test_update_user(self, auth_headers):
        # Get current user ID
        response = client.get("/api/v1/users/me", headers=auth_headers)
        user_id = response.json()["id"]
        
        # Update user
        response = client.patch(
            f"/api/v1/users/{user_id}",
            headers=auth_headers,
            json={"fullName": "Updated Name"}
        )
        assert response.status_code == 200
        data = response.json()
        assert data["fullName"] == "Updated Name"

    def test_unauthorized_access(self):
        response = client.get("/api/v1/users")
        assert response.status_code == 401

class TestValidation:
    def test_invalid_email(self):
        response = client.post(
            "/api/v1/auth/register",
            json={
                "email": "invalid-email",
                "username": "testuser",
                "password": "Test123!@#"
            }
        )
        assert response.status_code == 422

    def test_weak_password(self):
        response = client.post(
            "/api/v1/auth/register",
            json={
                "email": "test@example.com",
                "username": "testuser",
                "password": "weak"
            }
        )
        assert response.status_code == 422

    def test_short_username(self):
        response = client.post(
            "/api/v1/auth/register",
            json={
                "email": "test@example.com",
                "username": "ab",
                "password": "Test123!@#"
            }
        )
        assert response.status_code == 422
"""

# Create generator instance
generator = TemplateGenerator()
