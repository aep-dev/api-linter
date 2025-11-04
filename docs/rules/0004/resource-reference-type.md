---
rule:
  aep: 4
  name: [core, '4', resource-reference-type]
  summary: Resource reference annotations should only apply to strings.
permalink: /4/resource-reference-type
redirect_from:
  - /4/resource-reference-type
---

# Resource annotation presence

This rule enforces that any field with a `(aep.api.field_info).resource_reference`
annotation has a `string` type, as described in [AEP-4][].

## Details

This rule scans all fields with a `(aep.api.field_info).resource_reference` annotation.
If one is found, the type is checked, and the rule complains if the type is
anything other than `string`.

It suggests the removal of the annotation rather than fixing the type, because
what we have observed in real life is that the annotation is usually what is
in error rather than the selected type.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message Book {
  string path = 1;

  // This is not a resource reference; the annotation does not belong.
  Author author = 2 [(aep.api.field_info).resource_reference = "library.googleapis.com/Author"];
}
```

**Correct** code for this rule:

```proto
// Correct.
message Book {
  string path = 1;

  Author author = 2;
}
```

```proto
// Correct.
message Book {
  string path = 1;

  string author = 2 [(aep.api.field_info).resource_reference = "library.googleapis.com/Author"];
}
```

## Disabling

Do not violate this rule; it will break several tools.

[aep-4]: https://aep.dev/4
