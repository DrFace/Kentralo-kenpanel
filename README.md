# KenPanel

> **One panel. Your servers. Your rules.**  
> A free, open-source, self-hosted unified control plane for hosting, servers, applications, containers, virtualization, multi-server fleets, and cloud infrastructure.

[![License](https://img.shields.io/badge/license-AGPL--3.0-blue.svg)](LICENSE)
[![Architecture Spec](https://img.shields.io/badge/spec-v2.3.0-purple.svg)](https://kenpanel.kentralo.com/architecture)
[![Software Status](https://img.shields.io/badge/status-v0.1.0--alpha%20preview-amber.svg)](https://kenpanel.kentralo.com/roadmap)
[![Website](https://img.shields.io/badge/website-kenpanel.kentralo.com-blue)](https://kenpanel.kentralo.com)
[![Build & Test](https://github.com/DrFace/Kentralo-kenpanel/actions/workflows/ci.yml/badge.svg)](https://github.com/DrFace/Kentralo-kenpanel/actions)

> [!IMPORTANT]
> **Project Status: Early Development (Phase 1 In Progress)**  
> The KenPanel v2.3 Architecture Specification is complete. Core hosting functionality, OS adapters, and daemon subsystems are actively under construction in Phase 1 (Q4 2026). This software is not yet recommended for mission-critical production servers. Review the [Public Roadmap](https://kenpanel.kentralo.com/roadmap) for detailed milestone tracking.

---

## What is KenPanel?

KenPanel combines traditional web hosting management (cPanel, Plesk, CyberPanel, aaPanel) with modern container/PaaS orchestration (Coolify, Dokploy, Easypanel) and multi-cloud fleet management into a single, cohesive control plane.

- **100% Free & Open Source**: No licenses, subscriptions, paywalls, or activation servers.
- **Truly Self-Hosted**: Runs entirely on your own infrastructure; no forced dependency on KenPanel Cloud.
- **Privacy & Telemetry**: Zero forced tracking, ads, or analytics in the installed software.
- **Multi-Server Architecture**: Designed to control hundreds of Linux and Windows nodes from a unified dashboard.
- **Privilege Separated**: Web UI runs as an unprivileged user; the privileged node daemon communicates via secure, typed RPC over mTLS or local Unix socket.

---

## First-Class Subsystems

- **Migration Center**:
  - `Transplant`: Universal source-to-target migration engine across 15+ hosting stacks.
  - `Server Clone`: Whole-server and workload replication with automated identity/network remapping.
  - `MailBridge`: Complete email migration preserving IMAP flags, dates, folders, and Sieve filters.
- **Application Portability**:
  - `App Capsule`: Inspect, snapshot, and restore self-contained runtime applications with automated secret redaction.
- **Recovery Center**:
  - `RecoveryOS`: Standalone bootable live ISO for unbootable nodes and offline disaster recovery.
  - `Disaster Kit`: Offline human- and machine-readable recovery runbooks with optional encrypted secrets.
  - `Docker Rescue`: Container and volume reconstructor directly from disk metadata.
  - `Config Vault`: Versioned configuration history with atomic rollback on health check failure.
- **Health & Diagnostics**:
  - `WebStack Doctor`: Evidence-backed dependency tree failure diagnostics.
  - `Guardians`: Specialized monitors for SSL, Permissions, Ports, and Cron jobs.
- **Safe Changes**:
  - `ChangeGuard`: Transactional state machine (`Discover` $\rightarrow$ `Checkpoint` $\rightarrow$ `Apply` $\rightarrow$ `Verify` $\rightarrow$ `Commit`/`Rollback`).
  - `UpgradeLab` & `Compatibility Lab`: Disposable sandbox testing before production upgrades.
- **Intelligence & Performance**:
  - `Server Blueprint`: Reverse-engineers unmanaged servers into declarative manifests.
  - `Dependency Map`: In-depth visualization of domain, port, database, and process dependencies.
  - `Incident Autopsy`: Chronological correlation of system metrics, kernel events, and configuration diffs.
  - `Benchmark Lab`: Conservative, safe synthetic benchmarking for CPU, disk, memory, and network.

---

## Quick Installation

On any supported Linux server (Ubuntu, Debian, AlmaLinux, Rocky Linux):

```bash
curl -fsSL https://kenpanel.kentralo.com/install.sh | sudo bash
```

See the [Official Installation Guide](https://kenpanel.kentralo.com/install) for offline bundles, Windows Server, and unattended deployment flags.

---

## Monorepo Layout

- `cmd/`: Main entrypoints (`kenpanel`, `kenpanel-agent`, `kenpanel-cli`, standalone engines).
- `core/`: Authentication, authorization, database models, job scheduler, event bus, requirement registry.
- `pkg/`: Modular adapters (OS, Service, Database, DNS, Mail, Cloud) and core subsystem engines.
- `ui/`: Modern React 19 + TypeScript management interface.
- `packaging/`: Packaging specifications for DEB, RPM, and systemd units.
- `docs/`: Comprehensive technical architecture documentation.

---

## Community & Contributing

KenPanel is developed in the open by the community:
- [Suggest an Idea or Feature](https://kenpanel.kentralo.com/feedback/ideas)
- [Report a Bug](https://kenpanel.kentralo.com/feedback/bug-report)
- [Report a Compatibility Result](https://kenpanel.kentralo.com/feedback/compatibility-report)
- [Security Policy & Private Vulnerability Reporting](SECURITY.md)
- [Contributing Guide](CONTRIBUTING.md)

---

## License

KenPanel is licensed under the GNU Affero General Public License v3.0 (AGPL-3.0). Client CLI tools and SDKs are licensed under Apache-2.0.
