---
rule:
  aep: 132
  name: [core, '0132', request-unknown-fields]
  summary: List RPCs should not have unexpected fields in the request.
permalink: /132/request-unknown-fields
redirect_from:
  - /0132/request-unknown-fields
---

# List methods: Unknown fields (Request)

This rule enforces that all `List` standard methods do not have unexpected
fields, as mandated in [AEP-132][].

## Details

This rule looks at any message matching `List*Request` and complains if it
comes across any fields other than:

- `string parent` ([AEP-132][])
- `int32 max_page_size` ([AEP-158][])
- `string page_token` ([AEP-158][])
- `int32 skip` ([AEP-158][])
- `string filter` ([AEP-132][])
- `string order_by` ([AEP-132][])
- `bool show_deleted` ([AEP-132][])
- `string request_id` ([AEP-155][])
- `google.protobuf.FieldMask read_mask` ([AEP-157][])
- `View view` ([AEP-157][])

It only checks field names; it does not validate type correctness. This is
handled by other rules, such as
[request field types](./0132-request-field-types.md).

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message ListBooksRequest {
  string parent = 1;
  int32  max_page_size = 2;
  string page_token = 3;
  string library_id = 4;  // Non-standard field.
}
```

**Correct** code for this rule:

```proto
// Correct.
message ListBooksRequest {
  string parent = 1;
  int32  max_page_size = 2;
  string page_token = 3;
}
```

## Disabling

If you need to violate this rule, use a leading comment above the field.
Remember to also include an [aep.dev/not-precedent][] comment explaining why.

```proto
message ListBooksRequest {
  string parent = 1;
  int32  max_page_size = 2;
  string page_token = 3;

  // (-- api-linter: core::0132::request-unknown-fields=disabled
  //     aep.dev/not-precedent: We really need this field because reasons. --)
  string library_id = 4;
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aep-132]: https://aep.dev/132
[aep-135]: https://aep.dev/135
[aep-155]: https://aep.dev/155
[aep-157]: https://aep.dev/157
[aep-158]: https://aep.dev/158
[aep.dev/not-precedent]: https://aep.dev/not-precedent
