# Vision and Goals

## Overview

### Purpose

RegLint is a tool designed to help developers enforce custom code quality rules using regular expressions. By defining rules in a simple YAML format, teams can quickly identify and address issues in their codebase without the overhead of complex static analysis tools.

### Goals

- Ship a frictionless rule-to-result pipeline that turns regex intent into actionable findings in seconds.
- Make CI gates and PR annotations feel automatic by producing consistent, machine-ready outputs.
- Deliver precise, confidence-inspiring findings that are easy to trust and fast to resolve.
- Keep rule sharing simple so teams can enforce standards without extra process.
- Stay fast and predictable so scanning fits naturally into daily workflows.

### Non-Goals

- Compete as a full SAST/DAST platform with deep semantic analysis.
- Provide automatic code fixes or refactoring suggestions.
- Offer hosted scanning, dashboards, or multi-tenant reporting.
- Replace existing linters; the tool complements them.
- Support complex PCRE-only regex features beyond RE2.

### Scope

- Primary use cases:
  - Local scans in a repository.
  - Pre-commit checks.
  - CI gate checks.
  - PR annotations (via reviewdog in a separate repository).
- Platforms: Linux, macOS, Windows.

## Appendices

- PR annotation integration will be developed in a separate repository using reviewdog: https://github.com/reviewdog/reviewdog
