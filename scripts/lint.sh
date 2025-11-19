#!/bin/bash

# Comprehensive linting script for Cooperative ERP Lite
# Usage: ./scripts/lint.sh [--fix] [--backend] [--frontend]

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Parse arguments
FIX_MODE=false
RUN_BACKEND=true
RUN_FRONTEND=true

while [[ $# -gt 0 ]]; do
  case $1 in
    --fix)
      FIX_MODE=true
      shift
      ;;
    --backend)
      RUN_FRONTEND=false
      shift
      ;;
    --frontend)
      RUN_BACKEND=false
      shift
      ;;
    *)
      echo -e "${RED}Unknown argument: $1${NC}"
      exit 1
      ;;
  esac
done

echo -e "${BLUE}╔════════════════════════════════════════════════════╗${NC}"
echo -e "${BLUE}║   Cooperative ERP Lite - Code Quality Check       ║${NC}"
echo -e "${BLUE}╚════════════════════════════════════════════════════╝${NC}"
echo ""

ERRORS=0

# Backend linting
if [ "$RUN_BACKEND" = true ]; then
  echo -e "${YELLOW}→ Running Go linting...${NC}"

  if [ ! -d "backend" ]; then
    echo -e "${RED}✗ Backend directory not found${NC}"
    ERRORS=$((ERRORS + 1))
  else
    cd backend

    # Check if golangci-lint is installed
    if ! command -v golangci-lint &> /dev/null; then
      echo -e "${RED}✗ golangci-lint not installed${NC}"
      echo -e "${YELLOW}  Install with: brew install golangci-lint${NC}"
      ERRORS=$((ERRORS + 1))
    else
      if [ "$FIX_MODE" = true ]; then
        echo -e "${BLUE}  Running golangci-lint with auto-fix...${NC}"
        if golangci-lint run --fix; then
          echo -e "${GREEN}✓ Go linting passed (with fixes)${NC}"
        else
          echo -e "${RED}✗ Go linting failed${NC}"
          ERRORS=$((ERRORS + 1))
        fi
      else
        echo -e "${BLUE}  Running golangci-lint...${NC}"
        if golangci-lint run; then
          echo -e "${GREEN}✓ Go linting passed${NC}"
        else
          echo -e "${RED}✗ Go linting failed${NC}"
          ERRORS=$((ERRORS + 1))
        fi
      fi
    fi

    # Go formatting check
    echo -e "${YELLOW}→ Checking Go formatting...${NC}"
    UNFORMATTED=$(gofmt -l . 2>&1 | grep -v "vendor/" || true)
    if [ -n "$UNFORMATTED" ]; then
      if [ "$FIX_MODE" = true ]; then
        echo -e "${BLUE}  Formatting Go files...${NC}"
        gofmt -w .
        echo -e "${GREEN}✓ Go files formatted${NC}"
      else
        echo -e "${RED}✗ Go files need formatting:${NC}"
        echo "$UNFORMATTED"
        ERRORS=$((ERRORS + 1))
      fi
    else
      echo -e "${GREEN}✓ Go formatting passed${NC}"
    fi

    # Go vet
    echo -e "${YELLOW}→ Running go vet...${NC}"
    if go vet ./...; then
      echo -e "${GREEN}✓ Go vet passed${NC}"
    else
      echo -e "${RED}✗ Go vet failed${NC}"
      ERRORS=$((ERRORS + 1))
    fi

    # Go mod tidy check
    echo -e "${YELLOW}→ Checking go.mod...${NC}"
    if [ "$FIX_MODE" = true ]; then
      go mod tidy
      echo -e "${GREEN}✓ go.mod tidied${NC}"
    else
      cp go.mod go.mod.backup
      cp go.sum go.sum.backup
      go mod tidy
      if diff go.mod go.mod.backup && diff go.sum go.sum.backup; then
        echo -e "${GREEN}✓ go.mod is tidy${NC}"
      else
        echo -e "${RED}✗ go.mod needs tidying${NC}"
        ERRORS=$((ERRORS + 1))
      fi
      mv go.mod.backup go.mod
      mv go.sum.backup go.sum
    fi

    cd ..
  fi
  echo ""
fi

# Frontend linting
if [ "$RUN_FRONTEND" = true ]; then
  echo -e "${YELLOW}→ Running Frontend linting...${NC}"

  if [ ! -d "frontend" ]; then
    echo -e "${RED}✗ Frontend directory not found${NC}"
    ERRORS=$((ERRORS + 1))
  else
    cd frontend

    # Check if node_modules exists
    if [ ! -d "node_modules" ]; then
      echo -e "${RED}✗ node_modules not found. Run 'npm install' first.${NC}"
      ERRORS=$((ERRORS + 1))
    else
      # ESLint
      echo -e "${YELLOW}→ Running ESLint...${NC}"
      if [ "$FIX_MODE" = true ]; then
        if npm run lint:fix; then
          echo -e "${GREEN}✓ ESLint passed (with fixes)${NC}"
        else
          echo -e "${RED}✗ ESLint failed${NC}"
          ERRORS=$((ERRORS + 1))
        fi
      else
        if npm run lint; then
          echo -e "${GREEN}✓ ESLint passed${NC}"
        else
          echo -e "${RED}✗ ESLint failed${NC}"
          ERRORS=$((ERRORS + 1))
        fi
      fi

      # Prettier
      echo -e "${YELLOW}→ Running Prettier...${NC}"
      if [ "$FIX_MODE" = true ]; then
        npm run format
        echo -e "${GREEN}✓ Code formatted${NC}"
      else
        if npm run format:check; then
          echo -e "${GREEN}✓ Prettier check passed${NC}"
        else
          echo -e "${RED}✗ Code needs formatting${NC}"
          ERRORS=$((ERRORS + 1))
        fi
      fi

      # TypeScript type checking
      echo -e "${YELLOW}→ Running TypeScript type check...${NC}"
      if npm run type-check; then
        echo -e "${GREEN}✓ Type check passed${NC}"
      else
        echo -e "${RED}✗ Type check failed${NC}"
        ERRORS=$((ERRORS + 1))
      fi
    fi

    cd ..
  fi
  echo ""
fi

# Summary
echo -e "${BLUE}╔════════════════════════════════════════════════════╗${NC}"
if [ $ERRORS -eq 0 ]; then
  echo -e "${GREEN}║   ✓ All quality checks passed!                    ║${NC}"
  echo -e "${BLUE}╚════════════════════════════════════════════════════╝${NC}"
  exit 0
else
  echo -e "${RED}║   ✗ Quality checks failed with $ERRORS error(s)          ║${NC}"
  echo -e "${BLUE}╚════════════════════════════════════════════════════╝${NC}"
  if [ "$FIX_MODE" = false ]; then
    echo -e "${YELLOW}Tip: Run with --fix to automatically fix some issues${NC}"
  fi
  exit 1
fi
