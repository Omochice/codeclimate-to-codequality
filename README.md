# codeclimate-to-codequality

A Go command-line tool that converts [CodeClimate](https://github.com/codeclimate/platform/blob/master/spec/analyzers/SPEC.md) null-byte delimited output to [GitLab Code Quality](https://docs.gitlab.com/ee/ci/testing/code_quality.html) JSON format. It reads from a file or stdin and writes to standard output, making it easy to integrate any CodeClimate-compatible analyzer into GitLab CI/CD pipelines.

## Installation

```bash
go install github.com/Omochice/codeclimate-to-codequality@latest
```

## Usage

The tool requires exactly one argument: a file path or `-` for stdin.

```bash
codeclimate-to-codequality report.codeclimate > codequality.json
some-analyzer --format=codeclimate | codeclimate-to-codequality - > codequality.json
```

## CI/CD Integration

### GitLab CI Example

```yaml
analyze:
  stage: test
  script:
    - some-analyzer --format=codeclimate -o codeclimate-report.codeclimate
  artifacts:
    paths:
      - codeclimate-report.codeclimate
    when: always

codequality:
  stage: test
  image: ghcr.io/omochice/codeclimate-to-codequality:latest
  needs:
    - job: analyze
      artifacts: true
  when: always
  script:
    - codeclimate-to-codequality codeclimate-report.codeclimate > codequality.json
  artifacts:
    reports:
      codequality: codequality.json
```

## Format Conversion Details

### Severity Mapping

CodeClimate severity levels are mapped to GitLab severity levels:

- blocker: `critical`
- critical: `critical`
- major: `major`
- minor: `minor`
- info: `info`
- Unknown: `info`

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

## License

[zlib](./LICENSE)
