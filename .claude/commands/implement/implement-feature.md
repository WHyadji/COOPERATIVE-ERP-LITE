# `/implement-feature` Command - Detailed Implementation

## Command Structure
```bash
/implement-feature [plan-id] [phase] [mode] [scope] [--options]
```

## Parameters

### Required Parameters
- **plan-id**: Reference to existing feature plan (`plan-123` or `user-subscription-plan`)
- **phase**: `foundation` | `core` | `integration` | `polish` | `full` | `custom:<name>`
- **mode**: `development` | `testing` | `staging` | `production` | `dry-run`

### Optional Parameters
- **scope**: `backend` | `frontend` | `full-stack` | `infrastructure` | `database` | `api`

### Options
- `--incremental`: Implement in small, testable chunks
- `--parallel`: Enable parallel implementation tracks
- `--safe-mode`: Extra validation and rollback preparation
- `--fast-track`: Skip some validations for urgent implementations
- `--code-review`: Require review before each commit
- `--pair-programming`: Generate pair programming session plan
- `--documentation`: Auto-generate documentation during implementation
- `--testing-first`: TDD approach with tests before implementation
- `--performance-monitoring`: Add performance tracking to implementation
- `--feature-flags`: Implement with feature flag integration
- `--rollback-ready`: Prepare rollback mechanisms during implementation

### Convention Enforcement (References: /name-conventions)
- Clean names: !`echo "Applying clean, business-focused naming"`
- No jargon: !`echo "Avoiding technical prefixes and versioning"`
- Business focus: !`echo "Using domain-appropriate terminology"`
- Validation: !`echo "Checking names follow conventions"`

## Implementation Flow

### Phase 1: Gemini Context Ingestion & Code Analysis

```javascript
const geminiImplementationAnalysis = {
  // 1. Plan Context Loading
  planContextAnalysis: {
    // Load and validate the feature plan
    loadFeaturePlan: async (planId) => {
      const plan = await planStorage.load(planId);

      return {
        planValidation: {
          isValid: true,
          completeness: 0.95,
          staleness: "2 days old",
          dependencies: plan.dependencies,
          assumptions: plan.assumptions
        },

        implementationReadiness: {
          prerequisites: [
            {
              requirement: "Database schema updated",
              status: "completed",
              verifiedAt: "2025-07-10T10:00:00Z"
            },
            {
              requirement: "API interfaces defined",
              status: "pending",
              blockedBy: "Team review needed"
            }
          ],

          environmentReadiness: {
            development: "ready",
            testing: "ready",
            staging: "needs-deployment",
            production: "not-ready"
          }
        },

        riskFactors: [
          {
            risk: "Third-party API rate limits",
            currentStatus: "monitored",
            mitigation: "Implemented exponential backoff"
          }
        ]
      };
    },

    // Analyze current codebase state vs plan expectations
    codebaseStateAnalysis: async (plan, currentCodebase) => {
      return {
        // Check if codebase matches plan assumptions
        planCodebaseAlignment: {
          architecturalChanges: [
            {
              expected: "Service layer for subscriptions",
              actual: "Monolithic controller structure",
              impact: "Refactoring needed",
              effort: "2-3 hours"
            }
          ],

          dependencyChanges: [
            {
              dependency: "stripe",
              plannedVersion: "^10.0.0",
              currentVersion: "^9.5.0",
              upgradeRequired: true,
              breakingChanges: ["Payment intent API changes"]
            }
          ]
        },

        // Identify what's already implemented
        existingImplementation: {
          completedComponents: [
            {
              component: "User model",
              status: "complete",
              alignment: "matches plan",
              reusability: "high"
            }
          ],

          partialImplementations: [
            {
              component: "Payment service",
              completeness: 0.6,
              gaps: ["Webhook handling", "Error recovery"],
              reworkNeeded: false
            }
          ],

          conflictingCode: [
            {
              location: "src/auth/UserController.js",
              issue: "Different authentication pattern than planned",
              resolution: "Adapter pattern needed"
            }
          ]
        }
      };
    }
  },

  // 2. Implementation Context Preparation
  implementationContextPrep: {
    // Analyze implementation patterns in existing codebase
    patternAnalysis: async (codebase, targetPhase) => {
      return {
        // Identify coding patterns to follow
        codingPatterns: {
          fileStructure: {
            pattern: "Feature-based folders",
            example: "src/features/subscriptions/",
            conventions: ["services/", "controllers/", "models/", "tests/"]
          },

          namingConventions: {
            functions: "camelCase",
            classes: "PascalCase",
            constants: "UPPER_SNAKE_CASE",
            files: "kebab-case"
          },

          errorHandling: {
            pattern: "Custom error classes with error codes",
            example: "throw new ValidationError('INVALID_EMAIL', email)",
            logging: "Structured logging with context"
          },

          testing: {
            framework: "Jest",
            pattern: "AAA (Arrange, Act, Assert)",
            coverage: "90% minimum",
            mockingStrategy: "Dependency injection with test doubles"
          }
        },

        // Map integration points and dependencies
        integrationMapping: {
          existingServices: [
            {
              service: "AuthService",
              interface: "IAuthService",
              location: "src/auth/AuthService.js",
              integrationPoints: ["User validation", "Permission checking"]
            }
          ],

          externalAPIs: [
            {
              api: "Stripe",
              wrapper: "PaymentService",
              errorHandling: "Retry with exponential backoff",
              testingStrategy: "Mock with real response structures"
            }
          ],

          dataLayerIntegration: {
            orm: "Sequelize",
            migrationStrategy: "Incremental with rollback",
            seedingApproach: "Environment-specific fixtures"
          }
        }
      };
    },

    // Prepare detailed implementation instructions
    implementationInstructions: async (plan, patterns, phase) => {
      return {
        phaseSpecificInstructions: {
          foundation: {
            priority: ["Database schema", "Core models", "Service interfaces"],
            order: "Bottom-up (data layer first)",
            validationCriteria: ["All migrations pass", "Models validate correctly"]
          },

          core: {
            priority: ["Business logic", "API endpoints", "Integration services"],
            order: "Outside-in (API first, then implementation)",
            validationCriteria: ["API contracts fulfilled", "Business rules enforced"]
          },

          integration: {
            priority: ["UI components", "End-to-end flows", "Error handling"],
            order: "User journey driven",
            validationCriteria: ["User stories completed", "Error scenarios handled"]
          }
        },

        codeGeneration: {
          templates: "Use existing project templates",
          boilerplate: "Generate from established patterns",
          customization: "Adapt to specific feature requirements",
          quality: "Follow project linting and formatting rules"
        },

        incrementalDelivery: {
          chunkSize: "Single user story or component",
          integrationFrequency: "Every 2-3 commits",
          feedbackLoop: "Automated tests + manual verification",
          rollbackStrategy: "Feature flags + database transactions"
        }
      };
    }
  }
};
```

### Phase 2: Claude Precise Code Implementation

```javascript
const claudeFeatureImplementation = {
  // 1. Strategic Implementation Orchestration
  implementationOrchestration: {
    // Create detailed implementation plan for the phase
    createImplementationPlan: (geminiContext, phase, options) => {
      return {
        executionStrategy: {
          approach: options.incremental ? "incremental" : "batch",
          parallelTracks: options.parallel ? [
            {
              track: "backend",
              dependencies: [],
              estimatedTime: "3 days"
            },
            {
              track: "frontend",
              dependencies: ["backend API contracts"],
              estimatedTime: "2 days"
            }
          ] : null,

          qualityGates: [
            {
              gate: "Unit tests pass",
              frequency: "Every commit",
              blocking: true
            },
            {
              gate: "Integration tests pass",
              frequency: "Every feature completion",
              blocking: true
            },
            {
              gate: "Code review approved",
              frequency: "Every pull request",
              blocking: options.codeReview
            }
          ]
        },

        implementationOrder: [
          {
            step: 1,
            component: "Database migrations",
            reason: "Foundation for all other components",
            estimatedTime: "2 hours",
            riskLevel: "low"
          },
          {
            step: 2,
            component: "Domain models",
            reason: "Core business logic foundation",
            estimatedTime: "4 hours",
            riskLevel: "low"
          },
          {
            step: 3,
            component: "Service layer",
            reason: "Business logic implementation",
            estimatedTime: "8 hours",
            riskLevel: "medium"
          },
          {
            step: 4,
            component: "API controllers",
            reason: "External interface definition",
            estimatedTime: "6 hours",
            riskLevel: "medium"
          }
        ]
      };
    },

    // Generate specific implementation tasks
    generateImplementationTasks: (plan, context, phase) => {
      return {
        tasks: [
          {
            id: "SUB-001",
            title: "Implement Subscription Model",
            description: "Create Subscription domain model with validation and business rules",

            technicalSpecs: {
              files: [
                "src/features/subscriptions/models/Subscription.js",
                "src/features/subscriptions/models/SubscriptionStatus.js"
              ],

              implementation: {
                baseClass: "BaseModel",
                interfaces: ["ISubscription"],
                validations: [
                  "userId must be valid UUID",
                  "planId must exist in plans table",
                  "status must be valid enum value"
                ],
                relationships: [
                  "belongsTo User",
                  "belongsTo Plan",
                  "hasMany BillingEvents"
                ]
              },

              businessRules: [
                "Only one active subscription per user per plan type",
                "Status transitions must follow state machine rules",
                "Cancellation requires proper cleanup of related services"
              ]
            },

            acceptanceCriteria: [
              "Model validates all required fields",
              "Business rules enforced at model level",
              "All relationships work correctly",
              "State transitions validate properly"
            ],

            testRequirements: {
              unitTests: [
                "Valid subscription creation",
                "Invalid data rejection",
                "State transition validation",
                "Business rule enforcement"
              ],

              integrationTests: [
                "Database persistence",
                "Relationship loading",
                "Transaction rollback on errors"
              ]
            },

            estimatedTime: "4 hours",
            dependencies: ["Database migration complete"],
            riskLevel: "low"
          }
        ]
      };
    }
  },

  // 2. Code Generation Engine
  codeGeneration: {
    // Generate actual implementation code
    generateCode: (task, patterns, context) => {
      return {
        // Database Migration
        migration: {
          file: "migrations/20250712_create_subscriptions.js",
          content: `
            'use strict';

            module.exports = {
              up: async (queryInterface, Sequelize) => {
                await queryInterface.createTable('subscriptions', {
                  id: {
                    type: Sequelize.UUID,
                    defaultValue: Sequelize.UUIDV4,
                    primaryKey: true,
                    allowNull: false
                  },
                  userId: {
                    type: Sequelize.UUID,
                    allowNull: false,
                    references: {
                      model: 'users',
                      key: 'id'
                    },
                    onUpdate: 'CASCADE',
                    onDelete: 'CASCADE'
                  },
                  planId: {
                    type: Sequelize.UUID,
                    allowNull: false,
                    references: {
                      model: 'subscription_plans',
                      key: 'id'
                    }
                  },
                  status: {
                    type: Sequelize.ENUM('pending', 'active', 'cancelled', 'expired'),
                    allowNull: false,
                    defaultValue: 'pending'
                  },
                  startDate: {
                    type: Sequelize.DATE,
                    allowNull: true
                  },
                  endDate: {
                    type: Sequelize.DATE,
                    allowNull: true
                  },
                  createdAt: {
                    type: Sequelize.DATE,
                    allowNull: false
                  },
                  updatedAt: {
                    type: Sequelize.DATE,
                    allowNull: false
                  }
                });

                // Add indexes for performance
                await queryInterface.addIndex('subscriptions', ['userId']);
                await queryInterface.addIndex('subscriptions', ['status']);
                await queryInterface.addIndex('subscriptions', ['planId']);

                // Add unique constraint for active subscriptions
                await queryInterface.addConstraint('subscriptions', {
                  fields: ['userId', 'planId'],
                  type: 'unique',
                  name: 'unique_active_user_plan',
                  where: {
                    status: 'active'
                  }
                });
              },

              down: async (queryInterface, Sequelize) => {
                await queryInterface.dropTable('subscriptions');
              }
            };
          `
        },

        // Domain Model
        model: {
          file: "src/features/subscriptions/models/Subscription.js",
          content: `
            const { Model, DataTypes } = require('sequelize');
            const { ValidationError } = require('../../../core/errors');

            class Subscription extends Model {
              static init(sequelize) {
                super.init({
                  id: {
                    type: DataTypes.UUID,
                    defaultValue: DataTypes.UUIDV4,
                    primaryKey: true
                  },
                  userId: {
                    type: DataTypes.UUID,
                    allowNull: false,
                    validate: {
                      notNull: {
                        msg: 'User ID is required'
                      },
                      isUUID: {
                        msg: 'User ID must be a valid UUID'
                      }
                    }
                  },
                  planId: {
                    type: DataTypes.UUID,
                    allowNull: false,
                    validate: {
                      notNull: {
                        msg: 'Plan ID is required'
                      },
                      async planExists(value) {
                        const plan = await this.sequelize.models.SubscriptionPlan.findByPk(value);
                        if (!plan) {
                          throw new Error('Plan does not exist');
                        }
                      }
                    }
                  },
                  status: {
                    type: DataTypes.ENUM('pending', 'active', 'cancelled', 'expired'),
                    allowNull: false,
                    defaultValue: 'pending',
                    validate: {
                      isIn: {
                        args: [['pending', 'active', 'cancelled', 'expired']],
                        msg: 'Invalid subscription status'
                      }
                    }
                  },
                  startDate: {
                    type: DataTypes.DATE,
                    allowNull: true,
                    validate: {
                      isDate: true,
                      isAfterCreation() {
                        if (this.startDate && this.startDate < this.createdAt) {
                          throw new Error('Start date cannot be before creation date');
                        }
                      }
                    }
                  },
                  endDate: {
                    type: DataTypes.DATE,
                    allowNull: true,
                    validate: {
                      isDate: true,
                      isAfterStart() {
                        if (this.endDate && this.startDate && this.endDate <= this.startDate) {
                          throw new Error('End date must be after start date');
                        }
                      }
                    }
                  }
                }, {
                  sequelize,
                  modelName: 'Subscription',
                  tableName: 'subscriptions',
                  hooks: {
                    beforeValidate: async (subscription) => {
                      // Validate business rules before saving
                      await subscription.validateBusinessRules();
                    },
                    beforeCreate: async (subscription) => {
                      // Set start date for active subscriptions
                      if (subscription.status === 'active' && !subscription.startDate) {
                        subscription.startDate = new Date();
                      }
                    }
                  }
                });

                return Subscription;
              }

              static associate(models) {
                Subscription.belongsTo(models.User, {
                  foreignKey: 'userId',
                  as: 'user'
                });

                Subscription.belongsTo(models.SubscriptionPlan, {
                  foreignKey: 'planId',
                  as: 'plan'
                });

                Subscription.hasMany(models.BillingEvent, {
                  foreignKey: 'subscriptionId',
                  as: 'billingEvents'
                });
              }

              // Business Logic Methods
              async validateBusinessRules() {
                // Rule: Only one active subscription per user per plan type
                if (this.status === 'active') {
                  const existingActive = await Subscription.findOne({
                    where: {
                      userId: this.userId,
                      planId: this.planId,
                      status: 'active',
                      id: { [this.sequelize.Sequelize.Op.ne]: this.id }
                    }
                  });

                  if (existingActive) {
                    throw new ValidationError('User already has an active subscription for this plan');
                  }
                }
              }

              async activate() {
                await this.validateStatusTransition('active');
                this.status = 'active';
                this.startDate = new Date();
                await this.save();

                // Emit domain event
                await this.emitEvent('subscription.activated', {
                  subscriptionId: this.id,
                  userId: this.userId,
                  planId: this.planId
                });
              }

              async cancel(reason = null) {
                await this.validateStatusTransition('cancelled');
                this.status = 'cancelled';
                this.endDate = new Date();
                await this.save();

                // Emit domain event
                await this.emitEvent('subscription.cancelled', {
                  subscriptionId: this.id,
                  userId: this.userId,
                  reason
                });
              }

              async validateStatusTransition(newStatus) {
                const validTransitions = {
                  pending: ['active', 'cancelled'],
                  active: ['cancelled', 'expired'],
                  cancelled: [],
                  expired: []
                };

                if (!validTransitions[this.status]?.includes(newStatus)) {
                  throw new ValidationError(
                    \`Invalid status transition from \${this.status} to \${newStatus}\`
                  );
                }
              }

              async emitEvent(eventType, payload) {
                // Integration with event system
                const EventBus = require('../../../core/EventBus');
                await EventBus.emit(eventType, payload);
              }

              // Query Methods
              static async findActiveByUser(userId) {
                return await Subscription.findAll({
                  where: {
                    userId,
                    status: 'active'
                  },
                  include: ['plan']
                });
              }

              static async findExpiring(days = 7) {
                const expirationDate = new Date();
                expirationDate.setDate(expirationDate.getDate() + days);

                return await Subscription.findAll({
                  where: {
                    status: 'active',
                    endDate: {
                      [this.sequelize.Sequelize.Op.lte]: expirationDate
                    }
                  },
                  include: ['user', 'plan']
                });
              }
            }

            module.exports = Subscription;
          `
        },

        // Service Layer
        service: {
          file: "src/features/subscriptions/services/SubscriptionService.js",
          content: `
            const Subscription = require('../models/Subscription');
            const PaymentService = require('../../payment/PaymentService');
            const NotificationService = require('../../../core/NotificationService');
            const { ValidationError, BusinessLogicError } = require('../../../core/errors');

            class SubscriptionService {
              constructor(paymentService = new PaymentService(), notificationService = new NotificationService()) {
                this.paymentService = paymentService;
                this.notificationService = notificationService;
              }

              async createSubscription(userId, planId, paymentMethodId) {
                const transaction = await Subscription.sequelize.transaction();

                try {
                  // Validate user and plan
                  const user = await this.validateUser(userId);
                  const plan = await this.validatePlan(planId);

                  // Check for existing active subscriptions
                  const existingSubscription = await Subscription.findOne({
                    where: { userId, planId, status: 'active' },
                    transaction
                  });

                  if (existingSubscription) {
                    throw new BusinessLogicError('User already has an active subscription for this plan');
                  }

                  // Create pending subscription
                  const subscription = await Subscription.create({
                    userId,
                    planId,
                    status: 'pending'
                  }, { transaction });

                  // Process payment
                  const paymentResult = await this.paymentService.processSubscriptionPayment({
                    amount: plan.price,
                    currency: plan.currency,
                    paymentMethodId,
                    subscriptionId: subscription.id,
                    customerId: user.stripeCustomerId
                  });

                  if (paymentResult.status === 'succeeded') {
                    await subscription.activate();
                    await this.notificationService.sendSubscriptionWelcome(user, subscription, plan);
                  }

                  await transaction.commit();
                  return await subscription.reload({ include: ['plan', 'user'] });

                } catch (error) {
                  await transaction.rollback();
                  throw error;
                }
              }

              async cancelSubscription(subscriptionId, userId, reason = null) {
                const subscription = await this.findUserSubscription(subscriptionId, userId);

                if (subscription.status !== 'active') {
                  throw new BusinessLogicError('Can only cancel active subscriptions');
                }

                // Cancel recurring payments
                await this.paymentService.cancelRecurringPayment(subscription.stripeSubscriptionId);

                // Cancel subscription
                await subscription.cancel(reason);

                // Send cancellation notification
                await this.notificationService.sendSubscriptionCancellation(
                  subscription.user,
                  subscription,
                  reason
                );

                return subscription;
              }

              async upgradeSubscription(subscriptionId, userId, newPlanId) {
                const transaction = await Subscription.sequelize.transaction();

                try {
                  const subscription = await this.findUserSubscription(subscriptionId, userId);
                  const newPlan = await this.validatePlan(newPlanId);

                  // Calculate prorated amount
                  const proratedAmount = await this.calculateProration(subscription, newPlan);

                  // Process upgrade payment if needed
                  if (proratedAmount > 0) {
                    await this.paymentService.processUpgradePayment({
                      subscriptionId,
                      amount: proratedAmount,
                      customerId: subscription.user.stripeCustomerId
                    });
                  }

                  // Update subscription
                  subscription.planId = newPlanId;
                  await subscription.save({ transaction });

                  await transaction.commit();
                  return await subscription.reload({ include: ['plan', 'user'] });

                } catch (error) {
                  await transaction.rollback();
                  throw error;
                }
              }

              // Helper Methods
              async validateUser(userId) {
                const User = Subscription.sequelize.models.User;
                const user = await User.findByPk(userId);
                if (!user) {
                  throw new ValidationError('User not found');
                }
                return user;
              }

              async validatePlan(planId) {
                const SubscriptionPlan = Subscription.sequelize.models.SubscriptionPlan;
                const plan = await SubscriptionPlan.findByPk(planId);
                if (!plan || !plan.isActive) {
                  throw new ValidationError('Plan not found or inactive');
                }
                return plan;
              }

              async findUserSubscription(subscriptionId, userId) {
                const subscription = await Subscription.findOne({
                  where: { id: subscriptionId, userId },
                  include: ['user', 'plan']
                });

                if (!subscription) {
                  throw new ValidationError('Subscription not found');
                }

                return subscription;
              }

              async calculateProration(subscription, newPlan) {
                // Implement prorated billing calculation
                const currentPlan = subscription.plan;
                const daysRemaining = this.getDaysRemaining(subscription);
                const dailyOldRate = currentPlan.price / 30;
                const dailyNewRate = newPlan.price / 30;

                return Math.max(0, (dailyNewRate - dailyOldRate) * daysRemaining);
              }

              getDaysRemaining(subscription) {
                const now = new Date();
                const endDate = subscription.endDate || this.calculateEndDate(subscription);
                return Math.max(0, Math.ceil((endDate - now) / (1000 * 60 * 60 * 24)));
              }
            }

            module.exports = SubscriptionService;
          `
        },

        // Unit Tests
        tests: {
          file: "src/features/subscriptions/tests/Subscription.test.js",
          content: `
            const { Subscription, User, SubscriptionPlan } = require('../../../models');
            const { ValidationError } = require('../../../core/errors');
            const { setupTestDB, cleanupTestDB } = require('../../../test/helpers');

            describe('Subscription Model', () => {
              beforeAll(async () => {
                await setupTestDB();
              });

              afterAll(async () => {
                await cleanupTestDB();
              });

              beforeEach(async () => {
                await Subscription.destroy({ where: {}, force: true });
                await User.destroy({ where: {}, force: true });
                await SubscriptionPlan.destroy({ where: {}, force: true });
              });

              describe('Validation', () => {
                test('should create valid subscription', async () => {
                  const user = await User.create({
                    email: 'test@example.com',
                    password: 'password123'
                  });

                  const plan = await SubscriptionPlan.create({
                    name: 'Basic Plan',
                    price: 9.99,
                    currency: 'USD'
                  });

                  const subscription = await Subscription.create({
                    userId: user.id,
                    planId: plan.id,
                    status: 'pending'
                  });

                  expect(subscription.id).toBeDefined();
                  expect(subscription.status).toBe('pending');
                });

                test('should reject subscription without userId', async () => {
                  const plan = await SubscriptionPlan.create({
                    name: 'Basic Plan',
                    price: 9.99,
                    currency: 'USD'
                  });

                  await expect(Subscription.create({
                    planId: plan.id,
                    status: 'pending'
                  })).rejects.toThrow('User ID is required');
                });

                test('should reject subscription with invalid planId', async () => {
                  const user = await User.create({
                    email: 'test@example.com',
                    password: 'password123'
                  });

                  await expect(Subscription.create({
                    userId: user.id,
                    planId: 'invalid-uuid',
                    status: 'pending'
                  })).rejects.toThrow();
                });
              });

              describe('Business Rules', () => {
                test('should prevent multiple active subscriptions for same plan', async () => {
                  const user = await User.create({
                    email: 'test@example.com',
                    password: 'password123'
                  });

                  const plan = await SubscriptionPlan.create({
                    name: 'Basic Plan',
                    price: 9.99,
                    currency: 'USD'
                  });

                  // Create first active subscription
                  await Subscription.create({
                    userId: user.id,
                    planId: plan.id,
                    status: 'active',
                    startDate: new Date()
                  });

                  // Attempt to create second active subscription
                  await expect(Subscription.create({
                    userId: user.id,
                    planId: plan.id,
                    status: 'active',
                    startDate: new Date()
                  })).rejects.toThrow('User already has an active subscription for this plan');
                });
              });

              describe('Status Transitions', () => {
                test('should allow valid status transitions', async () => {
                  const user = await User.create({
                    email: 'test@example.com',
                    password: 'password123'
                  });

                  const plan = await SubscriptionPlan.create({
                    name: 'Basic Plan',
                    price: 9.99,
                    currency: 'USD'
                  });

                  const subscription = await Subscription.create({
                    userId: user.id,
                    planId: plan.id,
                    status: 'pending'
                  });

                  // Should allow pending -> active
                  await subscription.activate();
                  expect(subscription.status).toBe('active');
                  expect(subscription.startDate).toBeDefined();

                  // Should allow active -> cancelled
                  await subscription.cancel('User requested');
                  expect(subscription.status).toBe('cancelled');
                  expect(subscription.endDate).toBeDefined();
                });

                test('should reject invalid status transitions', async () => {
                  const user = await User.create({
                    email: 'test@example.com',
                    password: 'password123'
                  });

                  const plan = await SubscriptionPlan.create({
                    name: 'Basic Plan',
                    price: 9.99,
                    currency: 'USD'
                  });

                  const subscription = await Subscription.create({
                    userId: user.id,
                    planId: plan.id,
                    status: 'cancelled'
                  });

                  // Should reject cancelled -> active
                  await expect(subscription.activate())
                    .rejects.toThrow('Invalid status transition from cancelled to active');
                });
              });

              describe('Query Methods', () => {
                test('should find active subscriptions by user', async () => {
                  const user = await User.create({
                    email: 'test@example.com',
                    password: 'password123'
                  });

                  const plan = await SubscriptionPlan.create({
                    name: 'Basic Plan',
                    price: 9.99,
                    currency: 'USD'
                  });

                  await Subscription.create({
                    userId: user.id,
                    planId: plan.id,
                    status: 'active',
                    startDate: new Date()
                  });

                  await Subscription.create({
                    userId: user.id,
                    planId: plan.id,
                    status: 'cancelled'
                  });

                  const activeSubscriptions = await Subscription.findActiveByUser(user.id);
                  expect(activeSubscriptions).toHaveLength(1);
                  expect(activeSubscriptions[0].status).toBe('active');
                });
              });
            });
          `
        }
      };
    },

    // Generate integration and API code
    generateAPICode: (spec, patterns) => {
      return {
        controller: {
          file: "src/features/subscriptions/controllers/SubscriptionController.js",
          content: `
            const SubscriptionService = require('../services/SubscriptionService');
            const { ValidationError, BusinessLogicError } = require('../../../core/errors');
            const { authenticate, authorize } = require('../../../middleware/auth');

            class SubscriptionController {
              constructor(subscriptionService = new SubscriptionService()) {
                this.subscriptionService = subscriptionService;
              }

              async createSubscription(req, res, next) {
                try {
                  const { planId, paymentMethodId } = req.body;
                  const userId = req.user.id;

                  // Input validation
                  if (!planId || !paymentMethodId) {
                    throw new ValidationError('Plan ID and payment method ID are required');
                  }

                  const subscription = await this.subscriptionService.createSubscription(
                    userId,
                    planId,
                    paymentMethodId
                  );

                  res.status(201).json({
                    success: true,
                    data: subscription,
                    message: 'Subscription created successfully'
                  });
                } catch (error) {
                  next(error);
                }
              }

              async getMySubscriptions(req, res, next) {
                try {
                  const userId = req.user.id;
                  const subscriptions = await Subscription.findActiveByUser(userId);

                  res.json({
                    success: true,
                    data: subscriptions
                  });
                } catch (error) {
                  next(error);
                }
              }

              async cancelSubscription(req, res, next) {
                try {
                  const { subscriptionId } = req.params;
                  const { reason } = req.body;
                  const userId = req.user.id;

                  const subscription = await this.subscriptionService.cancelSubscription(
                    subscriptionId,
                    userId,
                    reason
                  );

                  res.json({
                    success: true,
                    data: subscription,
                    message: 'Subscription cancelled successfully'
                  });
                } catch (error) {
                  next(error);
                }
              }
            }

            module.exports = SubscriptionController;
          `
        },

        routes: {
          file: "src/features/subscriptions/routes/subscriptionRoutes.js",
          content: `
            const express = require('express');
            const SubscriptionController = require('../controllers/SubscriptionController');
            const { authenticate, authorize } = require('../../../middleware/auth');
            const { validateRequest } = require('../../../middleware/validation');
            const { createSubscriptionSchema, cancelSubscriptionSchema } = require('../validation/subscriptionSchemas');

            const router = express.Router();
            const subscriptionController = new SubscriptionController();

            // Apply authentication to all routes
            router.use(authenticate);

            // Create subscription
            router.post('/',
              validateRequest(createSubscriptionSchema),
              subscriptionController.createSubscription.bind(subscriptionController)
            );

            // Get user's subscriptions
            router.get('/me',
              subscriptionController.getMySubscriptions.bind(subscriptionController)
            );

            // Cancel subscription
            router.delete('/:subscriptionId',
              validateRequest(cancelSubscriptionSchema),
              subscriptionController.cancelSubscription.bind(subscriptionController)
            );

            // Upgrade subscription
            router.put('/:subscriptionId/upgrade',
              authorize(['user', 'admin']),
              subscriptionController.upgradeSubscription.bind(subscriptionController)
            );

            module.exports = router;
          `
        }
      };
    }
  },

  // 3. Quality Assurance and Validation
  qualityAssurance: {
    // Generate comprehensive test suites
    generateTestSuite: (implementation, requirements) => {
      return {
        unitTests: "Generated unit tests for all models and services",
        integrationTests: "API endpoint tests with database integration",
        endToEndTests: "Complete user journey testing",
        performanceTests: "Load testing for subscription operations",
        securityTests: "Authentication and authorization validation"
      };
    },

    // Code quality validation
    validateCodeQuality: (generatedCode, standards) => {
      return {
        linting: "ESLint validation passed",
        formatting: "Prettier formatting applied",
        complexity: "Cyclomatic complexity within limits",
        coverage: "Test coverage above 90%",
        documentation: "JSDoc comments added to all public methods"
      };
    }
  }
};
```

## Usage Examples

### Basic Feature Implementation
```bash
/implement-feature "user-subscription-plan" core development full-stack --incremental
```

**Expected Output:**
```markdown
# Implementation: User Subscription Management - Core Phase

## Implementation Plan
**Phase**: Core Implementation
**Mode**: Development
**Approach**: Incremental delivery with continuous integration

## Tasks Completed ✅

### 1. Database Schema Implementation
- ✅ Created subscription table migration
- ✅ Added indexes for performance optimization
- ✅ Implemented unique constraints for business rules
- **Files**: `migrations/20250712_create_subscriptions.js`

### 2. Domain Model Implementation
- ✅ Subscription model with validation
- ✅ Business rule enforcement
- ✅ Status transition state machine
- ✅ Event emission for domain events
- **Files**: `src/features/subscriptions/models/Subscription.js`

### 3. Service Layer Implementation
- ✅ SubscriptionService with business logic
- ✅ Payment integration
- ✅ Error handling and transaction management
- ✅ Notification integration
- **Files**: `src/features/subscriptions/services/SubscriptionService.js`

### 4. API Implementation
- ✅ RESTful endpoints for subscription management
- ✅ Input validation and error handling
- ✅ Authentication and authorization
- **Files**:
  - `src/features/subscriptions/controllers/SubscriptionController.js`
  - `src/features/subscriptions/routes/subscriptionRoutes.js`

## Quality Metrics

### Test Coverage
- **Unit Tests**: 94% coverage
- **Integration Tests**: 87% coverage
- **API Tests**: 100% endpoint coverage

### Code Quality
- **ESLint**: 0 errors, 0 warnings
- **Complexity**: Average 4.2 (target: <5)
- **Documentation**: 100% JSDoc coverage

### Performance
- **API Response Time**: 145ms average
- **Database Queries**: Optimized with proper indexing
- **Memory Usage**: Within acceptable limits

## Integration Points Completed

### Payment Service Integration
- ✅ Stripe API integration for subscription billing
- ✅ Webhook handling for payment events
- ✅ Error handling and retry logic

### Notification Service Integration
- ✅ Welcome email on subscription activation
- ✅ Cancellation notification
- ✅ Billing event notifications

### Event System Integration
- ✅ Domain events for subscription lifecycle
- ✅ Event handlers for analytics and reporting

## Next Steps
1. **Frontend Integration**: Connect UI components to API
2. **Error Monitoring**: Set up error tracking and alerting
3. **Performance Testing**: Load testing with realistic data
4. **Security Review**: Code review focused on security
5. **Documentation**: Update API documentation and user guides

## Rollback Plan
- **Feature Flags**: `subscription_management_enabled` (currently: true)
- **Database**: Migration rollback script available
- **Monitoring**: Alerts configured for error rates and performance
```

### Advanced Implementation with Testing First
```bash
/implement-feature "payment-webhook-processing" foundation development backend --testing-first --safe-mode --code-review
```

**Expected Output includes**:
- TDD approach with tests written before implementation
- Enhanced error handling and validation
- Automated code review requirements
- Comprehensive rollback mechanisms
- Security-focused implementation

### Parallel Track Implementation
```bash
/implement-feature "user-dashboard-upgrade" integration development full-stack --parallel --performance-monitoring
```

**Expected Output includes**:
- Separate backend and frontend implementation tracks
- Real-time performance monitoring during implementation
- Dependency management between tracks
- Integrated delivery timeline

## Advanced Implementation Features

### Feature Flag Integration
```javascript
const featureFlagImplementation = {
  // Automatic feature flag creation and management
  flags: [
    {
      name: "subscription_management_v2",
      description: "New subscription management system",
      defaultValue: false,
      environments: {
        development: true,
        staging: true,
        production: false
      }
    }
  ],

  // Code generation includes feature flag checks
  implementationPattern: `
    if (FeatureFlags.isEnabled('subscription_management_v2', req.user)) {
      return await newSubscriptionService.createSubscription(data);
    } else {
      return await legacySubscriptionService.createSubscription(data);
    }
  `
};
```

### Performance Monitoring Integration
```javascript
const performanceMonitoring = {
  // Automatic instrumentation during implementation
  metrics: [
    "API response times",
    "Database query performance",
    "Memory usage patterns",
    "Error rates and types"
  ],

  // Built-in performance optimization
  optimizations: [
    "Database query optimization",
    "Caching strategy implementation",
    "Async processing for heavy operations",
    "Resource pooling and management"
  ]
};
```

### Rollback and Recovery
```javascript
const rollbackStrategy = {
  // Automatic rollback trigger conditions
  triggers: [
    {
      condition: "Error rate > 5%",
      action: "Automatic rollback",
      notification: "Immediate alert"
    },
    {
      condition: "Response time > 500ms",
      action: "Investigation alert",
      escalation: "Manual rollback after 10min"
    }
  ],

  // Rollback implementation
  mechanisms: [
    "Feature flag toggle",
    "Database migration rollback",
    "Service version rollback",
    "Configuration revert"
  ]
};
```

## Integration with Development Workflow

### CI/CD Integration
```yaml
# .github/workflows/implement-feature.yml
name: Feature Implementation Pipeline
on:
  workflow_dispatch:
    inputs:
      plan_id:
        description: 'Feature plan ID'
        required: true
      phase:
        description: 'Implementation phase'
        required: true

jobs:
  implement:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout code
        uses: actions/checkout@v3

      - name: Run Plan Mode Implementation
        run: |
          /implement-feature ${{ github.event.inputs.plan_id }} ${{ github.event.inputs.phase }} development full-stack --safe-mode

      - name: Run tests
        run: npm test

      - name: Quality checks
        run: |
          npm run lint
          npm run audit

      - name: Deploy to staging
        if: github.event.inputs.phase == 'integration'
        run: |
          /deploy-feature ${{ github.event.inputs.plan_id }} staging --validate
```

### Real-time Collaboration
```bash
# Pair programming session
/implement-feature "auth-improvement" core development backend --pair-programming --code-review

# Team collaboration with review checkpoints
/implement-feature "dashboard-redesign" integration development frontend --incremental --code-review --documentation
```
