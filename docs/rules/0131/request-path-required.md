---
rule:
  aep: 131
  name: [core, '0131', request-path-required]
  summary: Get RPCs must have a `path` field in the request.
permalink: /131/request-path-required
redirect_from:
  - /0131/request-path-required
---

# Get methods: Name field

This rule enforces that all `Get` standard methods have a `string path` field
in the request message, as mandated in [AEP-131][].

## Details

This rule looks at any message matching `Get*Request` and complains if
the `path` field is missing.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message GetBookRequest {
  string book = 1 [  // Field path should be `path`.
    (aep.api.field_behavior) = FIELD_BEHAVIOR_REQUIRED,
    (aep.api.field_info).resource_reference = "library.googleapis.com/Book"
  ];
}
```

**Correct** code for this rule:

```proto
// Correct.
message GetBookRequest {
  string path = 1 [
    (aep.api.field_behavior) = FIELD_BEHAVIOR_REQUIRED,
    (aep.api.field_info).resource_reference = "library.googleapis.com/Book"
  ];
}
```

## Disabling

If you need to violate this rule, use a leading comment above the message.
Remember to also include an [aep.dev/not-precedent][] comment explaining why.

```proto
// (-- api-linter: core::0131::request-path-required=disabled
//     aep.dev/not-precedent: This is named "book" for historical reasons. --)
message GetBookRequest {
  string book = 1 [
    (aep.api.field_behavior) = FIELD_BEHAVIOR_REQUIRED,
    (aep.api.field_info).resource_reference = "library.googleapis.com/Book"
  ];
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aep-131]: https://aep.dev/131
[aep.dev/not-precedent]: https://aep.dev/not-precedent
