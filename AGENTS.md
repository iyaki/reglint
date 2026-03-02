# Agent Guidelines

## Specifications

**IMPORTANT:** Before implementing any feature, consult the specifications in `specs/README.md`.

- **Assume NOT implemented.** Many specs describe planned features that may not yet exist in the codebase.
- **Check the codebase first.** Before concluding something is or isn't implemented, search the actual code. Specs describe intent; code describes reality.
- **Use specs as guidance.** When implementing a feature, follow the design patterns, types, and architecture defined in the relevant spec.
- **Spec index:** `specs/README.md` lists all specifications organized by category (core, LLM, security, etc.).

## Commands

## Deployment (Optional)

## Database Migrations (Optional)

## Local Testing

## Architecture

## Code Style

##

- When multiple code paths do similar things with slight variations, create a shared service with a request struct that cpatures the variations, rather than having each caller implemnt its own logic.
