
# envsubst

A lightweight Go package for substituting environment variables in configuration strings, similar to Docker Compose's `${VAR}` syntax.

## Features

- ✅ Replaces `${VAR}` with the value of the environment variable `VAR`
- ❌ Does **not yet** support default values (`${VAR:-default}`) or required checks (`${VAR:?error}`) — coming soon!

## Installation

```bash
go get github.com/emperorsixpacks/envsubst
