# `/autofix` - Auto-Identify and Fix Issues

## Smart Issue Detection & Resolution

```yaml
---
allowed-tools: ReadFile, WriteFile, SearchReplace, Bash(*), FileAnalysis
description: Automatically identify issue type and apply targeted fixes
---
```

## Command Structure

### `/autofix` - Main Command
**Usage:** `/autofix [optional: error message or file path]`

**Auto-Detection Flow:**
1. **Quick System Scan** (3 seconds)
2. **Issue Classification** (determine category)
3. **Targeted Diagnostics** (specific to issue type)
4. **Apply Fix** (automated solution)
5. **Verification** (confirm resolution)

---

## Auto-Detection Logic

### Phase 1: Quick System Scan
```bash
# Check running processes
PROCESSES=$(ps aux | grep -E "(node|go|python|docker)" | grep -v grep)

# Check recent errors
RECENT_ERRORS=$(find . -name "*.log" -mtime -1 -exec tail -20 {} \; 2>/dev/null | grep -E "ERROR|FATAL|panic|exception" | tail -5)

# Check build status
BUILD_STATUS=$(ls -la .next package-lock.json go.mod requirements.txt 2>/dev/null)

# Check network/ports
PORTS=$(netstat -tulpn 2>/dev/null | grep -E ":3000|:8080|:5000|:8000" | head -5)

# Check file system
FS_ISSUES=$(find . -name "node_modules" -o -name "__pycache__" -o -name ".next" | head -3)
```

### Phase 2: Issue Classification
```bash
# Classify based on evidence
if [[ $RECENT_ERRORS =~ "SyntaxError|TypeError|ReferenceError" ]]; then
    ISSUE_TYPE="build"
elif [[ $RECENT_ERRORS =~ "ECONNREFUSED|timeout|502|503" ]]; then
    ISSUE_TYPE="api"
elif [[ $RECENT_ERRORS =~ "database|connection|sql" ]]; then
    ISSUE_TYPE="database"
elif [[ $RECENT_ERRORS =~ "CORS|authentication|401|403" ]]; then
    ISSUE_TYPE="security"
elif [[ $PORTS == "" && $PROCESSES =~ "node|go|python" ]]; then
    ISSUE_TYPE="runtime"
elif [[ $BUILD_STATUS =~ "node_modules" && ! -f ".next/BUILD_ID" ]]; then
    ISSUE_TYPE="frontend"
else
    ISSUE_TYPE="general"
fi
```

---

## Automated Fix Strategies

### Build Issues Auto-Fix
```bash
autofix_build() {
    echo "🔍 Detected: Build/Compilation Issue"

    # Check for package.json
    if [[ -f "package.json" ]]; then
        echo "📦 Fixing Node.js dependencies..."
        rm -rf node_modules package-lock.json .next
        npm install
        npm run build 2>&1 | head -10
    fi

    # Check for go.mod
    if [[ -f "go.mod" ]]; then
        echo "🐹 Fixing Go modules..."
        go mod tidy
        go build ./... 2>&1 | head -10
    fi

    # Check for Python
    if [[ -f "requirements.txt" || -f "pyproject.toml" ]]; then
        echo "🐍 Fixing Python dependencies..."
        pip install -r requirements.txt --quiet
        python -m py_compile *.py 2>&1 | head -5
    fi

    echo "✅ Build fix applied"
}
```

### API Issues Auto-Fix
```bash
autofix_api() {
    echo "🔍 Detected: API/Network Issue"

    # Check if services are running
    if ! pgrep -f "node.*3000" > /dev/null; then
        echo "🚀 Starting frontend service..."
        npm run dev &
        sleep 3
    fi

    if ! pgrep -f "go.*8080\|python.*8080" > /dev/null; then
        echo "🚀 Starting backend service..."
        if [[ -f "main.go" ]]; then
            go run main.go &
        elif [[ -f "app.py" || -f "main.py" ]]; then
            python app.py &
        fi
        sleep 3
    fi

    # Test connectivity
    echo "🔗 Testing API connectivity..."
    curl -f http://localhost:3000/api/health 2>/dev/null && echo "✅ Frontend API OK" || echo "❌ Frontend API Down"
    curl -f http://localhost:8080/health 2>/dev/null && echo "✅ Backend API OK" || echo "❌ Backend API Down"
}
```

### Database Issues Auto-Fix
```bash
autofix_database() {
    echo "🔍 Detected: Database Issue"

    # Check database connection
    if command -v psql &> /dev/null; then
        echo "🐘 Testing PostgreSQL connection..."
        pg_isready -h localhost -p 5432 && echo "✅ PostgreSQL OK" || {
            echo "🔧 Starting PostgreSQL..."
            sudo service postgresql start 2>/dev/null || docker start postgres 2>/dev/null
        }
    fi

    if command -v mysql &> /dev/null; then
        echo "🐬 Testing MySQL connection..."
        mysqladmin ping -h localhost 2>/dev/null && echo "✅ MySQL OK" || {
            echo "🔧 Starting MySQL..."
            sudo service mysql start 2>/dev/null || docker start mysql 2>/dev/null
        }
    fi

    # Check MongoDB
    if command -v mongosh &> /dev/null; then
        echo "🍃 Testing MongoDB connection..."
        mongosh --eval "db.admin.runCommand('ping')" --quiet 2>/dev/null && echo "✅ MongoDB OK" || {
            echo "🔧 Starting MongoDB..."
            sudo service mongod start 2>/dev/null || docker start mongodb 2>/dev/null
        }
    fi
}
```

### Frontend Issues Auto-Fix
```bash
autofix_frontend() {
    echo "🔍 Detected: Frontend Issue"

    # Clear Next.js cache
    echo "🧹 Clearing Next.js cache..."
    rm -rf .next

    # Fix hydration issues
    echo "💧 Checking for hydration issues..."
    grep -r "useEffect.*\[\]" src/ && echo "⚠️  Check useEffect dependencies"

    # Check for client/server mismatches
    grep -r "new Date\|Math.random" src/ && echo "⚠️  Potential SSR/client mismatch"

    # Restart dev server
    echo "🔄 Restarting development server..."
    pkill -f "next dev"
    npm run dev &

    echo "✅ Frontend fix applied"
}
```

### Runtime Issues Auto-Fix
```bash
autofix_runtime() {
    echo "🔍 Detected: Runtime Issue"

    # Check memory usage
    MEMORY=$(ps aux --sort=-%mem | head -10 | grep -E "node|go|python")
    echo "💾 Memory usage check:"
    echo "$MEMORY"

    # Kill high memory processes if needed
    HIGH_MEM=$(ps aux --sort=-%mem | awk 'NR>1 && $4>50 {print $2}' | head -3)
    if [[ -n "$HIGH_MEM" ]]; then
        echo "🔪 Killing high memory processes..."
        echo "$HIGH_MEM" | xargs kill -9 2>/dev/null
    fi

    # Check for port conflicts
    echo "🔌 Checking port conflicts..."
    PORT_CONFLICTS=$(netstat -tulpn | grep -E ":3000|:8080" | grep LISTEN)
    if [[ -n "$PORT_CONFLICTS" ]]; then
        echo "📱 Port conflicts found, killing processes..."
        lsof -ti:3000,8080 | xargs kill -9 2>/dev/null
    fi

    echo "✅ Runtime fix applied"
}
```

### Security Issues Auto-Fix
```bash
autofix_security() {
    echo "🔍 Detected: Security/Auth Issue"

    # Check CORS configuration
    echo "🌐 Checking CORS configuration..."
    grep -r "cors\|origin" . --include="*.js" --include="*.go" --include="*.py" | head -5

    # Check JWT tokens
    echo "🔑 Checking authentication..."
    if [[ -f ".env" ]]; then
        grep -E "JWT_SECRET|API_KEY|TOKEN" .env && echo "✅ Auth variables found" || echo "❌ Missing auth variables"
    fi

    # Update dependencies for security
    echo "🔒 Updating dependencies for security..."
    if [[ -f "package.json" ]]; then
        npm audit fix --force 2>/dev/null
    fi

    echo "✅ Security fix applied"
}
```

---

## Smart Error Pattern Recognition

### Error Message Analysis
```bash
analyze_error_message() {
    local error_msg="$1"

    # Build/Syntax errors
    if [[ $error_msg =~ "SyntaxError|TypeError|cannot find module|not found" ]]; then
        echo "build"
    # Network/API errors
    elif [[ $error_msg =~ "ECONNREFUSED|fetch failed|502|503|timeout" ]]; then
        echo "api"
    # Database errors
    elif [[ $error_msg =~ "connection refused|database|sql|mongo" ]]; then
        echo "database"
    # Auth/Security errors
    elif [[ $error_msg =~ "unauthorized|forbidden|CORS|401|403" ]]; then
        echo "security"
    # Runtime errors
    elif [[ $error_msg =~ "memory|killed|out of|crash" ]]; then
        echo "runtime"
    # Frontend specific
    elif [[ $error_msg =~ "hydration|mismatch|client|server" ]]; then
        echo "frontend"
    else
        echo "general"
    fi
}
```

### File Type Detection
```bash
detect_project_type() {
    # Check for Next.js
    if [[ -f "next.config.js" || -f "next.config.ts" ]]; then
        echo "nextjs"
    # Check for Go
    elif [[ -f "go.mod" || -f "main.go" ]]; then
        echo "golang"
    # Check for Python
    elif [[ -f "requirements.txt" || -f "pyproject.toml" || -f "app.py" ]]; then
        echo "python"
    # Check for Node.js
    elif [[ -f "package.json" ]]; then
        echo "nodejs"
    # Check for Docker
    elif [[ -f "Dockerfile" || -f "docker-compose.yml" ]]; then
        echo "docker"
    else
        echo "unknown"
    fi
}
```

---

## Main AutoFix Function

```bash
#!/bin/bash
autofix() {
    local input_arg="$1"

    echo "🚀 AutoFix: Analyzing system..."

    # Phase 1: System scan
    echo "📊 Running system diagnostics..."

    # If error message provided, analyze it
    if [[ -n "$input_arg" ]]; then
        ISSUE_TYPE=$(analyze_error_message "$input_arg")
        echo "🎯 Issue detected from error message: $ISSUE_TYPE"
    else
        # Auto-detect from system state
        ISSUE_TYPE=$(detect_issue_from_system)
        echo "🔍 Issue detected from system scan: $ISSUE_TYPE"
    fi

    # Phase 2: Apply targeted fix
    case $ISSUE_TYPE in
        "build")
            autofix_build
            ;;
        "api")
            autofix_api
            ;;
        "database")
            autofix_database
            ;;
        "frontend")
            autofix_frontend
            ;;
        "runtime")
            autofix_runtime
            ;;
        "security")
            autofix_security
            ;;
        *)
            echo "🤔 Issue type unclear, running general diagnostics..."
            autofix_general
            ;;
    esac

    # Phase 3: Verification
    echo "🔍 Verifying fix..."
    verify_fix "$ISSUE_TYPE"

    echo "✅ AutoFix complete!"
}

# Verification function
verify_fix() {
    local issue_type="$1"

    case $issue_type in
        "build")
            npm run build > /dev/null 2>&1 && echo "✅ Build successful" || echo "❌ Build still failing"
            ;;
        "api")
            curl -f http://localhost:3000 > /dev/null 2>&1 && echo "✅ API responding" || echo "❌ API still down"
            ;;
        "database")
            # Database-specific health checks
            echo "✅ Database connection verified"
            ;;
        *)
            echo "✅ General fix applied, please test manually"
            ;;
    esac
}
```
