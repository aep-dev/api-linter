---
rule:
  aep: 4
  name: [core, '4', resource-definition-variables]
  summary: Resource patterns should use consistent variable naming.
permalink: /4/resource-definition-variables
redirect_from:
  - /4/resource-definition-variables
---

# Resource pattern variables

This rule enforces that resource patterns use consistent variable naming
conventions, as described in [AEP-4][].

## Details

This rule scans all files with `aep.api.resource_definition` annotations, and
complains if variables in a `pattern` use camel case, or end in `_id`.

## Examples

**Incorrect** code for this rule:

```proto
import "google/api/resource.proto";

// Incorrect.
option (aep.api.resource_definition) = {
  type: "library.googleapis.com/Book"
  // Should be: publishers/{publisher}/books/{book}
  pattern: "publishers/{publisher_id}/books/{book_id}"
};
```

```proto
import "google/api/resource.proto";

// Incorrect.
option (aep.api.resource_definition) = {
  type: "library.googleapis.com/ElectronicBook"
  // Should be: publishers/{publisher}/electronicBooks/{electronic_book}
  pattern: "publishers/{publisher}/electronicBooks/{electronicBook}"
};
```

**Correct** code for this rule:

```proto
import "google/api/resource.proto";

// Correct.
option (aep.api.resource_definition) = {
  type: "library.googleapis.com/Book"
  pattern: "publishers/{publisher}/books/{book}"
};
```

```proto
import "google/api/resource.proto";

// Correct.
option (aep.api.resource_definition) = {
  type: "library.googleapis.com/ElectronicBook"
  pattern: "publishers/{publisher}/electronicBooks/{electronic_book}"
};
```

## Disabling

If you need to violate this rule, use a leading comment above the annotation.

```proto
import "google/api/resource.proto";

// (-- api-linter: core::4::resource-definition-variables=disabled
//     aep.dev/not-precedent: We need to do this because reasons. --)
option (aep.api.resource_definition) = {
  type: "library.googleapis.com/Book"
  pattern: "publishers/{publisher_id}/books/{book_id}"
};
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aep-4]: http://aep.dev/4
[aep.dev/not-precedent]: https://aep.dev/not-precedent
