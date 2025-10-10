## v0.4.1

### Fixed

- Update kubeVersion parsing to handle provider suffixes (e.g., GKE
  x.y.z-gke.a).
- Fixed image rendering in docs/index.rst.
- Conversion of GHFM admonitions from .md to .rst.

## v0.4.0

### Added

- Added GPU device plugin parsing into a GRES request

### Changed

- Document new slurm-operator Token CR and add it to hack/kind.sh

### Removed

- Removed the need to specify --bridge when running hack/kind.sh
