#!/usr/bin/env python3
"""
Docker Orchestrator
Generate optimized Dockerfiles, docker-compose configurations, and container strategies
"""

import os
import sys
import json
import yaml
import argparse
from pathlib import Path
from typing import Dict, List, Optional, Tuple
from dataclasses import dataclass
from enum import Enum

class AppType(Enum):
    NODE = "node"
    PYTHON = "python"
    GO = "go"
    JAVA = "java"
    DOTNET = "dotnet"
    RUBY = "ruby"
    PHP = "php"
    STATIC = "static"
    
class Framework(Enum):
    EXPRESS = "express"
    FASTIFY = "fastify"
    NESTJS = "nestjs"
    DJANGO = "django"
    FLASK = "flask"
    FASTAPI = "fastapi"
    SPRING = "spring"
    RAILS = "rails"
    LARAVEL = "laravel"
    REACT = "react"
    VUE = "vue"
    ANGULAR = "angular"
    NEXTJS = "nextjs"

@dataclass
class DockerConfig:
    app_type: AppType
    framework: Optional[Framework]
    port: int
    multi_stage: bool
    optimize: bool
    healthcheck: bool
    non_root: bool
    
class DockerOrchestrator:
    def __init__(self):
        self.templates = self._load_templates()
        
    def _load_templates(self) -> Dict:
        """Load Docker templates"""
        return {
            "node": self._node_template,
            "python": self._python_template,
            "go": self._go_template,
            "java": self._java_template,
            "static": self._static_template
        }
    
    def generate_dockerfile(self, config: DockerConfig) -> str:
        """Generate optimized Dockerfile based on configuration"""
        
        template_func = self.templates.get(config.app_type.value)
        if not template_func:
            raise ValueError(f"Unsupported app type: {config.app_type.value}")
        
        return template_func(config)
    
    def _node_template(self, config: DockerConfig) -> str:
        """Generate Node.js Dockerfile"""
        
        if config.multi_stage and config.optimize:
            dockerfile = '''# syntax=docker/dockerfile:1
# Stage 1: Dependencies
FROM node:18-alpine AS deps
WORKDIR /app

# Copy package files
COPY package*.json ./
RUN npm ci --only=production

# Stage 2: Build
FROM node:18-alpine AS builder
WORKDIR /app

# Copy package files
COPY package*.json ./
RUN npm ci

# Copy source code
COPY . .

# Build application
RUN npm run build

# Stage 3: Production
FROM node:18-alpine AS runner
WORKDIR /app

ENV NODE_ENV=production

# Create non-root user
RUN addgroup -g 1001 -S nodejs && \\
    adduser -S nodejs -u 1001

# Copy production dependencies
COPY --from=deps --chown=nodejs:nodejs /app/node_modules ./node_modules

# Copy built application
COPY --from=builder --chown=nodejs:nodejs /app/dist ./dist
COPY --from=builder --chown=nodejs:nodejs /app/package*.json ./

USER nodejs

EXPOSE ''' + str(config.port) + '''

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \\
    CMD node healthcheck.js || exit 1

CMD ["node", "dist/index.js"]'''
        
        else:
            # Simple Dockerfile
            dockerfile = '''FROM node:18-alpine

WORKDIR /app

# Copy and install dependencies
COPY package*.json ./
RUN npm ci --only=production

# Copy application code
COPY . .

EXPOSE ''' + str(config.port) + '''

CMD ["node", "index.js"]'''
        
        return dockerfile
    
    def _python_template(self, config: DockerConfig) -> str:
        """Generate Python Dockerfile"""
        
        if config.framework == Framework.DJANGO:
            cmd = 'CMD ["gunicorn", "--bind", "0.0.0.0:8000", "--workers", "4", "myproject.wsgi:application"]'
        elif config.framework == Framework.FASTAPI:
            cmd = 'CMD ["uvicorn", "main:app", "--host", "0.0.0.0", "--port", "8000"]'
        else:
            cmd = 'CMD ["python", "app.py"]'
        
        if config.multi_stage and config.optimize:
            dockerfile = '''# syntax=docker/dockerfile:1
# Stage 1: Build dependencies
FROM python:3.11-slim AS builder

WORKDIR /app

# Install build dependencies
RUN apt-get update && apt-get install -y \\
    gcc \\
    g++ \\
    && rm -rf /var/lib/apt/lists/*

# Install Python dependencies
COPY requirements.txt .
RUN pip install --user --no-cache-dir -r requirements.txt

# Stage 2: Production
FROM python:3.11-slim AS runner

WORKDIR /app

# Create non-root user
RUN groupadd -r appuser && useradd -r -g appuser appuser

# Copy dependencies from builder
COPY --from=builder --chown=appuser:appuser /root/.local /home/appuser/.local

# Copy application code
COPY --chown=appuser:appuser . .

# Update PATH
ENV PATH=/home/appuser/.local/bin:$PATH

USER appuser

EXPOSE ''' + str(config.port) + '''

''' + cmd
        
        else:
            # Simple Dockerfile
            dockerfile = '''FROM python:3.11-slim

WORKDIR /app

# Install dependencies
COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

# Copy application code
COPY . .

EXPOSE ''' + str(config.port) + '''

''' + cmd
        
        return dockerfile
    
    def _go_template(self, config: DockerConfig) -> str:
        """Generate Go Dockerfile"""
        
        dockerfile = '''# syntax=docker/dockerfile:1
# Stage 1: Build
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Stage 2: Production
FROM gcr.io/distroless/static:nonroot

WORKDIR /

# Copy binary from builder
COPY --from=builder /app/main .

USER nonroot:nonroot

EXPOSE ''' + str(config.port) + '''

ENTRYPOINT ["/main"]'''
        
        return dockerfile
    
    def _java_template(self, config: DockerConfig) -> str:
        """Generate Java Dockerfile"""
        
        if config.framework == Framework.SPRING:
            dockerfile = '''# syntax=docker/dockerfile:1
# Stage 1: Build
FROM maven:3.9-eclipse-temurin-17 AS builder

WORKDIR /app

# Copy POM
COPY pom.xml .
RUN mvn dependency:go-offline

# Copy source and build
COPY src ./src
RUN mvn clean package -DskipTests

# Stage 2: Extract layers
FROM eclipse-temurin:17-jre AS extractor

WORKDIR /app

COPY --from=builder /app/target/*.jar app.jar
RUN java -Djarmode=layertools -jar app.jar extract

# Stage 3: Production
FROM eclipse-temurin:17-jre-alpine

WORKDIR /app

# Create non-root user
RUN addgroup -S spring && adduser -S spring -G spring

# Copy layers
COPY --from=extractor --chown=spring:spring /app/dependencies/ ./
COPY --from=extractor --chown=spring:spring /app/spring-boot-loader/ ./
COPY --from=extractor --chown=spring:spring /app/snapshot-dependencies/ ./
COPY --from=extractor --chown=spring:spring /app/application/ ./

USER spring

EXPOSE ''' + str(config.port) + '''

ENTRYPOINT ["java", "org.springframework.boot.loader.JarLauncher"]'''
        else:
            dockerfile = '''FROM eclipse-temurin:17-jre-alpine

WORKDIR /app

COPY target/*.jar app.jar

EXPOSE ''' + str(config.port) + '''

ENTRYPOINT ["java", "-jar", "app.jar"]'''
        
        return dockerfile
    
    def _static_template(self, config: DockerConfig) -> str:
        """Generate static site Dockerfile with nginx"""
        
        dockerfile = '''# Stage 1: Build
FROM node:18-alpine AS builder

WORKDIR /app

# Install dependencies
COPY package*.json ./
RUN npm ci

# Copy source and build
COPY . .
RUN npm run build

# Stage 2: Serve with nginx
FROM nginx:alpine

# Copy custom nginx config
COPY nginx.conf /etc/nginx/nginx.conf

# Copy built assets from builder
COPY --from=builder /app/dist /usr/share/nginx/html

# Add non-root user
RUN chown -R nginx:nginx /usr/share/nginx/html && \\
    chmod -R 755 /usr/share/nginx/html && \\
    chown -R nginx:nginx /var/cache/nginx && \\
    chown -R nginx:nginx /var/log/nginx && \\
    chown -R nginx:nginx /etc/nginx/conf.d && \\
    touch /var/run/nginx.pid && \\
    chown -R nginx:nginx /var/run/nginx.pid

USER nginx

EXPOSE 80

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \\
    CMD curl -f http://localhost/ || exit 1

CMD ["nginx", "-g", "daemon off;"]'''
        
        return dockerfile
    
    def generate_docker_compose(self, services: List[str], scale: int = 1) -> str:
        """Generate docker-compose.yml configuration"""
        
        compose = {
            "version": "3.9",
            "services": {},
            "networks": {
                "app-network": {
                    "driver": "bridge"
                }
            },
            "volumes": {}
        }
        
        # Add services based on requested stack
        if "web" in services or "frontend" in services:
            compose["services"]["frontend"] = {
                "build": {
                    "context": "./frontend",
                    "dockerfile": "Dockerfile"
                },
                "ports": ["80:80"],
                "networks": ["app-network"],
                "depends_on": ["backend"] if "api" in services or "backend" in services else [],
                "restart": "unless-stopped"
            }
        
        if "api" in services or "backend" in services:
            compose["services"]["backend"] = {
                "build": {
                    "context": "./backend",
                    "dockerfile": "Dockerfile"
                },
                "ports": ["3000:3000"],
                "environment": [
                    "NODE_ENV=production",
                    "DATABASE_URL=postgres://user:pass@postgres:5432/appdb"
                ],
                "networks": ["app-network"],
                "depends_on": ["postgres"] if "db" in services or "postgres" in services else [],
                "restart": "unless-stopped"
            }
            
            if scale > 1:
                compose["services"]["backend"]["deploy"] = {
                    "replicas": scale
                }
        
        if "db" in services or "postgres" in services:
            compose["services"]["postgres"] = {
                "image": "postgres:15-alpine",
                "environment": {
                    "POSTGRES_DB": "appdb",
                    "POSTGRES_USER": "user",
                    "POSTGRES_PASSWORD": "${DB_PASSWORD}"
                },
                "volumes": ["postgres-data:/var/lib/postgresql/data"],
                "ports": ["5432:5432"],
                "networks": ["app-network"],
                "restart": "unless-stopped",
                "healthcheck": {
                    "test": ["CMD-SHELL", "pg_isready -U user -d appdb"],
                    "interval": "10s",
                    "timeout": "5s",
                    "retries": 5
                }
            }
            compose["volumes"]["postgres-data"] = {"driver": "local"}
        
        if "redis" in services:
            compose["services"]["redis"] = {
                "image": "redis:7-alpine",
                "command": "redis-server --appendonly yes",
                "volumes": ["redis-data:/data"],
                "ports": ["6379:6379"],
                "networks": ["app-network"],
                "restart": "unless-stopped"
            }
            compose["volumes"]["redis-data"] = {"driver": "local"}
        
        if "mongodb" in services or "mongo" in services:
            compose["services"]["mongodb"] = {
                "image": "mongo:6",
                "environment": {
                    "MONGO_INITDB_ROOT_USERNAME": "admin",
                    "MONGO_INITDB_ROOT_PASSWORD": "${MONGO_PASSWORD}"
                },
                "volumes": ["mongo-data:/data/db"],
                "ports": ["27017:27017"],
                "networks": ["app-network"],
                "restart": "unless-stopped"
            }
            compose["volumes"]["mongo-data"] = {"driver": "local"}
        
        if "rabbitmq" in services:
            compose["services"]["rabbitmq"] = {
                "image": "rabbitmq:3-management-alpine",
                "ports": ["5672:5672", "15672:15672"],
                "volumes": ["rabbitmq-data:/var/lib/rabbitmq"],
                "networks": ["app-network"],
                "restart": "unless-stopped"
            }
            compose["volumes"]["rabbitmq-data"] = {"driver": "local"}
        
        if "nginx" in services:
            compose["services"]["nginx"] = {
                "image": "nginx:alpine",
                "ports": ["443:443", "80:80"],
                "volumes": [
                    "./nginx/nginx.conf:/etc/nginx/nginx.conf:ro",
                    "./nginx/conf.d:/etc/nginx/conf.d:ro"
                ],
                "networks": ["app-network"],
                "depends_on": ["backend"] if "backend" in compose["services"] else [],
                "restart": "unless-stopped"
            }
        
        return yaml.dump(compose, default_flow_style=False, sort_keys=False)
    
    def generate_dockerignore(self, app_type: AppType) -> str:
        """Generate .dockerignore file"""
        
        common = [
            ".git",
            ".gitignore",
            "README.md",
            ".env",
            ".env.*",
            ".vscode",
            ".idea",
            "*.swp",
            ".DS_Store",
            "Thumbs.db"
        ]
        
        if app_type == AppType.NODE:
            specific = [
                "node_modules",
                "npm-debug.log",
                "yarn-debug.log",
                "yarn-error.log",
                ".npm",
                ".yarn",
                "coverage",
                ".nyc_output",
                "dist",
                "build"
            ]
        elif app_type == AppType.PYTHON:
            specific = [
                "__pycache__",
                "*.py[cod]",
                "*$py.class",
                "*.so",
                ".Python",
                "env",
                "venv",
                "ENV",
                ".venv",
                "pip-log.txt",
                "pip-delete-this-directory.txt",
                ".tox",
                ".coverage",
                ".pytest_cache",
                "htmlcov",
                ".mypy_cache",
                ".pytype",
                "*.egg-info"
            ]
        elif app_type == AppType.GO:
            specific = [
                "*.exe",
                "*.exe~",
                "*.dll",
                "*.so",
                "*.dylib",
                "*.test",
                "*.out",
                "vendor"
            ]
        elif app_type == AppType.JAVA:
            specific = [
                "target",
                "*.class",
                "*.jar",
                "*.war",
                "*.ear",
                ".gradle",
                "build",
                ".mvn",
                "mvnw",
                "mvnw.cmd"
            ]
        else:
            specific = []
        
        return "\n".join(common + specific)
    
    def generate_nginx_config(self) -> str:
        """Generate nginx configuration"""
        
        return '''events {
    worker_connections 1024;
}

http {
    include /etc/nginx/mime.types;
    default_type application/octet-stream;

    # Logging
    access_log /var/log/nginx/access.log;
    error_log /var/log/nginx/error.log;

    # Gzip compression
    gzip on;
    gzip_vary on;
    gzip_min_length 1024;
    gzip_types text/plain text/css text/xml text/javascript application/javascript application/xml+rss application/json;

    # Rate limiting
    limit_req_zone $binary_remote_addr zone=api:10m rate=10r/s;

    # Upstream servers
    upstream backend {
        least_conn;
        server backend:3000 max_fails=3 fail_timeout=30s;
    }

    server {
        listen 80;
        server_name _;

        # Security headers
        add_header X-Frame-Options "SAMEORIGIN" always;
        add_header X-Content-Type-Options "nosniff" always;
        add_header X-XSS-Protection "1; mode=block" always;

        # Frontend
        location / {
            root /usr/share/nginx/html;
            index index.html;
            try_files $uri $uri/ /index.html;
        }

        # API proxy
        location /api {
            limit_req zone=api burst=20 nodelay;
            
            proxy_pass http://backend;
            proxy_http_version 1.1;
            proxy_set_header Upgrade $http_upgrade;
            proxy_set_header Connection 'upgrade';
            proxy_set_header Host $host;
            proxy_cache_bypass $http_upgrade;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
            proxy_set_header X-Forwarded-Proto $scheme;
        }

        # Health check
        location /health {
            access_log off;
            return 200 "healthy\\n";
            add_header Content-Type text/plain;
        }
    }
}'''
    
    def analyze_project(self, project_path: str) -> DockerConfig:
        """Analyze project and suggest Docker configuration"""
        
        path = Path(project_path)
        
        # Detect app type based on files
        if (path / "package.json").exists():
            app_type = AppType.NODE
            
            # Detect framework
            with open(path / "package.json", 'r') as f:
                package_json = json.load(f)
                deps = package_json.get("dependencies", {})
                
                if "express" in deps:
                    framework = Framework.EXPRESS
                elif "fastify" in deps:
                    framework = Framework.FASTIFY
                elif "next" in deps:
                    framework = Framework.NEXTJS
                elif "react" in deps:
                    framework = Framework.REACT
                elif "vue" in deps:
                    framework = Framework.VUE
                else:
                    framework = None
                    
            port = 3000
            
        elif (path / "requirements.txt").exists() or (path / "setup.py").exists():
            app_type = AppType.PYTHON
            
            # Detect framework
            if (path / "requirements.txt").exists():
                with open(path / "requirements.txt", 'r') as f:
                    requirements = f.read()
                    
                if "django" in requirements.lower():
                    framework = Framework.DJANGO
                elif "flask" in requirements.lower():
                    framework = Framework.FLASK
                elif "fastapi" in requirements.lower():
                    framework = Framework.FASTAPI
                else:
                    framework = None
            else:
                framework = None
                
            port = 8000
            
        elif (path / "go.mod").exists():
            app_type = AppType.GO
            framework = None
            port = 8080
            
        elif (path / "pom.xml").exists() or (path / "build.gradle").exists():
            app_type = AppType.JAVA
            framework = Framework.SPRING if (path / "pom.xml").exists() else None
            port = 8080
            
        elif (path / "Gemfile").exists():
            app_type = AppType.RUBY
            framework = Framework.RAILS
            port = 3000
            
        elif (path / "composer.json").exists():
            app_type = AppType.PHP
            framework = Framework.LARAVEL
            port = 8000
            
        else:
            # Default to static site
            app_type = AppType.STATIC
            framework = None
            port = 80
        
        return DockerConfig(
            app_type=app_type,
            framework=framework,
            port=port,
            multi_stage=True,
            optimize=True,
            healthcheck=True,
            non_root=True
        )
    
    def optimize_dockerfile(self, dockerfile_path: str) -> str:
        """Analyze and optimize existing Dockerfile"""
        
        with open(dockerfile_path, 'r') as f:
            content = f.read()
        
        suggestions = []
        
        # Check for anti-patterns
        if "FROM.*:latest" in content:
            suggestions.append("❌ Use specific version tags instead of :latest")
        
        if content.count("RUN ") > 5:
            suggestions.append("❌ Combine RUN commands to reduce layers")
        
        if "COPY . ." in content and "COPY package" not in content:
            suggestions.append("❌ Copy dependency files first to leverage cache")
        
        if "USER " not in content:
            suggestions.append("❌ Run as non-root user for security")
        
        if "HEALTHCHECK" not in content:
            suggestions.append("❌ Add HEALTHCHECK for container monitoring")
        
        if "apt-get install" in content and "rm -rf /var/lib/apt/lists" not in content:
            suggestions.append("❌ Clean apt cache to reduce image size")
        
        # Positive patterns
        if "# syntax=docker/dockerfile:1" in content:
            suggestions.append("✅ Using BuildKit features")
        
        if "--mount=type=cache" in content:
            suggestions.append("✅ Using cache mounts for dependencies")
        
        if "FROM.*AS.*builder" in content:
            suggestions.append("✅ Using multi-stage build")
        
        return "\n".join(suggestions)

def main():
    parser = argparse.ArgumentParser(description="Docker Orchestrator")
    
    # Commands
    subparsers = parser.add_subparsers(dest="command", help="Command to run")
    
    # Generate Dockerfile
    dockerfile_parser = subparsers.add_parser("dockerfile", help="Generate Dockerfile")
    dockerfile_parser.add_argument("--app", choices=["node", "python", "go", "java", "static"],
                                  help="Application type")
    dockerfile_parser.add_argument("--framework", help="Framework (express, django, spring, etc.)")
    dockerfile_parser.add_argument("--port", type=int, default=3000, help="Application port")
    dockerfile_parser.add_argument("--optimize", action="store_true", help="Apply optimizations")
    dockerfile_parser.add_argument("--output", default="Dockerfile", help="Output file")
    
    # Generate docker-compose
    compose_parser = subparsers.add_parser("compose", help="Generate docker-compose.yml")
    compose_parser.add_argument("--services", required=True, help="Comma-separated services")
    compose_parser.add_argument("--scale", type=int, default=1, help="Service replicas")
    compose_parser.add_argument("--output", default="docker-compose.yml", help="Output file")
    
    # Analyze project
    analyze_parser = subparsers.add_parser("analyze", help="Analyze project")
    analyze_parser.add_argument("path", help="Project path")
    analyze_parser.add_argument("--generate", action="store_true", help="Generate files")
    
    # Optimize existing
    optimize_parser = subparsers.add_parser("optimize", help="Optimize Dockerfile")
    optimize_parser.add_argument("dockerfile", help="Dockerfile path")
    
    args = parser.parse_args()
    
    if not args.command:
        parser.print_help()
        sys.exit(1)
    
    orchestrator = DockerOrchestrator()
    
    if args.command == "dockerfile":
        # Generate Dockerfile
        app_type = AppType(args.app) if args.app else AppType.NODE
        framework = Framework(args.framework) if args.framework else None
        
        config = DockerConfig(
            app_type=app_type,
            framework=framework,
            port=args.port,
            multi_stage=args.optimize,
            optimize=args.optimize,
            healthcheck=True,
            non_root=True
        )
        
        dockerfile = orchestrator.generate_dockerfile(config)
        
        with open(args.output, 'w') as f:
            f.write(dockerfile)
        
        print(f"✅ Generated {args.output}")
        
        # Also generate .dockerignore
        dockerignore = orchestrator.generate_dockerignore(app_type)
        with open(".dockerignore", 'w') as f:
            f.write(dockerignore)
        print("✅ Generated .dockerignore")
        
    elif args.command == "compose":
        # Generate docker-compose.yml
        services = args.services.split(",")
        compose = orchestrator.generate_docker_compose(services, args.scale)
        
        with open(args.output, 'w') as f:
            f.write(compose)
        
        print(f"✅ Generated {args.output}")
        
        # Generate nginx config if needed
        if "nginx" in services:
            nginx_config = orchestrator.generate_nginx_config()
            os.makedirs("nginx", exist_ok=True)
            with open("nginx/nginx.conf", 'w') as f:
                f.write(nginx_config)
            print("✅ Generated nginx/nginx.conf")
        
    elif args.command == "analyze":
        # Analyze project
        config = orchestrator.analyze_project(args.path)
        
        print(f"📋 Project Analysis:")
        print(f"  App Type: {config.app_type.value}")
        print(f"  Framework: {config.framework.value if config.framework else 'None'}")
        print(f"  Suggested Port: {config.port}")
        
        if args.generate:
            # Generate Dockerfile
            dockerfile = orchestrator.generate_dockerfile(config)
            output_path = Path(args.path) / "Dockerfile"
            with open(output_path, 'w') as f:
                f.write(dockerfile)
            print(f"✅ Generated {output_path}")
            
            # Generate .dockerignore
            dockerignore = orchestrator.generate_dockerignore(config.app_type)
            ignore_path = Path(args.path) / ".dockerignore"
            with open(ignore_path, 'w') as f:
                f.write(dockerignore)
            print(f"✅ Generated {ignore_path}")
        
    elif args.command == "optimize":
        # Analyze and optimize Dockerfile
        suggestions = orchestrator.optimize_dockerfile(args.dockerfile)
        
        print(f"📋 Dockerfile Analysis:")
        print(suggestions)

if __name__ == "__main__":
    main()