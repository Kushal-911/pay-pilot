# Repository contribution rules

## Documentation changes

- Whenever documentation is changed, commit the documentation change using a Conventional Commit and push it to GitHub `main` (`origin/main`). This applies whether documentation is changed alone or alongside code.
- Keep the documentation commit focused. Do not include unrelated work.
- Before pushing, confirm the current branch is `main` and the push can fast-forward. Never force-push.
- If `main` has advanced or the push cannot safely fast-forward, fetch/reconcile with the latest `origin/main` and retry. Report any unresolved blocker instead of leaving the documentation change silently unpushed.

## Commits

- Use Conventional Commit format: `<type>(optional-scope): imperative summary`.
- Keep each commit focused; use `!` and a `BREAKING CHANGE:` footer for breaking changes.
