# Architecture Design

The `aws-kill` CLI tool follows a clean layer architecture:

```
CLI Command Layer (Cobra/Viper)
      │
      ▼
Workflow Engines (Scanner, Planner, Killer, Verifier)
      │
      ▼
Unified AWS Client Registry (AWS Session)
      │
      ▼
AWS Service Modules (EC2, VPC, RDS, S3, ECS, etc.)
      │
      ▼
AWS SDK for Go v2
```

## Layers Overview

1. **CLI Layer (`cmd/`)**: Manages inputs, environment configuration, command-line arguments, global flags, and structured output presentation.
   - `list.go`: Parses `status.json` and prints the human-readable summary. Runs locally without calling any network APIs.
   - `scan.go`, `plan.go`, `kill.go`, `verify.go`: Subcommand interfaces with sequence verification guards.
2. **Engine Layer (`engine/`)**: Contains orchestrator engines for each phase:
   - `Scanner`: Discovers infrastructure and maps raw data to standard Resource models, while calculating the status summary counts.
   - `Planner`: Analyzes references and compiles a Directed Acyclic Graph (DAG) for dependency resolution.
   - `Killer`: Traverses the planned steps in reverse topological order, calling service-specific deletion endpoints, polling status with waiters, and maintaining run states.
   - `Verifier`: Double-checks the environment to ensure zero surviving target resources.
3. **AWS Registry (`aws/`)**: Centralizes the AWS session configuration, credentials profile loading, retry limits, backoff middleware, and client wrapper handles.
4. **Service Logic (`services/`)**: Decentralized, modular code per AWS resource type (e.g., `s3.go`, `rds.go`, `ec2.go`) implementing concrete query, delete and verify endpoints.
5. **Models (`models/`)**: Standardizes interfaces, structures, and JSON serialization schemas for cross-layer data exchange.
   - `status.go`: Defines the structure for `status.json` containing `ServiceStatus` and `StatusReport`.
   - `resource.go`: Represents an AWS resource with metadata and its current `State`.

## Data Storage & Intermediate Reports

All commands transfer state between steps by writing local JSON files under the `reports/` folder:

*   `reports/inventory.json`: Raw discovered resources output by `scan`.
*   `reports/status.json`: Discovered resources count summary output by `scan` and read by `list`.
*   `reports/plan.json`: Safe sequential deletion order plan output by `plan`.
*   `reports/result.json`: Execution status outcome report output by `kill`.
*   `reports/verification.json`: Post-cleanup audit check output by `verify`.
