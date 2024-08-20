---
rule:
  aep: 158
  name: [core, '0158', response-repeated-field]
  summary: |
    One field of a paginated RPCs' response
    should be repeated.
permalink: /158/response-repeated-field
redirect_from:
  - /0158/response-repeated-field
---

# Paginated methods: Page token field

This rule enforces that all `List` and `Search` methods have a repeatable field
in the response message, as mandated in [AEP-158][].

## Details

This rule looks at any message matching `List*Response` or `Search*Response`
that has `next_page_token` field and complains if there does not exist a field that is repeated.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message ListBooksResponse {
  Book books = 1;  // Field should be repeated.
  string next_page_token = 2;
}
```

**Correct** code for this rule:

```proto
// Correct.
message ListBooksResponse {
  repeated Book books = 1;
  string next_page_token = 2;
}
```

## Disabling

If you need to violate this rule, use a leading comment above the message or
above the field. Remember to also include an [aep.dev/not-precedent][] comment
explaining why.

```proto
// (-- api-linter: core::0158::response-repeated-first-field=disabled
//     aep.dev/not-precedent: We need to do this because reasons. --)
message ListBooksResponse {
    Book books = 1;
    string next_page_token = 2;
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aep-158]: https://aep.dev/158
[aep.dev/not-precedent]: https://aep.dev/not-precedent
