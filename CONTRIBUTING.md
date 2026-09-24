# Contributing to KenPanel

Thank you for your interest in contributing to KenPanel! As a 100% free, open-source, and community-driven project, your contributions help shape the future of modern server and infrastructure management.

## Code of Conduct

All contributors and maintainers are expected to adhere to our [Code of Conduct](CODE_OF_CONDUCT.md). Please read it before participating.

---

## How Can I Contribute?

### 1. Reporting Bugs
- Verify the bug hasn't already been reported in [GitHub Issues](https://github.com/Kentralo/kenpanel/issues).
- Use the structured [Bug Report Template](https://github.com/Kentralo/kenpanel/issues/new?template=bug_report.yml).
- Include sanitized logs, OS details, reproduction steps, and expected vs. actual behavior.

### 2. Suggesting Features & Integrations
- We welcome feature suggestions! Please submit them using our [Feature Request Template](https://github.com/Kentralo/kenpanel/issues/new?template=feature_request.yml) or [Integration Request Template](https://github.com/Kentralo/kenpanel/issues/new?template=integration_request.yml).
- Explain the problem, proposed workflow, and target operating systems/providers.

### 3. Reporting OS & Software Compatibility
- If you have verified KenPanel or its adapters on a specific Linux distribution, kernel version, or cloud provider, submit a report using the [Compatibility Report Template](https://github.com/Kentralo/kenpanel/issues/new?template=compatibility_report.yml).

### 4. Improving Documentation & Guides
- Documentation lives in `docs/` and on [kenpanel.kentralo.com/guides](https://kenpanel.kentralo.com/guides).
- Fix typos, enhance runbook steps, or submit new task-based guides via pull request.

---

## Development Setup

### Prerequisites
- Go 1.23+
- Node.js 20+ (for UI development)
- Docker / Podman (for integration tests)
- `make`

### Building from Source
```bash
git clone https://github.com/Kentralo/kenpanel.git
cd kenpanel
make build
```
Binaries will be placed in `bin/`:
- `bin/kenpanel` (Control Plane)
- `bin/kenpanel-agent` (Node Agent)
- `bin/kenpanel-cli` (Operator CLI)

### Running Tests
```bash
make test
```

---

## Pull Request Guidelines

1. **Fork & Branch**: Create a feature branch from `main` (e.g., `feat/nginx-custom-log-format` or `fix/transplant-quota-check`).
2. **Coding Standards**:
   - Write clean, idiomatic Go with error handling and unit tests.
   - For UI, use TypeScript and responsive vanilla CSS components.
   - Do not commit secrets, tokens, or private environment files.
3. **Commit Messages**: Follow standard conventional commits: `feat:`, `fix:`, `docs:`, `test:`, `refactor:`.
4. **Documentation**: Update relevant documentation and OpenAPI definitions for API changes.
5. **PR Template**: Complete all sections of the [Pull Request Template](.github/PULL_REQUEST_TEMPLATE.md).
