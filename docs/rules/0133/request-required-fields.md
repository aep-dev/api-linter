---
rule:
  aep: 133
  name: [core, '0133', request-required-fields]
  summary: Create RPCs must not have unexpected required fields in the request.
permalink: /133/request-required-fields
redirect_from:
  - /0133/request-required-fields
---

# Create methods: Required fields

This rule enforces that all `Create` standard methods do not have unexpected
required fields, as mandated in [AEP-133][].

## Details

This rule looks at any message matching `Create*Request` and complains if it
comes across any required fields other than:

- `string parent` ([AEP-133][])
- `{Resource} {resource}` ([AEP-133][])
- `string {resource}_id` ([AEP-133][])

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message CreateBookRequest {
  string parent = 1;
  Book book = 2;
  string book_id = 3;
  // Non-standard required field.
  string validate_only = 4 [(aep.api.field_info).field_behavior = FIELD_BEHAVIOR_REQUIRED];
}
```

**Correct** code for this rule:

```proto
// Correct.
message CreateBookRequest {
  string parent = 1;
  Book book = 2;
  string id = 3;
  string validate_only = 4 [(aep.api.field_info).field_behavior = OPTIONAL];
}
```

## Disabling

If you need to violate this rule, use a leading comment above the field.
Remember to also include an [aep.dev/not-precedent][] comment explaining why.

```proto
message CreateBookRequest {
  string parent = 1;
  Book book = 2;
  string id = 3;

  // (-- api-linter: core::0133::request-required-fields=disabled
  //     aep.dev/not-precedent: We really need this field to be required because
  // reasons. --)
  string validate_only = 4 [(aep.api.field_info).field_behavior = FIELD_BEHAVIOR_REQUIRED];
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aep-133]: https://aep.dev/133
[aep.dev/not-precedent]: https://aep.dev/not-precedent
