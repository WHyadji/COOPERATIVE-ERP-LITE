---
name: Data Pipeline Builder
description: Generate production-ready data pipelines with ETL/ELT workflows, supporting multiple data sources, transformations, orchestration, and monitoring. Includes Apache Airflow, Prefect, dbt, and streaming pipelines.
---

# Data Pipeline Builder

## Overview

This skill generates complete data pipeline solutions with ETL/ELT workflows, data quality checks, orchestration, monitoring, and documentation. Supports batch processing, real-time streaming, and modern data stack architectures.

## Quick Start

### Generate ETL Pipeline
```bash
# Apache Airflow Pipeline
python scripts/generate_pipeline.py --type batch --orchestrator airflow --name "SalesETL"

# Prefect Pipeline
python scripts/generate_pipeline.py --type batch --orchestrator prefect --name "CustomerData"

# dbt + Airflow
python scripts/generate_pipeline.py --type transform --tool dbt --orchestrator airflow
```

### Generate Streaming Pipeline
```bash
# Kafka + Spark Streaming
python scripts/generate_pipeline.py --type stream --tool spark --source kafka --sink postgres

# Apache Flink Pipeline
python scripts/generate_pipeline.py --type stream --tool flink --source kinesis
```

## Supported Technologies

### Orchestrators
- **Apache Airflow** - Complex DAGs, enterprise-ready
- **Prefect** - Modern Python-native orchestration
- **Dagster** - Data-aware orchestration
- **Apache NiFi** - Visual data flow
- **Temporal** - Workflow orchestration

### Processing Engines
- **Apache Spark** - Large-scale batch/stream processing
- **Apache Flink** - Real-time stream processing
- **Apache Beam** - Unified batch/stream
- **dbt** - SQL-based transformations
- **Pandas** - Python data processing

### Data Sources
- **Databases**: PostgreSQL, MySQL, MongoDB, Cassandra, Redis
- **Cloud Storage**: S3, GCS, Azure Blob
- **APIs**: REST, GraphQL, Webhooks
- **Streaming**: Kafka, Kinesis, Pub/Sub
- **Files**: CSV, JSON, Parquet, Avro

### Data Destinations
- **Data Warehouses**: Snowflake, BigQuery, Redshift, Synapse
- **Databases**: PostgreSQL, MySQL, MongoDB
- **Data Lakes**: S3, ADLS, GCS
- **Analytics**: Elasticsearch, ClickHouse
- **Streaming**: Kafka, Event Hubs

## Pipeline Patterns

### ETL Pattern (Extract-Transform-Load)
```python
# Traditional ETL - Transform before loading
Extract → Transform → Load
Source → Processing → Destination

# Best for:
- Structured data
- Complex transformations
- Data cleansing before storage
```

### ELT Pattern (Extract-Load-Transform)
```python
# Modern ELT - Transform in destination
Extract → Load → Transform
Source → Raw Storage → Processing → Analytics

# Best for:
- Cloud data warehouses
- Big data volumes
- Flexible transformations
```

### Real-time Streaming
```python
# Continuous processing
Source → Stream Processing → Sink
Kafka → Spark/Flink → Database

# Best for:
- Real-time analytics
- Event processing
- IoT data
```

## Apache Airflow Pipeline

### DAG Structure
```python
from airflow import DAG
from airflow.operators.python import PythonOperator
from airflow.providers.postgres.operators.postgres import PostgresOperator
from airflow.providers.amazon.transfers.s3_to_redshift import S3ToRedshiftOperator
from datetime import datetime, timedelta

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
    'sales_etl_pipeline',
    default_args=default_args,
    description='Daily sales data ETL pipeline',
    schedule_interval='@daily',
    catchup=False,
    tags=['sales', 'etl']
)

# Tasks
def extract_from_source(**context):
    """Extract data from source systems"""
    import pandas as pd
    from sqlalchemy import create_engine
    
    engine = create_engine(CONNECTION_STRING)
    
    query = """
    SELECT *
    FROM sales
    WHERE date = '{{ ds }}'
    """
    
    df = pd.read_sql(query, engine)
    
    # Save to staging
    df.to_parquet(f'/tmp/sales_{{ ds }}.parquet')
    
    return f"Extracted {len(df)} records"

def transform_data(**context):
    """Apply transformations"""
    import pandas as pd
    
    # Read staged data
    df = pd.read_parquet(f'/tmp/sales_{{ ds }}.parquet')
    
    # Transformations
    df['revenue'] = df['quantity'] * df['unit_price']
    df['tax'] = df['revenue'] * 0.1
    df['total'] = df['revenue'] + df['tax']
    
    # Data quality checks
    assert df['revenue'].min() >= 0, "Negative revenue found"
    assert df['quantity'].min() > 0, "Invalid quantity"
    
    # Aggregations
    daily_summary = df.groupby(['product_id', 'region']).agg({
        'quantity': 'sum',
        'revenue': 'sum',
        'total': 'sum',
        'order_id': 'count'
    }).reset_index()
    
    # Save transformed data
    daily_summary.to_parquet(f'/tmp/sales_transformed_{{ ds }}.parquet')
    
    return f"Transformed {len(daily_summary)} records"

def load_to_warehouse(**context):
    """Load to data warehouse"""
    import pandas as pd
    from sqlalchemy import create_engine
    
    df = pd.read_parquet(f'/tmp/sales_transformed_{{ ds }}.parquet')
    
    engine = create_engine(WAREHOUSE_CONNECTION)
    
    # Load to warehouse
    df.to_sql(
        'fact_sales',
        engine,
        if_exists='append',
        index=False,
        method='multi'
    )
    
    return f"Loaded {len(df)} records to warehouse"

# Define tasks
extract_task = PythonOperator(
    task_id='extract_from_source',
    python_callable=extract_from_source,
    dag=dag
)

transform_task = PythonOperator(
    task_id='transform_data',
    python_callable=transform_data,
    dag=dag
)

load_task = PythonOperator(
    task_id='load_to_warehouse',
    python_callable=load_to_warehouse,
    dag=dag
)

# Data quality check
quality_check = PostgresOperator(
    task_id='data_quality_check',
    postgres_conn_id='warehouse',
    sql="""
        SELECT COUNT(*) as count
        FROM fact_sales
        WHERE date = '{{ ds }}'
        HAVING COUNT(*) > 0
    """,
    dag=dag
)

# Task dependencies
extract_task >> transform_task >> load_task >> quality_check
```

### Sensors and Hooks
```python
from airflow.sensors.s3_key_sensor import S3KeySensor
from airflow.hooks.postgres_hook import PostgresHook
from airflow.sensors.external_task import ExternalTaskSensor

# Wait for file in S3
wait_for_file = S3KeySensor(
    task_id='wait_for_file',
    bucket_name='data-bucket',
    bucket_key='incoming/{{ ds }}/sales.csv',
    aws_conn_id='aws_default',
    timeout=3600,
    poke_interval=60,
    dag=dag
)

# Wait for upstream DAG
wait_for_upstream = ExternalTaskSensor(
    task_id='wait_for_upstream',
    external_dag_id='customer_pipeline',
    external_task_id='export_complete',
    dag=dag
)

# Custom sensor
class DataQualitySensor(BaseSensorOperator):
    def poke(self, context):
        hook = PostgresHook(postgres_conn_id='warehouse')
        result = hook.get_first(
            "SELECT COUNT(*) FROM staging_table WHERE date = '{{ ds }}'"
        )
        return result[0] > 0
```

## Prefect Pipeline

### Flow Definition
```python
from prefect import flow, task
from prefect.tasks import task_input_hash
from datetime import timedelta
import pandas as pd
import numpy as np

@task(
    name="Extract Data",
    retries=3,
    retry_delay=timedelta(seconds=10),
    cache_key_fn=task_input_hash,
    cache_expiration=timedelta(hours=1)
)
def extract_data(source: str, date: str) -> pd.DataFrame:
    """Extract data from source"""
    logger = get_run_logger()
    logger.info(f"Extracting data from {source} for {date}")
    
    # Connection logic
    if source == "postgres":
        conn_string = get_secret("POSTGRES_CONN")
        df = pd.read_sql(
            f"SELECT * FROM orders WHERE date = '{date}'",
            conn_string
        )
    elif source == "api":
        response = requests.get(
            f"https://api.example.com/data?date={date}",
            headers={"Authorization": f"Bearer {get_secret('API_KEY')}"}
        )
        df = pd.DataFrame(response.json())
    
    logger.info(f"Extracted {len(df)} records")
    return df

@task(name="Transform Data")
def transform_data(df: pd.DataFrame) -> pd.DataFrame:
    """Apply transformations"""
    logger = get_run_logger()
    
    # Data cleaning
    df = df.dropna(subset=['customer_id', 'amount'])
    df['amount'] = pd.to_numeric(df['amount'], errors='coerce')
    
    # Feature engineering
    df['revenue'] = df['quantity'] * df['amount']
    df['tax'] = df['revenue'] * 0.1
    df['total'] = df['revenue'] + df['tax']
    
    # Aggregations
    summary = df.groupby(['product_id', 'category']).agg({
        'revenue': 'sum',
        'quantity': 'sum',
        'customer_id': 'nunique'
    }).reset_index()
    
    logger.info(f"Transformed data to {len(summary)} records")
    return summary

@task(name="Data Quality Check")
def validate_data(df: pd.DataFrame) -> bool:
    """Validate data quality"""
    checks = []
    
    # Schema validation
    required_columns = ['product_id', 'revenue', 'quantity']
    checks.append(all(col in df.columns for col in required_columns))
    
    # Data quality rules
    checks.append(df['revenue'].min() >= 0)
    checks.append(df['quantity'].min() > 0)
    checks.append(not df['product_id'].isna().any())
    
    if not all(checks):
        raise ValueError("Data quality checks failed")
    
    return True

@task(name="Load to Warehouse")
def load_data(df: pd.DataFrame, target: str) -> None:
    """Load data to target"""
    logger = get_run_logger()
    
    if target == "snowflake":
        conn = snowflake.connector.connect(
            user=get_secret("SNOWFLAKE_USER"),
            password=get_secret("SNOWFLAKE_PASSWORD"),
            account=get_secret("SNOWFLAKE_ACCOUNT"),
            warehouse="ETL_WH",
            database="ANALYTICS",
            schema="PUBLIC"
        )
        
        # Write to Snowflake
        success, nchunks, nrows, _ = write_pandas(
            conn, df, "FACT_SALES", auto_create_table=True
        )
        
        logger.info(f"Loaded {nrows} rows to Snowflake")
    
    elif target == "bigquery":
        client = bigquery.Client()
        table_id = "project.dataset.fact_sales"
        
        job_config = bigquery.LoadJobConfig(
            write_disposition="WRITE_APPEND",
            schema_autodetect=True
        )
        
        job = client.load_table_from_dataframe(
            df, table_id, job_config=job_config
        )
        job.result()
        
        logger.info(f"Loaded {len(df)} rows to BigQuery")

@flow(name="ETL Pipeline")
def etl_pipeline(
    source: str = "postgres",
    target: str = "snowflake",
    date: str = None
) -> None:
    """Main ETL pipeline flow"""
    if date is None:
        date = datetime.now().strftime("%Y-%m-%d")
    
    # Pipeline steps
    raw_data = extract_data(source, date)
    transformed_data = transform_data(raw_data)
    is_valid = validate_data(transformed_data)
    
    if is_valid:
        load_data(transformed_data, target)

# Schedule and deploy
if __name__ == "__main__":
    # Create deployment
    deployment = Deployment.build_from_flow(
        flow=etl_pipeline,
        name="daily-etl",
        schedule=IntervalSchedule(interval=timedelta(days=1)),
        parameters={
            "source": "postgres",
            "target": "snowflake"
        },
        tags=["etl", "production"]
    )
    
    deployment.apply()
```

## dbt (Data Build Tool) Pipeline

### Project Structure
```yaml
# dbt_project.yml
name: 'analytics'
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
```

### Models
```sql
-- models/staging/stg_orders.sql
{{ config(
    materialized='view',
    tags=['staging', 'daily']
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
        CAST(order_date AS DATE) as order_date,
        CAST(amount AS DECIMAL(10,2)) as amount,
        status,
        {{ dbt_utils.generate_surrogate_key(['order_id', 'customer_id']) }} as order_key
    FROM source
    WHERE status NOT IN ('cancelled', 'failed')
)

SELECT * FROM cleaned
```

```sql
-- models/marts/fct_sales.sql
{{ config(
    materialized='incremental',
    unique_key='sale_id',
    on_schema_change='fail',
    tags=['marts', 'sales']
) }}

WITH orders AS (
    SELECT * FROM {{ ref('stg_orders') }}
    {% if is_incremental() %}
        WHERE order_date > (SELECT MAX(order_date) FROM {{ this }})
    {% endif %}
),

customers AS (
    SELECT * FROM {{ ref('dim_customers') }}
),

products AS (
    SELECT * FROM {{ ref('dim_products') }}
),

final AS (
    SELECT
        {{ dbt_utils.generate_surrogate_key(['o.order_id', 'o.order_date']) }} as sale_id,
        o.order_date,
        o.customer_id,
        c.customer_segment,
        c.customer_lifetime_value,
        o.product_id,
        p.product_category,
        p.product_subcategory,
        o.quantity,
        o.amount,
        o.quantity * o.amount as revenue,
        (o.quantity * o.amount) * 0.1 as tax,
        (o.quantity * o.amount) * 1.1 as total_amount,
        CURRENT_TIMESTAMP as processed_at
    FROM orders o
    LEFT JOIN customers c ON o.customer_id = c.customer_id
    LEFT JOIN products p ON o.product_id = p.product_id
)

SELECT * FROM final
```

### Tests
```yaml
# models/schema.yml
version: 2

models:
  - name: fct_sales
    description: "Sales fact table"
    columns:
      - name: sale_id
        description: "Unique identifier for each sale"
        tests:
          - unique
          - not_null
      - name: revenue
        description: "Total revenue"
        tests:
          - not_null
          - dbt_expectations.expect_column_values_to_be_between:
              min_value: 0
              max_value: 1000000
      - name: customer_id
        tests:
          - relationships:
              to: ref('dim_customers')
              field: customer_id

sources:
  - name: raw
    database: raw_db
    schema: public
    tables:
      - name: orders
        description: "Raw orders data"
        columns:
          - name: order_id
            tests:
              - not_null
              - unique
        freshness:
          warn_after: {count: 12, period: hour}
          error_after: {count: 24, period: hour}
        loaded_at_field: created_at
```

## Spark Pipeline

### Batch Processing
```python
from pyspark.sql import SparkSession
from pyspark.sql.functions import *
from pyspark.sql.types import *
from delta.tables import DeltaTable

# Initialize Spark
spark = SparkSession.builder \
    .appName("ETL_Pipeline") \
    .config("spark.sql.extensions", "io.delta.sql.DeltaSparkSessionExtension") \
    .config("spark.sql.catalog.spark_catalog", "org.apache.spark.sql.delta.catalog.DeltaCatalog") \
    .getOrCreate()

# Read from multiple sources
def extract_data():
    """Extract data from various sources"""
    
    # Read from database
    jdbc_df = spark.read \
        .format("jdbc") \
        .option("url", "jdbc:postgresql://localhost:5432/sourcedb") \
        .option("dbtable", "orders") \
        .option("user", "username") \
        .option("password", "password") \
        .option("driver", "org.postgresql.Driver") \
        .option("fetchsize", "10000") \
        .load()
    
    # Read from S3 Parquet
    s3_df = spark.read \
        .parquet("s3a://bucket/data/orders/*.parquet")
    
    # Read from Kafka
    kafka_df = spark.read \
        .format("kafka") \
        .option("kafka.bootstrap.servers", "localhost:9092") \
        .option("subscribe", "orders-topic") \
        .option("startingOffsets", "earliest") \
        .load()
    
    # Parse Kafka JSON
    kafka_parsed = kafka_df \
        .select(from_json(col("value").cast("string"), orders_schema).alias("data")) \
        .select("data.*")
    
    # Union all sources
    all_orders = jdbc_df.unionByName(s3_df).unionByName(kafka_parsed)
    
    return all_orders

def transform_data(df):
    """Apply transformations"""
    
    # Data cleaning
    cleaned_df = df \
        .filter(col("amount") > 0) \
        .filter(col("status").isin(["completed", "shipped"])) \
        .dropDuplicates(["order_id"])
    
    # Add calculated fields
    transformed_df = cleaned_df \
        .withColumn("revenue", col("quantity") * col("amount")) \
        .withColumn("tax", col("revenue") * 0.1) \
        .withColumn("total", col("revenue") + col("tax")) \
        .withColumn("order_date", to_date(col("order_timestamp"))) \
        .withColumn("order_hour", hour(col("order_timestamp"))) \
        .withColumn("is_weekend", when(dayofweek(col("order_date")).isin([1, 7]), True).otherwise(False))
    
    # Aggregations
    daily_summary = transformed_df \
        .groupBy("order_date", "product_id", "customer_segment") \
        .agg(
            sum("revenue").alias("total_revenue"),
            sum("quantity").alias("total_quantity"),
            countDistinct("customer_id").alias("unique_customers"),
            avg("revenue").alias("avg_order_value"),
            max("revenue").alias("max_order_value")
        )
    
    # Window functions
    window_spec = Window.partitionBy("customer_id").orderBy("order_date")
    
    with_customer_metrics = transformed_df \
        .withColumn("customer_order_rank", row_number().over(window_spec)) \
        .withColumn("customer_total_spent", sum("total").over(window_spec)) \
        .withColumn("days_since_last_order", 
                   datediff(col("order_date"), lag("order_date", 1).over(window_spec)))
    
    return with_customer_metrics

def load_data(df, output_path):
    """Load data to destination"""
    
    # Write to Delta Lake
    df.write \
        .format("delta") \
        .mode("overwrite") \
        .option("overwriteSchema", "true") \
        .partitionBy("order_date") \
        .save(f"{output_path}/delta/orders")
    
    # Create or replace Delta table
    df.write \
        .format("delta") \
        .mode("overwrite") \
        .option("overwriteSchema", "true") \
        .saveAsTable("analytics.fact_orders")
    
    # Optimize Delta table
    delta_table = DeltaTable.forPath(spark, f"{output_path}/delta/orders")
    delta_table.optimize().executeCompaction()
    delta_table.vacuum(168)  # 7 days
    
    # Write to PostgreSQL
    df.write \
        .format("jdbc") \
        .option("url", "jdbc:postgresql://localhost:5432/warehouse") \
        .option("dbtable", "fact_orders") \
        .option("user", "username") \
        .option("password", "password") \
        .option("driver", "org.postgresql.Driver") \
        .mode("append") \
        .save()

# Main pipeline
def run_pipeline():
    """Execute ETL pipeline"""
    try:
        # Extract
        raw_data = extract_data()
        print(f"Extracted {raw_data.count()} records")
        
        # Transform
        transformed_data = transform_data(raw_data)
        transformed_data.cache()
        print(f"Transformed {transformed_data.count()} records")
        
        # Data quality checks
        assert transformed_data.filter(col("revenue") < 0).count() == 0, "Negative revenue found"
        
        # Load
        load_data(transformed_data, "/data/warehouse")
        print("Pipeline completed successfully")
        
    except Exception as e:
        print(f"Pipeline failed: {e}")
        raise
    finally:
        spark.stop()

if __name__ == "__main__":
    run_pipeline()
```

### Streaming Pipeline
```python
from pyspark.sql import SparkSession
from pyspark.sql.functions import *
from pyspark.sql.types import *

spark = SparkSession.builder \
    .appName("StreamingPipeline") \
    .config("spark.sql.streaming.checkpointLocation", "/checkpoint") \
    .getOrCreate()

# Define schema
event_schema = StructType([
    StructField("event_id", StringType()),
    StructField("user_id", StringType()),
    StructField("event_type", StringType()),
    StructField("timestamp", TimestampType()),
    StructField("properties", MapType(StringType(), StringType()))
])

# Read from Kafka stream
events_df = spark \
    .readStream \
    .format("kafka") \
    .option("kafka.bootstrap.servers", "localhost:9092") \
    .option("subscribe", "events") \
    .option("startingOffsets", "latest") \
    .load()

# Parse JSON events
parsed_events = events_df \
    .select(from_json(col("value").cast("string"), event_schema).alias("data")) \
    .select("data.*") \
    .withWatermark("timestamp", "10 minutes")

# Aggregate in windows
windowed_stats = parsed_events \
    .groupBy(
        window(col("timestamp"), "5 minutes", "1 minute"),
        col("event_type")
    ) \
    .agg(
        count("*").alias("event_count"),
        countDistinct("user_id").alias("unique_users"),
        min("timestamp").alias("window_start"),
        max("timestamp").alias("window_end")
    )

# Write to multiple sinks
query1 = windowed_stats \
    .writeStream \
    .outputMode("update") \
    .format("console") \
    .start()

query2 = parsed_events \
    .writeStream \
    .format("delta") \
    .outputMode("append") \
    .option("path", "/data/events") \
    .start()

query3 = windowed_stats \
    .selectExpr("to_json(struct(*)) as value") \
    .writeStream \
    .format("kafka") \
    .option("kafka.bootstrap.servers", "localhost:9092") \
    .option("topic", "aggregated-events") \
    .option("checkpointLocation", "/checkpoint/kafka") \
    .start()

# Wait for termination
spark.streams.awaitAnyTermination()
```

## Data Quality & Testing

### Great Expectations
```python
import great_expectations as ge
from great_expectations.checkpoint import SimpleCheckpoint
from great_expectations.core.batch import BatchRequest

# Create context
context = ge.get_context()

# Define expectations
def create_expectations():
    """Define data quality expectations"""
    
    # Create expectation suite
    suite = context.create_expectation_suite(
        "orders_validation",
        overwrite_existing=True
    )
    
    # Add expectations
    validator = context.get_validator(
        batch_request=BatchRequest(
            datasource_name="postgres",
            data_connector_name="default",
            data_asset_name="orders"
        ),
        expectation_suite_name="orders_validation"
    )
    
    # Column expectations
    validator.expect_column_to_exist("order_id")
    validator.expect_column_values_to_not_be_null("order_id")
    validator.expect_column_values_to_be_unique("order_id")
    
    validator.expect_column_values_to_be_between(
        "amount",
        min_value=0,
        max_value=100000
    )
    
    validator.expect_column_values_to_be_in_set(
        "status",
        ["pending", "completed", "shipped", "cancelled"]
    )
    
    validator.expect_column_values_to_match_regex(
        "email",
        "^[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+\.[a-zA-Z0-9-.]+$"
    )
    
    # Table expectations
    validator.expect_table_row_count_to_be_between(
        min_value=1000,
        max_value=1000000
    )
    
    validator.expect_table_column_count_to_equal(15)
    
    # Save suite
    validator.save_expectation_suite()
    
    return suite

# Run validation
def validate_data():
    """Run data validation"""
    
    checkpoint = SimpleCheckpoint(
        name="orders_checkpoint",
        data_context=context,
        validations=[
            {
                "batch_request": {
                    "datasource_name": "postgres",
                    "data_connector_name": "default",
                    "data_asset_name": "orders"
                },
                "expectation_suite_name": "orders_validation"
            }
        ]
    )
    
    results = checkpoint.run()
    
    if not results["success"]:
        raise ValueError("Data validation failed")
    
    return results
```

### Data Quality Monitoring
```python
import pandas as pd
from datetime import datetime, timedelta

class DataQualityMonitor:
    def __init__(self, connection_string):
        self.conn = connection_string
        self.metrics = []
    
    def check_freshness(self, table, timestamp_column, threshold_hours=24):
        """Check data freshness"""
        query = f"""
        SELECT MAX({timestamp_column}) as latest_record
        FROM {table}
        """
        
        result = pd.read_sql(query, self.conn)
        latest = result['latest_record'].iloc[0]
        
        if datetime.now() - latest > timedelta(hours=threshold_hours):
            self.alert(f"Data in {table} is stale. Latest record: {latest}")
            return False
        return True
    
    def check_completeness(self, table, required_columns):
        """Check for null values in required columns"""
        issues = []
        
        for column in required_columns:
            query = f"""
            SELECT 
                COUNT(*) as total_rows,
                SUM(CASE WHEN {column} IS NULL THEN 1 ELSE 0 END) as null_count,
                ROUND(100.0 * SUM(CASE WHEN {column} IS NULL THEN 1 ELSE 0 END) / COUNT(*), 2) as null_percentage
            FROM {table}
            """
            
            result = pd.read_sql(query, self.conn)
            
            if result['null_percentage'].iloc[0] > 5:
                issues.append({
                    'column': column,
                    'null_percentage': result['null_percentage'].iloc[0]
                })
        
        return issues
    
    def check_anomalies(self, table, metric_column, threshold_std=3):
        """Detect anomalies using statistical methods"""
        query = f"""
        WITH stats AS (
            SELECT 
                AVG({metric_column}) as mean_val,
                STDDEV({metric_column}) as std_val
            FROM {table}
            WHERE date >= CURRENT_DATE - INTERVAL '30 days'
        )
        SELECT 
            date,
            {metric_column} as value,
            (SELECT mean_val FROM stats) as mean,
            (SELECT std_val FROM stats) as std
        FROM {table}
        WHERE date = CURRENT_DATE
        """
        
        result = pd.read_sql(query, self.conn)
        
        if len(result) > 0:
            current_value = result['value'].iloc[0]
            mean = result['mean'].iloc[0]
            std = result['std'].iloc[0]
            
            z_score = abs((current_value - mean) / std)
            
            if z_score > threshold_std:
                self.alert(f"Anomaly detected in {table}.{metric_column}: Z-score = {z_score}")
                return False
        
        return True
    
    def alert(self, message):
        """Send alert"""
        print(f"⚠️ ALERT: {message}")
        # Send to Slack, PagerDuty, etc.
        
    def generate_report(self):
        """Generate data quality report"""
        report = {
            'timestamp': datetime.now().isoformat(),
            'checks_performed': len(self.metrics),
            'issues_found': sum(1 for m in self.metrics if not m['passed']),
            'details': self.metrics
        }
        
        return report
```

## Orchestration with Dagster

```python
from dagster import (
    asset, 
    op, 
    job, 
    schedule, 
    sensor,
    Output,
    AssetMaterialization,
    RunRequest,
    SkipReason
)
from dagster_dbt import dbt_cli_resource, dbt_run_op, dbt_test_op
from dagster_spark import spark_resource
import pandas as pd

# Define assets
@asset(
    compute_kind="python",
    description="Raw orders data from source system"
)
def raw_orders():
    """Extract raw orders"""
    df = pd.read_csv("s3://bucket/orders.csv")
    df.to_parquet("/data/raw/orders.parquet")
    
    return Output(
        value=df,
        metadata={
            "rows": len(df),
            "columns": list(df.columns),
            "preview": df.head().to_dict()
        }
    )

@asset(
    compute_kind="spark",
    required_resource_keys={"spark"}
)
def processed_orders(context, raw_orders):
    """Process orders with Spark"""
    spark = context.resources.spark
    
    df = spark.read.parquet("/data/raw/orders.parquet")
    
    # Processing logic
    processed = df.filter(df.amount > 0) \
                 .groupBy("product_id") \
                 .sum("amount")
    
    processed.write.mode("overwrite").parquet("/data/processed/orders.parquet")
    
    return AssetMaterialization(
        asset_key="processed_orders",
        metadata={
            "record_count": processed.count()
        }
    )

# Define ops
@op(required_resource_keys={"dbt"})
def run_dbt_models(context):
    """Run dbt models"""
    return dbt_run_op(context)

@op(required_resource_keys={"dbt"})
def test_dbt_models(context):
    """Test dbt models"""
    return dbt_test_op(context)

# Define job
@job(
    resource_defs={
        "dbt": dbt_cli_resource,
        "spark": spark_resource
    }
)
def etl_job():
    """Complete ETL job"""
    test_dbt_models(run_dbt_models())

# Schedule
@schedule(
    cron_schedule="0 2 * * *",
    job=etl_job,
    execution_timezone="UTC"
)
def daily_etl_schedule(context):
    """Daily ETL schedule"""
    return {}

# Sensor
@sensor(job=etl_job)
def file_sensor(context):
    """Trigger job when file arrives"""
    if check_file_exists("s3://bucket/trigger.txt"):
        return RunRequest(run_key=context.cursor)
    else:
        return SkipReason("File not found")
```

## Monitoring & Alerting

### Pipeline Monitoring
```python
import logging
from datadog import initialize, statsd
from prometheus_client import Counter, Histogram, Gauge, push_to_gateway
import time

# Setup monitoring
initialize(
    api_key="YOUR_API_KEY",
    app_key="YOUR_APP_KEY"
)

# Metrics
pipeline_runs = Counter('pipeline_runs_total', 'Total pipeline runs', ['pipeline', 'status'])
pipeline_duration = Histogram('pipeline_duration_seconds', 'Pipeline duration', ['pipeline'])
records_processed = Counter('records_processed_total', 'Records processed', ['pipeline', 'stage'])
data_quality_score = Gauge('data_quality_score', 'Data quality score', ['pipeline'])

class PipelineMonitor:
    def __init__(self, pipeline_name):
        self.pipeline_name = pipeline_name
        self.logger = self._setup_logger()
        
    def _setup_logger(self):
        """Setup structured logging"""
        logger = logging.getLogger(self.pipeline_name)
        handler = logging.StreamHandler()
        formatter = logging.Formatter(
            '{"timestamp":"%(asctime)s","pipeline":"%(name)s","level":"%(levelname)s","message":"%(message)s"}'
        )
        handler.setFormatter(formatter)
        logger.addHandler(handler)
        logger.setLevel(logging.INFO)
        return logger
    
    def start_pipeline(self):
        """Log pipeline start"""
        self.start_time = time.time()
        self.logger.info("Pipeline started")
        statsd.increment(f'{self.pipeline_name}.runs')
        pipeline_runs.labels(pipeline=self.pipeline_name, status='started').inc()
    
    def end_pipeline(self, status='success'):
        """Log pipeline end"""
        duration = time.time() - self.start_time
        
        self.logger.info(f"Pipeline completed: {status}", extra={'duration': duration})
        
        # Send metrics
        statsd.histogram(f'{self.pipeline_name}.duration', duration)
        statsd.increment(f'{self.pipeline_name}.{status}')
        
        pipeline_runs.labels(pipeline=self.pipeline_name, status=status).inc()
        pipeline_duration.labels(pipeline=self.pipeline_name).observe(duration)
        
        # Push to Prometheus
        push_to_gateway('localhost:9091', job=self.pipeline_name)
    
    def log_stage(self, stage, records_count, duration):
        """Log stage metrics"""
        self.logger.info(f"Stage completed: {stage}", extra={
            'records': records_count,
            'duration': duration
        })
        
        records_processed.labels(pipeline=self.pipeline_name, stage=stage).inc(records_count)
        statsd.gauge(f'{self.pipeline_name}.{stage}.records', records_count)
    
    def alert(self, message, severity='warning'):
        """Send alerts"""
        if severity == 'critical':
            self.logger.error(message)
            # Send to PagerDuty, Slack, etc.
        else:
            self.logger.warning(message)
```

## Configuration Management

### Pipeline Configuration
```yaml
# config/pipeline_config.yml
pipeline:
  name: sales_etl
  version: 1.0.0
  owner: data-team
  schedule: "0 2 * * *"
  
sources:
  postgres:
    host: ${DB_HOST}
    port: 5432
    database: sales_db
    table: orders
    connection_pool_size: 10
    
  s3:
    bucket: data-lake
    prefix: raw/orders/
    format: parquet
    
  kafka:
    bootstrap_servers: ${KAFKA_BROKERS}
    topic: orders-stream
    consumer_group: etl-consumer

transformations:
  - name: clean_data
    type: filter
    conditions:
      - amount > 0
      - status != 'cancelled'
      
  - name: add_metrics
    type: calculate
    fields:
      - name: revenue
        formula: quantity * amount
      - name: tax
        formula: revenue * 0.1
        
  - name: aggregate
    type: groupby
    keys: [product_id, date]
    aggregations:
      - sum(revenue)
      - count(order_id)

destinations:
  snowflake:
    account: ${SNOWFLAKE_ACCOUNT}
    warehouse: ETL_WH
    database: ANALYTICS
    schema: PUBLIC
    table: fact_sales
    
  bigquery:
    project: ${GCP_PROJECT}
    dataset: analytics
    table: fact_sales
    
monitoring:
  alerts:
    - type: freshness
      threshold: 24h
      channel: slack
      
    - type: volume
      min_records: 1000
      max_records: 1000000
      channel: email
      
    - type: quality
      min_score: 0.95
      channel: pagerduty

retry:
  max_attempts: 3
  backoff: exponential
  initial_delay: 60s
```

## Deployment

### Docker Configuration
```dockerfile
# Dockerfile for pipeline
FROM apache/airflow:2.5.0-python3.9

USER root

# Install system dependencies
RUN apt-get update && apt-get install -y \
    gcc \
    g++ \
    libpq-dev \
    && rm -rf /var/lib/apt/lists/*

USER airflow

# Install Python packages
COPY requirements.txt /tmp/
RUN pip install --no-cache-dir -r /tmp/requirements.txt

# Copy DAGs and scripts
COPY dags/ /opt/airflow/dags/
COPY plugins/ /opt/airflow/plugins/
COPY scripts/ /opt/airflow/scripts/

# Set environment variables
ENV AIRFLOW__CORE__LOAD_EXAMPLES=False
ENV AIRFLOW__CORE__EXECUTOR=CeleryExecutor

EXPOSE 8080

CMD ["webserver"]
```

### Kubernetes Deployment
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: data-pipeline
spec:
  replicas: 2
  selector:
    matchLabels:
      app: data-pipeline
  template:
    metadata:
      labels:
        app: data-pipeline
    spec:
      containers:
      - name: pipeline
        image: data-pipeline:latest
        env:
        - name: ENVIRONMENT
          value: production
        - name: DB_HOST
          valueFrom:
            secretKeyRef:
              name: pipeline-secrets
              key: db-host
        resources:
          requests:
            memory: "2Gi"
            cpu: "1000m"
          limits:
            memory: "4Gi"
            cpu: "2000m"
```

## Best Practices

1. **Idempotency**: Ensure pipelines can be rerun safely
2. **Incremental Processing**: Process only new/changed data
3. **Error Handling**: Implement retry logic and dead letter queues
4. **Data Validation**: Validate at each stage
5. **Monitoring**: Track metrics and set up alerts
6. **Documentation**: Document data lineage and transformations
7. **Version Control**: Version pipeline code and configurations
8. **Testing**: Unit test transformations, integration test pipelines
9. **Security**: Encrypt sensitive data, use secrets management
10. **Scalability**: Design for horizontal scaling