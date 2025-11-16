#!/usr/bin/env python3
"""
API Generator - Main Script
Generates complete REST or GraphQL APIs with chosen framework
"""

import os
import sys
import json
import argparse
from pathlib import Path
from typing import Dict, List, Optional
from datetime import datetime
import shutil

# Import template generator
sys.path.append(str(Path(__file__).parent))
from templates import TemplateGenerator

class APIGenerator:
    def __init__(self, api_type: str, framework: str, name: str, options: Dict = None):
        self.api_type = api_type.lower()
        self.framework = framework.lower()
        self.name = name
        self.options = options or {}
        self.output_dir = Path(self.options.get('output', f"./{name.lower().replace(' ', '-')}"))
        self.templates = TemplateGenerator()
        
        # Framework mappings
        self.rest_frameworks = {
            'express': self._generate_express_api,
            'fastapi': self._generate_fastapi_api,
            'django': self._generate_django_api,
            'flask': self._generate_flask_api,
            'nestjs': self._generate_nestjs_api
        }
        
        self.graphql_frameworks = {
            'apollo': self._generate_apollo_api,
            'yoga': self._generate_yoga_api,
            'strawberry': self._generate_strawberry_api,
            'graphene': self._generate_graphene_api
        }
    
    def generate(self):
        """Main generation process"""
        print(f"🚀 Generating {self.api_type.upper()} API with {self.framework}...")
        
        # Create output directory
        self.output_dir.mkdir(parents=True, exist_ok=True)
        
        # Generate based on type
        if self.api_type == 'rest':
            if self.framework in self.rest_frameworks:
                self.rest_frameworks[self.framework]()
            else:
                raise ValueError(f"Unknown REST framework: {self.framework}")
        elif self.api_type == 'graphql':
            if self.framework in self.graphql_frameworks:
                self.graphql_frameworks[self.framework]()
            else:
                raise ValueError(f"Unknown GraphQL framework: {self.framework}")
        else:
            raise ValueError(f"Unknown API type: {self.api_type}")
        
        print(f"✅ API generated successfully in {self.output_dir}")
        print(f"\n📖 Next steps:")
        print(f"   1. cd {self.output_dir}")
        print(f"   2. Review README.md for setup instructions")
        print(f"   3. Configure environment variables in .env")
        print(f"   4. Install dependencies and run")
    
    def _create_file(self, path: str, content: str):
        """Create a file with content"""
        file_path = self.output_dir / path
        file_path.parent.mkdir(parents=True, exist_ok=True)
        file_path.write_text(content)
        print(f"    ✓ Created {path}")
    
    def _generate_express_api(self):
        """Generate Express.js REST API"""
        print("  📦 Creating Express.js project structure...")
        
        # Create directory structure
        dirs = [
            'src/controllers', 'src/models', 'src/routes',
            'src/middleware', 'src/services', 'src/validators',
            'src/utils', 'src/config', 'tests/unit',
            'tests/integration', 'logs', 'docs'
        ]
        
        for dir_path in dirs:
            (self.output_dir / dir_path).mkdir(parents=True, exist_ok=True)
        
        # Generate files using templates
        self._create_file('package.json', self._get_express_package_json())
        self._create_file('src/server.js', self._get_express_server())
        self._create_file('src/app.js', self._get_express_app())
        self._create_file('src/controllers/userController.js', self.templates.get_express_user_controller())
        self._create_file('src/middleware/validation.js', self.templates.get_express_validation_rules())
        self._create_file('.env.example', self._get_env_example())
        self._create_file('README.md', self._get_readme())
        self._create_file('Dockerfile', self._get_dockerfile())
        self._create_file('docker-compose.yml', self._get_docker_compose())
        self._create_file('.gitignore', self._get_gitignore())
        
        print("  ✅ Express.js API generated!")
    
    def _generate_fastapi_api(self):
        """Generate FastAPI REST API"""
        print("  📦 Creating FastAPI project structure...")
        
        # Create directory structure
        dirs = [
            'app/api/v1/endpoints', 'app/api/v1/middleware',
            'app/core', 'app/models', 'app/schemas',
            'app/services', 'app/db', 'app/utils',
            'tests', 'alembic/versions', 'docs'
        ]
        
        for dir_path in dirs:
            (self.output_dir / dir_path).mkdir(parents=True, exist_ok=True)
        
        # Generate files
        self._create_file('requirements.txt', self._get_fastapi_requirements())
        self._create_file('app/main.py', self._get_fastapi_main())
        self._create_file('app/api/v1/endpoints/users.py', self.templates.get_fastapi_user_endpoint())
        self._create_file('tests/test_api.py', self.templates.get_api_tests())
        self._create_file('.env.example', self._get_env_example())
        self._create_file('README.md', self._get_readme())
        self._create_file('Dockerfile', self._get_dockerfile())
        self._create_file('docker-compose.yml', self._get_docker_compose())
        self._create_file('.gitignore', self._get_gitignore())
        
        print("  ✅ FastAPI API generated!")
    
    def _generate_apollo_api(self):
        """Generate Apollo GraphQL API"""
        print("  📦 Creating Apollo GraphQL project structure...")
        
        # Create directory structure
        dirs = [
            'src/schema', 'src/resolvers', 'src/dataSources',
            'src/models', 'src/utils', 'src/directives',
            'src/middleware', 'tests', 'docs'
        ]
        
        for dir_path in dirs:
            (self.output_dir / dir_path).mkdir(parents=True, exist_ok=True)
        
        # Generate files
        self._create_file('package.json', self._get_apollo_package_json())
        self._create_file('src/schema/typeDefs.graphql', self.templates.get_graphql_schema())
        self._create_file('.env.example', self._get_env_example())
        self._create_file('README.md', self._get_readme())
        self._create_file('Dockerfile', self._get_dockerfile())
        self._create_file('docker-compose.yml', self._get_docker_compose())
        self._create_file('.gitignore', self._get_gitignore())
        
        print("  ✅ Apollo GraphQL API generated!")
    
    # Template methods
    def _get_express_package_json(self) -> str:
        return json.dumps({
            "name": self.name.lower().replace(' ', '-'),
            "version": "1.0.0",
            "description": f"{self.name} REST API",
            "main": "src/server.js",
            "scripts": {
                "start": "node src/server.js",
                "dev": "nodemon src/server.js",
                "test": "jest",
                "lint": "eslint src/"
            },
            "dependencies": {
                "express": "^4.18.2",
                "express-validator": "^7.0.1",
                "helmet": "^7.0.0",
                "cors": "^2.8.5",
                "compression": "^1.7.4",
                "bcryptjs": "^2.4.3",
                "jsonwebtoken": "^9.0.0",
                "dotenv": "^16.0.3",
                "sequelize": "^6.32.0",
                "pg": "^8.11.0",
                "redis": "^4.6.5",
                "winston": "^3.8.2",
                "swagger-ui-express": "^4.6.3"
            },
            "devDependencies": {
                "nodemon": "^2.0.22",
                "jest": "^29.5.0",
                "supertest": "^6.3.3",
                "eslint": "^8.41.0"
            }
        }, null, 2)
    
    def _get_apollo_package_json(self) -> str:
        return json.dumps({
            "name": self.name.lower().replace(' ', '-'),
            "version": "1.0.0",
            "description": f"{self.name} GraphQL API",
            "main": "src/index.js",
            "scripts": {
                "start": "node src/index.js",
                "dev": "nodemon src/index.js",
                "test": "jest"
            },
            "dependencies": {
                "@apollo/server": "^4.9.0",
                "graphql": "^16.7.1",
                "graphql-subscriptions": "^2.0.0",
                "jsonwebtoken": "^9.0.0",
                "bcryptjs": "^2.4.3",
                "dotenv": "^16.0.3",
                "mongoose": "^7.4.0",
                "dataloader": "^2.2.0"
            },
            "devDependencies": {
                "nodemon": "^2.0.22",
                "jest": "^29.5.0"
            }
        }, null, 2)
    
    def _get_fastapi_requirements(self) -> str:
        return """fastapi==0.104.1
uvicorn[standard]==0.24.0
python-jose[cryptography]==3.3.0
passlib[bcrypt]==1.7.4
python-multipart==0.0.6
email-validator==2.1.0
sqlalchemy==2.0.23
alembic==1.12.1
psycopg2-binary==2.9.9
redis==5.0.1
pydantic==2.5.0
python-dotenv==1.0.0
pytest==7.4.3
pytest-asyncio==0.21.1"""
    
    def _get_express_server(self) -> str:
        return """require('dotenv').config();
const app = require('./app');

const PORT = process.env.PORT || 3000;

app.listen(PORT, () => {
  console.log(`🚀 Server running on port ${PORT}`);
  console.log(`📖 API Docs: http://localhost:${PORT}/api-docs`);
});"""
    
    def _get_express_app(self) -> str:
        return """const express = require('express');
const helmet = require('helmet');
const cors = require('cors');
const compression = require('compression');

const app = express();

// Security middleware
app.use(helmet());
app.use(cors());

// Body parsing
app.use(express.json());
app.use(express.urlencoded({ extended: true }));

// Compression
app.use(compression());

// Health check
app.get('/health', (req, res) => {
  res.json({ status: 'healthy', timestamp: new Date().toISOString() });
});

// Routes would go here
app.get('/api/v1', (req, res) => {
  res.json({ message: 'API is running!' });
});

module.exports = app;"""
    
    def _get_fastapi_main(self) -> str:
        return """from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware

app = FastAPI(
    title="API",
    version="1.0.0",
    docs_url="/api-docs"
)

app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"]
)

@app.get("/health")
async def health_check():
    return {"status": "healthy"}

@app.get("/api/v1")
async def root():
    return {"message": "API is running!"}"""
    
    def _get_env_example(self) -> str:
        return """# Application
NODE_ENV=development
PORT=3000

# Database
DB_HOST=localhost
DB_PORT=5432
DB_NAME=api_db
DB_USER=postgres
DB_PASSWORD=postgres

# Authentication
JWT_SECRET=your-secret-key-change-this
JWT_EXPIRY=24h

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379"""
    
    def _get_readme(self) -> str:
        return f"""# {self.name}

A production-ready {self.api_type.upper()} API built with {self.framework}.

## Quick Start

1. Install dependencies:
   {'npm install' if self.framework in ['express', 'apollo', 'nestjs'] else 'pip install -r requirements.txt'}

2. Configure environment:
   cp .env.example .env

3. Run the server:
   {'npm run dev' if self.framework in ['express', 'apollo', 'nestjs'] else 'uvicorn app.main:app --reload'}

4. View API documentation:
   http://localhost:{'3000' if self.framework == 'express' else '8000'}/api-docs

## Docker

Build and run:
docker-compose up -d

## Testing

{'npm test' if self.framework in ['express', 'apollo', 'nestjs'] else 'pytest'}"""
    
    def _get_dockerfile(self) -> str:
        if self.framework in ['express', 'apollo', 'nestjs']:
            return """FROM node:18-alpine
WORKDIR /app
COPY package*.json ./
RUN npm ci --only=production
COPY . .
EXPOSE 3000
CMD ["node", "src/server.js"]"""
        else:
            return """FROM python:3.11-slim
WORKDIR /app
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt
COPY . .
EXPOSE 8000
CMD ["uvicorn", "app.main:app", "--host", "0.0.0.0", "--port", "8000"]"""
    
    def _get_docker_compose(self) -> str:
        return """version: '3.8'

services:
  api:
    build: .
    ports:
      - "3000:3000"
    environment:
      - NODE_ENV=production
      - DB_HOST=postgres
      - REDIS_HOST=redis
    depends_on:
      - postgres
      - redis

  postgres:
    image: postgres:15-alpine
    environment:
      - POSTGRES_DB=api_db
      - POSTGRES_USER=postgres
      - POSTGRES_PASSWORD=postgres
    volumes:
      - postgres_data:/var/lib/postgresql/data

  redis:
    image: redis:7-alpine
    volumes:
      - redis_data:/data

volumes:
  postgres_data:
  redis_data:"""
    
    def _get_gitignore(self) -> str:
        return """# Dependencies
node_modules/
venv/
__pycache__/
*.pyc

# Environment
.env
.env.local

# Logs
logs/
*.log

# Testing
coverage/
.coverage
.pytest_cache/

# IDE
.vscode/
.idea/
.DS_Store

# Build
dist/
build/"""
    
    # Stub methods for other frameworks
    def _generate_django_api(self):
        print("  📦 Django REST Framework support coming soon...")
        
    def _generate_flask_api(self):
        print("  📦 Flask support coming soon...")
        
    def _generate_nestjs_api(self):
        print("  📦 NestJS support coming soon...")
        
    def _generate_yoga_api(self):
        print("  📦 GraphQL Yoga support coming soon...")
        
    def _generate_strawberry_api(self):
        print("  📦 Strawberry GraphQL support coming soon...")
        
    def _generate_graphene_api(self):
        print("  📦 Graphene support coming soon...")

def main():
    parser = argparse.ArgumentParser(description="Generate production-ready APIs")
    parser.add_argument(
        "--type",
        choices=["rest", "graphql"],
        required=True,
        help="Type of API to generate"
    )
    parser.add_argument(
        "--framework",
        required=True,
        help="Framework to use (express, fastapi, django, apollo, etc.)"
    )
    parser.add_argument(
        "--name",
        required=True,
        help="Name of the API project"
    )
    parser.add_argument(
        "--output",
        help="Output directory (default: ./<name>)"
    )
    parser.add_argument(
        "--features",
        nargs="+",
        choices=["auth", "cache", "rate-limit", "websocket", "file-upload", "email", "payment"],
        help="Additional features to include"
    )
    
    args = parser.parse_args()
    
    options = {
        "output": args.output,
        "features": args.features or []
    }
    
    generator = APIGenerator(
        api_type=args.type,
        framework=args.framework,
        name=args.name,
        options=options
    )
    
    try:
        generator.generate()
    except Exception as e:
        print(f"❌ Error: {e}", file=sys.stderr)
        sys.exit(1)

if __name__ == "__main__":
    main()