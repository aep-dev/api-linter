---
rule:
  aep: 140
  name: [core, '0140', abbreviations]
  summary: Field names should use common abbreviations.
permalink: /140/abbreviations
redirect_from:
  - /0140/abbreviations
---

# Field names: Abbreviations

This rule enforces that field names use common abbreviations, as mandated in
[AEP-140][].

## Details

This rule checks every descriptor in the proto and complains if the long form
of any of the following words are used instead of the abbreviation:

- configuration
- identifier
- information
- specification
- statistics

## Examples

### Single word method

**Incorrect** code for this rule:

```proto
// Incorrect.
message Book {
  string name = 1;
  string identifier = 2;  // Should be `id`.
}
```

**Correct** code for this rule:

```proto
// Correct.
message Book {
  string name = 1;
  string id = 2;
}
```

## Disabling

If you need to violate this rule, use a leading comment above the method.
Remember to also include an [aep.dev/not-precedent][] comment explaining why.

```proto
// (-- api-linter: core::0140::abbreviations=disabled
//     aep.dev/not-precedent: We need to do this because reasons. --)
message Book {
  string name = 1;
  string identifier = 2;  // Should be `id`.
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aep-140]: https://aep.dev/140
[aep.dev/not-precedent]: https://aep.dev/not-precedent
