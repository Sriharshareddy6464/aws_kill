# Feature: Dynamic AWS Inventory Status

## Overview

Enhanced the **Scan** and **List** workflow by introducing a dynamic infrastructure summary generated directly from the latest AWS scan.

The scan process now produces two independent reports:

- `reports/inventory.json` — Complete inventory containing raw AWS resource metadata.
- `reports/status.json` — Aggregated infrastructure summary used by the CLI.

This separates resource discovery from presentation and establishes a reusable data layer for future planning and verification phases.

---

## Objectives

- Generate a reusable infrastructure inventory from a single AWS scan.
- Aggregate discovered resources into logical AWS service groups.
- Eliminate unnecessary AWS API calls from the `list` command.
- Present users with a concise overview of their AWS account before planning resource deletion.

---

## Implementation

### Data Model

#### Added

- `models/status.go`

Introduced the following models:

- `ServiceStatus`
- `StatusReport`

These models define the structure of `status.json` and provide a standardized format for aggregated infrastructure summaries.

#### Updated

- `models/resource.go`

Added the `State` property to support resource-specific state tracking (e.g., running and stopped EC2 instances).

---

## Scan Engine

Updated the scan workflow to generate two outputs from a single AWS scan.

### `inventory.json`

Contains complete resource metadata including identifiers, regions, resource types, and states.

Purpose:
- Application data source
- Planning engine input
- Verification reference

### `status.json`

Contains aggregated resource counts grouped by AWS service.

Purpose:
- CLI presentation
- Infrastructure overview
- User-readable summary

---

## Service Layer

Updated all supported AWS service modules to return raw resource metadata and aggregated resource counts.

Current supported services include:
- **EC2** (Instances, Running/Stopped Instances, Elastic IPs, Volumes, Snapshots, Key Pairs, Security Groups, Network Interfaces, Launch Templates, Placement Groups, Dedicated Hosts, Capacity Reservations)
- **VPC** (VPCs, Subnets, Route Tables, Internet Gateways, NAT Gateways)
- **Application Load Balancer** (Load Balancers, Target Groups)
- **ECS** (Clusters, Services, Task Definitions, Running Tasks)
- **ECR** (Repositories, Images)
- **RDS** (DB Instances, DB Snapshots, Subnet Groups)
- **S3** (Buckets)
- **CloudFront** (Distributions)

Each service reports both discovered resources and categorized counts for infrastructure summarization.

---

## Scanner Orchestration

Updated `engine/scanner.go` to:
- Execute all service scanners.
- Aggregate resource counts.
- Group resources under parent AWS services.
- Exclude service groups with zero discovered resources.
- Generate both `inventory.json` and `status.json`.

The scan process remains resilient by continuing execution even if individual AWS service scans encounter permission or API errors.

---

## List Command

The `list` command has been redesigned.

*   **Previous behavior**: Displayed a predefined list of supported services.
*   **Current behavior**: Reads only `reports/status.json` (performs no AWS API requests), displays only services discovered during the latest scan, and presents resource counts in a formatted summary.

This makes the command deterministic, lightweight, and independent of live AWS connectivity.

---

## Architecture

```text
AWS APIs
    │
    ▼
Scan Engine
    │
    ├── reports/inventory.json
    │
    └── reports/status.json
             │
             ▼
        List Command
             │
             ▼
 AWS Infrastructure Summary
```

---

## Validation

### Build Verification

```bash
go build
```
Project compiled successfully without errors.

### Scan Verification

```bash
go run . scan
```
Result:
- Successfully connected to AWS.
- Queried all supported services.
- Generated `reports/inventory.json` and `reports/status.json`.
- Logged each service scan independently.
- Reported total active resources discovered.

### List Verification

```bash
go run . list
```
Result:
- Successfully loaded `status.json`.
- Displayed only discovered AWS services.
- Printed aggregated resource counts grouped by service.
- Correctly calculated total resources without double-counting derived metrics such as:
  - Running Instances
  - Stopped Instances
  - Running Tasks
  - Image Counts

---

## Outcome

The Scan phase is now the single source of truth for infrastructure discovery. Subsequent phases will consume generated reports instead of repeatedly querying AWS.

Current workflow:
```text
scan
 │
 ├── inventory.json
 └── status.json
       │
       ▼
list ──> plan ──> kill ──> verify
```
This enhancement establishes the foundation for implementing the dependency planner while providing users with an immediate and accurate overview of their AWS infrastructure.
