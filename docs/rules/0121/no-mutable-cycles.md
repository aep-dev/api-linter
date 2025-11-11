---
rule:
  aep: 121
  name: [core, '0121', no-mutable-cycles]
  summary: Resources must not form a resource reference cycle.
permalink: /121/no-mutable-cycles
redirect_from:
  - /0121/no-mutable-cycles
---

# Resources must not form a resource reference cycle

This rule enforces that resources do not create reference cycles of mutable 
references as mandated in [AEP-121][].

## Details

This rule scans the fields of every resource and ensures that any references to
other resources do not create a mutable cycle between them.

## Examples

**Incorrect** code for this rule:

```proto
message Book {
  option (aep.api.resource) = {
    type: "library.googleapis.com/Book"
    pattern: "books/{book}"
  };

  string name = 1;

  // Incorrect. Creates potential reference cycle.
  string author = 2 [
    (aep.api.field_info).resource_reference = "library.googleapis.com/Author"
  ];
}

message Author {
  option (aep.api.resource) = {
    type: "library.googleapis.com/Author"
    pattern: "authors/{author}"
  };

  string name = 1;

  // Incorrect. Creates potential reference cycle.
  string book = 2 [
    (aep.api.field_info).resource_reference = "library.googleapis.com/Book"
  ];
}
```

**Correct** code for this rule:

```proto
message Book {
  option (aep.api.resource) = {
    type: "library.googleapis.com/Book"
    pattern: "books/{book}"
  };

  string name = 1;

  // Correct because the other reference is OUTPUT_ONLY.
  string author = 2 [
    (aep.api.field_info).resource_reference = "library.googleapis.com/Author"
  ];
}

message Author {
  option (aep.api.resource) = {
    type: "library.googleapis.com/Author"
    pattern: "authors/{author}"
  };

  string name = 1;

  // Correct because an OUTPUT_ONLY reference breaks the mutation cycle.
  string book = 2 [
    (aep.api.field_info).resource_reference = "library.googleapis.com/Book",
    (aep.api.field_behavior) = FIELD_BEHAVIOR_OUTPUT_ONLY
  ];
}
```

## Disabling

If you need to violate this rule, use a leading comment above the service.
Remember to also include an [aep.dev/not-precedent][] comment explaining why.

```proto
message Book {
  option (aep.api.resource) = {
    type: "library.googleapis.com/Book"
    pattern: "books/{book}"
  };

  string name = 1;

  // (-- api-linter: core::0121::no-mutable-cycles=disabled
  //     aep.dev/not-precedent: We need to do this because reasons. --)
  string author = 2 [
    (aep.api.field_info).resource_reference = "library.googleapis.com/Author"
  ];
}

message Author {
  option (aep.api.resource) = {
    type: "library.googleapis.com/Author"
    pattern: "authors/{author}"
  };

  string name = 1;

  // (-- api-linter: core::0121::no-mutable-cycles=disabled
  //     aep.dev/not-precedent: We need to do this because reasons. --)
  string book = 2 [
    (aep.api.field_info).resource_reference = "library.googleapis.com/Book"
  ];
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aep-121]: https://aep.dev/121
[aep.dev/not-precedent]: https://aep.dev/not-precedent
