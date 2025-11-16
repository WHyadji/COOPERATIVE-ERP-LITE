# Naming Conventions - Clean and Simple Names

## 📋 **IMPORTANT NAMING GUIDELINES**

All slash commands must follow clean, simple naming conventions. Avoid technical jargon, versioning, or implementation details in names.

---

## ✅ **DO - Use Clean, Business-Focused Names**

### Database Tables
```sql
-- ✅ GOOD: Clean, descriptive names
users
roles  
permissions
role_permissions
user_sessions
payment_methods
subscription_plans
billing_events
audit_logs
security_events

-- ✅ GOOD: Simple relationship tables
user_roles
plan_features
order_items
```

### Services and Classes
```javascript
// ✅ GOOD: Clear, purpose-driven names
UserService
AuthService
PaymentService
NotificationService
SubscriptionService
ValidationService
SecurityService
AuditService

// ✅ GOOD: Model names
User
Role
Permission
Subscription
BillingEvent
```

### API Endpoints
```bash
# ✅ GOOD: RESTful, resource-focused
/api/users
/api/roles
/api/permissions
/api/subscriptions
/api/payments
/api/notifications
/api/audit-logs
```

### Configuration and Environment Variables
```bash
# ✅ GOOD: Clear purpose
JWT_SECRET
DATABASE_URL
STRIPE_API_KEY
REDIS_URL
LOG_LEVEL
APP_PORT
```

---

## ❌ **DON'T - Avoid Technical Jargon and Prefixes**

### Database Tables - Avoid These Patterns
```sql
-- ❌ BAD: Technical prefixes and versioning
enhanced_users
rbac_v2_roles
auth_system_permissions
new_role_permissions
improved_user_sessions
advanced_payment_methods
v3_subscription_plans
updated_billing_events
security_audit_logs_v2
enhanced_security_events

-- ❌ BAD: Implementation details in names
jwt_user_sessions
oauth2_auth_tokens
bcrypt_password_hashes
redis_cache_entries
```

### Services and Classes - Avoid These Patterns
```javascript
// ❌ BAD: Technical implementation details
EnhancedUserService
AdvancedAuthService
ImprovedPaymentService
ModernNotificationService
NewSubscriptionService
UpdatedValidationService
SecureSecurityService
ComprehensiveAuditService

// ❌ BAD: Version numbers and technical terms
UserServiceV2
OAuth2AuthService
JWTTokenService
BCryptHashService
RedisSessionService
```

### API Endpoints - Avoid These Patterns
```bash
# ❌ BAD: Technical prefixes and versions
/api/v2/enhanced-users
/api/new/rbac-roles
/api/improved/permissions
/api/advanced/subscriptions
/api/modern/payments
/api/updated/notifications
/api/secure/audit-logs
```

---

## 🎯 **NAMING PRINCIPLES**

### 1. **Business Domain First**
- Name things by **what they represent** in the business domain
- Use **user-facing terminology** when possible
- Think about how a business person would refer to it

```javascript
// ✅ GOOD: Business domain focused
class SubscriptionService {
  createSubscription(userId, planId) { }
  cancelSubscription(subscriptionId) { }
  upgradeSubscription(subscriptionId, newPlanId) { }
}

// ❌ BAD: Technical implementation focused
class EnhancedSubscriptionManagementService {
  createNewSubscriptionRecord(userId, planId) { }
  performSubscriptionCancellation(subscriptionId) { }
  executeSubscriptionUpgradeProcess(subscriptionId, newPlanId) { }
}
```

### 2. **Simple and Clear**
- Use the **simplest name** that clearly conveys purpose
- Avoid unnecessary adjectives like "enhanced", "improved", "advanced"
- One concept = one clear name

```sql
-- ✅ GOOD: Simple and clear
CREATE TABLE users (
  id UUID PRIMARY KEY,
  email VARCHAR(255) UNIQUE,
  password_hash VARCHAR(255),
  created_at TIMESTAMP
);

-- ❌ BAD: Overcomplicated
CREATE TABLE enhanced_user_management_records (
  advanced_user_id UUID PRIMARY KEY,
  improved_email_address VARCHAR(255) UNIQUE,
  secure_password_hash VARCHAR(255),
  enhanced_creation_timestamp TIMESTAMP
);
```

### 3. **Consistent Patterns**
- Follow the same naming pattern across the entire system
- Use consistent relationship naming (`user_roles`, not `users_to_roles`)
- Maintain consistency in pluralization

```javascript
// ✅ GOOD: Consistent patterns
users          → UserService       → /api/users
roles          → RoleService       → /api/roles  
permissions    → PermissionService → /api/permissions
subscriptions  → SubscriptionService → /api/subscriptions

// ❌ BAD: Inconsistent patterns
users           → EnhancedUserMgmt     → /api/v2/user-management
role_definitions → AdvancedRoleService → /api/modern/role-system
permission_sets  → PermService         → /api/perms
subscriptions   → SubMgmtService      → /api/sub-mgmt
```

### 4. **Avoid Implementation Details**
- Don't include technology choices in names
- Don't include architectural patterns
- Focus on **what** it does, not **how** it does it

```javascript
// ✅ GOOD: Implementation agnostic
class SessionService {
  createSession(userId) { }
  validateSession(sessionId) { }
  destroySession(sessionId) { }
}

// ❌ BAD: Implementation details exposed
class JWTRedisSessionManagementService {
  createJWTSessionWithRedisBackend(userId) { }
  validateJWTTokenAgainstRedisStore(sessionId) { }
  destroyJWTSessionAndClearRedisCache(sessionId) { }
}
```

---

## 🔧 **PRACTICAL EXAMPLES**

### Database Schema Design
```sql
-- ✅ GOOD: Clean schema names
CREATE TABLE users (
  id UUID PRIMARY KEY,
  email VARCHAR(255) UNIQUE,
  first_name VARCHAR(100),
  last_name VARCHAR(100),
  created_at TIMESTAMP,
  updated_at TIMESTAMP
);

CREATE TABLE roles (
  id UUID PRIMARY KEY,
  name VARCHAR(50) UNIQUE,
  description TEXT,
  created_at TIMESTAMP
);

CREATE TABLE permissions (
  id UUID PRIMARY KEY,
  name VARCHAR(50) UNIQUE,
  resource VARCHAR(50),
  action VARCHAR(50),
  created_at TIMESTAMP
);

CREATE TABLE user_roles (
  user_id UUID REFERENCES users(id),
  role_id UUID REFERENCES roles(id),
  assigned_at TIMESTAMP,
  assigned_by UUID REFERENCES users(id),
  PRIMARY KEY (user_id, role_id)
);

CREATE TABLE role_permissions (
  role_id UUID REFERENCES roles(id),
  permission_id UUID REFERENCES permissions(id),
  granted_at TIMESTAMP,
  granted_by UUID REFERENCES users(id),
  PRIMARY KEY (role_id, permission_id)
);
```

### Service Layer Implementation
```javascript
// ✅ GOOD: Clean service names and methods
class UserService {
  async createUser(userData) { }
  async getUserById(userId) { }
  async updateUser(userId, updateData) { }
  async deleteUser(userId) { }
  async getUsersByRole(roleId) { }
}

class RoleService {
  async createRole(roleData) { }
  async assignRoleToUser(userId, roleId) { }
  async removeRoleFromUser(userId, roleId) { }
  async getUserRoles(userId) { }
}

class PermissionService {
  async checkPermission(userId, resource, action) { }
  async grantPermission(roleId, permissionId) { }
  async revokePermission(roleId, permissionId) { }
  async getRolePermissions(roleId) { }
}
```

### API Endpoint Design
```javascript
// ✅ GOOD: Clean RESTful endpoints
// Users
GET    /api/users
POST   /api/users
GET    /api/users/:id
PUT    /api/users/:id
DELETE /api/users/:id

// Roles
GET    /api/roles
POST   /api/roles
GET    /api/roles/:id
PUT    /api/roles/:id
DELETE /api/roles/:id

// Role assignments
POST   /api/users/:userId/roles/:roleId
DELETE /api/users/:userId/roles/:roleId
GET    /api/users/:userId/roles

// Permissions
GET    /api/permissions
POST   /api/permissions
GET    /api/roles/:roleId/permissions
POST   /api/roles/:roleId/permissions/:permissionId
DELETE /api/roles/:roleId/permissions/:permissionId
```

---

## 🚀 **COMMAND INTEGRATION EXAMPLES**

### Feature Planning with Clean Names
```bash
# ✅ GOOD: Clean feature names
/plan-feature "user-management" medium high weeks full-stack
/plan-feature "role-permissions" low medium weeks backend
/plan-feature "subscription-billing" high urgent months full-stack

# ❌ AVOID: Technical jargon in feature names
/plan-feature "enhanced-rbac-v2-system" medium high weeks full-stack
/plan-feature "advanced-permission-framework" low medium weeks backend
/plan-feature "improved-subscription-management" high urgent months full-stack
```

### Implementation with Clean Names
```bash
# ✅ GOOD: Clear implementation scope
/implement-feature "user-authentication" core development backend
/implement-feature "payment-processing" foundation production full-stack
/implement-feature "audit-logging" integration testing backend

# ❌ AVOID: Technical implementation details
/implement-feature "jwt-oauth2-auth-system" core development backend
/implement-feature "stripe-payment-gateway-integration" foundation production full-stack
/implement-feature "comprehensive-security-audit-framework" integration testing backend
```

### Security Audit with Clean Focus
```bash
# ✅ GOOD: Clear security scope
/security-audit authentication deep owasp access-control
/security-audit permissions standard nist authorization
/security-audit payments comprehensive pci-dss data-protection

# ❌ AVOID: Technical security jargon
/security-audit jwt-token-validation deep owasp crypto-analysis
/security-audit rbac-authorization-framework standard nist access-audit
/security-audit payment-gateway-integration comprehensive pci-dss encryption-analysis
```

---

## 📝 **QUICK REFERENCE CHECKLIST**

Before naming anything, ask yourself:

- [ ] **Is this name business-focused?** (What it represents, not how it works)
- [ ] **Is this the simplest clear name?** (No unnecessary adjectives)
- [ ] **Is this consistent with existing patterns?** (Follows established conventions)
- [ ] **Does this avoid implementation details?** (Technology-agnostic)
- [ ] **Would a business person understand this name?** (Domain-appropriate)

---

## 🎯 **SUMMARY**

**DO USE:**
- `users`, `roles`, `permissions`, `subscriptions`
- `UserService`, `PaymentService`, `NotificationService`
- `/api/users`, `/api/payments`, `/api/subscriptions`

**DON'T USE:**
- `enhanced_users`, `rbac_v2_roles`, `advanced_permissions`
- `EnhancedUserService`, `ImprovedPaymentService`, `ModernNotificationService`
- `/api/v2/enhanced-users`, `/api/advanced/payments`, `/api/modern/subscriptions`

**Remember:** Clean, simple names make code more maintainable, understandable, and professional. Focus on **what** things do, not **how** they do it.