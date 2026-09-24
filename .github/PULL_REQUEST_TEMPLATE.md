## Description
Provide a concise explanation of what this pull request changes and why.

Fixes #(issue number)

## Type of Change
- [ ] Bug fix (non-breaking change which fixes an issue)
- [ ] New feature (non-breaking change adding functionality)
- [ ] Subsystem / Engine improvement (Transplant, Capsule, Recovery, Doctor, ChangeGuard)
- [ ] Breaking change (fix or feature that causes existing functionality to change)
- [ ] Documentation / Guide update
- [ ] Security hardening

## Architecture & Safety Checklist
- [ ] Tested on target operating system(s) (Ubuntu / Debian / RHEL-compatible / Windows Server).
- [ ] Code follows KenPanel privilege separation: web UI/API does not run unrestricted root shell commands.
- [ ] Parameterized database queries used throughout (no raw SQL interpolation).
- [ ] No secrets, keys, or sensitive customer data committed.
- [ ] Unit tests added / updated and passing (`make test`).
- [ ] Requirement Registry updated if applicable (`requirements.json`).
- [ ] Relevant documentation or task runbook guide added / updated.

## Verification & Screenshots
Describe the exact test steps taken to verify the change. Attach terminal output or UI screenshots if applicable.
