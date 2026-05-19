# Changelog

All notable changes to this project will be documented here.

The format follows [Keep a Changelog](https://keepachangelog.com/en/1.1.0/)
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## v0.1.1 - 2026-05-19

### Changed

- Pin `charm.land/fantasy` requirement to `v0.23.2` to match Crush's
  current dependency baseline. The API surface used by this module is
  identical across `v0.23.2`-`v0.25.0`, so Go's minimum-version selection
  will still pick whatever the consuming module needs.

## v0.1.0 - 2026-05-19

### Added

- Initial portable `AskUserQuestion` implementation as a
  `charm.land/fantasy` `AgentTool`.
- Wire-format types (`Params`, `Question`, `Option`, `Answer`, `Request`)
  matching the Claude Agent SDK exactly.
- `Validate` enforcing 1-4 questions, 2-4 options, non-empty fields, and
  unique labels per question.
- `Format` rendering the canonical tool-result string the Claude Agent
  SDK expects.
- `Resolver` interface plus `StaticResolver`, `ResolverFunc`, and
  `StdinResolver` implementations.
- `NewTool` and `NewToolNamed` constructors for both `ask_user_question`
  (default) and `AskUserQuestion` (SDK-compatible) naming.
- `ErrResolverCancelled` sentinel surfaced as a non-error tool result so
  the conversation stays valid when the user dismisses a request.
- CLI demo in `examples/cli`.
- CI workflow (`build`, `vet`, `test -race`) on push and PR.
- Release workflow that creates a GitHub Release from semver tags.
