# AWS Kill Switch (`aws-kill`)

A lightweight command-line tool that automatically discovers, plans, destroys, and verifies AWS infrastructure created during development.

---

## Workflow Flow

The CLI enforces a strict execution sequence to prevent accidental deletion and ensure consistency. The **Scan** phase is the single source of truth:

```
      [Scan] ──> [status.json] ──> [List]
        │
        ▼
[inventory.json] ──> [Plan] ──> [Kill] ──> [Verify]
```

1. **`aws-kill scan`**: Discovers all supported resources via live AWS API calls. Saves raw metadata to `reports/inventory.json` and aggregated counts to `reports/status.json`.
2. **`aws-kill list`**: Reads `reports/status.json` (no live API calls) to display a structured overview of the active services and resource counts.
3. **`aws-kill plan`**: Reads `reports/inventory.json`, analyzes relationships, maps dependencies, and generates a safe deletion order in `reports/plan.json`.
4. **`aws-kill kill`**: Reads `reports/plan.json` and destroys the planned resources in order, polling until deleted, saving outcomes to `reports/result.json`.
5. **`aws-kill verify`**: Re-scans AWS post-deletion to confirm all planned resources are deleted and saves the report to `reports/verification.json`.

---

## Directory Structure

```text
aws-kill/
├── main.go                          # Entry point
├── cmd/                             # Cobra CLI Command Handlers (scan, list, plan, kill, verify)
├── engine/                          # Workflow Engines (Scan, Plan, Kill, Verify)
├── aws/                             # AWS Config, Client Sessions, Wait Helpers
├── services/                        # AWS Service-Specific Deletion Logic (15 services supported)
├── models/                          # Shared Resource, Status, Inventory & Result Structs
├── reports/                         # Generated JSON Reports (inventory.json, status.json, plan.json, etc.)
├── utils/                           # Loggers, helpers, JSON utilities
└── docs/                            # Markdown Documentation & Feature Specs
```

## Setup & Run

### Prerequisites
* Go 1.24+
* AWS Credentials configured locally (e.g. `~/.aws/credentials`)

### Installation
```bash
go build -o aws-kill
```

### Usage
```bash
# 1. Scan your resources (filtered by target tag)
./aws-kill scan --tag Environment=dev

# 2. List the active services and resource counts discovered
./aws-kill list

# 3. Plan the deletion sequence
./aws-kill plan

# 4. Destroy the resources
./aws-kill kill

# 5. Verify everything was removed
./aws-kill verify
```
