# Orc8r Local Docker Setup Notes

This directory contains Docker Compose files for running Magma Orc8r services
locally for development and testing purposes.

## Notes on Local Setup (Docker Compose v2 & Ubuntu 24.04)

Running Orc8r locally using Docker Compose may fail on newer environments
(e.g. Ubuntu 24.04 with Docker Compose v2).

During local setup, the following issues were observed:

- Docker Compose v2 ignores the `version` field (harmless warning)
- Orc8r controller, test, nginx, and fluentd images may fail to build
- Build failures occur due to missing expected build-context paths such as:
  - `configs/`
  - `src/`
  - `gomod/`

Example error:

`
COPY configs /etc/magma/configs
failed to calculate checksum: "/configs": not found
`


Because of these build failures:
- Orc8r services may not fully start
- The NMS (Admin UI) login page may not be reachable locally
- Admin login issues discussed in the community may not be reproducible
  without a fully running Orc8r environment

## Recommendation

- Local Orc8r Docker setup should be treated as **best-effort**
- Contributors working on AGW, linting, documentation, or CI do **not**
  need a locally running Orc8r
- For admin login or NMS testing, a cloud-based or CI-backed Orc8r deployment
  may be more reliable

This documentation reflects observed behavior using official instructions
and is intended to help contributors avoid setup blockers.

