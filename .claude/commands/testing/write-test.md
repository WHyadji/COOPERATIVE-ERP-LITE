# Enhanced Write Test - Multi-Stack Testing Framework (SOLID + TDD)

```yaml
---
allowed-tools: ReadFile, WriteFile, SearchReplace, Bash(npm:*), Bash(go:*), Bash(python:*), Bash(pip:*), Bash(git:*)
description: Comprehensive test generation with SOLID principles, TDD, and multi-technology support
---

## Pre-Test Analysis & Discovery

### Codebase Architecture Discovery
- Project structure: !`find . -type f -name "*.ts" -o -name "*.js" -o -name "*.go" -o -name "*.py" | head -20`
- Test frameworks: !`grep -r "jest\\|vitest\\|testing\\|pytest\\|testify\\|ginkgo" package.json go.mod requirements.txt pyproject.toml 2>/dev/null || echo "Detecting test frameworks"`
- Existing tests: !`find . -name "*test*" -o -name "*spec*" -type f | head -15`
- Test coverage: !`find . -name "coverage" -o -name ".nyc_output" -o -name "htmlcov" -type d | head -10`

### Frontend Testing Environment (Next.js/React)
- React Testing Library: !`grep -r "@testing-library" package.json || echo "RTL not detected"`
- Jest configuration: !`find . -name "jest.config.*" -o -name ".jestrc*" | head -5`
- Test utilities: !`find . -path "*/test*" -name "*.ts" -o -name "*.js" | head-10`
- Component tests: !`find . -name "*.test.tsx" -o -name "*.spec.tsx" | head-10`
- Hook tests: !`grep -r "renderHook\\|act" . --include="*.test.*" --include="*.spec.*" | head-5`

### Go Testing Environment
- Go test files: !`find . -name "*_test.go" | head-15`
- Test frameworks: !`grep -r "testify\\|ginkgo\\|gomega" go.mod go.sum 2>/dev/null || echo "Standard testing"`
- Mocking libraries: !`grep -r "gomock\\|testify/mock\\|counterfeiter" go.mod go.sum 2>/dev/null || echo "No mocking libs"`
- Benchmark tests: !`grep -r "func Benchmark" . --include="*_test.go" | head-5`
- Integration tests: !`find . -name "*integration*test*.go" | head-5`

### Python Testing Environment
- Test frameworks: !`grep -r "pytest\\|unittest\\|nose" requirements.txt pyproject.toml setup.py 2>/dev/null || echo "Standard unittest"`
- Mocking libraries: !`grep -r "mock\\|pytest-mock\\|factory-boy" requirements.txt pyproject.toml 2>/dev/null || echo "No mocking libs"`
- Test structure: !`find . -name "test_*.py" -o -name "*_test.py" | head-15`
- Fixtures: !`grep -r "@pytest.fixture\\|@fixture" . --include="*.py" | head-5`
- Property testing: !`grep -r "hypothesis\\|property" . --include="*.py" | head-3`

### Service Architecture Testing
- API contracts: !`find . -name "*contract*" -o -name "*schema*" -o -name "openapi*" | head-10`
- Integration points: !`grep -r "http://\\|https://\\|localhost:" . --include="*.go" --include="*.py" --include="*.ts" | head-10`
- Database tests: !`grep -r "testcontainers\\|docker.*test\\|memory.*db" . --include="*.go" --include="*.py" | head-5`
- Message queue tests: !`grep -r "kafka.*test\\|redis.*test\\|rabbitmq.*test" . --include="*.go" --include="*.py" | head-5`

### SOLID Principle Test Patterns
- Interface testing: !`grep -r "interface\\|Protocol\\|ABC" . --include="*.go" --include="*.py" --include="*.ts" | head-10`
- Dependency injection: !`grep -r "inject\\|DI\\|container" . --include="*.go" --include="*.py" --include="*.ts" | head-5`
- Mock implementations: !`grep -r "mock\\|stub\\|fake" . --include="*test*" | head-10`
- Strategy pattern tests: !`grep -r "strategy\\|algorithm" . --include="*test*" | head-5`

### Convention Enforcement (References: /name-conventions)
- Clean names: !`echo "Applying clean, business-focused naming"`
- No jargon: !`echo "Avoiding technical prefixes and versioning"`
- Business focus: !`echo "Using domain-appropriate terminology"`
- Validation: !`echo "Checking names follow conventions"`

## Your Task

**Code to Test:** $ARGUMENTS

## Test Generation Strategy (SOLID + TDD)

### 1. SOLID-Based Test Architecture

#### Single Responsibility Principle (SRP) Testing
**Each test should verify one specific behavior**

```typescript
// ❌ Violates SRP - Testing multiple concerns
describe('UserService', () => {
  it('should validate, save, and send email for user', () => {
    // Tests validation, persistence, and email sending
  });
});

// ✅ Follows SRP - Focused tests
describe('UserValidator', () => {
  it('should reject invalid email format', () => {
    // Tests only validation logic
  });
});

describe('UserRepository', () => {
  it('should save user to database', () => {
    // Tests only persistence logic
  });
});

describe('EmailService', () => {
  it('should send welcome email', () => {
    // Tests only email sending
  });
});
```

#### Open/Closed Principle (OCP) Testing
**Tests should be extensible without modification**

```python
# ✅ Strategy Pattern Testing
class TestPaymentProcessor:
    @pytest.fixture
    def processors(self):
        return [
            CreditCardProcessor(),
            PayPalProcessor(),
            BitcoinProcessor()  # New processor, no test changes needed
        ]

    @pytest.mark.parametrize("processor", processors)
    def test_payment_processing(self, processor):
        result = processor.process(100.0)
        assert result.success is True
        assert result.amount == 100.0
```

#### Liskov Substitution Principle (LSP) Testing
**Subclass tests should pass for base class contracts**

```go
// ✅ Contract-based testing
func TestBirdContract(t *testing.T) {
    birds := []Bird{
        &Eagle{},
        &Sparrow{},
        // All implementations must satisfy Bird contract
    }

    for _, bird := range birds {
        t.Run(fmt.Sprintf("%T", bird), func(t *testing.T) {
            // Test common behavior
            assert.NotNil(t, bird.GetSpecies())
            assert.True(t, bird.CanMove())
        })
    }
}

func TestFlyableBirds(t *testing.T) {
    flyableBirds := []FlyableBird{
        &Eagle{},
        &Sparrow{},
        // Only flying birds implement this interface
    }

    for _, bird := range flyableBirds {
        t.Run(fmt.Sprintf("%T", bird), func(t *testing.T) {
            distance := bird.Fly(10)
            assert.Greater(t, distance, 0.0)
        })
    }
}
```

#### Interface Segregation Principle (ISP) Testing
**Test focused interfaces separately**

```typescript
// ✅ Segregated interface testing
interface Readable {
  read(): string;
}

interface Writable {
  write(data: string): void;
}

describe('FileReader (Readable)', () => {
  it('should read file content', () => {
    const reader: Readable = new FileReader('test.txt');
    expect(reader.read()).toBe('file content');
  });
});

describe('FileWriter (Writable)', () => {
  it('should write file content', () => {
    const writer: Writable = new FileWriter('test.txt');
    writer.write('new content');
    // Verify write operation
  });
});

describe('FileManager (Readable & Writable)', () => {
  it('should read and write files', () => {
    const manager: Readable & Writable = new FileManager('test.txt');
    manager.write('test data');
    expect(manager.read()).toBe('test data');
  });
});
```

#### Dependency Inversion Principle (DIP) Testing
**Test with injected dependencies (mocks/stubs)**

```python
# ✅ Dependency injection testing
class TestUserService:
    @pytest.fixture
    def mock_repository(self):
        return Mock(spec=UserRepository)

    @pytest.fixture
    def mock_email_service(self):
        return Mock(spec=EmailService)

    @pytest.fixture
    def user_service(self, mock_repository, mock_email_service):
        return UserService(
            repository=mock_repository,
            email_service=mock_email_service
        )

    def test_create_user_success(self, user_service, mock_repository, mock_email_service):
        # Arrange
        user = User(email="test@example.com")
        mock_repository.save.return_value = user

        # Act
        result = user_service.create_user(user)

        # Assert
        mock_repository.save.assert_called_once_with(user)
        mock_email_service.send_welcome.assert_called_once_with(user)
        assert result == user
```

### 2. Test-Driven Development (TDD) Workflow

#### Red-Green-Refactor Cycle
1. **RED**: Write failing test
2. **GREEN**: Write minimal code to pass
3. **REFACTOR**: Improve code while keeping tests green

#### TDD Test Templates

**Unit Test Template (Go)**
```go
func TestNewFeature(t *testing.T) {
    // Arrange
    given := setupTestData()
    expected := expectedResult()

    // Act
    actual := subject.DoSomething(given)

    // Assert
    assert.Equal(t, expected, actual)
}

func TestNewFeature_ErrorCase(t *testing.T) {
    // Arrange
    invalidInput := setupInvalidData()

    // Act & Assert
    _, err := subject.DoSomething(invalidInput)
    assert.Error(t, err)
    assert.Contains(t, err.Error(), "expected error message")
}
```

**Component Test Template (React/TypeScript)**
```typescript
describe('ComponentName', () => {
  it('should render with required props', () => {
    // Arrange
    const props = {
      title: 'Test Title',
      onSubmit: jest.fn()
    };

    // Act
    render(<ComponentName {...props} />);

    // Assert
    expect(screen.getByText('Test Title')).toBeInTheDocument();
  });

  it('should call onSubmit when form is submitted', async () => {
    // Arrange
    const mockSubmit = jest.fn();
    const user = userEvent.setup();

    render(<ComponentName onSubmit={mockSubmit} />);

    // Act
    await user.click(screen.getByRole('button', { name: /submit/i }));

    // Assert
    expect(mockSubmit).toHaveBeenCalledTimes(1);
  });
});
```

**Service Test Template (Python)**
```python
class TestServiceName:
    @pytest.fixture
    def service(self):
        return ServiceName()

    def test_method_success_case(self, service):
        # Arrange
        input_data = {"key": "value"}
        expected = ExpectedResult()

        # Act
        result = service.method(input_data)

        # Assert
        assert result == expected

    def test_method_validation_error(self, service):
        # Arrange
        invalid_data = {"invalid": "data"}

        # Act & Assert
        with pytest.raises(ValidationError) as exc_info:
            service.method(invalid_data)

        assert "validation message" in str(exc_info.value)
```

### 3. Multi-Layer Testing Strategy

#### Unit Tests (Isolation)
- **Scope**: Single function/method/class
- **Dependencies**: Mocked/stubbed
- **Speed**: Very fast (< 10ms)
- **Reliability**: High

```typescript
// Fast, isolated unit test
describe('calculateTax', () => {
  it('should calculate correct tax for given amount', () => {
    expect(calculateTax(100, 0.1)).toBe(10);
  });
});
```

#### Integration Tests (Component Interaction)
- **Scope**: Multiple components working together
- **Dependencies**: Real implementations where possible
- **Speed**: Medium (< 1s)
- **Reliability**: Medium

```go
func TestUserServiceIntegration(t *testing.T) {
    // Use real database with test data
    db := setupTestDatabase()
    defer db.Close()

    userRepo := NewUserRepository(db)
    emailService := NewMockEmailService()
    userService := NewUserService(userRepo, emailService)

    user := &User{Email: "test@example.com"}
    createdUser, err := userService.CreateUser(user)

    assert.NoError(t, err)
    assert.NotEmpty(t, createdUser.ID)
}
```

#### End-to-End Tests (Full System)
- **Scope**: Complete user workflows
- **Dependencies**: Real services/databases
- **Speed**: Slow (> 1s)
- **Reliability**: Lower

```python
def test_user_registration_flow():
    # Start real services
    client = TestClient(app)

    # Complete user journey
    response = client.post("/api/register", json={
        "email": "test@example.com",
        "password": "secure123"
    })

    assert response.status_code == 201
    assert "id" in response.json()

    # Verify email was sent
    assert email_service.sent_emails[-1].to == "test@example.com"
```

### 4. Advanced Testing Patterns

#### Property-Based Testing
```python
from hypothesis import given, strategies as st

@given(st.integers(min_value=0), st.floats(min_value=0, max_value=1))
def test_tax_calculation_properties(amount, rate):
    result = calculate_tax(amount, rate)

    # Properties that should always hold
    assert result >= 0  # Tax can't be negative
    assert result <= amount  # Tax can't exceed amount
    if rate == 0:
        assert result == 0  # No tax when rate is zero
```

#### Contract Testing (API)
```typescript
// Using Pact for contract testing
describe('User API Contract', () => {
  const provider = new Pact({
    consumer: 'Frontend',
    provider: 'UserService'
  });

  it('should get user by ID', async () => {
    await provider
      .given('user exists')
      .uponReceiving('a request for user')
      .withRequest({
        method: 'GET',
        path: '/api/users/123'
      })
      .willRespondWith({
        status: 200,
        body: { id: 123, name: 'John Doe' }
      });

    const response = await userApi.getUser(123);
    expect(response.name).toBe('John Doe');
  });
});
```

#### Mutation Testing
```bash
# Go mutation testing
go install github.com/go-mutesting/mutesting
mutesting ./...

# JavaScript mutation testing
npm install --save-dev stryker-cli
stryker run
```

#### Performance Testing
```go
func BenchmarkUserCreation(b *testing.B) {
    service := setupUserService()
    user := &User{Email: "test@example.com"}

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        service.CreateUser(user)
    }
}

func TestUserCreationPerformance(t *testing.T) {
    service := setupUserService()
    user := &User{Email: "test@example.com"}

    start := time.Now()
    service.CreateUser(user)
    duration := time.Since(start)

    assert.Less(t, duration, 100*time.Millisecond)
}
```

### 5. Test Organization & Best Practices

#### Test Structure (AAA Pattern)
```python
def test_method_name():
    # Arrange - Set up test data and dependencies
    user = User(email="test@example.com")
    mock_repo = Mock()
    service = UserService(mock_repo)

    # Act - Execute the code under test
    result = service.create_user(user)

    # Assert - Verify the outcome
    assert result.id is not None
    mock_repo.save.assert_called_once()
```

#### Test Naming Conventions
- **Method**: `test_should_[expected_behavior]_when_[condition]`
- **Describe**: Component/class being tested
- **It**: Specific behavior being verified

```typescript
describe('UserValidator', () => {
  describe('validateEmail', () => {
    it('should return true when email format is valid', () => {
      // Test implementation
    });

    it('should return false when email format is invalid', () => {
      // Test implementation
    });

    it('should throw error when email is null', () => {
      // Test implementation
    });
  });
});
```

#### Test Data Management
```go
// Test fixtures
type UserFixture struct {
    ValidUser   *User
    InvalidUser *User
}

func NewUserFixture() *UserFixture {
    return &UserFixture{
        ValidUser: &User{
            Email: "valid@example.com",
            Name:  "Valid User",
        },
        InvalidUser: &User{
            Email: "invalid-email",
            Name:  "",
        },
    }
}

func TestUserValidation(t *testing.T) {
    fixture := NewUserFixture()
    validator := NewUserValidator()

    t.Run("valid user", func(t *testing.T) {
        err := validator.Validate(fixture.ValidUser)
        assert.NoError(t, err)
    })

    t.Run("invalid user", func(t *testing.T) {
        err := validator.Validate(fixture.InvalidUser)
        assert.Error(t, err)
    })
}
```

### 6. Test Automation & CI/CD

#### Test Pipeline Configuration
```yaml
# .github/workflows/test.yml
name: Test Suite
on: [push, pull_request]

jobs:
  unit-tests:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - name: Run unit tests
        run: |
          npm test -- --coverage
          go test ./... -v -race -coverprofile=coverage.out
          pytest --cov=src tests/

  integration-tests:
    runs-on: ubuntu-latest
    services:
      postgres:
        image: postgres:13
        env:
          POSTGRES_PASSWORD: test
    steps:
      - name: Run integration tests
        run: |
          npm run test:integration
          go test ./... -tags=integration
          pytest tests/integration/

  e2e-tests:
    runs-on: ubuntu-latest
    steps:
      - name: Run e2e tests
        run: |
          npm run test:e2e
          pytest tests/e2e/
```

#### Coverage Requirements
```json
// jest.config.js
{
  "coverageThreshold": {
    "global": {
      "branches": 80,
      "functions": 80,
      "lines": 80,
      "statements": 80
    }
  }
}
```

### 7. Microservice Testing Strategies

#### Service-Level Testing
```python
# Test service boundaries
class TestUserService:
    def test_service_isolation(self):
        # Service should work independently
        user_service = UserService()

        # Mock external dependencies
        with patch('external_email_service') as mock_email:
            user = user_service.create_user(valid_user_data)
            assert user.id is not None
            mock_email.send.assert_called_once()
```

#### Contract Testing Between Services
```typescript
// API contract tests
describe('User Service API Contract', () => {
  it('should maintain backward compatibility', async () => {
    const response = await request(app)
      .get('/api/v1/users/123')
      .expect(200);

    // Verify response structure hasn't changed
    expect(response.body).toMatchObject({
      id: expect.any(Number),
      email: expect.any(String),
      created_at: expect.any(String)
    });
  });
});
```

#### Event-Driven Testing
```go
func TestEventPublishing(t *testing.T) {
    eventBus := NewMockEventBus()
    userService := NewUserService(eventBus)

    user := &User{Email: "test@example.com"}
    userService.CreateUser(user)

    // Verify event was published
    events := eventBus.GetPublishedEvents()
    assert.Len(t, events, 1)
    assert.Equal(t, "UserCreated", events[0].Type)
}
```

### 8. Test Maintenance & Quality

#### Test Code Quality
- **DRY**: Don't repeat test setup
- **Clear**: Tests should be self-documenting
- **Fast**: Unit tests should run quickly
- **Isolated**: Tests shouldn't depend on each other
- **Deterministic**: Same input, same output

#### Test Refactoring
```python
# Before: Duplicated setup
class TestUserService:
    def test_create_user(self):
        db = Database()
        email_service = EmailService()
        user_service = UserService(db, email_service)
        # test logic

    def test_update_user(self):
        db = Database()
        email_service = EmailService()
        user_service = UserService(db, email_service)
        # test logic

# After: Shared fixtures
class TestUserService:
    @pytest.fixture
    def user_service(self):
        db = Database()
        email_service = EmailService()
        return UserService(db, email_service)

    def test_create_user(self, user_service):
        # test logic

    def test_update_user(self, user_service):
        # test logic
```

### 9. Advanced Testing Tools & Frameworks

#### Frontend Testing Stack
- **Jest/Vitest**: Unit testing framework
- **React Testing Library**: Component testing
- **Playwright/Cypress**: E2E testing
- **MSW**: API mocking
- **Storybook**: Component documentation & testing

#### Backend Testing Stack
- **Go**: testify, ginkgo, gomega, testcontainers
- **Python**: pytest, unittest, factory-boy, responses
- **Database**: testcontainers, in-memory databases
- **API**: httptest, requests-mock

#### Testing Utilities
```typescript
// Custom testing utilities
export const renderWithProviders = (
  ui: React.ReactElement,
  options: RenderOptions = {}
) => {
  const AllTheProviders = ({ children }: { children: React.ReactNode }) => (
    <QueryProvider>
      <ThemeProvider>
        <AuthProvider>
          {children}
        </AuthProvider>
      </ThemeProvider>
    </QueryProvider>
  );

  return render(ui, { wrapper: AllTheProviders, ...options });
};

// Test factory
export const createUser = (overrides: Partial<User> = {}): User => ({
  id: Math.random().toString(),
  email: 'test@example.com',
  name: 'Test User',
  created_at: new Date().toISOString(),
  ...overrides
});
```

Remember: Good tests are an investment in code quality, maintainability, and team confidence. They should be treated as first-class citizens in your codebase, following the same quality standards as production code.
```
