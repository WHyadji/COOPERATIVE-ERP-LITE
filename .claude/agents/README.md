# Agent Organization System

This directory contains specialized Claude Code agents organized by job function to improve workflow efficiency and agent discovery.

## Directory Structure

### 🏗️ Architecture (`architecture/`)
Agents focused on system design and architectural decisions:
- **architecture-design-expert.md** - Overall system architecture design
- **backend-architect.md** - Backend architecture and microservices patterns
- **database-architect.md** - Database design, data modeling, and polyglot persistence
- **domain-architect.md** - Domain-driven design implementation
- **microservices-architect.md** - Microservices architecture patterns

### 💻 Development (`development/`)
Agents for code generation and development best practices:
- **api-integration-specialist.md** - API integration, authentication flows, and real-time communication
- **code-generator.md** - Clean code generation following SOLID principles
- **code-quality-guardian.md** - Code review and quality analysis
- **code-reviewer.md** - Code review specialist for quality and best practices
- **frontend-architect.md** - Frontend architecture and patterns
- **frontend-developer.md** - Frontend development specialist
- **go-security-architect.md** - Go-specific security and architecture
- **ui-specialist.md** - Expert UI developer for React components, HTML/CSS, Tailwind, animations, and responsive design implementation

### 🧪 Testing (`testing/`)
Agents for comprehensive testing strategies:
- **e2e-test-specialist.md** - End-to-end testing implementation
- **load-tester.md** - Performance and load testing
- **qa-specialist-agent.md** - Quality assurance specialist for comprehensive testing and validation
- **test-architect.md** - Test strategy and architecture
- **test-data-manager.md** - Test data management and generation
- **tdd-integration-agent.md** - Test-driven development workflow coordination

### 🚀 DevOps (`devops/`)
Agents for deployment, infrastructure, and operations:
- **cicd-pipeline-architect.md** - CI/CD pipeline design and optimization
- **database-optimization.md** - Database performance optimization and query tuning
- **devops-engineer.md** - DevOps engineer for CI/CD, infrastructure automation, and cloud operations
- **devops-infrastructure.md** - DevOps and infrastructure specialist handling deployment automation, cloud infrastructure, containerization, and monitoring
- **integration-specialist.md** - System integration patterns
- **metrics-monitoring-specialist.md** - Monitoring and observability
- **performance-engineer.md** - Performance analysis, load testing, database optimization, and system scalability
- **performance-optimizer.md** - Performance optimization specialist

### 📚 Documentation (`documentation/`)
Agents for project documentation and analysis:
- **api-documenter.md** - API documentation specialist
- **doc-generator.md** - Technical documentation generation
- **progress-tracker.md** - Project progress monitoring
- **requirements-analyst.md** - Requirements analysis and management
- **technical-writer.md** - Technical writing specialist

### 🎨 Design (`design/`)
Agents for UI/UX design and optimization:
- **react-performance-optimization.md** - React performance optimization specialist
- **ui-ux-designer.md** - UI/UX design specialist

### 💬 Languages (`languages/`)
Language-specific programming experts:
- **javascript-pro.md** - JavaScript programming expert
- **python-pro.md** - Python programming expert
- **sql-pro.md** - SQL and database query expert
- **typescript-pro.md** - TypeScript programming expert

### 📋 Management (`management/`)
Agents for project coordination and strategic oversight:
- **pricing-strategist.md** - Pricing strategy development, revenue optimization, and competitive positioning
- **product-manager.md** - Proactively manages the entire product lifecycle from vision to launch with specialist agent orchestration
- **subscription-specialist.md** - SaaS monetization, billing systems, subscription lifecycle management, and revenue optimization strategies

### 🔒 Security (`security/`)
Agents focused on security and compliance:
- **api-contract-guardian.md** - API contract validation and security
- **security-specialist.md** - Comprehensive security analysis, vulnerability assessment, and secure implementation
- **supabase-security-expert.md** - Supabase-specific security implementation

### 🛠️ Tools (`tools/`)
Specialized utility agents:
- **accessibility-specialist.md** - Web accessibility audits and WCAG compliance validation
- **bugfix-agents.md** - Multi-agent bug analysis and resolution coordination
- **content-writer.md** - Content creation, copywriting, technical documentation, and SEO optimization
- **conversion-optimization-agent.md** - Conversion rate optimization through A/B testing and user behavior analysis
- **error-handler.md** - Error handling and debugging
- **messaging-systems-specialist.md** - Expert in Kafka, Redis, RabbitMQ, ActiveMQ, and NATS for distributed messaging architectures
- **sdlc-orchestrator.md** - Meta-agent that coordinates multiple specialized agents throughout the software development lifecycle
- **seo-specialist.md** - Expert SEO consultant specializing in technical SEO audits, on-page optimization, Core Web Vitals, and schema markup
- **ux-specialist.md** - UI/UX design, design systems, wireframes, prototypes, information architecture, and accessibility
- **visual-designer.md** - Visual design and branding specialist focusing on aesthetics, visual hierarchy, design systems, and brand expression

### 📦 Other (`other/`)
Specialized agents that don't fit standard categories:
- **mcp-expert.md** - MCP (Model Context Protocol) expert
- **nextjs-architecture-expert.md** - Next.js architecture specialist
- **seo-analyzer.md** - SEO analysis specialist
- **task-decomposition-expert.md** - Task decomposition and planning specialist

## Usage Guidelines

### Agent Selection
Choose agents based on your current task:
- **Strategic Planning**: Use `management/` and `architecture/` agents
- **Requirements Phase**: Use `documentation/` and `management/` agents
- **Development Phase**: Use `development/` and `testing/` agents
- **Deployment Phase**: Use `devops/` and `security/` agents
- **Maintenance Phase**: Use `tools/` and `scripts/` agents

### Multi-Agent Workflows
For complex tasks, consider using agents in sequence:
1. **product-manager** → **requirements-analyst** → **architecture-design-expert** → **code-generator**
2. **test-architect** → **e2e-test-specialist** → **load-tester**
3. **cicd-pipeline-architect** → **security/** agents → **performance-engineer**
4. **product-manager** → orchestrate multiple specialist agents throughout development lifecycle

### Context-Aware Selection
For accounting webapp specific tasks:
- **Go Backend**: Use `development/go-security-architect`, `architecture/backend-architect`, `architecture/microservices-architect`
- **Frontend**: Use `development/frontend-architect`, `design/ui-ux-designer`, `languages/typescript-pro`
- **Database**: Use `architecture/database-architect`, `devops/database-optimization`, `testing/test-data-manager`
- **Compliance**: Use `security/` agents, `security/api-contract-guardian`
- **Performance**: Use `devops/performance-engineer`, `devops/performance-optimizer`, `design/react-performance-optimization`
- **Documentation**: Use `documentation/` agents, especially `technical-writer` and `api-documenter`

## Maintenance

Keep this organization current by:
- Adding new agents to appropriate categories
- Updating this README when adding new categories
- Reviewing agent effectiveness and reorganizing as needed
- Considering cross-cutting concerns that span multiple categories