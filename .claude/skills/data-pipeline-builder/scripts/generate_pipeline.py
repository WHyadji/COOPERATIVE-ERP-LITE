#!/usr/bin/env python3
"""
Data Pipeline Generator
Generates complete ETL/ELT pipelines with chosen orchestrator and processing engine
"""

import os
import sys
import json
import yaml
import argparse
from pathlib import Path
from typing import Dict, List, Optional
from datetime import datetime

class PipelineGenerator:
    def __init__(self, pipeline_type: str, orchestrator: str, name: str, options: Dict = None):
        self.pipeline_type = pipeline_type.lower()
        self.orchestrator = orchestrator.lower()
        self.name = name
        self.options = options or {}
        self.output_dir = Path(self.options.get('output', f"./{name.lower().replace(' ', '-')}"))
        
        # Supported configurations
        self.batch_orchestrators = {
            'airflow': self._generate_airflow,
            'prefect': self._generate_prefect,
            'dagster': self._generate_dagster,
            'luigi': self._generate_luigi
        }
        
        self.stream_processors = {
            'spark': self._generate_spark_streaming,
            'flink': self._generate_flink,
            'kafka-streams': self._generate_kafka_streams,
            'beam': self._generate_beam
        }
        
        self.transform_tools = {
            'dbt': self._generate_dbt,
            'spark': self._generate_spark_batch,
            'pandas': self._generate_pandas
        }
    
    def generate(self):
        """Main generation process"""
        print(f"🚀 Generating {self.pipeline_type} pipeline with {self.orchestrator}...")
        
        # Create output directory
        self.output_dir.mkdir(parents=True, exist_ok=True)
        
        # Generate based on type
        if self.pipeline_type == 'batch':
            if self.orchestrator in self.batch_orchestrators:
                self.batch_orchestrators[self.orchestrator]()
            else:
                raise ValueError(f"Unknown batch orchestrator: {self.orchestrator}")
                
        elif self.pipeline_type == 'stream':
            tool = self.options.get('tool', 'spark')
            if tool in self.stream_processors:
                self.stream_processors[tool]()
            else:
                raise ValueError(f"Unknown stream processor: {tool}")
                
        elif self.pipeline_type == 'transform':
            tool = self.options.get('tool', 'dbt')
            if tool in self.transform_tools:
                self.transform_tools[tool]()
            else:
                raise ValueError(f"Unknown transform tool: {tool}")
        else:
            raise ValueError(f"Unknown pipeline type: {self.pipeline_type}")
        
        # Generate common files
        self._generate_common_files()
        
        print(f"✅ Pipeline generated successfully in {self.output_dir}")
        print(f"\n📖 Next steps:")
        print(f"   1. cd {self.output_dir}")
        print(f"   2. Review README.md for setup")
        print(f"   3. Configure connections in config/")
        print(f"   4. Run the pipeline")
    
    def _create_file(self, path: str, content: str):
        """Create a file with content"""
        file_path = self.output_dir / path
        file_path.parent.mkdir(parents=True, exist_ok=True)
        file_path.write_text(content)
        print(f"    ✓ Created {path}")
    
    def _generate_airflow(self):
        """Generate Airflow pipeline"""
        print("  📦 Creating Apache Airflow pipeline...")
        
        # Create directory structure
        dirs = [
            'dags',
            'plugins/operators',
            'plugins/hooks',
            'plugins/sensors',
            'scripts',
            'sql',
            'tests',
            'config',
            'data'
        ]
        
        for dir_path in dirs:
            (self.output_dir / dir_path).mkdir(parents=True, exist_ok=True)
        
        # Generate files
        self._create_file('dags/main_pipeline.py', self._get_airflow_dag())
        self._create_file('dags/data_quality.py', self._get_airflow_quality_dag())
        self._create_file('plugins/operators/custom_operator.py', self._get_custom_operator())
        self._create_file('docker-compose.yml', self._get_airflow_docker_compose())
        self._create_file('requirements.txt', self._get_airflow_requirements())
        self._create_file('config/connections.yml', self._get_connections_config())
        self._create_file('tests/test_dag.py', self._get_dag_tests())
        
        print("  ✅ Airflow pipeline generated!")
    
    def _generate_prefect(self):
        """Generate Prefect pipeline"""
        print("  📦 Creating Prefect pipeline...")
        
        # Create directory structure
        dirs = [
            'flows',
            'tasks',
            'blocks',
            'deployments',
            'tests',
            'config'
        ]
        
        for dir_path in dirs:
            (self.output_dir / dir_path).mkdir(parents=True, exist_ok=True)
        
        # Generate files
        self._create_file('flows/main_flow.py', self._get_prefect_flow())
        self._create_file('tasks/extract.py', self._get_prefect_extract_tasks())
        self._create_file('tasks/transform.py', self._get_prefect_transform_tasks())
        self._create_file('tasks/load.py', self._get_prefect_load_tasks())
        self._create_file('deployments/deploy.py', self._get_prefect_deployment())
        self._create_file('requirements.txt', self._get_prefect_requirements())
        
        print("  ✅ Prefect pipeline generated!")
    
    def _generate_dbt(self):
        """Generate dbt project"""
        print("  📦 Creating dbt project...")
        
        # Create dbt project structure
        dirs = [
            'models/staging',
            'models/marts/core',
            'models/marts/finance',
            'models/intermediate',
            'tests',
            'macros',
            'analyses',
            'data',
            'snapshots'
        ]
        
        for dir_path in dirs:
            (self.output_dir / dir_path).mkdir(parents=True, exist_ok=True)
        
        # Generate files
        self._create_file('dbt_project.yml', self._get_dbt_project_config())
        self._create_file('profiles.yml', self._get_dbt_profiles())
        self._create_file('models/schema.yml', self._get_dbt_schema())
        self._create_file('models/staging/stg_orders.sql', self._get_dbt_staging_model())
        self._create_file('models/marts/core/fct_sales.sql', self._get_dbt_fact_model())
        self._create_file('models/marts/core/dim_customers.sql', self._get_dbt_dimension_model())
        self._create_file('macros/get_custom_schema.sql', self._get_dbt_macro())
        self._create_file('tests/assert_positive_values.sql', self._get_dbt_test())
        
        print("  ✅ dbt project generated!")
    
    def _generate_spark_streaming(self):
        """Generate Spark Streaming pipeline"""
        print("  📦 Creating Spark Streaming pipeline...")
        
        dirs = [
            'src/main/scala',
            'src/main/python',
            'src/main/resources',
            'src/test',
            'config',
            'scripts'
        ]
        
        for dir_path in dirs:
            (self.output_dir / dir_path).mkdir(parents=True, exist_ok=True)
        
        # Generate files
        self._create_file('src/main/python/streaming_pipeline.py', self._get_spark_streaming())
        self._create_file('src/main/python/batch_pipeline.py', self._get_spark_batch())
        self._create_file('build.sbt', self._get_sbt_config())
        self._create_file('requirements.txt', self._get_spark_requirements())
        self._create_file('submit_job.sh', self._get_spark_submit_script())
        
        print("  ✅ Spark Streaming pipeline generated!")
    
    def _generate_common_files(self):
        """Generate common files"""
        self._create_file('README.md', self._get_readme())
        self._create_file('.env.example', self._get_env_example())
        self._create_file('Dockerfile', self._get_dockerfile())
        self._create_file('.gitignore', self._get_gitignore())
        self._create_file('Makefile', self._get_makefile())
    
    # Template methods
    def _get_airflow_dag(self) -> str:
        return """from airflow import DAG
from airflow.operators.python import PythonOperator
from airflow.operators.bash import BashOperator
from airflow.providers.postgres.operators.postgres import PostgresOperator
from airflow.providers.amazon.sensors.s3 import S3KeySensor
from datetime import datetime, timedelta
import pandas as pd

default_args = {
    'owner': 'data-team',
    'depends_on_past': False,
    'start_date': datetime(2024, 1, 1),
    'email_on_failure': True,
    'email_on_retry': False,
    'retries': 2,
    'retry_delay': timedelta(minutes=5)
}

dag = DAG(
    'main_etl_pipeline',
    default_args=default_args,
    description='Main ETL pipeline',
    schedule_interval='@daily',
    catchup=False,
    tags=['etl', 'production']
)

def extract_data(**context):
    \"\"\"Extract data from source\"\"\"
    # Your extraction logic here
    df = pd.read_csv('/data/source/input.csv')
    df.to_parquet(f'/data/staging/data_{{{{ ds }}}}.parquet')
    return f"Extracted {{len(df)}} records"

def transform_data(**context):
    \"\"\"Transform data\"\"\"
    df = pd.read_parquet(f'/data/staging/data_{{{{ ds }}}}.parquet')
    # Transformation logic
    df['processed_at'] = datetime.now()
    df.to_parquet(f'/data/processed/data_{{{{ ds }}}}.parquet')
    return f"Transformed {{len(df)}} records"

def load_data(**context):
    \"\"\"Load data to destination\"\"\"
    df = pd.read_parquet(f'/data/processed/data_{{{{ ds }}}}.parquet')
    # Load logic
    return f"Loaded {{len(df)}} records"

# Define tasks
wait_for_file = S3KeySensor(
    task_id='wait_for_file',
    bucket_name='data-bucket',
    bucket_key='input/{{ ds }}/data.csv',
    aws_conn_id='aws_default',
    timeout=600,
    poke_interval=60,
    dag=dag
)

extract = PythonOperator(
    task_id='extract_data',
    python_callable=extract_data,
    dag=dag
)

transform = PythonOperator(
    task_id='transform_data',
    python_callable=transform_data,
    dag=dag
)

load = PythonOperator(
    task_id='load_data',
    python_callable=load_data,
    dag=dag
)

quality_check = PostgresOperator(
    task_id='quality_check',
    postgres_conn_id='warehouse',
    sql=\"\"\"
        SELECT COUNT(*) FROM processed_data
        WHERE date = '{{ ds }}'
        HAVING COUNT(*) > 0
    \"\"\",
    dag=dag
)

# Define dependencies
wait_for_file >> extract >> transform >> load >> quality_check"""
    
    def _get_prefect_flow(self) -> str:
        return """from prefect import flow, task
from prefect.tasks import task_input_hash
from datetime import timedelta
import pandas as pd

@task(
    name="Extract Data",
    retries=3,
    retry_delay=timedelta(seconds=10),
    cache_key_fn=task_input_hash,
    cache_expiration=timedelta(hours=1)
)
def extract_data(source: str) -> pd.DataFrame:
    \"\"\"Extract data from source\"\"\"
    # Extract logic
    df = pd.read_csv(source)
    return df

@task(name="Transform Data")
def transform_data(df: pd.DataFrame) -> pd.DataFrame:
    \"\"\"Transform data\"\"\"
    # Transformation logic
    df['processed'] = True
    return df

@task(name="Load Data")
def load_data(df: pd.DataFrame, destination: str) -> None:
    \"\"\"Load data to destination\"\"\"
    df.to_parquet(destination)

@flow(name="ETL Pipeline")
def etl_pipeline(source: str, destination: str):
    \"\"\"Main ETL pipeline\"\"\"
    raw_data = extract_data(source)
    transformed_data = transform_data(raw_data)
    load_data(transformed_data, destination)

if __name__ == "__main__":
    etl_pipeline(
        source="data/input.csv",
        destination="data/output.parquet"
    )"""
    
    def _get_dbt_project_config(self) -> str:
        return """name: 'analytics'
version: '1.0.0'
profile: 'analytics'

model-paths: ["models"]
analysis-paths: ["analyses"]
test-paths: ["tests"]
seed-paths: ["data"]
macro-paths: ["macros"]
snapshot-paths: ["snapshots"]

target-path: "target"
clean-targets:
  - "target"
  - "dbt_packages"

models:
  analytics:
    staging:
      +materialized: view
      +schema: staging
    marts:
      +materialized: table
      +schema: analytics
    intermediate:
      +materialized: ephemeral

vars:
  start_date: '2024-01-01'
  
on-run-start:
  - "{{ log('Starting dbt run', info=True) }}"

on-run-end:
  - "{{ log('Completed dbt run', info=True) }}"
"""
    
    def _get_spark_streaming(self) -> str:
        return """from pyspark.sql import SparkSession
from pyspark.sql.functions import *
from pyspark.sql.types import *

spark = SparkSession.builder \\
    .appName("StreamingPipeline") \\
    .config("spark.sql.adaptive.enabled", "true") \\
    .getOrCreate()

# Read from Kafka
df = spark \\
    .readStream \\
    .format("kafka") \\
    .option("kafka.bootstrap.servers", "localhost:9092") \\
    .option("subscribe", "events") \\
    .option("startingOffsets", "latest") \\
    .load()

# Parse JSON
schema = StructType([
    StructField("event_id", StringType()),
    StructField("user_id", StringType()),
    StructField("event_type", StringType()),
    StructField("timestamp", TimestampType()),
    StructField("amount", DoubleType())
])

events = df.select(
    from_json(col("value").cast("string"), schema).alias("data")
).select("data.*")

# Process with watermark
processed = events \\
    .withWatermark("timestamp", "10 minutes") \\
    .groupBy(
        window(col("timestamp"), "5 minutes"),
        col("event_type")
    ) \\
    .agg(
        count("*").alias("event_count"),
        sum("amount").alias("total_amount")
    )

# Write to console (for testing)
query = processed \\
    .writeStream \\
    .outputMode("update") \\
    .format("console") \\
    .start()

query.awaitTermination()"""
    
    def _get_airflow_docker_compose(self) -> str:
        return """version: '3.8'

x-airflow-common:
  &airflow-common
  image: apache/airflow:2.8.0
  environment:
    &airflow-common-env
    AIRFLOW__CORE__EXECUTOR: LocalExecutor
    AIRFLOW__DATABASE__SQL_ALCHEMY_CONN: postgresql+psycopg2://airflow:airflow@postgres/airflow
    AIRFLOW__CORE__FERNET_KEY: ''
    AIRFLOW__CORE__DAGS_ARE_PAUSED_AT_CREATION: 'true'
    AIRFLOW__CORE__LOAD_EXAMPLES: 'false'
  volumes:
    - ./dags:/opt/airflow/dags
    - ./logs:/opt/airflow/logs
    - ./plugins:/opt/airflow/plugins
    - ./data:/opt/airflow/data
  user: "${AIRFLOW_UID:-50000}:0"
  depends_on:
    &airflow-common-depends-on
    postgres:
      condition: service_healthy

services:
  postgres:
    image: postgres:13
    environment:
      POSTGRES_USER: airflow
      POSTGRES_PASSWORD: airflow
      POSTGRES_DB: airflow
    volumes:
      - postgres-db-volume:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD", "pg_isready", "-U", "airflow"]
      interval: 5s
      retries: 5
    restart: always

  airflow-webserver:
    <<: *airflow-common
    command: webserver
    ports:
      - 8080:8080
    healthcheck:
      test: ["CMD", "curl", "--fail", "http://localhost:8080/health"]
      interval: 10s
      timeout: 10s
      retries: 5
    restart: always
    depends_on:
      <<: *airflow-common-depends-on
      airflow-init:
        condition: service_completed_successfully

  airflow-scheduler:
    <<: *airflow-common
    command: scheduler
    healthcheck:
      test: ["CMD-SHELL", 'airflow jobs check --job-type SchedulerJob --hostname "$${HOSTNAME}"']
      interval: 10s
      timeout: 10s
      retries: 5
    restart: always
    depends_on:
      <<: *airflow-common-depends-on
      airflow-init:
        condition: service_completed_successfully

  airflow-init:
    <<: *airflow-common
    entrypoint: /bin/bash
    command:
      - -c
      - |
        airflow db init
        airflow users create \\
          --username admin \\
          --password admin \\
          --firstname Admin \\
          --lastname User \\
          --role Admin \\
          --email admin@example.com
    environment:
      <<: *airflow-common-env
      _AIRFLOW_DB_UPGRADE: 'true'
      _AIRFLOW_WWW_USER_CREATE: 'true'
      _AIRFLOW_WWW_USER_USERNAME: admin
      _AIRFLOW_WWW_USER_PASSWORD: admin

volumes:
  postgres-db-volume:"""
    
    def _get_readme(self) -> str:
        return f"""# {self.name} Data Pipeline

A production-ready data pipeline built with {self.orchestrator}.

## Architecture

```
Source → Extract → Transform → Load → Destination
         ↓          ↓          ↓
    Validation  Processing  Quality
```

## Quick Start

### Prerequisites
- Python 3.8+
- Docker & Docker Compose
- {"Apache Airflow" if self.orchestrator == "airflow" else self.orchestrator.title()}

### Installation

1. Clone the repository
```bash
git clone <repository>
cd {self.output_dir.name}
```

2. Set up environment
```bash
cp .env.example .env
# Edit .env with your configuration
```

3. Install dependencies
```bash
pip install -r requirements.txt
```

4. Start services
```bash
docker-compose up -d
```

5. Access UI
- {"Airflow UI: http://localhost:8080" if self.orchestrator == "airflow" else ""}
- {"Prefect UI: http://localhost:4200" if self.orchestrator == "prefect" else ""}

## Project Structure

```
.
├── {"dags/" if self.orchestrator == "airflow" else "flows/"}     # Pipeline definitions
├── {"plugins/" if self.orchestrator == "airflow" else "tasks/"}   # Custom components
├── config/      # Configuration files
├── tests/       # Test files
├── data/        # Local data directory
└── scripts/     # Utility scripts
```

## Configuration

Edit `config/connections.yml` to set up your data sources and destinations.

## Running the Pipeline

{"airflow dags trigger main_etl_pipeline" if self.orchestrator == "airflow" else ""}
{"prefect deployment run etl-pipeline/production" if self.orchestrator == "prefect" else ""}

## Testing

```bash
pytest tests/
```

## Monitoring

- Check pipeline status in the UI
- Logs are available in `logs/` directory
- Metrics are sent to monitoring service

## License

MIT"""
    
    def _get_env_example(self) -> str:
        return """# Environment Configuration

# Database
DB_HOST=localhost
DB_PORT=5432
DB_NAME=pipeline_db
DB_USER=pipeline
DB_PASSWORD=password

# Data Warehouse
WAREHOUSE_HOST=localhost
WAREHOUSE_PORT=5432
WAREHOUSE_NAME=warehouse
WAREHOUSE_USER=warehouse
WAREHOUSE_PASSWORD=password

# Cloud Storage
AWS_ACCESS_KEY_ID=
AWS_SECRET_ACCESS_KEY=
S3_BUCKET=data-bucket

# Streaming
KAFKA_BOOTSTRAP_SERVERS=localhost:9092
KAFKA_TOPIC=events

# Monitoring
DATADOG_API_KEY=
SLACK_WEBHOOK_URL=

# Orchestrator
AIRFLOW_UID=50000
PREFECT_API_URL=http://localhost:4200/api"""
    
    def _get_dockerfile(self) -> str:
        if self.orchestrator == 'airflow':
            return """FROM apache/airflow:2.8.0-python3.9

USER root
RUN apt-get update && apt-get install -y gcc g++ && rm -rf /var/lib/apt/lists/*

USER airflow
COPY requirements.txt /tmp/
RUN pip install --no-cache-dir -r /tmp/requirements.txt

COPY dags/ /opt/airflow/dags/
COPY plugins/ /opt/airflow/plugins/
COPY scripts/ /opt/airflow/scripts/"""
        else:
            return """FROM python:3.9-slim

WORKDIR /app

RUN apt-get update && apt-get install -y gcc g++ && rm -rf /var/lib/apt/lists/*

COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

COPY . .

CMD ["python", "main.py"]"""
    
    def _get_makefile(self) -> str:
        return """# Makefile for Data Pipeline

.PHONY: help install test run clean

help:
	@echo "Available commands:"
	@echo "  make install    Install dependencies"
	@echo "  make test       Run tests"
	@echo "  make run        Run pipeline"
	@echo "  make clean      Clean up"

install:
	pip install -r requirements.txt

test:
	pytest tests/ -v

run:
	docker-compose up -d

clean:
	docker-compose down
	find . -type d -name __pycache__ -exec rm -rf {} +
	rm -rf .pytest_cache

start:
	docker-compose up -d

stop:
	docker-compose stop

logs:
	docker-compose logs -f

reset:
	docker-compose down -v
	docker-compose up -d"""
    
    def _get_gitignore(self) -> str:
        return """# Python
__pycache__/
*.py[cod]
*$py.class
*.so
.Python
venv/
env/
ENV/

# Data
data/
*.csv
*.parquet
*.json
*.db

# Logs
logs/
*.log

# Environment
.env
.env.local

# IDE
.vscode/
.idea/
*.swp
.DS_Store

# Testing
.pytest_cache/
.coverage
htmlcov/

# Airflow
airflow.db
airflow-webserver.pid

# dbt
target/
dbt_packages/
logs/

# Prefect
.prefect/

# Docker
.docker/"""
    
    # Additional template methods for other files...
    def _get_airflow_requirements(self) -> str:
        return """apache-airflow==2.8.0
pandas==2.0.3
numpy==1.24.3
sqlalchemy==2.0.23
psycopg2-binary==2.9.9
boto3==1.29.7
requests==2.31.0
pyarrow==14.0.1
great-expectations==0.18.3
pytest==7.4.3"""
    
    def _get_prefect_requirements(self) -> str:
        return """prefect==2.14.0
pandas==2.0.3
numpy==1.24.3
sqlalchemy==2.0.23
psycopg2-binary==2.9.9
boto3==1.29.7
requests==2.31.0
pyarrow==14.0.1
dask==2023.12.0
pytest==7.4.3"""
    
    def _get_spark_requirements(self) -> str:
        return """pyspark==3.5.0
pandas==2.0.3
numpy==1.24.3
pyarrow==14.0.1
delta-spark==3.0.0
kafka-python==2.0.2
boto3==1.29.7"""
    
    def _get_dbt_staging_model(self) -> str:
        return """{{ config(
    materialized='view',
    tags=['staging']
) }}

WITH source AS (
    SELECT *
    FROM {{ source('raw', 'orders') }}
    WHERE created_at >= '{{ var("start_date") }}'
),

cleaned AS (
    SELECT
        order_id,
        customer_id,
        product_id,
        CAST(order_date AS DATE) as order_date,
        CAST(quantity AS INTEGER) as quantity,
        CAST(amount AS DECIMAL(10,2)) as amount,
        status,
        created_at,
        updated_at
    FROM source
    WHERE status NOT IN ('test', 'deleted')
)

SELECT * FROM cleaned"""
    
    def _get_custom_operator(self) -> str:
        return """from airflow.models.baseoperator import BaseOperator
from airflow.utils.decorators import apply_defaults

class DataQualityOperator(BaseOperator):
    \"\"\"Custom operator for data quality checks\"\"\"
    
    @apply_defaults
    def __init__(self, checks, conn_id, *args, **kwargs):
        super().__init__(*args, **kwargs)
        self.checks = checks
        self.conn_id = conn_id
    
    def execute(self, context):
        \"\"\"Execute data quality checks\"\"\"
        for check in self.checks:
            # Run check logic
            pass
        return True"""
    
    # Stub methods for frameworks not fully implemented
    def _generate_dagster(self):
        print("  📦 Dagster support coming soon...")
    
    def _generate_luigi(self):
        print("  📦 Luigi support coming soon...")
    
    def _generate_flink(self):
        print("  📦 Apache Flink support coming soon...")
    
    def _generate_kafka_streams(self):
        print("  📦 Kafka Streams support coming soon...")
    
    def _generate_beam(self):
        print("  📦 Apache Beam support coming soon...")
    
    def _generate_spark_batch(self):
        print("  📦 Spark batch processing included...")
    
    def _generate_pandas(self):
        print("  📦 Pandas processing included...")

def main():
    parser = argparse.ArgumentParser(description="Generate data pipelines")
    parser.add_argument(
        "--type",
        choices=["batch", "stream", "transform"],
        required=True,
        help="Type of pipeline"
    )
    parser.add_argument(
        "--orchestrator",
        help="Orchestrator to use (airflow, prefect, dagster)"
    )
    parser.add_argument(
        "--tool",
        help="Processing tool (spark, flink, dbt, pandas)"
    )
    parser.add_argument(
        "--name",
        required=True,
        help="Pipeline name"
    )
    parser.add_argument(
        "--source",
        help="Data source (postgres, mysql, s3, kafka)"
    )
    parser.add_argument(
        "--sink",
        help="Data destination (warehouse, s3, kafka)"
    )
    parser.add_argument(
        "--output",
        help="Output directory"
    )
    
    args = parser.parse_args()
    
    # Set defaults based on type
    if args.type == "batch" and not args.orchestrator:
        args.orchestrator = "airflow"
    elif args.type == "stream" and not args.tool:
        args.tool = "spark"
    elif args.type == "transform" and not args.tool:
        args.tool = "dbt"
    
    options = {
        "output": args.output,
        "tool": args.tool,
        "source": args.source,
        "sink": args.sink
    }
    
    generator = PipelineGenerator(
        pipeline_type=args.type,
        orchestrator=args.orchestrator or args.tool,
        name=args.name,
        options=options
    )
    
    try:
        generator.generate()
    except Exception as e:
        print(f"❌ Error: {e}", file=sys.stderr)
        sys.exit(1)

if __name__ == "__main__":
    main()