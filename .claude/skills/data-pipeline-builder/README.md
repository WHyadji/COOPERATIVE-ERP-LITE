# Data Pipeline Builder Skill

A comprehensive data pipeline generation skill that creates production-ready ETL/ELT pipelines with orchestration, monitoring, and data quality checks.

## Features

### 🚀 Core Capabilities
- **Batch Processing**: Apache Airflow, Prefect, Dagster workflows
- **Stream Processing**: Spark Streaming, Flink, Kafka Streams
- **Data Transformation**: dbt, Spark, Pandas transformations
- **Orchestration**: DAG generation with dependencies and scheduling
- **Data Quality**: Great Expectations, custom validation checks
- **Monitoring**: Logging, metrics, alerting integration
- **Multiple Sources**: Databases, APIs, Files, Streaming platforms
- **Multiple Destinations**: Data warehouses, Lakes, Databases
- **Incremental Processing**: Change data capture, watermarking
- **Error Handling**: Retry logic, dead letter queues
- **Testing**: Unit tests, integration tests, data validation

## Quick Start

### Generate Batch ETL Pipeline

```bash
# Apache Airflow Pipeline
python scripts/generate_pipeline.py --type batch --orchestrator airflow --name "SalesETL"

# Prefect Pipeline
python scripts/generate_pipeline.py --type batch --orchestrator prefect --name "CustomerETL"

# Dagster Pipeline
python scripts/generate_pipeline.py --type batch --orchestrator dagster --name "ProductETL"
```

### Generate Streaming Pipeline

```bash
# Spark Streaming
python scripts/generate_pipeline.py --type stream --tool spark --name "EventStream" --source kafka --sink postgres

# Apache Flink
python scripts/generate_pipeline.py --type stream --tool flink --name "RealtimeAnalytics"

# Kafka Streams
python scripts/generate_pipeline.py --type stream --tool kafka-streams --name "OrderProcessing"
```

### Generate Transform Pipeline

```bash
# dbt Project
python scripts/generate_pipeline.py --type transform --tool dbt --name "Analytics"

# Spark SQL Transformations
python scripts/generate_pipeline.py --type transform --tool spark --name "DataTransform"
```

## Generated Project Structure

### Airflow Pipeline Structure
```
sales-etl/
├── dags/
│   ├── main_pipeline.py      # Main DAG definition
│   └── data_quality.py       # Quality check DAG
├── plugins/
│   ├── operators/            # Custom operators
│   ├── hooks/               # Custom hooks
│   └── sensors/             # Custom sensors
├── scripts/                  # Helper scripts
├── sql/                     # SQL queries
├── config/                  # Configuration files
├── tests/                   # Test files
├── docker-compose.yml       # Local development setup
└── requirements.txt         # Python dependencies
```

### Prefect Pipeline Structure
```
customer-etl/
├── flows/                   # Flow definitions
│   └── main_flow.py
├── tasks/                   # Reusable tasks
│   ├── extract.py
│   ├── transform.py
│   └── load.py
├── blocks/                  # Prefect blocks
├── deployments/            # Deployment configs
├── tests/                  # Test files
└── requirements.txt        # Dependencies
```

### dbt Project Structure
```
analytics/
├── models/
│   ├── staging/            # Staging models
│   ├── marts/             # Business logic models
│   └── intermediate/      # Intermediate models
├── tests/                 # Data tests
├── macros/               # Reusable SQL macros
├── analyses/             # Ad-hoc analyses
├── snapshots/            # SCD Type 2 history
├── dbt_project.yml       # Project configuration
└── profiles.yml          # Connection profiles
```

## Supported Technologies

### Orchestrators
| Tool | Use Case | Strengths |
|------|----------|-----------|
| **Apache Airflow** | Complex workflows | Mature, extensive operators, UI |
| **Prefect** | Modern Python workflows | Cloud-native, dynamic DAGs |
| **Dagster** | Data-aware orchestration | Asset-based, testing focus |
| **Luigi** | Simple pipelines | Lightweight, Python-native |

### Processing Engines
| Tool | Use Case | Strengths |
|------|----------|-----------|
| **Apache Spark** | Big data processing | Distributed, batch & stream |
| **Apache Flink** | Real-time streaming | Low latency, exactly-once |
| **dbt** | SQL transformations | Version control, testing |
| **Pandas** | Small-medium data | Rich functionality, easy |

### Data Sources
- **Databases**: PostgreSQL, MySQL, MongoDB, Cassandra
- **Cloud Storage**: S3, GCS, Azure Blob Storage
- **APIs**: REST, GraphQL, Webhooks
- **Streaming**: Kafka, Kinesis, Pub/Sub, Event Hubs
- **Files**: CSV, JSON, Parquet, Avro, ORC

### Data Destinations
- **Data Warehouses**: Snowflake, BigQuery, Redshift, Synapse
- **Databases**: PostgreSQL, MySQL, MongoDB
- **Data Lakes**: S3, ADLS Gen2, GCS
- **Analytics**: Elasticsearch, ClickHouse, Druid
- **Streaming**: Kafka, Event Hubs, Kinesis

## Pipeline Patterns

### ETL Pattern
```
Extract → Transform → Load
```
- Transform data before loading
- Best for: structured data, complex transformations
- Example: Database → Python processing → Data warehouse

### ELT Pattern
```
Extract → Load → Transform
```
- Load raw data, transform in destination
- Best for: cloud warehouses, big data
- Example: API → S3 → Snowflake → dbt

### Streaming Pattern
```
Source → Stream Processing → Sink
```
- Continuous real-time processing
- Best for: event data, real-time analytics
- Example: Kafka → Spark Streaming → Database

### Lambda Architecture
```
Batch Layer → Serving Layer
Speed Layer ↗
```
- Combines batch and streaming
- Best for: complete and fast results
- Example: Historical + real-time analytics

## Configuration Options

### Pipeline Configuration
```yaml
pipeline:
  name: sales_etl
  schedule: "0 2 * * *"  # Daily at 2 AM
  retries: 3
  retry_delay: 300  # seconds

sources:
  postgres:
    host: localhost
    port: 5432
    database: sales

transformations:
  - type: filter
    condition: "amount > 0"
  - type: aggregate
    group_by: [product_id]
    metrics: [sum(amount)]

destinations:
  snowflake:
    account: myaccount
    warehouse: ETL_WH
    database: ANALYTICS
```

## Data Quality Checks

Generated pipelines include:
- **Schema validation** - Column types, required fields
- **Completeness checks** - Null values, missing data
- **Uniqueness checks** - Primary keys, duplicates
- **Referential integrity** - Foreign key relationships
- **Business rules** - Custom validations
- **Freshness checks** - Data recency
- **Volume checks** - Row count thresholds
- **Distribution checks** - Statistical anomalies

## Monitoring & Alerting

Built-in monitoring:
- **Pipeline metrics** - Duration, success rate
- **Data metrics** - Row counts, processing volume
- **Quality metrics** - Validation pass rate
- **System metrics** - CPU, memory, disk usage
- **Alerting channels** - Email, Slack, PagerDuty
- **Dashboards** - Grafana, DataDog integration

## Testing

Each pipeline includes:
- **Unit tests** - Individual component testing
- **Integration tests** - End-to-end pipeline testing
- **Data validation tests** - Quality assertions
- **Performance tests** - Load and stress testing
- **Mock data generators** - Test data creation

## Deployment

### Local Development
```bash
# Start services
docker-compose up -d

# Run pipeline
make run

# Run tests
make test
```

### Production Deployment
- **Docker** - Containerized deployment
- **Kubernetes** - Orchestrated scaling
- **Cloud services** - AWS MWAA, GCP Composer
- **CI/CD** - GitHub Actions, GitLab CI

## Environment Variables

```env
# Database Connections
DB_HOST=localhost
DB_PORT=5432
DB_NAME=pipeline_db
DB_USER=user
DB_PASSWORD=password

# Cloud Storage
AWS_ACCESS_KEY_ID=xxx
AWS_SECRET_ACCESS_KEY=xxx
S3_BUCKET=data-bucket

# Data Warehouse
SNOWFLAKE_ACCOUNT=xxx
SNOWFLAKE_USER=xxx
SNOWFLAKE_PASSWORD=xxx

# Streaming
KAFKA_BOOTSTRAP_SERVERS=localhost:9092

# Monitoring
DATADOG_API_KEY=xxx
SLACK_WEBHOOK_URL=xxx
```

## Best Practices Implemented

1. **Idempotency** - Safe reruns
2. **Incremental processing** - Process only new data
3. **Error handling** - Comprehensive retry logic
4. **Data validation** - Quality checks at each stage
5. **Monitoring** - Metrics and alerting
6. **Documentation** - Self-documenting pipelines
7. **Version control** - Git-friendly structure
8. **Testing** - Comprehensive test coverage
9. **Security** - Credential management
10. **Scalability** - Distributed processing ready

## Command Reference

```bash
# Generate Airflow batch pipeline
python scripts/generate_pipeline.py \
  --type batch \
  --orchestrator airflow \
  --name "MyETL" \
  --source postgres \
  --sink snowflake

# Generate Spark streaming pipeline
python scripts/generate_pipeline.py \
  --type stream \
  --tool spark \
  --name "EventProcessor" \
  --source kafka \
  --sink elasticsearch

# Generate dbt transformation project
python scripts/generate_pipeline.py \
  --type transform \
  --tool dbt \
  --name "Analytics" \
  --output ./analytics-transform

# See all options
python scripts/generate_pipeline.py --help
```

## Customization

After generation, customize:
- Add business logic to transformations
- Configure specific data sources
- Add custom quality checks
- Extend monitoring metrics
- Integrate with your infrastructure

## Requirements

- Python 3.8+
- Docker & Docker Compose
- Specific tool requirements (Airflow, Spark, etc.)
- Database/warehouse access
- Cloud credentials (if using cloud services)

## License

This skill is provided for use with Claude AI. Generated pipelines are yours to use and modify.