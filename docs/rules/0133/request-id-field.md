---
rule:
  aep: 133
  name: [core, '0133', request-id-field]
  summary: create methods should have a client-specified ID field.
permalink: /133/request-id-field
redirect_from:
  - /0133/request-id-field
---

# Client-specified IDs

This rule enforces that declarative-friendly create methods include a
client-specified ID field, as mandated in [AEP-133][].

## Details

This rule looks at any `Create` method connected to a resource and complains if
it lacks a client-specified ID (e.g. `string book_id`) field.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message CreateBookRequest {
  string parent = 1 [(aep.api.field_info).resource_reference = {
    child_type: "library.googleapis.com/Book"
  }];

  Book book = 2;

  // A `string book_id` field should exist.
}
```

**Correct** code for this rule:

```proto
// Correct.
message CreateBookRequest {
  string parent = 1 [(aep.api.field_info).resource_reference = {
    child_type: "library.googleapis.com/Book"
  }];

  string id = 2;

  Book book = 3;

}
```

## Disabling

If you need to violate this rule, use a leading comment above the message.
Remember to also include an [aep.dev/not-precedent][] comment explaining why.

```proto
// (-- api-linter: core::0133::request-id-field=disabled
//     aep.dev/not-precedent: We need to do this because reasons. --)
message CreateBookRequest {
  string parent = 1 [(aep.api.field_info).resource_reference = {
    child_type: "library.googleapis.com/Book"
  }];

  Book book = 2;
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aep-133]: https://aep.dev/133
[aep.dev/not-precedent]: https://aep.dev/not-precedent