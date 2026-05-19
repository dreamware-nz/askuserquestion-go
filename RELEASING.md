# Releasing

Releases are cut by pushing a semver tag. The
[`release` workflow](.github/workflows/release.yml) takes it from there:
runs `build`, `vet`, and `test -race`, then creates a GitHub Release.

## Steps

1. Update [`CHANGELOG.md`](./CHANGELOG.md) with a new `## vX.Y.Z` section.
   Group entries under `Added`, `Changed`, `Fixed`, `Removed`,
   `Deprecated`, or `Security`. Commit on `main`.
2. Tag the release commit:

   ```sh
   git tag -a vX.Y.Z -m "vX.Y.Z"
   git push origin vX.Y.Z
   ```

3. The release workflow runs. If `CHANGELOG.md` has a matching
   `## vX.Y.Z` heading the workflow uses that section as the release
   notes; otherwise it falls back to GitHub's auto-generated notes.
4. Go module proxy picks up the new tag within a few minutes. Consumers
   can immediately:

   ```sh
   go get github.com/dreamware-nz/askuserquestion-go@vX.Y.Z
   ```

## Versioning policy

- `v0.x.y` — pre-1.0. Breaking changes can land in any minor bump but
  must be called out in the changelog.
- `v1.0.0` will be cut once the schema and `Resolver` interface are
  considered stable across at least two host integrations (Crush plus one
  other).
- Patch releases (`v0.x.Z`) are reserved for bug fixes and documentation;
  no API changes.
