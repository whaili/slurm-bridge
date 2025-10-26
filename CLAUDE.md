# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Slurm Bridge is a Kubernetes scheduler that enables Slurm (HPC workload manager) to work as a Kubernetes scheduler, allowing co-location of traditional HPC workloads and cloud-like workloads within the same Kubernetes cluster. This is a Slinky project developed by SchedMD LLC.

## Development Commands

### Build and Development
- `make help` - Display available commands
- `make build` - Build OCI packages (images + charts)
- `make build-images` - Build container images
- `make build-chart` - Build Helm charts
- `make fmt` - Format Go code
- `make vet` - Run Go vet
- `make test` - Run tests with coverage (70% threshold)
- `make golangci-lint` - Run linting with auto-fix
- `make golangci-lint-fmt` - Run linting (no auto-fix)
- `make govulncheck` - Security vulnerability checking
- `make manifests` - Generate Kubernetes manifests
- `make install-dev` - Install development dependencies

### Helm Commands
- `make helm-validate` - Validate Helm charts
- `make helm-docs` - Generate Helm documentation
- `make helm-lint` - Lint Helm charts
- `make helm-dependency-update` - Update chart dependencies

### Docker Build Commands
- `make push` - Push OCI packages to registry
- `make clean` - Clean build artifacts

## Architecture

### Core Components
1. **Scheduler Plugin** (`cmd/scheduler/`)
   - Kubernetes scheduler that integrates Slurm scheduling capabilities
   - Uses Slurm's external job API to create placeholder jobs
   - Manages pod binding to Slurm-allocated nodes

2. **Admission Webhook** (`cmd/admission/`)
   - Validates and mutates pods before they are scheduled
   - Adds Slurm-specific annotations and configurations

3. **Controllers** (`cmd/controllers/`)
   - Manage Kubernetes resources and Slurm job lifecycle
   - Handle job state synchronization

4. **Slurm Bridge Plugin** (`internal/scheduler/plugins/slurmbridge/`)
   - Core plugin that bridges Kubernetes and Slurm scheduling
   - Manages placeholder job creation and pod binding

### Workflow
1. **Pod Submission**: User submits pod to configured namespace
2. **Admission Webhook**: Pod is validated and modified by slurm-bridge
3. **Scheduling**: slurm-bridge scheduler creates Slurm placeholder job
4. **Resource Allocation**: Slurm schedules the placeholder job to a node
5. **Pod Binding**: Scheduler binds the pod to the allocated node
6. **Pod Execution**: Kubelet launches the pod as normal

## Code Quality

### Linting and Formatting
- Uses golangci-lint with comprehensive rules focusing on security, error handling, and code quality
- Auto-fixes issues when possible with `make golangci-lint`
- Format Go code with `make fmt`
- Run vet with `make vet`

### Testing
- Tests use Ginkgo/Gomega framework
- Minimum coverage threshold: 70%
- Run tests with `make test`
- Test files are co-located with source files

### Security
- Run govulncheck with `make govulncheck`
- Focus on common Go security vulnerabilities

## Project Structure

### Key Directories
- `cmd/` - Main application entry points (scheduler, admission, controllers)
- `internal/` - Internal application code organized by component
  - `scheduler/plugins/slurmbridge/` - Core Slurm Bridge plugin logic
  - `admission/` - Webhook validation and mutation logic
  - `controller/` - Kubernetes controller implementations
  - `config/` - Configuration management
  - `utils/` - Utility functions
- `config/` - Kubernetes manifests and RBAC configurations
- `helm/` - Helm charts for deployment
- `docs/` - Project documentation
- `hack/` - Build and development scripts

## Development Environment

### Tools and Dependencies
- Go 1.24.0+ with modern toolchain
- Kubernetes v1.34+ dependencies via controller-runtime
- Slurm client library v0.4.1
- Helm for packaging
- Docker for container builds
- Kind for local Kubernetes testing

### VS Code Configuration
- Debug configurations available in `.vscode/launch.json`
- Pre-launch tasks for setup and configuration
- Use `make install-dev` to install development dependencies

## Important Notes

### Build Configuration
- Uses Docker buildx for multi-architecture builds (linux/amd64, linux/arm64)
- Default registry: `ghcr.io/slinkyproject`
- Version managed via `VERSION` variable in Makefile

### Kubernetes Integration
- Designed for Kubernetes v1.34+ with specific dependency versions
- Generates manifests via controller-gen
- Implements admission webhook patterns

### Slurm Integration
- Uses Slurm 25.05+ external job API
- Creates placeholder jobs for Kubernetes workloads
- Manages exclusive whole node allocations

## Documentation

### Generated Documentation
- Use `make helm-docs` to generate Helm chart documentation
- Use `make generate-docs` to convert README.md to Sphinx format
- Documentation lives in `docs/` directory with comprehensive guides

### Reading Material
- `docs/architecture.md` - Detailed architecture overview
- `docs/quickstart.md` - Installation and getting started
- `docs/workload.md` - Working with different workload types
- `docs/scheduler.md` - Scheduler configuration
- `docs/controllers.md` - Controller documentation
- `docs/testing.md` - Testing instructions


## Code Architecture
The detailed project understanding documents are organized into five parts under the `docs/claude/` directory:  
- 01-overview.md – Project Overview (directories, responsibilities, build/run methods, external dependencies, newcomer reading order)  
- 02-entrypoint.md – Program Entry & Startup Flow (entry functions, CLI commands, initialization and startup sequence)  
- 03-callchains.md – Core Call Chains (function call tree, key logic explanations, main sequence diagram)  
- 04-modules.md – Module Dependencies & Data Flow (module relationships, data structures, request/response processing, APIs)  
- 05-architecture.md – System Architecture (overall structure, startup flow, key call chains, module dependencies, external systems, configuration)  
When answering any questions related to source code structure, module relationships, or execution flow, **always refer to these five documents first**, and include file paths and function names for clarity.

## Reply Guidelines
- Always reference **file path + function name** when explaining code.
- Use **Mermaid diagrams** for flows, call chains, and module dependencies.
- If context is missing, ask explicitly which files to `/add`.
- Never hallucinate non-existing functions or files.
- Always reply in **Chinese**

## Excluded Paths
- vendor/
- build/
- dist/
- .git/
- third_party/


## 系统设计文档
- 原系统设计文档在docs/ 目录下，阅读并参考