# Project Roadmap

Milestones and upcoming features for `aws-kill`:

## MVP Goals (Current)
- [x] Initial codebase and package structure setup.
- [x] CLI Subcommands implementation.
- [x] Enforce sequential Scan -> Plan -> Kill -> Verify transition constraints.
- [x] Define abstract model schemas for inventories, steps, and results.
- [x] AWS service module interfaces and boilerplate.

## Phase 2: Core Scanning & Dependency Engine
- [ ] Implement AWS Resource Groups Tagging API integration.
- [ ] Build DAG (Directed Acyclic Graph) engine in `engine/planner.go`.
- [ ] Establish topological sort with standard service dependency rules (e.g. Subnet -> VPC).

## Phase 3: Destructive Operations (Kill Engine)
- [ ] Connect AWS service modules to SDK delete APIs.
- [ ] Add parallel worker pool execution with dependency synchronization.
- [ ] Add interactive prompt confirmations and global dry-run checks (`--dry-run`).
- [ ] Implement robust retries with jitter for eventual consistency delays.

## Phase 4: Audit & Post-Verification
- [ ] Complete `Verifier` logic to query resources post-destruction.
- [ ] Generate comprehensive JSON report details in `reports/verification.json`.
