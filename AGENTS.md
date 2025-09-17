# Repository Guidelines

## Project Structure & Module Organization
Core library lives in the root Go package (`opengraph.go`, `intent.go`, `structured.go`) and exposes Open Graph parsing helpers. The CLI sample lives at `ogp/main.go` and compiles to the `ogp` binary; keep it minimal and rely on exported APIs. HTML fixtures for parser tests reside under `test/html`, and vendorised dependencies are locked in `vendor/` for reproducible builds.

## Build, Test, and Development Commands
Run unit and integration tests before sending patches: `go test ./...`. Use `go test -run TestFetch` to target a scenario when iterating. Build the CLI with `go build ./ogp` or install it locally via `go install ./ogp`. Keep module metadata tidy with `go mod tidy` and format sources with `go fmt ./...` before pushing.

## Coding Style & Naming Conventions
Follow standard Go style: tabs for indentation, camel-case identifiers, and exported symbols documented with `// Name ...` comments. Keep files small and cohesive; group parser helpers by concern (e.g. tag handling in `tags.go`). Run `gofmt` or `goimports` on touched files; CI will reject unformatted code. Prefer clear intent-driven names such as `parseVideo` over abbreviations.

## Testing Guidelines
Tests rely on the `github.com/otiai10/mint` BDD helpers, so structure new cases with `When`/`Because` blocks for readability. Store new HTML fixtures beside existing ones in `test/html` and reference them through the in-process HTTP server. Name tests `TestFunction` or `TestType_Method` to keep `go test` discovery straightforward. Maintain coverage—Codecov monitors `main`—and add regression tests whenever fixing a bug.

## Commit & Pull Request Guidelines
Write commit subjects in the imperative mood (“Add helper for absolute URLs”) and keep bodies focused on the observable change. Reference issues or pull requests with `#NN` when applicable. Pull requests should describe the problem, outline the solution, and note any test commands run. Include screenshots or sample JSON output when behaviour changes the CLI.
