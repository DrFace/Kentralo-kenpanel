# KenPanel Development Roadmap

> This roadmap reflects our phased engineering milestones for the full KenPanel unified control plane.  
> Statuses: `Completed` | `In Progress` | `Planned`

---

## Phase 0: Foundations & Architecture Baseline (:white_check_mark: Completed)
- [x] Canonical monorepo structure (`DrFace/Kentralo-kenpanel`) and private static website (`kenpanel.kentralo.com`).
- [x] Machine-readable requirement and capability registry schema.
- [x] Domain models for Servers, Organizations, Websites, Databases, DNS, and Mail.
- [x] Secure mTLS agent communication protocol and Unix socket privilege separation.
- [x] Tamper-evident cryptographic audit log model.
- [x] ChangeGuard transactional state machine baseline.
- [x] GitHub Community Intake (Issue forms, PR template, SECURITY.md, Code of Conduct).

---

## Phase 1: Core Host, WebStack & Essential Hosting (:rocket: In Progress)
- [ ] Privileged node agent bootstrap and automated hardware/OS capability detection.
- [ ] OS Adapters: Ubuntu 22.04/24.04, Debian 12, AlmaLinux 9, Rocky Linux 9.
- [ ] Web Server Adapters: Nginx, Apache HTTP Server, Caddy.
- [ ] PHP Runtime Manager: Multiple concurrent PHP-FPM versions (7.4, 8.1, 8.2, 8.3, 8.4) with extension selector.
- [ ] Database Engines: MariaDB/MySQL and PostgreSQL instance provisioning and user/grant manager.
- [ ] SSL/TLS Automation: ACME / Let's Encrypt automated challenge validation (HTTP-01 & DNS-01) and custom certificate upload.
- [ ] WebStack Doctor v1: Automated diagnostic checks for DNS, proxy upstream, and socket connectivity.

---

## Phase 2: Migration Center & Application Portability
- [ ] **Transplant Engine**: Read-only discovery adapters for cPanel, Plesk, and CyberPanel.
- [ ] **Normalized Migration Manifest (`ManifestV1`)**: Serialization, disk preflight, dry-run simulation, and resumable streaming.
- [ ] **Server Clone Mode**: Full node replication with automated machine-id, hostname, SSH host key, and network remapping.
- [ ] **MailBridge**: IMAP synchronization preserving folders, flags, and internal dates, with DNS cutover assistant.
- [ ] **App Capsule**: Application packaging into `.capsule.tar.zst` with automated secret redaction.

---

## Phase 3: PaaS, Containers & Modern Developer Experience
- [ ] Git repository integration with automatic webhook triggers.
- [ ] Nixpacks / Dockerfile automated build pipeline.
- [ ] Docker & Docker Compose stack orchestrator.
- [ ] Application marketplace / one-click templates (WordPress Toolkit, Nextcloud, Discourse, Ghost).
- [ ] In-browser Web Terminal with strict authorization and session recording.

---

## Phase 4: Recovery Center & Resilience
- [ ] **RecoveryOS**: Bootable minimal live ISO with read-only root and filesystem discovery (ext4, XFS, Btrfs, ZFS).
- [ ] **Disaster Kit**: Self-contained offline HTML runbook and encrypted secrets export.
- [ ] **Docker Rescue**: Cold-container state excavator from dead dockerd daemon directories.
- [ ] **Config Vault**: Versioned configuration snapshots with automatic health-failure rollback.

---

## Phase 5: Multi-Server, Fleet & Cloud Infrastructure
- [ ] Multi-server node enrollment with automated token and fingerprint exchange.
- [ ] Multi-cloud provider adapters (AWS EC2/Route53/S3, Hetzner Cloud, DigitalOcean, Linode).
- [ ] Centralized cross-node metrics, unified logs, and health alerting.
- [ ] Incident Autopsy: Multi-node temporal log, metric, and kernel event correlation.
- [ ] Benchmark Lab: Safe CPU, disk, memory, and network performance verification.
