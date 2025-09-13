# Contract: Mock Generation

## Contract Name
mockery-generate-mocks

## Request
- Command: `make mocks` or `mockery --all`
- Inputs: Go interfaces in all modules

## Response
- Output: Generated mock files for each interface
- Location: `mocks/` or `modules/[domain]/mock/`
- Format: Go source files compatible with testify/gomock

## Error Cases
- Interface not found: error message
- Mock generation failed: error message

## Example
```
make mocks
# or
mockery --all --output=mocks/
```

---
