# Security Policy

## Supported Versions

Chessfut is a single rolling release tracked on `main`. Only the latest deployed version is supported with security fixes.

## Reporting a Vulnerability

If you find a security vulnerability, please **do not open a public GitHub issue**.

Instead, email **enisyaman4@gmail.com** with:

- A description of the vulnerability and its potential impact
- Steps to reproduce (or a proof of concept)
- Any relevant logs, requests, or screenshots

You should expect an initial response within a few days. Once confirmed, a fix will be prioritized and a coordinated disclosure timeline agreed on before any public write-up.

## Scope

This policy covers:

- `chessfut-be` — the Go API and its infrastructure (Postgres, Redis, chess.com integration)
- `chessfut-fe` — the Next.js frontend

Out of scope:

- Vulnerabilities in third-party dependencies without a demonstrated exploit path in this codebase (please report those upstream instead)
- The chess.com API itself