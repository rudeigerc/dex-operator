# Release

## Documentation

We use [CRD Reference Documentation Generator](https://github.com/elastic/crd-ref-docs) to generate documentation for custom resources.

```shell
make api-docs
```

### Configuration

In `hack/config.yaml`:

```yaml
processor:
  ignoreTypes:
    - "List$"
  ignoreFields:
    - "status$"
    - "TypeMeta$"
```
