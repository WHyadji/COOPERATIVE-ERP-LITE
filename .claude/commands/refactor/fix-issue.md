# Enhanced Debug & Fix Issues - Next.js & Services (SOLID Principles)

```yaml
---
allowed-tools: ReadFile, WriteFile, SearchReplace, Bash(npm:*), Bash(go:*), Bash(python:*), Bash(pip:*), Bash(git:*)
description: Comprehensive debugging and issue resolution for Next.js and microservices with SOLID principles
---

## Pre-Diagnostic Health Check

### Next.js Frontend Analysis
- TypeScript errors: !`cd frontend && npm run type-check 2>&1 | head -30 || echo "No Next.js frontend found"`
- Next.js build issues: !`cd frontend && npm run build 2>&1 | grep -E "ERROR|FAILED|error" | head -15 || echo "No build errors"`
- Linter warnings: !`cd frontend && npm run lint 2>&1 | head -30 || echo "No lint errors"`
- Runtime errors: !`grep -r "console.error\\|throw new Error\\|Error(" frontend/src --include="*.ts*" --include="*.js*" -n | tail -20 || echo "No frontend errors found"`
- SOLID violations: !`grep -r "class.*extends.*extends\\|new.*new.*new" frontend/src --include="*.ts*" -n | head-10 || echo "No apparent SOLID violations"`

### Go Services Analysis
- Go compilation: !`find services -name "*.go" -path "*/cmd/*" -o -path "*/main.go" | xargs -I {} sh -c 'cd $(dirname {}) && go build' 2>&1 | head -20 || echo "No Go services found"`
- Go tests: !`find services -name "*_test.go" | xargs -I {} sh -c 'cd $(dirname {}) && go test ./...' 2>&1 | grep -E "FAIL|ERROR|panic" | head -15 || echo "No Go tests found"`
- Go vet issues: !`find services -name "*.go" | xargs -I {} sh -c 'cd $(dirname {}) && go vet ./...' 2>&1 | head -15 || echo "No Go vet issues"`
- Interface segregation: !`grep -r "type.*interface" services --include="*.go" -A 10 | grep -E "func.*func.*func" | head-10 || echo "No large interfaces found"`
- Dependency injection: !`grep -r "\\*.*\\*.*\\*" services --include="*.go" | head-10 || echo "No apparent tight coupling"`

### Python Services Analysis
- Python syntax: !`find services -name "*.py" -exec python -m py_compile {} \; 2>&1 | head -15 || echo "No Python services found"`
- Python tests: !`find services -name "test_*.py" -o -name "*_test.py" | xargs -I {} sh -c 'cd $(dirname {}) && python -m pytest {}' 2>&1 | grep -E "FAILED|ERROR" | head-15 || echo "No Python tests found"`
- Python linting: !`find services -name "*.py" | xargs python -m flake8 2>&1 | head -20 || echo "No Python linting errors"`
- Abstract base classes: !`grep -r "ABC\\|abstractmethod" services --include="*.py" -n | head-15 || echo "No abstract classes found"`
- Single responsibility: !`grep -r "class.*:" services --include="*.py" -A 20 | grep -E "def.*def.*def.*def.*def" | head-10 || echo "No apparent SRP violations"`

### Service Architecture Analysis
- Service boundaries: !`find services -type d -maxdepth 1 | wc -l && echo "services found"`
- API contracts: !`grep -r "'/api\\|'/v1\\|swagger\\|openapi" services --include="*.py" --include="*.go" --include="*.yaml" -n | head-15 || echo "No API contracts found"`
- Shared dependencies: !`find services -name "go.mod" -o -name "requirements.txt" -o -name "pyproject.toml" | head-10 || echo "No service dependencies found"`
- Interface definitions: !`grep -r "interface\\|Protocol\\|ABC" services --include="*.py" --include="*.go" -n | head-15 || echo "No interfaces found"`

### SOLID Principles Validation
- **S**RP violations: !`grep -r "class\\|type.*struct" services --include="*.py" --include="*.go" -A 15 | grep -c "func\\|def" | head-10 || echo "Checking single responsibility"`
- **O**CP adherence: !`grep -r "interface\\|Protocol\\|abstract" services --include="*.py" --include="*.go" -n | head-15 || echo "Checking open/closed principle"`
- **L**SP violations: !`grep -r "isinstance\\|type.*assertion" services --include="*.py" --include="*.go" -n | head-10 || echo "Checking Liskov substitution"`
- **I**SP adherence: !`grep -r "interface.*{" services --include="*.go" -A 20 | grep -c "func" | head-10 || echo "Checking interface segregation"`
- **D**IP violations: !`grep -r "import.*\\.\\." services --include="*.py" --include="*.go" -n | head-15 || echo "Checking dependency inversion"`

### Cross-Service Communication
- API consistency: !`grep -r "http://\\|https://\\|localhost:" services --include="*.py" --include="*.go" --include="*.ts*" -n | head-15 || echo "No API endpoints found"`
- Message queues: !`grep -r "kafka\\|rabbitmq\\|redis\\|nats" services --include="*.py" --include="*.go" --include="*.yaml" -n | head-10 || echo "No message queues found"`
- Service discovery: !`grep -r "consul\\|etcd\\|eureka" services --include="*.py" --include="*.go" --include="*.yaml" -n | head-10 || echo "No service discovery found"`
- Circuit breakers: !`grep -r "circuit\\|breaker\\|retry\\|timeout" services --include="*.py" --include="*.go" -n | head-10 || echo "No resilience patterns found"`

### Dependencies & Security
- Node.js security: !`cd frontend && npm audit --audit-level=moderate 2>&1 | head -20 || echo "No npm security issues"`
- Go security: !`find services -name "go.mod" | xargs -I {} sh -c 'cd $(dirname {}) && go list -json -deps ./...' | grep -E "Module|ImportPath" | head-20 || echo "No Go dependencies"`
- Python security: !`find services -name "requirements.txt" | xargs -I {} pip-audit -r {} 2>&1 | head-20 || echo "No Python security scan available"`
- Shared vulnerabilities: !`grep -r "TODO\\|FIXME\\|HACK" services --include="*.py" --include="*.go" | head-15 || echo "No code smell comments found"`

## SOLID Principles Implementation Guide

### Single Responsibility Principle (SRP)
**Each class/module should have only one reason to change**

#### Detection Patterns:
- Classes with multiple unrelated methods
- Functions doing more than one thing
- Mixed concerns (business logic + data access + presentation)

#### Refactoring Strategies:
- **Extract Class**: Split large classes into focused ones
- **Extract Method**: Break down large functions
- **Separate Concerns**: Business logic, data access, presentation
- **Use Composition**: Prefer composition over inheritance

#### Implementation Examples:
```go
// ❌ Violates SRP
type UserService struct{}
func (s *UserService) ValidateUser(user User) error { /* validation */ }
func (s *UserService) SaveUser(user User) error { /* database */ }
func (s *UserService) SendEmail(user User) error { /* email */ }

// ✅ Follows SRP
type UserValidator struct{}
func (v *UserValidator) Validate(user User) error { /* validation only */ }

type UserRepository struct{}
func (r *UserRepository) Save(user User) error { /* database only */ }

type EmailService struct{}
func (e *EmailService) SendWelcomeEmail(user User) error { /* email only */ }
```

### Open/Closed Principle (OCP)
**Software entities should be open for extension, closed for modification**

#### Detection Patterns:
- Switch statements on types
- If-else chains for different behaviors
- Direct instantiation of concrete classes

#### Refactoring Strategies:
- **Strategy Pattern**: Encapsulate algorithms
- **Template Method**: Define skeleton, allow customization
- **Dependency Injection**: Inject behaviors
- **Interface Segregation**: Define contracts

#### Implementation Examples:
```python
# ❌ Violates OCP
class PaymentProcessor:
    def process(self, payment_type: str, amount: float):
        if payment_type == "credit_card":
            # credit card logic
        elif payment_type == "paypal":
            # paypal logic
        # Adding new type requires modification

# ✅ Follows OCP
from abc import ABC, abstractmethod

class PaymentProcessor(ABC):
    @abstractmethod
    def process(self, amount: float) -> bool:
        pass

class CreditCardProcessor(PaymentProcessor):
    def process(self, amount: float) -> bool:
        # credit card logic
        pass

class PayPalProcessor(PaymentProcessor):
    def process(self, amount: float) -> bool:
        # paypal logic
        pass
```

### Liskov Substitution Principle (LSP)
**Objects of a superclass should be replaceable with objects of subclasses**

#### Detection Patterns:
- Subclasses throwing "not implemented" exceptions
- Instanceof/type checking before method calls
- Subclasses requiring additional preconditions

#### Refactoring Strategies:
- **Contract by Design**: Ensure consistent interfaces
- **Behavioral Compatibility**: Maintain expected behavior
- **Precondition Weakening**: Subclasses accept more, not less
- **Postcondition Strengthening**: Subclasses guarantee more, not less

#### Implementation Examples:
```typescript
// ❌ Violates LSP
class Bird {
  fly(): void { /* flying logic */ }
}

class Penguin extends Bird {
  fly(): void {
    throw new Error("Penguins can't fly!"); // Violates LSP
  }
}

// ✅ Follows LSP
interface Flyable {
  fly(): void;
}

interface Swimmable {
  swim(): void;
}

class Eagle implements Flyable {
  fly(): void { /* flying logic */ }
}

class Penguin implements Swimmable {
  swim(): void { /* swimming logic */ }
}
```

### Interface Segregation Principle (ISP)
**Clients should not be forced to depend on interfaces they don't use**

#### Detection Patterns:
- Large interfaces with many methods
- Classes implementing empty methods
- Clients depending on unused functionality

#### Refactoring Strategies:
- **Split Interfaces**: Create focused, cohesive interfaces
- **Role Interfaces**: Define interfaces per client need
- **Composition**: Combine small interfaces when needed
- **Adapter Pattern**: Bridge incompatible interfaces

#### Implementation Examples:
```go
// ❌ Violates ISP
type Worker interface {
    Work()
    Eat()
    Sleep()
}

type Robot struct{}
func (r Robot) Work() { /* work */ }
func (r Robot) Eat() { /* robots don't eat! */ }
func (r Robot) Sleep() { /* robots don't sleep! */ }

// ✅ Follows ISP
type Worker interface {
    Work()
}

type Eater interface {
    Eat()
}

type Sleeper interface {
    Sleep()
}

type Human struct{}
func (h Human) Work() { /* work */ }
func (h Human) Eat() { /* eat */ }
func (h Human) Sleep() { /* sleep */ }

type Robot struct{}
func (r Robot) Work() { /* work */ }
```

### Dependency Inversion Principle (DIP)
**High-level modules should not depend on low-level modules. Both should depend on abstractions**

#### Detection Patterns:
- Direct instantiation of concrete classes
- Import statements to concrete implementations
- Hard-coded dependencies

#### Refactoring Strategies:
- **Dependency Injection**: Inject dependencies from outside
- **Inversion of Control**: Use containers/frameworks
- **Abstract Factory**: Create families of related objects
- **Service Locator**: Centralize dependency resolution

#### Implementation Examples:
```python
# ❌ Violates DIP
class EmailService:
    def send(self, message: str):
        # concrete email implementation
        pass

class UserService:
    def __init__(self):
        self.email_service = EmailService()  # Direct dependency

# ✅ Follows DIP
from abc import ABC, abstractmethod

class NotificationService(ABC):
    @abstractmethod
    def send(self, message: str):
        pass

class EmailService(NotificationService):
    def send(self, message: str):
        # email implementation
        pass

class UserService:
    def __init__(self, notification_service: NotificationService):
        self.notification_service = notification_service  # Abstraction
```

## Your Task

**Issue to Fix:** $ARGUMENTS

## SOLID-Based Debugging Approach

### 1. Architecture Assessment
- **Service Boundaries**: Are services properly segregated by business domain?
- **Interface Design**: Do interfaces follow ISP with focused responsibilities?
- **Dependency Flow**: Does the dependency graph follow DIP principles?
- **Extension Points**: Can new features be added without modifying existing code (OCP)?

### 2. Code Quality Analysis

#### Single Responsibility Violations
- **Large Classes**: Classes with multiple concerns
- **God Objects**: Classes that know too much
- **Mixed Layers**: Business logic mixed with data access
- **Utility Classes**: Classes with unrelated static methods

#### Open/Closed Violations
- **Switch Statements**: Type-based branching logic
- **Modification Patterns**: Frequent changes to existing classes
- **Hard-coded Behaviors**: Non-configurable business rules
- **Tight Coupling**: Classes that change together

#### Liskov Substitution Violations
- **Behavioral Differences**: Subclasses with different contracts
- **Exception Throwing**: Subclasses throwing unexpected exceptions
- **Precondition Strengthening**: Subclasses requiring more
- **Postcondition Weakening**: Subclasses guaranteeing less

#### Interface Segregation Violations
- **Fat Interfaces**: Interfaces with too many methods
- **Client Coupling**: Clients depending on unused methods
- **Role Confusion**: Mixed concerns in single interface
- **Implementation Burden**: Forced empty implementations

#### Dependency Inversion Violations
- **Concrete Dependencies**: Direct instantiation of implementations
- **Layering Issues**: High-level modules depending on low-level
- **Testing Difficulties**: Hard to mock dependencies
- **Configuration Coupling**: Hard-coded configuration

### 3. Service-Oriented Debugging

#### Microservice Patterns
- **Service Discovery**: How services find each other
- **Circuit Breakers**: Resilience and fault tolerance
- **API Gateway**: Centralized routing and cross-cutting concerns
- **Event Sourcing**: Event-driven architecture patterns

#### Communication Patterns
- **Synchronous**: REST, GraphQL, gRPC
- **Asynchronous**: Message queues, event streams
- **Saga Pattern**: Distributed transaction management
- **CQRS**: Command Query Responsibility Segregation

#### Data Management
- **Database per Service**: Data ownership boundaries
- **Shared Data**: Anti-pattern identification
- **Event Sourcing**: Audit trails and state reconstruction
- **Data Consistency**: Eventual consistency patterns

### 4. Technology-Specific SOLID Implementation

#### Next.js Frontend SOLID
- **Components**: Single-purpose, composable components
- **Hooks**: Custom hooks for specific concerns
- **Context**: Focused context providers
- **Pages**: Separation of concerns between data and presentation

#### Go Services SOLID
- **Interfaces**: Small, focused interfaces
- **Dependency Injection**: Constructor injection patterns
- **Middleware**: Cross-cutting concerns
- **Error Handling**: Consistent error interfaces

#### Python Services SOLID
- **Abstract Base Classes**: Define contracts
- **Dependency Injection**: Use frameworks like dependency-injector
- **Type Hints**: Define clear interfaces
- **Decorators**: Cross-cutting concerns

### 5. Refactoring Strategies

#### Extract Patterns
- **Extract Interface**: Define contracts
- **Extract Class**: Split responsibilities
- **Extract Method**: Break down complexity
- **Extract Service**: Separate concerns

#### Composition Patterns
- **Strategy Pattern**: Interchangeable algorithms
- **Decorator Pattern**: Add behavior without modification
- **Factory Pattern**: Create objects without coupling
- **Observer Pattern**: Loose coupling between components

#### Injection Patterns
- **Constructor Injection**: Required dependencies
- **Setter Injection**: Optional dependencies
- **Interface Injection**: Framework-driven injection
- **Service Locator**: Centralized dependency management

### 6. Testing with SOLID Principles

#### Unit Testing
- **Mock Interfaces**: Test in isolation
- **Dependency Injection**: Inject test doubles
- **Single Responsibility**: Test one thing at a time
- **Interface Contracts**: Test behavior, not implementation

#### Integration Testing
- **Service Boundaries**: Test service interactions
- **Contract Testing**: Verify API contracts
- **End-to-End**: Test complete user journeys
- **Performance**: Test under realistic conditions

### 7. Monitoring and Observability

#### Service Health
- **Health Checks**: Readiness and liveness probes
- **Metrics**: Business and technical metrics
- **Distributed Tracing**: Request flow across services
- **Log Aggregation**: Centralized logging with correlation IDs

#### SOLID Metrics
- **Cyclomatic Complexity**: Measure method complexity
- **Coupling Metrics**: Afferent and efferent coupling
- **Cohesion Metrics**: LCOM (Lack of Cohesion of Methods)
- **Dependency Metrics**: Stability and abstractness

## Emergency Response with SOLID Principles

### Rapid Assessment
1. **Identify Affected Services**: Which services are impacted?
2. **Check Interface Contracts**: Are interfaces being violated?
3. **Validate Dependencies**: Are injected dependencies available?
4. **Test Substitutability**: Can we swap implementations?

### Quick Fixes
- **Strategy Pattern**: Quickly swap algorithms
- **Feature Flags**: Toggle behavior without deployment
- **Circuit Breakers**: Isolate failing services
- **Fallback Patterns**: Graceful degradation

### Long-term Solutions
- **Refactor to SOLID**: Apply principles to prevent recurrence
- **Add Monitoring**: Detect violations early
- **Improve Testing**: Test contracts and behaviors
- **Documentation**: Document architectural decisions

Remember: SOLID principles lead to more maintainable, testable, and extensible code that's easier to debug and fix when issues arise.
```
