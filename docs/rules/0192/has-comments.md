---
rule:
  aep: 192
  name: [core, '0192', has-comments]
  summary: Everything in a proto file must be commented.
permalink: /192/has-comments
redirect_from:
  - /0192/has-comments
---

# Omnipresent comments

This rule enforces that every descriptor in every proto file has a _public_
leading comment, as mandated in [AEP-192][].

## Details

This rule looks at each descriptor in each proto file (exempting oneofs and the
file itself) and complains if no public comment is found _above_ the
descriptor.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
// A representation of a book.
message Book {
  string name = 1;  // No leading comment.
}
```

**Correct** code for this rule:

```proto
// Correct.
// A representation of a book.
message Book {
  // The resource name of the book.
  string name = 1;
}
```

## Disabling

If you need to violate this rule, use a leading comment above the descriptor
(and revel in the irony). Remember to also include an [aep.dev/not-precedent][]
comment explaining why.

```proto
// A representation of a book.
message Book {
  // (-- api-linter: core::0192::has-comments=disabled
  //     aep.dev/not-precedent: We need to do this because reasons. --)
  string name = 1;
}
```

[aep-192]: https://aep.dev/192
[aep.dev/not-precedent]: https://aep.dev/not-precedent
