# CodeClimate Format Specification

This document describes the CodeClimate engine output format that this tool consumes.

## Delimiter

The CodeClimate engine specification defines that each issue is an independent JSON object delimited by null bytes (`\0`).
This is neither a JSON array nor newline-delimited JSON (NDJSON).

```text
{"type":"issue","check_name":"example",...}\0{"type":"issue","check_name":"another",...}\0
```

The trailing `\0` after the last issue is optional.
Some tools (e.g. Brakeman) use `\0` as a separator between issues rather than a terminator.

Reference: [CodeClimate Engine Spec](https://github.com/codeclimate/platform/blob/master/spec/analyzers/SPEC.md)

## Issue Schema

Each JSON object has the following fields.

| Field                | Type     | Required | Description                                                                                                             |
| -------------------- | -------- | -------- | ----------------------------------------------------------------------------------------------------------------------- |
| `type`               | string   | yes      | Must be `"issue"`                                                                                                       |
| `check_name`         | string   | yes      | Unique name of the check                                                                                                |
| `description`        | string   | yes      | Single-line explanation                                                                                                 |
| `content`            | object   | no       | `{"body": "markdown text"}` with details                                                                                |
| `categories`         | string[] | yes      | One or more of: `Bug Risk`, `Clarity`, `Compatibility`, `Complexity`, `Duplication`, `Performance`, `Security`, `Style` |
| `severity`           | string   | no       | `info`, `minor`, `major`, `critical`, or `blocker`                                                                      |
| `fingerprint`        | string   | no       | Unique identifier for deduplication                                                                                     |
| `remediation_points` | integer  | no       | Estimated fix effort                                                                                                    |
| `location`           | object   | yes      | See below                                                                                                               |

## Location

The `location` field supports two variants.

### Line-based

```json
{
  "path": "app/models/user.rb",
  "lines": {
    "begin": 14,
    "end": 14
  }
}
```

### Position-based

```json
{
  "path": "app/models/user.rb",
  "positions": {
    "begin": { "line": 14, "column": 5 },
    "end": { "line": 14, "column": 20 }
  }
}
```

Line numbers are 1-based.

## Difference from GitLab Code Quality

GitLab Code Quality expects a single JSON array as an artifact file, not null-byte delimited objects.
This tool bridges the gap by reading CodeClimate engine output and producing GitLab Code Quality JSON array output.

Reference: [GitLab Code Quality](https://docs.gitlab.com/ci/testing/code_quality/)
