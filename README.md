# AWS Kill Switch (`aws-kill`)

A lightweight command-line tool that automatically discovers, plans, destroys, and verifies AWS infrastructure created during development.

---

## Workflow Flow

The CLI enforces a strict execution sequence to prevent accidental deletion and ensure consistency:

```
[Scan] ──> [Plan] ──> [Kill] ──> [Verify]
```

1. **`aws-kill scan`**: Discovers all supported resources and saves them to `reports/inventory.json`.
2. **`aws-kill plan`**: Analyzes relationships, maps dependencies, and generates a safe deletion order in `reports/plan.json`.
3. **`aws-kill kill`**: Destroys the planned resources in order, polling until deleted, and saves the outcome to `reports/result.json`.
4. **`aws-kill verify`**: Re-scans AWS post-deletion to confirm all planned resources are deleted and saves the report to `reports/verification.json`.

---

## Directory Structure

```text
aws-kill/
├── main.go                          # Entry point
├── cmd/                             # Cobra CLI Command Handlers
├── engine/                          # Workflow Engines (Scan, Plan, Kill, Verify)
├── aws/                             # AWS Config, Client Sessions, Wait Helpers
├── services/                        # AWS Service-Specific Deletion Logic
├── models/                          # Shared Resource & Inventory Structs
├── reports/                         # Generated JSON Reports
├── utils/                           # Loggers, helpers, JSON utilities
└── docs/                            # Markdown Documentation
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
# 1. Scan your resources (e.g., filtered by target tag)
./aws-kill scan --tag Environment=dev

# 2. Plan the deletion sequence
./aws-kill plan

# 3. Destroy the resources
./aws-kill kill

# 4. Verify everything was removed
./aws-kill verify
```
