# codeclimate-to-codequality

A Go command-line tool that converts [CodeClimate](https://github.com/codeclimate/platform/blob/master/spec/analyzers/SPEC.md) JSON output to [GitLab Code Quality](https://docs.gitlab.com/ee/ci/testing/code_quality.html) format.

## Overview

This tool reads CodeClimate JSON output from a file or stdin and outputs GitLab Code Quality JSON format to standard output, making it easy to integrate any CodeClimate-compatible analyzer into GitLab CI/CD pipelines.

## Installation

### From Source

```bash
go install github.com/Omochice/codeclimate-to-codequality@latest
```

### Build Locally

```bash
git clone https://github.com/Omochice/codeclimate-to-codequality.git
cd codeclimate-to-codequality
go build
```

## Usage

The tool requires exactly one argument: a file path or `-` for stdin.

```bash
codeclimate-to-codequality report.json > codequality.json
some-analyzer --format=codeclimate | codeclimate-to-codequality - > codequality.json
```

## CI/CD Integration

### GitLab CI Example

```yaml
analyze:
  stage: test
  script:
    - some-analyzer --format=codeclimate -o codeclimate-report.json
  artifacts:
    paths:
      - codeclimate-report.json
    when: always

codequality:
  stage: test
  image: ghcr.io/omochice/codeclimate-to-codequality:latest
  needs:
    - job: analyze
      artifacts: true
  when: always
  script:
    - codeclimate-to-codequality codeclimate-report.json > codequality.json
  artifacts:
    reports:
      codequality: codequality.json
```

## Format Conversion Details

### Severity Mapping

CodeClimate severity levels are mapped to GitLab severity levels:

- **blocker** -> `critical`
- **critical** -> `critical`
- **major** -> `major`
- **minor** -> `minor`
- **info** -> `info`
- Unknown -> `info`

### Skipped Issues

Issues that lack any of the following required fields are skipped:

- `location.path`
- Line number (from `location.lines.begin` or `location.positions.begin.line`)
- `check_name`
- `description`
- `fingerprint`

## Exit Codes

- `0`: Success
- `1`: Error (invalid JSON, I/O error, etc.)

## Error Handling

- Invalid or missing required fields in CodeClimate issues are skipped
- Error messages are written to standard error
- Empty issue arrays produce valid empty GitLab Code Quality output

## Requirements

- Go 1.25 or later

## License

[zlib](./LICENSE)

