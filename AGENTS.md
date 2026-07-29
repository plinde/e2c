# AGENTS.md — plinde/e2c

Agent instructions for this repository. `CLAUDE.md` is a symlink to this file.

`e2c` is a terminal UI for AWS EC2 ("k9s for EC2"). This checkout is
**`plinde/e2c`, a personal fork** of `nlamirault/e2c`.

## 🚫 Work on the fork only — never touch upstream

**`origin` is `plinde/e2c`. `upstream` is `nlamirault/e2c`, and it is read-only.**

Allowed against `upstream`: `git fetch upstream`, and reading its issues/PRs/code.

**Never**, unless the user explicitly asks for that specific action in that
message:

- `git push upstream ...` — anything at all
- `gh pr create --repo nlamirault/e2c` (or with `--repo` omitted while `upstream`
  is the default, which is the easy way to do this by accident)
- `gh issue create` / `gh issue comment` / any comment or review on an upstream
  issue or PR
- editing or closing anything in the upstream repo

Every PR belongs to `plinde/e2c` — always pass `--repo plinde/e2c` and
`--base main` explicitly. `gh` resolves a bare `gh pr create` against the
*upstream* of a fork, so the explicit flags are what keep an accidental
upstream PR from happening. Merging a PR on this fork is fine when asked.

If something genuinely looks worth upstreaming, **say so and stop** — the user
decides whether anything is ever sent to `nlamirault/e2c`.

## Branch and worktree layout

Work in a linked worktree beside the main checkout, never in the main checkout:

```
~/workspace/github.com/plinde/e2c/                 ← main checkout, keep clean
~/workspace/github.com/plinde/e2c--<description>/  ← linked worktrees
```

Branch off `origin/main` after a fetch, and pull the main checkout forward
(`git pull --ff-only`) once a PR merges.

## Keep fork-local changes isolated in their own commits

Fork-local preferences must not be mixed into commits that would otherwise be
clean upstream changes. Anything that is *only* a local preference gets its own
commit, so the rest stays cherry-pickable if the user ever chooses to offer it.

Live example: theme support (`aeccc32`) keeps upstream's `nord` default and is a
no-behaviour-change feature; flipping `DefaultTheme` to `gruvbox` (`2ca910c`) is
a separate commit, and `DefaultTheme` plus the docs quoting it are the only
lines that differ from an upstreamable version.

## Build, test, install

```bash
go build ./... && go vet ./... && gofmt -l . && go test ./...
make install
```

The Makefile installs to `INSTALL_DIR = $(GOBIN)`, falling back to
`$(GOPATH)/bin`. **Those are not the same directory under asdf** — here `GOBIN`
is `~/.asdf/installs/golang/<ver>/bin` while `GOPATH` is
`~/.asdf/installs/golang/<ver>/packages`, so reaching for `$(go env GOPATH)/bin`
gets you a path that does not exist. Resolve it the way the Makefile does:

```bash
E2C_BIN="$(go env GOBIN)"; [ -n "$E2C_BIN" ] || E2C_BIN="$(go env GOPATH)/bin"
```

⚠ **After every `make install`, re-sign the binary**:

```bash
codesign --force --sign - "$E2C_BIN/e2c"
```

`make install` `cp`s over the existing binary path. macOS keeps the previous
code signature cached for that path, so the new binary fails validation and is
SIGKILLed — **exit 137 with no output whatsoever**, which reads like the binary
silently doing nothing. The real fix is an `rm -f` before the `cp`; until that
lands, re-sign.

## Known upstream bugs (present in `main`, fix here first)

- **`--config` is declared but never wired up** — `cfgFile` is parsed in
  `internal/cmd/root.go` and then unused, so viper never sees it. To exercise an
  alternate config you must override `$HOME` *and* invoke the real binary
  (`$(go env GOPATH)/bin/e2c`) rather than the asdf shim, which needs the true
  `$HOME` itself.
- **Unknown config keys are ignored silently** by viper. Upstream's
  `examples/config.yaml` documented a `ui.theme: dark` key that `UIConfig` never
  had, so setting it looked like it worked and did nothing. When adding a config
  key, add it to the struct *and* validate it — the theme code errors on unknown
  names rather than falling back, deliberately.

## Local context

The user's `/e2c` skill (`~/.agents/skills/e2c/SKILL.md`) documents the
per-AWS-account `e2c-*` shell wrappers built on this binary, and is the place to
update when behaviour visible from those wrappers changes.
