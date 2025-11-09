---
rule:
  aep: 121
  name: [core, '0121', resource-must-support-list]
  summary: All resources must have a Standard List method.
permalink: /121/resource-must-support-list
redirect_from:
  - /0121/resource-must-support-list
---

# Resource must support list

This rule enforces that all, non-Singleton resources support the List operation
as mandated in [AEP-121][].

## Details

This rule scans a service for Create, Update, and Get methods for resources
(that are not Singletons), and ensures each one has a List method.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
service Foo {
  // Book has a create, but no List method.
  rpc CreateBook(CreateBookRequest) returns (Book) {};
}

message Book {
  option (aep.api.resource) = {
    type: "library.googleapis.com/Book"
    pattern: "books/{book}"
  };
}
```

**Correct** code for this rule:

```proto
// Correct.
service Foo {
  rpc CreateBook(CreateBookRequest) returns (Book) {};
  rpc ListBooks(ListBookRequest) returns (ListBooksResponse) {};
}

message Book {
  option (aep.api.resource) = {
    type: "library.googleapis.com/Book"
    pattern: "books/{book}"
  };
}
```

## Disabling

If you need to violate this rule, use a leading comment above the service.
Remember to also include an [aep.dev/not-precedent][] comment explaining why.

```proto
// (-- api-linter: core::0121::resource-must-support-list=disabled
//     aep.dev/not-precedent: We need to do this because reasons. --)
service Foo {
  // Book has a create, but no List method.
  rpc CreateBook(CreateBookRequest) returns (Book) {};
}

message Book {
  option (aep.api.resource) = {
    type: "library.googleapis.com/Book"
    pattern: "books/{book}"
  };
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aep-121]: https://aep.dev/121
[aep.dev/not-precedent]: https://aep.dev/not-precedent