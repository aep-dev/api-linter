---
rule:
  aep: 156
  name: [core, '0156', forbidden-methods]
  summary: Singletons must not define Create, or Delete methods.
permalink: /156/forbidden-methods
redirect_from:
  - /0156/forbidden-methods
---

# Singletons: Forbidden methods

This rule enforces that singleton resources do not define `Create`, or `Delete`
methods, as mandated in [AEP-156][].

## Details

This rule looks at any message with a `path` variable in the URI where the path
ends in anything other than `*`. It assumes that this is a method operating on
a singleton resource, and complains if the method is a `Create`, or `Delete`
standard method.

## Examples

**Incorrect** code for this rule:

```proto
// Incorrect.
rpc GetSettings(GetSettingsRequest) returns (Settings) {
  option (google.api.http) = {
    get: "/v1/{path=publishers/*/settings}"
  };
}

rpc UpdateSettings(UpdateSettingsRequest) returns (Settings) {
  option (google.api.http) = {
    patch: "/v1/{settings.path=publishers/*/settings}"
    body: "settings"
  };
}

// This method should not exist. The settings should be implicitly created
// when the publisher is created.
rpc CreateSettings(CreateSettingsRequest) returns (Settings) {
  option (google.api.http) = {
    post: "/v1/{path=publishers/*/settings}"
    body: "settings"
  };
}

// This method should not exist. The settings should always implicitly exist.
rpc DeleteSettings(DeleteSettingsRequest) returns (google.protobuf.Empty) {
  option (google.api.http) = {
    delete: "/v1/{path=publishers/*/settings}"
  };
}
```

**Correct** code for this rule:

```proto
// Correct.
rpc GetSettings(GetSettingsRequest) returns (Settings) {
  option (google.api.http) = {
    get: "/v1/{path=publishers/*/settings}"
  };
}

rpc UpdateSettings(UpdateSettingsRequest) returns (Settings) {
  option (google.api.http) = {
    patch: "/v1/{settings.path=publishers/*/settings}"
    body: "settings"
  };
}
```

## Disabling

If you need to violate this rule, use a leading comment above the method.
Remember to also include an [aep.dev/not-precedent][] comment explaining why.

```proto
rpc GetSettings(GetSettingsRequest) returns (Settings) {
  option (google.api.http) = {
    get: "/v1/{path=publishers/*/settings}"
  };
}

rpc UpdateSettings(UpdateSettingsRequest) returns (Settings) {
  option (google.api.http) = {
    patch: "/v1/{settings.path=publishers/*/settings}"
    body: "settings"
  };
}

// (-- api-linter: core::0156::forbidden-methods=disabled
//     aep.dev/not-precedent: We need to do this because reasons. --)
rpc CreateSettings(CreateSettingsRequest) returns (Settings) {
  option (google.api.http) = {
    post: "/v1/{path=publishers/*/settings}"
    body: "settings"
  };
}

// (-- api-linter: core::0156::forbidden-methods=disabled
//     aep.dev/not-precedent: We need to do this because reasons. --)
rpc DeleteSettings(DeleteSettingsRequest) returns (google.protobuf.Empty) {
  option (google.api.http) = {
    delete: "/v1/{path=publishers/*/settings}"
  };
}
```

If you need to violate this rule for an entire file, place the comment at the
top of the file.

[aep-156]: https://aep.dev/156
[aep.dev/not-precedent]: https://aep.dev/not-precedent
