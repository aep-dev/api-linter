---
rule:
  aep: 122
  name: [core, '0122', resource-reference-type]
  summary: All resource references must be strings.
permalink: /122/resource-reference-type
redirect_from:
  - /0122/resource-reference-type
---

# Resource reference type

This rule enforces that all fields with the `(aep.api.field_info).resource_reference`
annotation are strings, as mandated in [AEP-122][].

## Details

This rule complains if it sees a field with a `(aep.api.field_info).resource_reference`
that has a type other than `string`.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message Book {
  string path = 1;

  // Resource references should be strings.
  Author author = 2 [(aep.api.field_info).resource_reference = "library.googleapis.com/Author"];
}
```

**Correct** code for this rule:

```proto
// Correct.
message Book {
  string path = 1;

  string author = 2 [(aep.api.field_info).resource_reference = "library.googleapis.com/Author"];
}
```

```proto
// Correct.
message Book {
  string path = 1;

  // If "author" is not a first-class resource, then it may be a composite
  // field within the book.
  Author author = 2;
}
```

## Disabling

If you need to violate this rule, use a leading comment above the method.
Remember to also include an [aep.dev/not-precedent][] comment explaining why.

```proto
message Book {
  string path = 1;

  // (-- api-linter: core::0122::resource-reference-type=disabled
  //     aep.dev/not-precedent: We need to do this because reasons. --)
  Author author = 2 [(aep.api.field_info).resource_reference = "library.googleapis.com/Author"];
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aep-122]: https://aep.dev/122
[aep.dev/not-precedent]: https://aep.dev/not-precedent
