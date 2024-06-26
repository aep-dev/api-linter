---
rule:
  aep: 126
  name: [core, '0126', unspecified]
  summary: All enums must have a default unspecified value.
permalink: /126/unspecified
redirect_from:
  - /0126/unspecified
---

# Enum unspecified value

This rule enforces that all enums have a default unspecified value, as mandated
in [AEP-126][].

Because our APIs create automatically-generated client libraries, we need to
consider languages that have varying behavior around default values. To avoid
any ambiguity or confusion across languages, all enumerations should use an
"unspecified" value beginning with the name of the enum itself as the first
(`0`) value.

## Details

This rule finds all enumerations and ensures that the first one is named after
the enum itself with an `_UNSPECIFIED` suffix appended.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
enum Format {
  HARDCOVER = 0;  // Should have "FORMAT_UNSPECIFIED" first.
}
```

```proto
// Incorrect.
enum Format {
  UNSPECIFIED = 0;  // Should be "FORMAT_UNSPECIFIED".
  HARDCOVER = 1;
}
```

**Correct** code for this rule:

```proto
// Correct.
enum Format {
  FORMAT_UNSPECIFIED = 0;
  HARDCOVER = 1;
}
```

## Disabling

If you need to violate this rule, use a leading comment above the enum value.
Remember to also include an [aep.dev/not-precedent][] comment explaining why.

```proto
enum Format {
  // (-- api-linter: core::0126::unspecified=disabled
  //     aep.dev/not-precedent: We need to do this because reasons. --)
  HARDCOVER = 0;
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aep-126]: https://aep.dev/126
[aep.dev/not-precedent]: https://aep.dev/not-precedent
