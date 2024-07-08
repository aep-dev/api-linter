---
rule:
  aep: 216
  name: [core, '0216', value-synonyms]
  summary: Enforce common state values.
permalink: /216/value-synonyms
redirect_from:
  - /0216/value-synonyms
---

# Value synonyms

This rule enforces the use of state values enumerated in [AEP-216][] over
common synonyms.

## Details

This rule iterates over enumerations that end in `State` and looks for enum
values that are common synonyms or alternate spellings of the states listed at
the end of [AEP-216][], and suggests the use of the canonical value instead.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
enum State {
  STATE_UNSPECIFIED = 0;
  SUCCESSFUL = 1;  // Should be `SUCCEEDED`.
}
```

**Correct** code for this rule:

```proto
// Correct.
enum State {
  STATE_UNSPECIFIED = 0;
  SUCCEEDED = 1;
}
```

## Disabling

If you need to violate this rule, use a leading comment above the enum value.
Remember to also include an [aep.dev/not-precedent][] comment explaining why.

```proto
// Incorrect.
enum State {
  STATE_UNSPECIFIED = 0;
  // (-- api-linter: core::0216::value-synonyms=disabled
  //     aep.dev/not-precedent: We need to do this because reasons. --)
  SUCCESSFUL = 1;  // Should be `SUCCEEDED`.
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aep-216]: https://aep.dev/216
[aep.dev/not-precedent]: https://aep.dev/not-precedent
