# Smart Split Commits

Intelligently analyze and split large changes into multiple atomic commits with advanced grouping strategies.

Instructions:

1. **Deep Analysis Phase**:
   ```bash
   # Get all changes including untracked files
   git status --porcelain

   # Analyze each file's changes in detail
   git diff --name-status
   git diff --stat

   # For each changed file, understand:
   - Type of changes (feature, fix, refactor, style, docs)
   - Dependencies with other files
   - Logical grouping possibilities
   ```

2. **Intelligent Grouping Strategies**:

   a) **By Feature/Module**:
      - Group all files related to a specific feature
      - Include tests with their implementation
      - Include docs with their feature

   b) **By Layer**:
      - Backend changes (models, controllers, services)
      - Frontend changes (components, styles, assets)
      - Infrastructure (config, build, CI/CD)
      - Database (migrations, schemas)

   c) **By Change Type**:
      - Bug fixes (group related fixes)
      - Refactoring (keep separate from features)
      - Style/formatting changes
      - Documentation updates
      - Test additions/updates

   d) **By Impact**:
      - Breaking changes (separate commit)
      - Performance improvements
      - Security fixes (prioritize these)

3. **Commit Plan Generation**:
   ```
   Proposed Commit Breakdown:

   Commit 1: fix(auth): resolve login validation issue
   - src/auth/validator.js
   - tests/auth/validator.test.js
   Reason: Bug fix with its test

   Commit 2: refactor(api): extract common middleware
   - src/middleware/common.js
   - src/api/users.js
   - src/api/posts.js
   Reason: Refactoring should be separate

   Commit 3: feat(users): add profile picture upload
   - src/users/profile.js
   - src/users/upload.js
   - tests/users/upload.test.js
   - docs/api/users.md
   Reason: Complete feature with tests and docs
   ```

4. **Handle Complex Scenarios**:

   a) **Single File with Multiple Changes**:
      - Detect if a file has unrelated changes
      - Suggest using `git add -p` for partial staging
      - Guide through interactive staging if needed

   b) **Circular Dependencies**:
      - Identify files that must be committed together
      - Warn about potential build breaks

   c) **Large Refactors**:
      - Break into multiple passes if possible
      - Ensure each commit compiles/tests pass

5. **Execution Modes**:

   Based on arguments:
   - **"analyze"**: Deep analysis with detailed reasoning
   - **"suggest"**: Quick suggestions without execution
   - **"interactive"**: Step-by-step with user choices
   - **"auto"**: Execute all commits (with safety checks)
   - **"partial"**: Use `git add -p` for fine-grained control
   - **"by:TYPE"**: Group by specific strategy (feature/layer/type)

6. **Safety Checks**:
   - Verify no uncommitted changes will be lost
   - Option to create backup branch
   - Show diff before each commit
   - Allow editing commit messages

7. **Output Format**:
   ```
   📊 Analysis Complete: Found 15 files with changes

   🎯 Suggested Breakdown: 4 commits

   ┌─ Commit 1/4 ─────────────────────────┐
   │ Type: fix                            │
   │ Scope: auth                          │
   │ Files: 2                             │
   │ Impact: Resolves critical bug #123   │
   └──────────────────────────────────────┘

   Continue? [Y/n/edit/skip]
   ```

Arguments: $ARGUMENTS

Advanced Examples:
- "analyze by:feature" - Group by features
- "interactive partial" - Interactive with partial file staging
- "auto max:5" - Auto-create up to 5 commits
- "suggest by:layer" - Suggest layer-based grouping
- "dry-run verbose" - Detailed analysis without commits

IMPORTANT:
- Each commit should build and pass tests
- Never push automatically
- Preserve the ability to amend/squash later
