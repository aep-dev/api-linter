---
rule:
  aep: 132
  name: [core, '0132', request-required-fields]
  summary: List RPCs must not have unexpected required fields in the request.
permalink: /132/request-required-fields
redirect_from:
  - /0132/request-required-fields
---

# List methods: Required fields

This rule enforces that all `List` standard methods do not have unexpected
required fields, as mandated in [AEP-132][].

## Details

This rule looks at any message matching `List*Request` and complains if it
comes across any required fields other than:

- `string parent` ([AEP-132][])

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
message ListBooksRequest {
	// The parent, which owns this collection of books.
	// Format: publishers/{publisher}
	string parent = 1 [
	    (aep.api.field_behavior) = FIELD_BEHAVIOR_REQUIRED,
	    (google.api.resource_reference) = {
	  		child_type: "library.googleapis.com/Book"
	    }];

  // Non-standard required field.
  int32 page_size = 2 [(aep.api.field_info).field_behavior = FIELD_BEHAVIOR_REQUIRED]
}
```

**Correct** code for this rule:

```proto
// Correct.
message ListBooksRequest {
	// The parent, which owns this collection of books.
	// Format: publishers/{publisher}
	string parent = 1 [
	    (aep.api.field_behavior) = FIELD_BEHAVIOR_REQUIRED,
	    (google.api.resource_reference) = {
	  		child_type: "library.googleapis.com/Book"
	    }];

  int32 page_size = 2 [(aep.api.field_info).field_behavior = OPTIONAL]
}
```

## Disabling

If you need to violate this rule, use a leading comment above the field.
Remember to also include an [aep.dev/not-precedent][] comment explaining why.

```proto
message ListBooksRequest {
	// The parent, which owns this collection of books.
	// Format: publishers/{publisher}
	string parent = 1 [
	    (aep.api.field_behavior) = FIELD_BEHAVIOR_REQUIRED,
	    (google.api.resource_reference) = {
	  		child_type: "library.googleapis.com/Book"
	    }];

  // (-- api-linter: core::0132::request-required-fields=disabled
  //     aep.dev/not-precedent: We really need this field to be required because
  // reasons. --)
  int32 page_size = 2 [(aep.api.field_info).field_behavior = FIELD_BEHAVIOR_REQUIRED]
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aep-132]: https://aep.dev/132
[aep.dev/not-precedent]: https://aep.dev/not-precedent
