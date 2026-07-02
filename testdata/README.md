# testdata fixture convention

Fixtures are organized by rule ID, one subdirectory per rule:

```
testdata/
  S001/
    valid.md      # AGENTS.md that passes rule S001
    invalid.md    # AGENTS.md that fails rule S001
  S002/
    valid.md
    invalid.md
```

Every rule must have at least one passing and one failing fixture
(NFR-TEST-01). Subdirectories are added by the change that implements each
rule.
