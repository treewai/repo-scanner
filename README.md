# Secret Scanner

This project can add repository and detect secret key on each repository.

## Overview

Build a simple code scanning application that detects sensitive keywords in public git repos.
The application must fulfil the following requirements:
- A user can CRUD repositories. A repository contains a name and a link to the repo.
- A user can trigger a scan against a repository.
- A user can view the Security Scan Result ("Result") List

How to do a scan:
- Just keep it simple by iterating the words on the codebase to detect SECRET_KEY findings.
- SECRET_KEY start with prefix public_key || private_key.

The Result entity should have the following properties and be stored in a database of your choice:
- Id: any type of unique id
- Status: "Queued" | "In Progress" | "Success" | "Failure"
- RepositoryName: string
- RepositoryUrl: string
- Findings: JSONB, see [example](example-findings.json)
- QueuedAt: timestamp
- ScanningAt: timestamp
- FinishedAt: timestamp

Wherever you'd have to add something that requires product subscriptions or significant extra time, just mention it in your documentation.

## How to run
`git clone https://github.com/treewai/secret-scanner`

`cd secret-scanner`

`docker-compose up`
Add repository
``

**Remark:**

- Only has OpenAPI 3.0 spec (not swagger UI interface)
- Unit tests (not completed)
- Exported function comment (not completed)
