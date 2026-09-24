# Security Policy

KenPanel is infrastructure-grade control plane software designed to administer mission-critical servers, databases, and network services. Security is our highest priority.

## Supported Versions

Only the latest release channel and the most recent patch releases receive security updates.

| Version | Supported          |
| ------- | ------------------ |
| 2.3.x   | :white_check_mark: |
| 2.2.x   | :white_check_mark: |
| < 2.2   | :x:                |

---

## Reporting a Vulnerability

**DO NOT report security vulnerabilities via public GitHub issues, discussions, or community forums.**

To report a vulnerability responsibly:

1. **GitHub Private Vulnerability Reporting (Preferred)**:
   Navigate to the [Security Advisory page](https://github.com/Kentralo/kenpanel/security/advisories) on GitHub and click **"Report a vulnerability"**. This initiates a confidential disclosure thread with our core security team.
2. **Security Email**:
   If GitHub Private Vulnerability Reporting is unavailable, email our security team directly at:
   `security@kentralo.com`
   Include the subject line: `[KenPanel Security] Vulnerability Report`

### What to Include in Your Report
- Summary and severity of the vulnerability.
- Affected component(s) (Control plane, Agent, CLI, Transplant, Capsule, RecoveryOS, WebStack Doctor).
- Operating system, architecture, and exact KenPanel build/version.
- Step-by-step reproduction instructions or sanitized Proof of Concept (PoC).
- Potential impact and whether privileges can be escalated.
- **Do not include production credentials, private keys, or customer data.**

### Our Security Response Process
- **Initial Acknowledgment**: Within 24 hours of receipt.
- **Triage & Assessment**: Within 72 hours with priority classification.
- **Fix & Disclosure**: Coordinated patch release within 14 to 30 days depending on severity.
- **Credit**: Vulnerability researchers are credited in our public release notes (unless anonymity is requested).

---

## Security Principles & Architecture Baseline

- **Privilege Separation**: Control-plane web processes run as an unprivileged user (`kenpanel`). All mutating system operations are executed via typed, validated RPC requests to the privileged node daemon (`kenpanel-agent`).
- **No Arbitrary Shell Execution**: The API/UI does not expose unvalidated root shell execution primitives.
- **mTLS Agent Protocol**: Remote agent communication enforces mutual TLS with cryptographic identity checks.
- **Tamper-Evident Audit Logging**: All configuration changes are recorded in an append-only cryptographic ledger with HMAC/SHA-256 chain verification.
- **Zero Remote Telemetry**: KenPanel never transmits logs, configs, source code, or email data to external services.
