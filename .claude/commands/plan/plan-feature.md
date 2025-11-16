# `/plan-feature` Command - Detailed Implementation

## Command Structure
```bash
/plan-feature [feature-name] [complexity] [priority] [timeline] [scope] [--options]
```

## Parameters

### Required Parameters
- **feature-name**: Descriptive name for the feature (quoted if contains spaces)
- **complexity**: `trivial` | `low` | `medium` | `high` | `epic` | `unknown`
- **priority**: `critical` | `urgent` | `high` | `medium` | `low` | `backlog`

### Optional Parameters
- **timeline**: `hours` | `days` | `weeks` | `months` | `quarters` | `flexible`
- **scope**: `component` | `module` | `system` | `integration` | `full-stack`

### Options
- `--stakeholders`: Define stakeholder groups (`--stakeholders="dev,qa,product"`)
- `--dependencies`: Specify known dependencies (`--dependencies="auth,payment"`)
- `--constraints`: Technical constraints (`--constraints="mobile-first,offline-capable"`)
- `--compliance`: Regulatory requirements (`--compliance="gdpr,sox,hipaa"`)
- `--budget`: Resource constraints (`--budget="2-devs,1-month"`)
- `--research`: Include research phase (`--research="competitive,technical,user"`)
- `--prototype`: Generate prototype specifications (`--prototype="ui,api,data"`)
- `--rollback`: Define rollback strategy (`--rollback="feature-flags,database"`)

## Implementation Flow

### Phase 1: Gemini Context Analysis & Discovery

```javascript
const geminiFeatureAnalysis = {
  // 1. Contextual Understanding
  contextAnalysis: {
    // Analyze current codebase for integration points
    findIntegrationPoints: async (featureName, codebase) => {
      return {
        // Identify where feature will touch existing code
        touchPoints: [
          {
            module: "auth",
            reason: "User permissions for feature",
            complexity: "medium",
            riskLevel: "low"
          },
          {
            module: "database",
            reason: "New tables and relationships",
            complexity: "high", 
            riskLevel: "medium"
          }
        ],
        
        // Find similar existing features for pattern reference
        similarFeatures: [
          {
            name: "user-profile",
            similarity: 0.8,
            patterns: ["CRUD operations", "validation", "caching"],
            reusableComponents: ["FormBuilder", "ValidationService"]
          }
        ],
        
        // Identify architectural constraints
        constraints: {
          performance: ["Must handle 10k concurrent users"],
          security: ["PII data handling", "encryption at rest"],
          scalability: ["Horizontal scaling required"],
          compatibility: ["Must work with legacy auth system"]
        }
      };
    },

    // Analyze business domain and requirements
    domainAnalysis: async (featureName, businessContext) => {
      return {
        // Map business concepts to technical implementation
        businessConcepts: [
          {
            concept: "User Subscription",
            technicalMapping: "User.subscription relationship",
            complexity: "medium",
            businessRules: ["Trial period logic", "Billing cycles"]
          }
        ],
        
        // Identify stakeholder concerns
        stakeholderRequirements: {
          business: ["Revenue tracking", "Analytics integration"],
          users: ["Simple onboarding", "Clear pricing"],
          compliance: ["Data retention policies", "Audit trails"],
          support: ["Admin tools", "Debug capabilities"]
        },
        
        // Risk assessment
        businessRisks: [
          {
            risk: "User adoption",
            mitigation: "A/B testing framework",
            priority: "high"
          }
        ]
      };
    }
  },

  // 2. Technical Feasibility Analysis
  technicalAnalysis: {
    // Assess current system capabilities
    capabilityGapAnalysis: async (requirements, currentSystem) => {
      return {
        existingCapabilities: [
          {
            capability: "User Authentication",
            readiness: "ready",
            reusability: "high"
          },
          {
            capability: "Payment Processing", 
            readiness: "needs-extension",
            reusability: "medium",
            requiredChanges: ["Subscription billing", "Dunning management"]
          }
        ],
        
        missingCapabilities: [
          {
            capability: "Subscription Management",
            complexity: "high",
            estimatedEffort: "3-4 weeks",
            dependencies: ["Payment system", "Notification service"]
          }
        ],
        
        infrastructureNeeds: {
          newServices: ["Subscription service", "Billing service"],
          databaseChanges: ["Subscription tables", "Billing history"],
          thirdPartyIntegrations: ["Stripe webhooks", "Analytics APIs"]
        }
      };
    },

    // Identify potential technical challenges
    challengeIdentification: async (featureSpec, systemContext) => {
      return {
        dataFlowChallenges: [
          {
            challenge: "Cross-service data consistency",
            impact: "high",
            solutions: ["Event sourcing", "Saga pattern", "2PC"]
          }
        ],
        
        performanceChallenges: [
          {
            challenge: "Real-time billing calculations",
            impact: "medium", 
            solutions: ["Background jobs", "Caching", "Precomputation"]
          }
        ],
        
        securityChallenges: [
          {
            challenge: "PCI compliance for billing data",
            impact: "critical",
            solutions: ["Tokenization", "Vault integration", "Audit logging"]
          }
        ]
      };
    }
  },

  // 3. Pattern and Architecture Recommendation
  architecturalPlanning: {
    // Recommend implementation patterns based on codebase analysis
    patternRecommendation: async (featureRequirements, existingPatterns) => {
      return {
        recommendedPatterns: [
          {
            pattern: "Domain-Driven Design",
            reason: "Complex business logic with subscription rules",
            implementation: "Subscription aggregate with value objects",
            confidence: "high"
          },
          {
            pattern: "Event-Driven Architecture", 
            reason: "Multiple services need subscription updates",
            implementation: "Domain events for subscription state changes",
            confidence: "medium"
          }
        ],
        
        antiPatterns: [
          {
            pattern: "Distributed transactions",
            reason: "Complex error handling and poor performance",
            alternative: "Eventually consistent with compensation"
          }
        ]
      };
    },

    // Design component architecture
    componentDesign: async (featureScope, integrationPoints) => {
      return {
        coreComponents: [
          {
            name: "SubscriptionService",
            responsibility: "Business logic for subscription lifecycle",
            interfaces: ["ISubscriptionManager", "ISubscriptionQuery"],
            dependencies: ["PaymentService", "NotificationService"],
            testingStrategy: "Unit + Integration tests"
          }
        ],
        
        dataModel: {
          entities: [
            {
              name: "Subscription",
              attributes: ["id", "userId", "planId", "status", "billingCycle"],
              relationships: ["belongsTo User", "hasMany BillingEvents"],
              constraints: ["Unique per user-plan", "Status transitions"]
            }
          ],
          
          valueObjects: [
            {
              name: "BillingCycle", 
              attributes: ["interval", "intervalCount", "startDate"],
              validations: ["Valid interval types", "Positive count"]
            }
          ]
        },
        
        apiDesign: {
          endpoints: [
            {
              path: "/subscriptions",
              method: "POST",
              purpose: "Create subscription",
              security: ["Authentication", "Authorization"],
              validation: ["Plan exists", "User eligible"]
            }
          ],
          
          events: [
            {
              name: "SubscriptionCreated",
              payload: ["subscriptionId", "userId", "planId", "timestamp"],
              consumers: ["BillingService", "AnalyticsService", "NotificationService"]
            }
          ]
        }
      };
    }
  }
};
```

### Phase 2: Claude Strategic Planning & Execution

```javascript
const claudeFeaturePlanning = {
  // 1. Comprehensive Implementation Strategy
  strategicPlanning: {
    // Create detailed implementation roadmap
    createRoadmap: (geminiAnalysis, requirements) => {
      return {
        phases: [
          {
            phase: "Foundation",
            duration: "1-2 weeks",
            description: "Core infrastructure and data models",
            deliverables: [
              "Database schema design and migration",
              "Core domain models and value objects", 
              "Basic service interfaces",
              "Authentication/authorization framework"
            ],
            dependencies: [],
            risks: ["Schema changes might require downtime"],
            successCriteria: ["All models pass validation", "Migration completes successfully"]
          },
          
          {
            phase: "Core Implementation", 
            duration: "2-3 weeks",
            description: "Business logic and primary workflows",
            deliverables: [
              "Subscription service implementation",
              "Payment integration and webhooks",
              "Core API endpoints",
              "Basic admin interface"
            ],
            dependencies: ["Foundation phase complete"],
            risks: ["Third-party API integration complexity"],
            successCriteria: ["End-to-end subscription flow works", "Payment webhooks processed correctly"]
          },
          
          {
            phase: "Integration & Polish",
            duration: "1 week", 
            description: "UI/UX, edge cases, and system integration",
            deliverables: [
              "User-facing subscription management UI",
              "Error handling and edge cases",
              "Monitoring and alerting",
              "Documentation and runbooks"
            ],
            dependencies: ["Core Implementation complete"],
            risks: ["UI/UX feedback requires significant changes"],
            successCriteria: ["User acceptance testing passes", "Performance benchmarks met"]
          }
        ],
        
        milestones: [
          {
            name: "MVP Ready",
            date: "Week 4",
            criteria: ["Basic subscription creation/cancellation", "Payment processing", "Admin oversight"]
          },
          {
            name: "Production Ready",
            date: "Week 6", 
            criteria: ["Full feature set", "Performance optimized", "Monitoring in place"]
          }
        ]
      };
    },

    // Define technical specifications
    technicalSpecification: (architecturalDesign, requirements) => {
      return {
        systemArchitecture: {
          services: [
            {
              name: "subscription-service",
              type: "microservice",
              technology: "Node.js/Express",
              database: "PostgreSQL",
              caching: "Redis",
              messaging: "RabbitMQ",
              monitoring: "Prometheus/Grafana"
            }
          ],
          
          dataFlow: [
            {
              trigger: "User clicks 'Subscribe'",
              flow: [
                "UI validates plan selection",
                "API creates pending subscription", 
                "Payment service processes payment",
                "Webhook confirms payment success",
                "Subscription activated",
                "Welcome email sent",
                "Analytics event recorded"
              ],
              errorHandling: "Compensation pattern with user notification"
            }
          ],
          
          securityModel: {
            authentication: "JWT with refresh tokens",
            authorization: "RBAC with subscription-based permissions",
            dataProtection: "AES-256 encryption for PII",
            auditLogging: "All subscription changes logged"
          }
        },
        
        performanceRequirements: {
          responseTime: "API endpoints < 200ms P95",
          throughput: "1000 concurrent subscription operations",
          availability: "99.9% uptime SLA",
          scalability: "Horizontal scaling with load balancing"
        },
        
        qualityStandards: {
          testCoverage: "90% code coverage minimum",
          codeQuality: "ESLint + Prettier, SonarQube gates",
          documentation: "API docs, architecture diagrams, runbooks",
          monitoring: "Error rates, performance metrics, business KPIs"
        }
      };
    }
  },

  // 2. Implementation Planning
  implementationPlanning: {
    // Create detailed task breakdown
    createTaskBreakdown: (roadmap, technicalSpec) => {
      return {
        epics: [
          {
            epic: "Subscription Management Core",
            story_points: 34,
            stories: [
              {
                title: "As a user, I can view available subscription plans",
                points: 3,
                acceptance_criteria: [
                  "Plans display with pricing and features",
                  "Current subscription highlighted",
                  "Upgrade/downgrade options shown"
                ],
                technical_tasks: [
                  "Create Plan model and repository",
                  "Implement plan comparison API",
                  "Build responsive plan selection UI",
                  "Add analytics tracking"
                ],
                dependencies: ["User authentication"],
                definition_of_done: [
                  "Code reviewed and merged",
                  "Unit tests passing (90% coverage)",
                  "Integration tests verify full flow",
                  "UI/UX approved by design team",
                  "Performance tested under load"
                ]
              }
            ]
          }
        ],
        
        technical_debt: [
          {
            item: "Legacy payment integration cleanup",
            effort: "2 days",
            impact: "Reduced maintenance burden",
            timing: "After core implementation"
          }
        ],
        
        research_spikes: [
          {
            spike: "Payment provider comparison",
            effort: "3 days", 
            outcome: "Recommendation for Stripe vs PayPal integration",
            criteria: ["Cost", "Features", "Developer experience", "Compliance"]
          }
        ]
      };
    },

    // Resource allocation and timeline
    resourcePlanning: (taskBreakdown, constraints) => {
      return {
        teamComposition: {
          roles: [
            {
              role: "Senior Full-stack Developer",
              allocation: "100%",
              responsibilities: ["Architecture decisions", "Core implementation", "Code reviews"],
              duration: "6 weeks"
            },
            {
              role: "Frontend Developer", 
              allocation: "60%",
              responsibilities: ["UI components", "User experience", "Responsive design"],
              duration: "4 weeks"
            },
            {
              role: "DevOps Engineer",
              allocation: "20%", 
              responsibilities: ["Infrastructure", "CI/CD", "Monitoring"],
              duration: "6 weeks"
            }
          ]
        },
        
        timeline: {
          start_date: "2025-07-15",
          end_date: "2025-08-26",
          buffer: "1 week for unforeseen issues",
          checkpoints: [
            {
              date: "2025-07-22",
              deliverable: "Database schema and core models",
              review_criteria: ["Schema review", "Model validation"]
            },
            {
              date: "2025-08-05",
              deliverable: "MVP functionality complete",
              review_criteria: ["Feature demo", "Performance testing"]
            }
          ]
        },
        
        dependencies: [
          {
            dependency: "Payment provider API access",
            owner: "Business team",
            deadline: "2025-07-18",
            impact_if_delayed: "2-3 day implementation delay"
          }
        ]
      };
    }
  },

  // 3. Risk Management and Contingency Planning
  riskManagement: {
    // Identify and mitigate risks
    riskAssessment: (planData, historicalData) => {
      return {
        technical_risks: [
          {
            risk: "Third-party payment API changes",
            probability: "medium",
            impact: "high", 
            mitigation: [
              "Use stable API versions",
              "Implement adapter pattern for easy switching",
              "Monitor API deprecation notices"
            ],
            contingency: "Fallback to secondary payment provider"
          },
          {
            risk: "Database performance under load",
            probability: "low",
            impact: "high",
            mitigation: [
              "Performance testing in staging",
              "Database indexing optimization", 
              "Connection pooling"
            ],
            contingency: "Implement caching layer and read replicas"
          }
        ],
        
        business_risks: [
          {
            risk: "User adoption lower than expected",
            probability: "medium",
            impact: "medium",
            mitigation: [
              "A/B testing for UI variations",
              "User feedback collection",
              "Gradual rollout strategy"
            ],
            contingency: "Feature flag for quick rollback"
          }
        ],
        
        timeline_risks: [
          {
            risk: "Key developer unavailability",
            probability: "low",
            impact: "high",
            mitigation: [
              "Cross-training team members",
              "Comprehensive documentation",
              "Pair programming sessions"
            ],
            contingency: "Contractor backup identified"
          }
        ]
      };
    },

    // Define rollback and recovery strategies
    contingencyPlanning: (riskAssessment, systemDesign) => {
      return {
        rollback_strategy: {
          feature_flags: [
            {
              flag: "subscription_management_enabled",
              scope: "Global toggle for entire feature",
              rollback_time: "< 5 minutes"
            },
            {
              flag: "new_billing_flow_enabled", 
              scope: "Toggle for new vs legacy billing",
              rollback_time: "< 2 minutes"
            }
          ],
          
          database_rollback: {
            strategy: "Blue-green deployment with snapshot",
            rollback_time: "< 30 minutes",
            data_loss_risk: "Minimal (transactional)",
            testing: "Automated rollback testing in staging"
          },
          
          monitoring_triggers: [
            {
              metric: "Error rate > 5%",
              action: "Automatic rollback trigger",
              notification: "Immediate alert to on-call"
            },
            {
              metric: "Response time > 1s",
              action: "Investigation alert",
              escalation: "Manual rollback if not resolved in 10 min"
            }
          ]
        },
        
        disaster_recovery: {
          data_backup: "Continuous replication to secondary region",
          rto: "15 minutes", // Recovery Time Objective
          rpo: "1 minute",   // Recovery Point Objective
          testing_schedule: "Monthly disaster recovery drills"
        }
      };
    }
  }
};
```

## Usage Examples

### Basic Feature Planning
```bash
/plan-feature "user-subscription-management" medium high weeks full-stack
```

**Expected Output:**
```markdown
# Feature Plan: User Subscription Management

## Executive Summary
- **Complexity**: Medium (estimated 4-6 weeks)
- **Priority**: High (business-critical revenue feature)
- **Scope**: Full-stack implementation
- **Team**: 2.5 FTE for 6 weeks

## Implementation Strategy

### Phase 1: Foundation (Weeks 1-2)
**Goal**: Establish core infrastructure and data models

**Deliverables**:
- Database schema for subscriptions, plans, and billing
- Core domain models (Subscription, Plan, BillingCycle)
- Service interfaces and dependency injection setup
- Basic authentication/authorization framework

**Key Tasks**:
1. Design and implement subscription database schema
2. Create domain models with proper validation
3. Set up service layer architecture
4. Implement basic CRUD operations

**Success Criteria**:
- ✅ All database migrations run successfully
- ✅ Core models pass validation tests
- ✅ Service layer integration tests pass

### Phase 2: Core Implementation (Weeks 3-4)
**Goal**: Implement business logic and payment integration

**Deliverables**:
- Subscription lifecycle management
- Payment provider integration (Stripe)
- Webhook handling for payment events
- Basic admin interface for subscription management

**Key Tasks**:
1. Implement subscription state machine
2. Integrate with Stripe API for payments
3. Build webhook processing system
4. Create admin dashboard components

**Success Criteria**:
- ✅ End-to-end subscription flow functional
- ✅ Payment webhooks processed correctly
- ✅ Admin can manage user subscriptions

### Phase 3: Integration & Polish (Weeks 5-6)
**Goal**: Complete user experience and production readiness

**Deliverables**:
- User-facing subscription management UI
- Error handling and edge case coverage
- Monitoring, alerting, and analytics
- Documentation and deployment procedures

**Key Tasks**:
1. Build responsive user interface
2. Implement comprehensive error handling
3. Set up monitoring and alerting
4. Create user and developer documentation

**Success Criteria**:
- ✅ User acceptance testing passes
- ✅ Performance benchmarks met
- ✅ Production monitoring operational

## Technical Architecture

### Services & Components
```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Frontend UI   │────│ Subscription    │────│  Payment        │
│                 │    │ Service         │    │  Service        │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                              │                        │
                              │                        │
                       ┌─────────────────┐    ┌─────────────────┐
                       │   Database      │    │  Stripe API     │
                       │   (PostgreSQL)  │    │                 │
                       └─────────────────┘    └─────────────────┘
```

### Data Model
- **Subscription**: Core entity managing user subscription lifecycle
- **Plan**: Subscription tiers with pricing and feature definitions
- **BillingEvent**: Audit trail of all billing-related activities
- **PaymentMethod**: User payment information (tokenized)

### API Endpoints
- `POST /api/subscriptions` - Create new subscription
- `GET /api/subscriptions/me` - Get current user's subscription
- `PUT /api/subscriptions/{id}` - Update subscription (upgrade/downgrade)
- `DELETE /api/subscriptions/{id}` - Cancel subscription

## Risk Assessment & Mitigation

### Technical Risks
- **🔴 High**: Payment provider API changes
  - *Mitigation*: Use adapter pattern for easy provider switching
- **🟡 Medium**: Database performance under load
  - *Mitigation*: Implement caching and connection pooling

### Business Risks  
- **🟡 Medium**: Lower than expected user adoption
  - *Mitigation*: A/B testing and gradual rollout

### Timeline Risks
- **🟡 Medium**: Key developer unavailability
  - *Mitigation*: Cross-training and comprehensive documentation

## Success Metrics
- **Technical**: 99.9% uptime, <200ms API response time
- **Business**: 15% user conversion to paid subscriptions
- **Quality**: 90% test coverage, zero critical security issues

## Next Steps
1. **Immediate**: Stakeholder review and approval
2. **Week 1**: Kick-off meeting and environment setup  
3. **Ongoing**: Weekly progress reviews and risk assessment
```

### Advanced Feature Planning with Research
```bash
/plan-feature "ai-powered-content-recommendations" epic critical months system --research="competitive,technical,user" --prototype="ml-model,api,ui"
```

**Expected Output:**
```markdown
# Epic Feature Plan: AI-Powered Content Recommendations

## Research Phase (Weeks 1-2)

### Competitive Analysis
**Goal**: Understand market landscape and best practices

**Research Areas**:
- Netflix recommendation algorithm approach
- Spotify's music discovery features  
- Medium's content curation system
- Amazon's collaborative filtering

**Deliverables**:
- Competitive feature matrix
- Best practice recommendations
- Differentiation opportunities
- Performance benchmarks

### Technical Research  
**Goal**: Evaluate ML approaches and infrastructure needs

**Research Areas**:
- Collaborative filtering vs content-based filtering
- Deep learning models for recommendations
- Real-time vs batch processing tradeoffs
- Scalability considerations for 100K+ users

**Deliverables**:
- Technical architecture options
- ML model evaluation criteria
- Infrastructure cost analysis
- Performance requirements specification

### User Research
**Goal**: Validate assumptions and gather requirements

**Research Methods**:
- User interviews (20 participants)
- Surveys to existing user base
- A/B testing framework design
- Analytics review of current content consumption

**Deliverables**:
- User persona updates
- Feature prioritization matrix
- Success metrics definition
- UX wireframes and prototypes

## Prototype Phase (Weeks 3-4)

### ML Model Prototype
- Simple collaborative filtering implementation
- Content-based filtering using existing metadata  
- Hybrid approach combining both methods
- Performance evaluation on historical data

### API Prototype
- RESTful recommendation endpoints
- Real-time recommendation serving
- Batch recommendation generation
- A/B testing integration

### UI Prototype
- Recommendation display components
- User feedback collection interface
- Personalization settings panel
- Mobile-responsive design

## Implementation Strategy
[Detailed 6-month implementation plan...]
```

### Feature Planning with Compliance Requirements
```bash
/plan-feature "healthcare-patient-portal" high urgent months full-stack --compliance="hipaa,gdpr" --constraints="mobile-first,offline-capable"
```

**Expected Output includes**:
- HIPAA compliance checklist and implementation strategy
- GDPR data handling and consent management
- Mobile-first responsive design approach
- Offline capability with data synchronization
- Security architecture for healthcare data
- Audit trail and compliance reporting

## Integration with Development Workflow

### Pre-Planning Validation
```bash
# Validate feature feasibility before planning
/analyze-codebase full dependencies comprehensive json | \
/plan-feature "new-feature" medium high --validate-against-analysis
```

### Planning Template Generation
```bash
# Generate planning templates for common feature types
/plan-feature template crud-feature medium --generate-template
/plan-feature template integration-feature high --generate-template  
/plan-feature template ui-component low --generate-template
```

### Stakeholder Communication
```bash
# Generate stakeholder-specific summaries
/plan-feature "user-dashboard" medium high --audience="executive" --format="summary"
/plan-feature "user-dashboard" medium high --audience="technical" --format="detailed"
/plan-feature "user-dashboard" medium high --audience="product" --format="roadmap"
```

## Error Handling and Validation

### Input Validation
```javascript
const validatePlanFeatureCommand = (params) => {
  // Validate feature name
  if (!params.featureName || params.featureName.length < 3) {
    throw new Error('Feature name must be at least 3 characters long');
  }
  
  // Validate complexity vs timeline consistency
  const complexityTimelineMatrix = {
    trivial: ['hours', 'days'],
    low: ['days', 'weeks'], 
    medium: ['weeks', 'months'],
    high: ['months', 'quarters'],
    epic: ['quarters']
  };
  
  if (!complexityTimelineMatrix[params.complexity]?.includes(params.timeline)) {
    console.warn(`Timeline ${params.timeline} unusual for ${params.complexity} complexity`);
  }
  
  // Validate constraint combinations
  if (params.constraints?.includes('offline-capable') && !params.constraints?.includes('mobile-first')) {
    console.warn('Offline capability typically requires mobile-first approach');
  }
};
```

### Adaptive Planning
```javascript
const adaptivePlanning = async (baseParams, constraints) => {
  // Adjust plan based on team capacity
  const teamCapacity = await getTeamCapacity();
  if (teamCapacity.availableDevHours < estimatedEffort * 0.8) {
    return {
      ...basePlan,
      recommendation: 'Consider reducing scope or extending timeline',
      alternatives: [
        'Implement as MVP with reduced features',
        'Extend timeline by 2 weeks',
        'Add additional team member'
      ]
    };
  }
  
  // Adjust for technical constraints
  const technicalConstraints = await analyzeTechnicalConstraints();
  if (technicalConstraints.blockers.length > 0) {
    return {
      ...basePlan,
      blockers: technicalConstraints.blockers,
      prerequisiteWork: technicalConstraints.prerequisites
    };
  }
};
```
