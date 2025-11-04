---
rule:
  aep: 164
  name: [core, '0164', request-unknown-fields]
  summary: Undelete RPCs should not have unexpected fields in the request.
permalink: /164/request-unknown-fields
redirect_from:
  - /0164/request-unknown-fields
---

# Undelete methods: Unknown fields

This rule enforces that all `Undelete` requests do not have unexpected
fields, as mandated in [AEP-164][].

## Details

This rule looks at any message matching `Undelete*Request` and complains if it
comes across any fields other than:

- `string name` ([AEP-164][])
- `string etag` ([AEP-154][])
- `string request_id` ([AEP-155][])
- `bool validate_only` ([AEP-163][])

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message UndeleteBookRequest {
  string name = 1 [
    (aep.api.field_behavior) = FIELD_BEHAVIOR_REQUIRED,
    (aep.api.field_info).resource_reference = "library.googleapis.com/Book",
  ];
  string library_id = 2;  // Non-standard field.
}
```

**Correct** code for this rule:

```proto
// Correct.
message UndeleteBookRequest {
  string name = 1 [
    (aep.api.field_behavior) = FIELD_BEHAVIOR_REQUIRED,
    (aep.api.field_info).resource_reference = "library.googleapis.com/Book",
  ];
}
```

## Disabling

If you need to violate this rule, use a leading comment above the field.
Remember to also include an [aep.dev/not-precedent][] comment explaining why.

```proto
message UndeleteBookRequest {
  string name = 1 [
    (aep.api.field_behavior) = FIELD_BEHAVIOR_REQUIRED,
    (aep.api.field_info).resource_reference = "library.googleapis.com/Book",
  ];

  // (-- api-linter: core::0164::request-unknown-fields=disabled
  //     aep.dev/not-precedent: We really need this field because reasons. --)
  string library_id = 2;
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aep-154]: https://aep.dev/154
[aep-155]: https://aep.dev/155
[aep-163]: https://aep.dev/163
[aep-164]: https://aep.dev/164
[aep.dev/not-precedent]: https://aep.dev/not-precedent
