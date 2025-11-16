# `/implement-fix` Command - Detailed Implementation

## Command Structure
```bash
/implement-fix [issue-id] [approach] [scope] [priority] [--options]
```

## Parameters

### Required Parameters
- **issue-id**: Reference to bug report, ticket, or issue (`BUG-123`, `github-456`, `sentry-789`, `user-report-101`)
- **approach**: `minimal` | `comprehensive` | `preventive` | `hotfix` | `root-cause` | `temporary`
- **scope**: `local` | `module` | `system` | `integration` | `global` | `downstream`

### Optional Parameters
- **priority**: `critical` | `urgent` | `high` | `medium` | `low` | `maintenance`

### Options
- `--test-driven`: Write failing tests first, then implement fix
- `--regression-safe`: Include extensive regression testing
- `--backwards-compatible`: Ensure no breaking changes
- `--performance-impact`: Analyze and minimize performance impact
- `--security-review`: Include security impact assessment
- `--rollback-plan`: Prepare comprehensive rollback strategy
- `--monitoring`: Add monitoring and alerting for the fix
- `--documentation`: Document the fix and prevention measures
- `--upstream-fix`: Address root cause in upstream dependencies
- `--batch-fix`: Apply similar fixes across codebase
- `--emergency-deploy`: Fast-track deployment for critical issues

### Convention Enforcement (References: /name-conventions)
- Clean names: !`echo "Applying clean, business-focused naming"`
- No jargon: !`echo "Avoiding technical prefixes and versioning"`
- Business focus: !`echo "Using domain-appropriate terminology"`
- Validation: !`echo "Checking names follow conventions"`

## Implementation Flow

### Phase 1: Gemini Issue Analysis & Root Cause Investigation

```javascript
const geminiIssueAnalysis = {
  // 1. Comprehensive Issue Understanding
  issueInvestigation: {
    // Deep analysis of the reported issue
    analyzeIssueReport: async (issueId, issueData) => {
      return {
        // Issue classification and severity assessment
        issueClassification: {
          type: "NullPointerException",
          category: "Runtime Error",
          severity: "High",
          impact: "User-facing functionality broken",
          affected_users: "15% of active users",
          business_impact: "Revenue loss estimated at $2K/day",
          reproducibility: "Consistent (100%)",
          environment: ["production", "staging"]
        },

        // Detailed issue analysis
        issueDetails: {
          error_message: "Cannot read property 'id' of null at UserService.js:45",
          stack_trace: [
            "at UserService.updateProfile (src/services/UserService.js:45:12)",
            "at UserController.updateUser (src/controllers/UserController.js:28:18)",
            "at middleware (src/middleware/auth.js:15:5)"
          ],
          triggered_by: "User profile update with missing avatar field",
          preconditions: [
            "User must be logged in",
            "Profile form submitted without avatar selection",
            "Legacy users with null avatar_url field"
          ],
          frequency: "287 occurrences in last 24 hours",
          first_seen: "2025-07-10T14:30:00Z",
          last_seen: "2025-07-12T09:15:00Z"
        },

        // User impact analysis
        userImpact: {
          affected_functionality: ["Profile updates", "Avatar changes", "Settings save"],
          user_scenarios: [
            "New user trying to complete profile",
            "Existing user updating contact info",
            "Legacy user from pre-avatar feature"
          ],
          workarounds: [
            "User can skip avatar selection",
            "Admin can update profile directly"
          ],
          user_sentiment: "High frustration - multiple support tickets"
        }
      };
    },

    // Root cause analysis using codebase context
    rootCauseAnalysis: async (issueDetails, codebaseContext) => {
      return {
        // Primary root cause identification
        rootCause: {
          immediate_cause: "Null check missing for user.avatar_url property",
          underlying_cause: "Database migration added avatar_url column with NULL default for existing users",
          contributing_factors: [
            "No validation for required fields in profile update",
            "Frontend doesn't enforce avatar selection",
            "Legacy data migration incomplete"
          ],
          code_location: {
            file: "src/services/UserService.js",
            line: 45,
            function: "updateProfile",
            problematic_code: "const avatarId = user.avatar.id"
          }
        },

        // Historical context and pattern analysis
        historicalContext: {
          similar_issues: [
            {
              issue: "BUG-098 - Similar null pointer in payment service",
              resolution: "Added null checks and default values",
              date: "2025-06-15",
              prevention_measure: "Input validation middleware"
            }
          ],
          recent_changes: [
            {
              change: "Avatar feature rollout",
              date: "2025-07-08",
              author: "john.doe",
              pr: "#456",
              risk_factors: ["Schema changes", "Legacy data handling"]
            }
          ],
          code_quality_indicators: {
            test_coverage: "65% (below 90% target)",
            null_safety_patterns: "Inconsistent usage",
            error_handling: "Basic try-catch, no specific null handling"
          }
        },

        // Blast radius assessment
        blastRadius: {
          direct_impact: [
            "UserService.updateProfile method",
            "Profile update API endpoint",
            "User settings page"
          ],
          indirect_impact: [
            "User onboarding flow",
            "Admin user management",
            "Analytics user tracking"
          ],
          potential_cascading_effects: [
            "User session corruption",
            "Database inconsistency",
            "Cache invalidation issues"
          ],
          system_dependencies: [
            "Authentication service (depends on user data)",
            "Notification service (user profile changes)",
            "Analytics service (user behavior tracking)"
          ]
        }
      };
    },

    // Fix strategy assessment
    fixStrategyAnalysis: async (rootCause, approach, constraints) => {
      return {
        // Multiple fix approaches with trade-offs
        fixApproaches: [
          {
            approach: "minimal",
            description: "Add null check and default avatar handling",
            pros: [
              "Quick implementation (2-4 hours)",
              "Low risk of introducing new bugs",
              "Immediate problem resolution"
            ],
            cons: [
              "Doesn't address underlying data issue",
              "May need follow-up fixes",
              "Symptom fix rather than root cause"
            ],
            implementation: {
              changes: ["Null check in UserService.updateProfile"],
              testing_required: "Unit tests for null scenarios",
              deployment_risk: "Low"
            }
          },
          {
            approach: "comprehensive",
            description: "Fix data migration, add validation, improve error handling",
            pros: [
              "Addresses root cause completely",
              "Prevents similar issues in future",
              "Improves overall system robustness"
            ],
            cons: [
              "Longer implementation time (1-2 days)",
              "Higher risk of introducing changes",
              "Requires database migration"
            ],
            implementation: {
              changes: [
                "Data migration for legacy users",
                "Input validation middleware",
                "Comprehensive null safety",
                "Error handling improvements"
              ],
              testing_required: "Full regression test suite",
              deployment_risk: "Medium"
            }
          },
          {
            approach: "preventive",
            description: "Implement system-wide null safety and validation patterns",
            pros: [
              "Prevents entire class of similar bugs",
              "Improves code quality significantly",
              "Long-term technical debt reduction"
            ],
            cons: [
              "Extensive changes across codebase",
              "Requires significant testing",
              "Potential for breaking changes"
            ],
            implementation: {
              changes: [
                "TypeScript migration for null safety",
                "Schema validation middleware",
                "Standardized error handling",
                "Automated null safety linting"
              ],
              testing_required: "Complete system testing",
              deployment_risk: "High"
            }
          }
        ],

        // Risk assessment for each approach
        riskAssessment: {
          deployment_risks: [
            {
              risk: "Fix introduces new regression",
              probability: "Low-Medium",
              impact: "High",
              mitigation: "Comprehensive testing and gradual rollout"
            },
            {
              risk: "Database migration fails",
              probability: "Low",
              impact: "Critical",
              mitigation: "Test migration in staging, backup strategy"
            }
          ],
          business_risks: [
            {
              risk: "Extended downtime during fix deployment",
              probability: "Low",
              impact: "Medium",
              mitigation: "Blue-green deployment strategy"
            }
          ]
        }
      };
    }
  },

  // 2. Codebase Impact Analysis
  codebaseImpactAnalysis: {
    // Analyze what code will be affected by the fix
    identifyAffectedCode: async (rootCause, fixApproach) => {
      return {
        // Direct code changes required
        primaryChanges: [
          {
            file: "src/services/UserService.js",
            function: "updateProfile",
            change_type: "null_safety_check",
            current_code: `
              async updateProfile(userId, profileData) {
                const user = await User.findById(userId);
                const avatarId = user.avatar.id; // ❌ Potential null pointer
                return await user.update({ ...profileData, avatarId });
              }
            `,
            proposed_fix: `
              async updateProfile(userId, profileData) {
                const user = await User.findById(userId);
                const avatarId = user.avatar?.id || null; // ✅ Safe null handling
                return await user.update({ ...profileData, avatarId });
              }
            `,
            risk_level: "Low",
            testing_strategy: "Unit tests with null avatar scenarios"
          }
        ],

        // Related code that might need updates
        relatedChanges: [
          {
            file: "src/controllers/UserController.js",
            reason: "Input validation to prevent null avatar scenarios",
            change_type: "validation_enhancement",
            priority: "Medium"
          },
          {
            file: "src/models/User.js",
            reason: "Default value handling for avatar field",
            change_type: "model_enhancement",
            priority: "Low"
          }
        ],

        // Test files that need updates
        testUpdates: [
          {
            file: "tests/services/UserService.test.js",
            new_tests: [
              "updateProfile with null avatar",
              "updateProfile with undefined avatar",
              "updateProfile with legacy user data"
            ],
            updated_tests: [
              "Existing updateProfile tests to handle new logic"
            ]
          }
        ],

        // Migration scripts needed
        migrationRequirements: [
          {
            type: "data_migration",
            description: "Set default avatar for legacy users",
            sql: `
              UPDATE users
              SET avatar_url = '/images/default-avatar.png'
              WHERE avatar_url IS NULL AND created_at < '2025-07-08'
            `,
            rollback_sql: `
              UPDATE users
              SET avatar_url = NULL
              WHERE avatar_url = '/images/default-avatar.png'
              AND created_at < '2025-07-08'
            `
          }
        ]
      };
    },

    // Analyze testing requirements
    testingStrategy: async (affectedCode, fixApproach) => {
      return {
        // Test scenarios to cover
        testScenarios: [
          {
            scenario: "Null avatar handling",
            test_cases: [
              {
                name: "User with null avatar updates profile",
                setup: "Create user with avatar_url = null",
                action: "Call updateProfile with valid data",
                expected: "Profile updates successfully without error",
                assertion: "No null pointer exception thrown"
              },
              {
                name: "User with undefined avatar updates profile",
                setup: "Create user without avatar property",
                action: "Call updateProfile with valid data",
                expected: "Profile updates with default avatar handling",
                assertion: "avatarId set to null gracefully"
              }
            ]
          },
          {
            scenario: "Regression prevention",
            test_cases: [
              {
                name: "Existing users with valid avatars continue to work",
                setup: "User with proper avatar data",
                action: "Update profile normally",
                expected: "Existing functionality unchanged",
                assertion: "Avatar ID properly extracted and used"
              }
            ]
          }
        ],

        // Testing levels required
        testingLevels: {
          unit_tests: {
            required: true,
            focus: ["Null safety logic", "Edge case handling"],
            coverage_target: "100% for modified functions"
          },
          integration_tests: {
            required: true,
            focus: ["API endpoint behavior", "Database interactions"],
            scenarios: ["End-to-end profile update flow"]
          },
          regression_tests: {
            required: true,
            focus: ["Existing functionality preservation"],
            scope: "All user management features"
          },
          performance_tests: {
            required: false,
            reason: "Minimal performance impact expected"
          }
        }
      };
    }
  }
};
```

### Phase 2: Claude Precise Fix Implementation

```javascript
const claudeFixImplementation = {
  // 1. Fix Strategy Execution
  fixExecution: {
    // Implement the chosen fix approach
    implementFix: (geminiAnalysis, approach, options) => {
      return {
        // Implementation plan based on analysis
        implementationPlan: {
          priority_order: [
            {
              step: 1,
              task: "Implement immediate null safety fix",
              rationale: "Stop the bleeding - prevent further user impact",
              estimated_time: "1 hour",
              risk: "Low"
            },
            {
              step: 2,
              task: "Add comprehensive tests for edge cases",
              rationale: "Ensure fix works and prevent regression",
              estimated_time: "2 hours",
              risk: "Low"
            },
            {
              step: 3,
              task: "Data migration for legacy users",
              rationale: "Address root cause in existing data",
              estimated_time: "1 hour",
              risk: "Medium"
            },
            {
              step: 4,
              task: "Input validation enhancement",
              rationale: "Prevent similar issues in future",
              estimated_time: "2 hours",
              risk: "Low"
            }
          ],

          rollback_strategy: {
            immediate_rollback: "Revert code changes via git",
            data_rollback: "Execute rollback migration script",
            monitoring: "Watch error rates and user complaints",
            success_criteria: "Error rate returns to baseline within 1 hour"
          },

          deployment_plan: {
            environment_sequence: ["development", "staging", "production"],
            validation_gates: [
              "All tests pass",
              "Manual smoke testing",
              "Performance benchmarks maintained",
              "Security scan clean"
            ],
            rollout_strategy: "Blue-green deployment with 5% traffic test"
          }
        }
      };
    },

    // Generate the actual fix code
    generateFixCode: (analysisData, implementation_plan) => {
      return {
        // Primary fix implementation
        primaryFix: {
          file: "src/services/UserService.js",
          changes: [
            {
              type: "method_modification",
              method: "updateProfile",
              before: `
                async updateProfile(userId, profileData) {
                  try {
                    const user = await User.findById(userId);
                    if (!user) {
                      throw new NotFoundError('User not found');
                    }

                    // ❌ BUG: This will throw if user.avatar is null
                    const avatarId = user.avatar.id;

                    const updatedData = {
                      ...profileData,
                      avatarId,
                      updatedAt: new Date()
                    };

                    return await user.update(updatedData);
                  } catch (error) {
                    logger.error('Profile update failed:', error);
                    throw error;
                  }
                }
              `,
              after: `
                async updateProfile(userId, profileData) {
                  try {
                    const user = await User.findById(userId);
                    if (!user) {
                      throw new NotFoundError('User not found');
                    }

                    // ✅ FIX: Safe null handling with optional chaining and fallback
                    const avatarId = user.avatar?.id || this.getDefaultAvatarId();

                    // ✅ ENHANCEMENT: Validate profile data before update
                    const validatedData = await this.validateProfileData(profileData);

                    const updatedData = {
                      ...validatedData,
                      avatarId,
                      updatedAt: new Date()
                    };

                    // ✅ IMPROVEMENT: Use transaction for data consistency
                    const result = await User.sequelize.transaction(async (transaction) => {
                      return await user.update(updatedData, { transaction });
                    });

                    // ✅ MONITORING: Log successful updates for tracking
                    logger.info('Profile updated successfully', {
                      userId,
                      hasAvatar: !!user.avatar,
                      fields: Object.keys(validatedData)
                    });

                    return result;
                  } catch (error) {
                    // ✅ IMPROVEMENT: Better error context and classification
                    logger.error('Profile update failed:', {
                      userId,
                      error: error.message,
                      stack: error.stack,
                      profileData: this.sanitizeForLogging(profileData)
                    });

                    // Re-throw with user-friendly message for null avatar issues
                    if (error.message.includes('Cannot read property') && error.message.includes('avatar')) {
                      throw new ValidationError('Profile update failed due to avatar data issue. Please try again or contact support.');
                    }

                    throw error;
                  }
                }
              `,
              explanation: "Added null safety, input validation, transaction handling, and improved error messaging"
            }
          ],

          // Helper methods added to the class
          newMethods: [
            {
              method: "getDefaultAvatarId",
              implementation: `
                /**
                 * Get default avatar ID for users without avatars
                 * @returns {string|null} Default avatar ID or null
                 */
                getDefaultAvatarId() {
                  // Return null to indicate no avatar - frontend will show default
                  return null;
                }
              `
            },
            {
              method: "validateProfileData",
              implementation: `
                /**
                 * Validate profile data before update
                 * @param {Object} profileData - Profile data to validate
                 * @returns {Object} Validated and sanitized profile data
                 * @throws {ValidationError} If validation fails
                 */
                async validateProfileData(profileData) {
                  const allowedFields = ['firstName', 'lastName', 'email', 'phone', 'bio'];
                  const validatedData = {};

                  for (const [key, value] of Object.entries(profileData)) {
                    if (!allowedFields.includes(key)) {
                      continue; // Skip unknown fields
                    }

                    // Basic validation
                    if (value !== null && value !== undefined) {
                      if (typeof value === 'string' && value.trim().length > 0) {
                        validatedData[key] = value.trim();
                      } else if (typeof value !== 'string') {
                        validatedData[key] = value;
                      }
                    }
                  }

                  // Email validation if provided
                  if (validatedData.email && !this.isValidEmail(validatedData.email)) {
                    throw new ValidationError('Invalid email format');
                  }

                  return validatedData;
                }
              `
            },
            {
              method: "sanitizeForLogging",
              implementation: `
                /**
                 * Remove sensitive data from profile data for safe logging
                 * @param {Object} data - Data to sanitize
                 * @returns {Object} Sanitized data safe for logging
                 */
                sanitizeForLogging(data) {
                  const sanitized = { ...data };
                  // Remove or mask sensitive fields
                  if (sanitized.email) {
                    sanitized.email = sanitized.email.replace(/(.{2}).*(@.*)/, '$1***$2');
                  }
                  if (sanitized.phone) {
                    sanitized.phone = '***-***-' + sanitized.phone.slice(-4);
                  }
                  return sanitized;
                }
              `
            }
          ]
        },

        // Data migration script
        dataMigration: {
          file: "migrations/20250712_fix_null_avatars.js",
          content: `
            'use strict';

            module.exports = {
              up: async (queryInterface, Sequelize) => {
                console.log('Starting avatar null fix migration...');

                // Find users with null avatars created before avatar feature
                const usersWithNullAvatars = await queryInterface.sequelize.query(
                  \`SELECT id, email, created_at FROM users
                   WHERE avatar_url IS NULL
                   AND created_at < '2025-07-08'\`,
                  { type: Sequelize.QueryTypes.SELECT }
                );

                console.log(\`Found \${usersWithNullAvatars.length} users with null avatars\`);

                if (usersWithNullAvatars.length === 0) {
                  console.log('No users to update, migration complete');
                  return;
                }

                // Update users with default avatar URL
                const updateResult = await queryInterface.sequelize.query(
                  \`UPDATE users
                   SET avatar_url = '/images/default-avatar.png',
                       updated_at = NOW()
                   WHERE avatar_url IS NULL
                   AND created_at < '2025-07-08'\`
                );

                console.log(\`Updated \${updateResult[1]} users with default avatars\`);

                // Add index for avatar_url if it doesn't exist
                try {
                  await queryInterface.addIndex('users', ['avatar_url'], {
                    name: 'idx_users_avatar_url',
                    where: {
                      avatar_url: {
                        [Sequelize.Op.ne]: null
                      }
                    }
                  });
                  console.log('Added avatar_url index');
                } catch (error) {
                  if (!error.message.includes('already exists')) {
                    throw error;
                  }
                  console.log('Avatar_url index already exists');
                }
              },

              down: async (queryInterface, Sequelize) => {
                console.log('Rolling back avatar null fix migration...');

                // Revert default avatars for legacy users
                const revertResult = await queryInterface.sequelize.query(
                  \`UPDATE users
                   SET avatar_url = NULL,
                       updated_at = NOW()
                   WHERE avatar_url = '/images/default-avatar.png'
                   AND created_at < '2025-07-08'\`
                );

                console.log(\`Reverted \${revertResult[1]} users to null avatars\`);

                // Remove index
                try {
                  await queryInterface.removeIndex('users', 'idx_users_avatar_url');
                  console.log('Removed avatar_url index');
                } catch (error) {
                  console.log('Avatar_url index removal failed or already removed');
                }
              }
            };
          `
        },

        // Comprehensive test suite
        testSuite: {
          file: "tests/services/UserService.avatar-fix.test.js",
          content: `
            const UserService = require('../../src/services/UserService');
            const User = require('../../src/models/User');
            const { ValidationError, NotFoundError } = require('../../src/errors');
            const { setupTestDB, cleanupTestDB, createTestUser } = require('../helpers');

            describe('UserService - Avatar Null Fix', () => {
              let userService;

              beforeAll(async () => {
                await setupTestDB();
                userService = new UserService();
              });

              afterAll(async () => {
                await cleanupTestDB();
              });

              beforeEach(async () => {
                await User.destroy({ where: {}, force: true });
              });

              describe('Null Avatar Handling', () => {
                test('should handle user with null avatar gracefully', async () => {
                  // Arrange: Create user with null avatar (legacy user scenario)
                  const user = await createTestUser({
                    email: 'legacy@example.com',
                    avatar_url: null,
                    avatar: null
                  });

                  const profileData = {
                    firstName: 'John',
                    lastName: 'Doe'
                  };

                  // Act: Update profile (this used to throw null pointer error)
                  const result = await userService.updateProfile(user.id, profileData);

                  // Assert: Profile updated successfully without errors
                  expect(result.firstName).toBe('John');
                  expect(result.lastName).toBe('Doe');
                  expect(result.avatarId).toBeNull(); // Should be null, not throw error
                });

                test('should handle user with undefined avatar property', async () => {
                  // Arrange: Create user without avatar property at all
                  const user = await User.create({
                    email: 'test@example.com',
                    password: 'password123'
                    // No avatar property set
                  });

                  const profileData = { firstName: 'Jane' };

                  // Act & Assert: Should not throw null pointer error
                  await expect(userService.updateProfile(user.id, profileData))
                    .resolves.toMatchObject({ firstName: 'Jane' });
                });

                test('should preserve existing avatar for users who have one', async () => {
                  // Arrange: User with valid avatar
                  const user = await createTestUser({
                    email: 'hasavatar@example.com',
                    avatar: { id: 'avatar-123' }
                  });

                  const profileData = { bio: 'Updated bio' };

                  // Act
                  const result = await userService.updateProfile(user.id, profileData);

                  // Assert: Avatar should be preserved
                  expect(result.bio).toBe('Updated bio');
                  // Avatar handling logic should work correctly for existing avatars
                });
              });

              describe('Input Validation', () => {
                test('should validate email format in profile data', async () => {
                  const user = await createTestUser();

                  const invalidProfileData = {
                    email: 'invalid-email-format'
                  };

                  await expect(userService.updateProfile(user.id, invalidProfileData))
                    .rejects.toThrow(ValidationError);
                });

                test('should sanitize and trim string inputs', async () => {
                  const user = await createTestUser();

                  const profileData = {
                    firstName: '  John  ',
                    lastName: '  Doe  '
                  };

                  const result = await userService.updateProfile(user.id, profileData);

                  expect(result.firstName).toBe('John');
                  expect(result.lastName).toBe('Doe');
                });

                test('should ignore unknown fields', async () => {
                  const user = await createTestUser();

                  const profileData = {
                    firstName: 'John',
                    maliciousField: 'should be ignored',
                    anotherBadField: { nested: 'object' }
                  };

                  const result = await userService.updateProfile(user.id, profileData);

                  expect(result.firstName).toBe('John');
                  expect(result.maliciousField).toBeUndefined();
                  expect(result.anotherBadField).toBeUndefined();
                });
              });

              describe('Error Handling', () => {
                test('should throw NotFoundError for non-existent user', async () => {
                  const nonExistentUserId = 'non-existent-id';
                  const profileData = { firstName: 'John' };

                  await expect(userService.updateProfile(nonExistentUserId, profileData))
                    .rejects.toThrow(NotFoundError);
                });

                test('should provide user-friendly error message for avatar-related issues', async () => {
                  // This test ensures our error handling improvement works
                  const user = await createTestUser();

                  // Mock a scenario where avatar access might still fail
                  jest.spyOn(user, 'update').mockImplementation(() => {
                    const error = new Error("Cannot read property 'id' of null (avatar)");
                    throw error;
                  });

                  await expect(userService.updateProfile(user.id, { firstName: 'John' }))
                    .rejects.toThrow('Profile update failed due to avatar data issue');
                });
              });

              describe('Logging and Monitoring', () => {
                test('should log successful profile updates', async () => {
                  const user = await createTestUser();
                  const logSpy = jest.spyOn(console, 'log').mockImplementation();

                  await userService.updateProfile(user.id, { firstName: 'John' });

                  // Verify logging occurred (in real implementation, use proper logger)
                  expect(logSpy).toHaveBeenCalled();
                  logSpy.mockRestore();
                });

                test('should sanitize sensitive data in error logs', async () => {
                  const user = await createTestUser();

                  const sensitiveData = {
                    email: 'user@example.com',
                    phone: '555-123-4567'
                  };

                  const sanitized = userService.sanitizeForLogging(sensitiveData);

                  expect(sanitized.email).toBe('us***@example.com');
                  expect(sanitized.phone).toBe('***-***-4567');
                });
              });

              describe('Performance and Transactions', () => {
                test('should use database transactions for consistency', async () => {
                  const user = await createTestUser();

                  // Mock transaction to verify it's being used
                  const transactionSpy = jest.spyOn(User.sequelize, 'transaction');

                  await userService.updateProfile(user.id, { firstName: 'John' });

                  expect(transactionSpy).toHaveBeenCalled();
                  transactionSpy.mockRestore();
                });
              });
            });
          `
        },

        // Integration test for API endpoint
        integrationTest: {
          file: "tests/integration/user-profile-fix.test.js",
          content: `
            const request = require('supertest');
            const app = require('../../src/app');
            const { createTestUser, getAuthToken } = require('../helpers');

            describe('User Profile API - Avatar Fix Integration', () => {
              test('PUT /api/users/profile should handle null avatar users', async () => {
                // Create legacy user with null avatar
                const user = await createTestUser({
                  avatar_url: null,
                  avatar: null
                });
                const token = await getAuthToken(user);

                // Update profile via API
                const response = await request(app)
                  .put('/api/users/profile')
                  .set('Authorization', \`Bearer \${token}\`)
                  .send({
                    firstName: 'Fixed',
                    lastName: 'User'
                  });

                expect(response.status).toBe(200);
                expect(response.body.success).toBe(true);
                expect(response.body.data.firstName).toBe('Fixed');
                expect(response.body.data.lastName).toBe('User');
              });

              test('should return proper error response for validation failures', async () => {
                const user = await createTestUser();
                const token = await getAuthToken(user);

                const response = await request(app)
                  .put('/api/users/profile')
                  .set('Authorization', \`Bearer \${token}\`)
                  .send({
                    email: 'invalid-email'
                  });

                expect(response.status).toBe(400);
                expect(response.body.success).toBe(false);
                expect(response.body.error).toContain('Invalid email format');
              });
            });
          `
        }
      };
    }
  },

  // 2. Fix Validation and Quality Assurance
  fixValidation: {
    // Comprehensive testing and validation
    validateFix: (fixImplementation, originalIssue) => {
      return {
        // Test execution results
        testResults: {
          unit_tests: {
            total: 12,
            passed: 12,
            failed: 0,
            coverage: "100% for modified functions",
            execution_time: "2.3 seconds"
          },
          integration_tests: {
            total: 4,
            passed: 4,
            failed: 0,
            scenarios_covered: [
              "API endpoint with null avatar user",
              "Database transaction rollback",
              "Error response formatting",
              "Authentication with profile update"
            ]
          },
          regression_tests: {
            total: 28,
            passed: 28,
            failed: 0,
            scope: "All user management functionality",
            critical_paths_verified: true
          }
        },

        // Performance impact assessment
        performanceImpact: {
          response_time_change: "+2ms average (negligible)",
          memory_usage_change: "No significant change",
          database_query_impact: "Same number of queries, improved safety",
          load_testing_results: "No degradation under normal load"
        },

        // Security review results
        securityReview: {
          new_vulnerabilities: "None introduced",
          security_improvements: [
            "Input validation added",
            "Sensitive data sanitization in logs",
            "Better error message exposure control"
          ],
          security_scan_results: "Clean - no issues detected"
        },

        // Fix effectiveness validation
        fixEffectiveness: {
          original_issue_resolved: true,
          edge_cases_covered: [
            "Null avatar users",
            "Undefined avatar property",
            "Legacy users from before avatar feature",
            "Users with valid existing avatars"
          ],
          monitoring_recommendations: [
            "Track profile update success rate",
            "Monitor for new null pointer exceptions",
            "Alert on increased validation errors"
          ]
        }
      };
    },

    // Production readiness checklist
    productionReadiness: (fixValidation, rolloutPlan) => {
      return {
        readiness_checklist: [
          {
            item: "All tests passing",
            status: "✅ Complete",
            details: "100% test success rate across all levels"
          },
          {
            item: "Code review approved",
            status: "✅ Complete",
            reviewer: "Senior developer + Security team",
            approval_date: "2025-07-12"
          },
          {
            item: "Database migration tested",
            status: "✅ Complete",
            details: "Successfully tested on staging with production data copy"
          },
          {
            item: "Rollback plan verified",
            status: "✅ Complete",
            details: "Rollback tested and verified in staging environment"
          },
          {
            item: "Monitoring alerts configured",
            status: "✅ Complete",
            details: "Error rate and performance alerts set up"
          },
          {
            item: "Documentation updated",
            status: "✅ Complete",
            details: "Technical docs and incident post-mortem completed"
          }
        ],

        deployment_plan: {
          timeline: "Immediate deployment approved for critical fix",
          strategy: "Blue-green deployment with 5% traffic test",
          monitoring_period: "24 hours intensive monitoring",
          success_criteria: [
            "Zero null pointer exceptions in user profile updates",
            "Error rate remains below 0.1%",
            "User profile update success rate > 99.5%",
            "No increase in support tickets"
          ],
          rollback_triggers: [
            "Error rate > 1%",
            "Any new critical issues detected",
            "User profile update success rate < 95%"
          ]
        }
      };
    }
  }
};
```

## Usage Examples

### Critical Bug Fix
```bash
/implement-fix "BUG-456" minimal local critical --test-driven --rollback-plan
```

**Expected Output:**
```markdown
# Fix Implementation: User Profile Null Pointer Exception

## Issue Analysis Summary
**Bug ID**: BUG-456
**Severity**: Critical
**Impact**: 15% of users unable to update profiles
**Root Cause**: Null avatar handling in UserService.updateProfile()

## Fix Strategy: Minimal Approach
**Rationale**: Quick resolution to stop user impact with low risk
**Implementation Time**: 2-4 hours
**Deployment Risk**: Low

## Implementation Completed ✅

### 1. Primary Fix Applied
- ✅ Added null safety check for user.avatar property
- ✅ Implemented optional chaining: `user.avatar?.id || null`
- ✅ Enhanced error handling with user-friendly messages
- **File**: `src/services/UserService.js`

### 2. Test Coverage Added
- ✅ Unit tests for null avatar scenarios (100% coverage)
- ✅ Integration tests for API endpoint behavior
- ✅ Regression tests for existing functionality
- **Files**: `tests/services/UserService.avatar-fix.test.js`

### 3. Data Migration Prepared
- ✅ Migration script for legacy users with null avatars
- ✅ Rollback script tested and verified
- **File**: `migrations/20250712_fix_null_avatars.js`

## Quality Validation Results

### Test Results ✅
- **Unit Tests**: 12/12 passed (100% coverage)
- **Integration Tests**: 4/4 passed
- **Regression Tests**: 28/28 passed
- **Performance**: No degradation detected

### Security Review ✅
- **Vulnerabilities**: None introduced
- **Improvements**: Input validation, log sanitization
- **Scan Results**: Clean

## Production Deployment Plan

### Deployment Strategy
1. **Blue-Green Deployment** with 5% traffic test
2. **Monitoring Period**: 24 hours intensive monitoring
3. **Rollback Time**: < 5 minutes if needed

### Success Criteria
- ✅ Zero null pointer exceptions in profile updates
- ✅ Error rate < 0.1%
- ✅ Profile update success rate > 99.5%

### Rollback Plan
- **Trigger**: Error rate > 1% or new critical issues
- **Method**: Git revert + database migration rollback
- **Estimated Time**: < 5 minutes

## Monitoring & Alerts Configured
- **Error Rate Monitoring**: Alert if > 0.5%
- **Profile Update Success Rate**: Alert if < 99%
- **New Exception Types**: Immediate alert
- **User Feedback**: Enhanced support ticket categorization

## Next Steps
1. **Deploy to Staging**: Validate with production data copy ✅
2. **Production Deployment**: Deploy during low-traffic window
3. **Post-Deployment**: 24-hour monitoring period
4. **Post-Mortem**: Document lessons learned and prevention measures
```

### Comprehensive Root Cause Fix
```bash
/implement-fix "SENTRY-789" comprehensive system high --preventive --upstream-fix --documentation
```

**Expected Output includes**:
- Complete system-wide null safety improvements
- Input validation middleware implementation
- TypeScript migration planning for type safety
- Automated linting rules for null safety
- Documentation of patterns and best practices

### Emergency Hotfix
```bash
/implement-fix "PROD-CRITICAL-001" hotfix local critical --emergency-deploy --monitoring
```

**Expected Output includes**:
- Immediate temporary fix implementation
- Fast-track deployment with minimal testing
- Enhanced monitoring and alerting
- Follow-up comprehensive fix planning

## Advanced Fix Features

### Batch Fix Application
```bash
/implement-fix "NULL-POINTER-PATTERN" comprehensive global medium --batch-fix
```

**Identifies and fixes similar patterns across the entire codebase**:
- Scans for similar null pointer vulnerabilities
- Applies consistent fix patterns
- Updates related test suites
- Documents the systematic improvements

### Performance-Safe Fixes
```bash
/implement-fix "PERF-ISSUE-123" minimal local high --performance-impact --monitoring
```

**Includes performance impact analysis**:
- Before/after performance benchmarks
- Memory usage analysis
- Database query optimization
- Load testing validation

### Security-Focused Fixes
```bash
/implement-fix "SEC-VULN-456" comprehensive system critical --security-review --upstream-fix
```

**Comprehensive security remediation**:
- Security vulnerability assessment
- Input validation and sanitization
- Authentication and authorization review
- Compliance requirement validation

## Integration with Development Workflow

### CI/CD Integration
```yaml
# .github/workflows/implement-fix.yml
name: Bug Fix Implementation Pipeline
on:
  workflow_dispatch:
    inputs:
      issue_id:
        description: 'Issue/Bug ID'
        required: true
      approach:
        description: 'Fix approach'
        required: true

jobs:
  implement-fix:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout code
        uses: actions/checkout@v3

      - name: Run Plan Mode Fix Implementation
        run: |
          /implement-fix ${{ github.event.inputs.issue_id }} ${{ github.event.inputs.approach }} local high --test-driven --rollback-plan

      - name: Run comprehensive tests
        run: |
          npm test
          npm run test:integration
          npm run test:security

      - name: Performance validation
        run: npm run test:performance

      - name: Deploy to staging
        run: |
          /deploy-feature "fix-${{ github.event.inputs.issue_id }}" staging --validate --monitor
```

### Issue Tracking Integration
```bash
# Automatically update issue status and add fix details
/implement-fix "JIRA-ABC-123" comprehensive module high --documentation | \
update-issue-tracker --status="fixed" --add-implementation-details

# Link fix to related issues
/implement-fix "github-456" preventive system medium --batch-fix | \
link-related-issues --pattern="null-pointer-exceptions"
```
