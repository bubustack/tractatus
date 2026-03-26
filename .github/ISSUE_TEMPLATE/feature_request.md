---
name: "Feature request"
about: "Propose a new transport contract, helper API, or generation capability"
labels: ["kind/feature"]
---

## Problem statement
What contract or tooling gap are you trying to solve? Include latency, interoperability, or backward-compatibility constraints if relevant.

## Proposed change
Describe the behavior you’d like to see. If this affects `.proto` definitions, generated APIs, or helper packages, list the exact messages/fields/methods involved.

```
message ExampleRequest {
  string new_field = 1;
}
```

## Affected area
- [ ] Proto schema
- [ ] Generated Go bindings
- [ ] Release artifacts (Go / TypeScript / Python)
- [ ] Envelope helper package
- [ ] Transport stream constants
- [ ] Documentation

## Alternatives considered
What did you try already? Examples: metadata conventions, envelope payloads, custom downstream wrappers, or repo-local forks.

## Compatibility plan
- Can this ship in `transport.v1` without breaking existing consumers?
- If not, what coordinated rollout or versioning plan do you propose?

## Additional context
Links, design docs, screenshots, or related issues/discussions.
