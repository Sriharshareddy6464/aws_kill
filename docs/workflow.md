# AWS Kill Switch Workflow & Order Enforcement

The application enforces a strict, state-machine driven execution flow. The commands must be run in the following sequence:

```
[Scan] ──> [Plan] ──> [Kill] ──> [Verify]
```

Unless the preceding step has successfully completed and generated the required state artifact, the subsequent command will refuse to execute.

---

## 1. Flow Stages

| Stage | Command | Required Input File | Produced Output File | Description |
| :--- | :--- | :--- | :--- | :--- |
| **1. Scan** | `aws-kill scan` | *None* | `reports/inventory.json` | Scans target AWS account/region/tags for resources and writes them to inventory. |
| **2. Plan** | `aws-kill plan` | `reports/inventory.json` | `reports/plan.json` | Analyzes dependencies and maps out the optimal deletion sequence. |
| **3. Kill** | `aws-kill kill` | `reports/plan.json` | `reports/result.json` | Destroys resources in order, wait, retry, and track execution state. |
| **4. Verify** | `aws-kill verify` | `reports/result.json` | `reports/verification.json` | Post-deletion check confirming no planned resources remain in AWS. |

---

## 2. Enforcement Logic

Each command contains validation guards to prevent executing tasks out of order:

### `aws-kill scan`
*   Cleans up any existing `reports/plan.json`, `reports/result.json`, and `reports/verification.json` to prevent state mismatch or stale plans.
*   Discovers resources and outputs a new `reports/inventory.json`.

### `aws-kill plan`
*   **Guard Check**: Verifies that `reports/inventory.json` exists in the local directory and is not empty.
*   **Failure Behavior**: Aborts with code `1` and outputs:
    `Error: No scan inventory found at reports/inventory.json. Please run 'aws-kill scan' first.`
*   Generates a structured dependency graph and outputs `reports/plan.json`.

### `aws-kill kill`
*   **Guard Check**: Verifies that `reports/plan.json` exists.
*   **Failure Behavior**: Aborts with code `1` and outputs:
    `Error: No execution plan found at reports/plan.json. Please run 'aws-kill plan' first.`
*   Executes resource deletion and saves execution progression into `reports/result.json`.

### `aws-kill verify`
*   **Guard Check**: Verifies that `reports/result.json` exists, indicating a kill has occurred.
*   **Failure Behavior**: Aborts with code `1` and outputs:
    `Error: No kill execution state found at reports/result.json. Please run 'aws-kill kill' first.`
*   Queries AWS one final time to confirm deletions and writes `reports/verification.json`.

---

## 3. Helper Commands

### `aws-kill list`
*   **Action**: Prints a structured view of all 15 supported AWS resources categorized by service domain (Compute, Networking, Load Balancing, Containers, Storage, Database, CDN).
*   **State Constraints**: None. This command can be executed at any time to inspect capability support.

