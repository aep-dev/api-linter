---
rule:
  aep: 132
  name: [core, '0132', response-unknown-fields]
  summary: List RPCs should not have unexpected fields in the response.
permalink: /132/response-unknown-fields
redirect_from:
  - /0132/response-unknown-fields
---

# List methods: Unknown fields (Response)

This rule enforces that all `List` standard methods do not have unexpected
fields, as mandated in [AEP-132][].

## Details

This rule looks at any message matching `List*Response` and complains if it
comes across any fields other than:

- The resource.
- `int32/int64 total_size` ([AEP-132][])
- `string next_page_token` ([AEP-158][])
- `repeated string unavailable` ([AEP-217][])

It only checks field names; it does not validate type correctness or
repeated-ness.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message ListBooksResponse {
  repeated Book results = 1;
  string next_page_token = 2;
  string publisher_id = 3;  // Unrecognized field.
}
```

**Correct** code for this rule:

```proto
// Correct.
message ListBooksResponse {
  repeated Book results = 1;
  string next_page_token = 2;
}
```

## Disabling

If you need to violate this rule, use a leading comment above the field.
Remember to also include an [aep.dev/not-precedent][] comment explaining why.

```proto
message ListBooksResponse {
  repeated Book results = 1;
  string next_page_token = 2;
  // (-- api-linter: core::0132::response-unknown-fields=disabled
  //     aep.dev/not-precedent: We really need this field because reasons. --)
  string publisher_id = 3;
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aep-132]: https://aep.dev/132
[aep-135]: https://aep.dev/135
[aep-157]: https://aep.dev/157
[aep-158]: https://aep.dev/158
[aep.dev/not-precedent]: https://aep.dev/not-precedent
