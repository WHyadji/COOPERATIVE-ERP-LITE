---
description: Mock data and hardcoded response detector for APIs and services
argument-hint: [service/component] (optional focus area)
allowed-tools: grep_search, codebase_search, read_file
---

# Mock Data and Hardcoded Response Detection

Identify mock data, hardcoded responses, and placeholder implementations that need to be replaced with real data sources and proper business logic.

## Detection Patterns

Search for the following mock data patterns:

### 1. Hardcoded JSON Responses
Look for static data returns in API endpoints:
```javascript
// Patterns to find:
return { id: 1, name: "Test User", email: "test@example.com" }
res.json({ success: true, data: "placeholder" })
response.send({ message: "TODO: Implement this" })
return { status: "ok", result: [] }
```

### 2. Static Array and Object Returns
Identify hardcoded collections and data structures:
```javascript
// Static arrays:
return [1, 2, 3, 4, 5] // Hardcoded IDs
getUserList() { return ["user1", "user2"] }
const users = ["admin", "guest"] // Static user lists

// Static objects:
const config = { apiUrl: "http://localhost:3000" }
return { total: 100, items: mockItems }
```

### 3. Mock Data Files and Generators
Find mock data imports and generation:
```javascript
// Mock imports:
const mockData = require('./mockData.json')
import { mockUsers } from './fixtures'
import faker from 'faker'

// Mock generators:
const mockUsers = generateMockData()
return Mock.of(User).generate()
faker.random.words()
```

### 4. Development-Only Code Paths
Detect environment-specific mock responses:
```javascript
// Environment checks:
if (process.env.NODE_ENV === 'development') {
  return getMockResponse()
}
if (isDevelopment) return mockData

// Test/debug endpoints:
app.get('/test', ...)
router.post('/dummy', ...)
app.use('/debug', debugRouter)
```

### 5. Database Query Placeholders
Find incomplete or hardcoded database operations:
```sql
-- Hardcoded queries:
SELECT * FROM users LIMIT 5; -- No proper filtering
INSERT INTO orders VALUES (1, 'test', 100); -- Static values
UPDATE products SET price = 99.99; -- No WHERE clause

-- TODO placeholders:
-- TODO: Add proper JOIN logic
-- FIXME: Replace with dynamic query
```

### 6. Time and ID Placeholders
Identify placeholder timestamps and IDs:
```javascript
// Time placeholders:
setTimeout(() => resolve(mockData), 1000) // Fake delays
return Promise.resolve({ fake: "data" })
createdAt: "2023-01-01T00:00:00Z" // Hardcoded dates

// ID placeholders:
Date.now() // Used as entity IDs
Math.random() // Used for IDs
id: "temp-" + Math.random()
```

### 7. Authentication and Security Mocks
Find bypassed security and auth placeholders:
```javascript
// Auth bypasses:
const isAuthenticated = () => true // Always returns true
if (true || checkAuth()) // Auth check bypassed
userId: "12345" // Hardcoded user ID
token: "mock-jwt-token"

// Permission mocks:
hasPermission: () => true
isAdmin: true // Hardcoded admin status
```

### 8. External Service Mocks
Identify mocked external API calls:
```javascript
// Mocked API calls:
const mockApiResponse = { success: true }
// return fetch(url) // Commented out real API call
return Promise.resolve(mockResponse)

// Service mocks:
const paymentService = { charge: () => ({ success: true }) }
emailService.send = () => console.log("Email sent") // Mock implementation
```

## Analysis Scope

${ARGUMENTS ? `**Focused analysis**: ${ARGUMENTS}` : '**Full codebase scan** for mock data across all services and components'}

## Expected Output

Provide a detailed report with:

### Executive Summary
```
Mock Data Analysis Report
========================
🎭 Total Mock Instances: X
📊 By Risk Level:
  🔴 Critical (Production Impact): X
  🟡 High (Business Logic): X
  🟠 Medium (Data Quality): X
  🟢 Low (Development Aid): X

📁 By Component:
  🌐 API Endpoints: X
  🗄️  Database Layer: X
  🔐 Authentication: X
  🛠️  External Services: X
  🧪 Test/Debug Code: X
```

### Detailed Findings
For each mock data instance:
- **Location**: `file:line` with surrounding context
- **Type**: Category (hardcoded response, mock service, etc.)
- **Code Snippet**: The actual mock implementation
- **Risk Assessment**: Why this needs to be addressed
- **Business Impact**: What functionality is affected
- **Replacement Strategy**: Specific steps to implement real data source

### Implementation Roadmap

#### Phase 1: Critical Replacements (Security & Production)
- Remove auth bypasses and hardcoded credentials
- Replace payment/financial mock implementations
- Fix hardcoded production configurations

#### Phase 2: Business Logic Implementation
- Replace mock business rules with real implementations
- Implement proper data validation and processing
- Connect to real external services

#### Phase 3: Data Quality Improvements
- Replace static data with dynamic database queries
- Implement proper filtering and pagination
- Add real-time data sources where needed

#### Phase 4: Development Cleanup
- Remove debug endpoints from production code
- Clean up development-only mock data
- Maintain legitimate test fixtures in appropriate locations

### Migration Recommendations
- **Data Source Mapping**: Which real data sources should replace each mock
- **API Integration Plan**: Steps to connect to real external services
- **Testing Strategy**: How to maintain test coverage during migration
- **Rollback Plan**: Safe migration approach with rollback options

Begin the comprehensive mock data scan and provide actionable replacement strategies.
