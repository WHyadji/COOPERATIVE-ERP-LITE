# ============================================================================
# Multi-Stage Dockerfile for Next.js 15 Frontend
# Production-ready with security hardening and optimal image size
#
# Build stages:
#   1. base    - Base Node.js Alpine image
#   2. deps    - Install dependencies only
#   3. builder - Build Next.js application
#   4. runner  - Minimal production runtime
#
# Expected image size: ~229MB (75% reduction from non-standalone)
# ============================================================================

# ============================================================================
# Stage 1: Base - Foundation for all stages
# ============================================================================
FROM node:20-alpine AS base

# Install libc6-compat for Alpine Linux compatibility with Next.js
RUN apk add --no-cache libc6-compat

WORKDIR /app

# ============================================================================
# Stage 2: Dependencies - Install packages with optimal caching
# ============================================================================
FROM base AS deps

# Copy package files for dependency installation
# These files change less frequently, so Docker can cache this layer
COPY package.json package-lock.json* ./

# Install dependencies using npm ci for reproducible builds
# --omit=dev skips devDependencies (not needed in production)
# We install all deps here because builder needs devDependencies for build
RUN npm ci

# ============================================================================
# Stage 3: Builder - Build Next.js application
# ============================================================================
FROM base AS builder

WORKDIR /app

# Copy dependencies from deps stage
COPY --from=deps /app/node_modules ./node_modules

# Copy all source code
COPY . .

# Set build-time environment variables
# NEXT_PUBLIC_* variables must be set at build time (baked into JS bundles)
ARG NEXT_PUBLIC_API_BASE_URL=http://localhost:8080/api/v1
ENV NEXT_PUBLIC_API_BASE_URL=$NEXT_PUBLIC_API_BASE_URL

# Disable Next.js telemetry during build
ENV NEXT_TELEMETRY_DISABLED=1

# Build Next.js application with Turbopack
# This creates .next/standalone directory with minimal dependencies
RUN npm run build

# ============================================================================
# Stage 4: Runner - Minimal production runtime
# ============================================================================
FROM base AS runner

WORKDIR /app

# Set production environment
ENV NODE_ENV=production
ENV NEXT_TELEMETRY_DISABLED=1

# Create non-root user for security
# Using specific UID/GID for consistency across environments
RUN addgroup --system --gid 1001 nodejs && \
    adduser --system --uid 1001 nextjs

# Copy public assets (static files served by Next.js)
COPY --from=builder /app/public ./public

# Create .next directory with proper permissions
RUN mkdir .next && chown nextjs:nodejs .next

# Copy standalone output (minimal server + dependencies)
# This is the magic of Next.js standalone mode - only necessary files
COPY --from=builder --chown=nextjs:nodejs /app/.next/standalone ./

# Copy static files (built CSS, JS, images)
COPY --from=builder --chown=nextjs:nodejs /app/.next/static ./.next/static

# Switch to non-root user
USER nextjs

# Expose port (can be overridden by $PORT env var for Cloud Run)
EXPOSE 3000

# Set hostname to accept connections from any interface
# Required for Docker networking
ENV HOSTNAME=0.0.0.0

# Health check (optional but recommended for production)
# Verifies that the server is responding
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD node -e "require('http').get('http://localhost:3000/api/health', (r) => {process.exit(r.statusCode === 200 ? 0 : 1)})" || exit 1

# Start the Next.js standalone server
# The standalone output includes a minimal Node.js server
CMD ["node", "server.js"]
