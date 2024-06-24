---
rule:
  aep: 155
  name: [core, '0155', request-id-type]
  summary: request_id should have type aep.api.IdempotencyKey
permalink: /155/request-id-format
redirect_from:
  - /0155/request-id-type
---

# `request_id` type

This rule mandates that all fields named `request_id` have type
aep.api.IdempotencyKey.

## Details

This rule looks on for fields named `request_id` and complains if it does not
have the type aep.api.IdempotencyKey.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message CreateBookRequest {
  string parent = 1;

  Book book = 2;

  string request_id = 3; // wrong type.
}
```

**Correct** code for this rule:

```proto
// Correct.
message CreateBookRequest {
  string parent = 1;

  Book book = 2;

  aep.api.IdempotencyKey request_id = 3;
}
```

## Disabling

If you need to violate this rule, use a leading comment above the field or its
enclosing message. Remember to also include an [aep.dev/not-precedent][]
comment explaining why.

```proto
message CreateBookRequest {
  string parent = 1;

  Book book = 2;

  // (-- api-linter: core::0155::request-id-format=disabled
  //     aep.dev/not-precedent: We need to do this because reasons. --)
  string request_id = 3;
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aep-155]: https://aep.dev/155
[aep.dev/not-precedent]: https://aep.dev/not-precedent