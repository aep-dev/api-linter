---
rule:
  aep: 131
  name: [core, '0131', request-unknown-fields]
  summary: Get RPCs should not have unexpected fields in the request.
permalink: /131/request-unknown-fields
redirect_from:
  - /0131/request-unknown-fields
---

# Get methods: Unknown fields

This rule enforces that all `Get` standard methods do not have unexpected
fields, as mandated in [AEP-131][].

## Details

This rule looks at any message matching `Get*Request` and complains if it comes
across any fields other than:

- `string path` ([AEP-131][])
- `string request_id` ([AEP-155][])
- `google.protobuf.FieldMask read_mask` ([AEP-157][])
- `View view` ([AEP-157][])

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message GetBookRequest {
  string path = 1;
  string library_id = 2;  // Non-standard field.
}
```

**Correct** code for this rule:

```proto
// Correct.
message GetBookRequest {
  string path = 1;
}
```

## Disabling

If you need to violate this rule, use a leading comment above the field.
Remember to also include an [aep.dev/not-precedent][] comment explaining why.

```proto
message GetBookRequest {
  string path = 1;

  // (-- api-linter: core::0131::request-unknown-fields=disabled
  //     aep.dev/not-precedent: We really need this field because reasons. --)
  string library_id = 2;
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aep-131]: https://aep.dev/131
[aep-155]: https://aep.dev/155
[aep-157]: https://aep.dev/157
[aep.dev/not-precedent]: https://aep.dev/not-precedent
