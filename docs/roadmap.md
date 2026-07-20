# Project Roadmap

Milestones and upcoming features for `aws-kill`:

## MVP Goals (Completed)
- [x] Initial codebase and package structure setup.
- [x] CLI Subcommands implementation (`scan`, `list`, `plan`, `kill`, `verify`).
- [x] Enforce sequential Scan -> List -> Plan -> Kill -> Verify transition constraints.
- [x] Integrate live AWS SDK v2 resource discovery queries for all 15 services.
- [x] Implement structured status summary report (`reports/status.json`) from a single scan.
- [x] Build local-only `list` command to print grouped resource counts from `status.json`.

## Phase 2: Core Dependency Planning Engine
- [ ] Build DAG (Directed Acyclic Graph) engine in `engine/planner.go`.
- [ ] Establish topological sort with standard service dependency rules (e.g. Subnet -> VPC).
- [ ] Implement support for tag-based filtering (`--tag`) within the scanner resource aggregation.

## Phase 3: Destructive Operations (Kill Engine)
- [ ] Connect AWS service modules to SDK delete APIs.
- [ ] Add parallel worker pool execution with dependency synchronization.
- [ ] Add interactive prompt confirmations and global dry-run checks (`--dry-run`).
- [ ] Implement robust retries with jitter for eventual consistency delays.

## Phase 4: Audit & Post-Verification
- [ ] Complete `Verifier` logic to query resources post-destruction.
- [ ] Generate comprehensive JSON report details in `reports/verification.json`.
