# Getting Support for KenPanel

KenPanel is free and open-source software built and maintained by Kentralo and community contributors.

## Community Support Channels

- **Documentation & Guides**: [kenpanel.kentralo.com/docs](https://kenpanel.kentralo.com/docs) and [kenpanel.kentralo.com/guides](https://kenpanel.kentralo.com/guides)
- **GitHub Discussions**: [Kentralo/kenpanel Discussions](https://github.com/Kentralo/kenpanel/discussions) for general questions, troubleshooting, and architectural discussions.
- **GitHub Issues**: Use [GitHub Issues](https://github.com/Kentralo/kenpanel/issues) for verified bugs, compatibility reports, and feature requests.

## Self-Diagnostics

Before requesting help, run KenPanel's built-in diagnostics:
```bash
kenpanel-cli doctor run
kenpanel-cli doctor permissions --check
kenpanel-cli doctor ssl --all
```
To generate a sanitized diagnostic bundle for review:
```bash
kenpanel-cli support-bundle --sanitize
```
*Note: Always review the output to verify no private credentials or internal IPs are exposed before sharing.*
