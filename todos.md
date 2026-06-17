# go_tools / gopls (fork)

Commands for building and testing gopls against the extended Go fork.
See also `doc/new_features/gopls.md` in the GOROOT fork repo.

## Environment (PowerShell)

```powershell
$env:GOROOT = "C:\Users\pc\MyCodeOtherC\go"
$env:GOPATH = "C:\Users\pc\MyCodeOtherC\gopath"   # must not equal GOROOT
$env:GOEXPERIMENT = "genericmethods"
$env:PATH = "$env:GOPATH\bin;$env:GOROOT\bin;$env:PATH"
```

Rebuild the fork toolchain after pulling compiler/stdlib changes:

```powershell
cd C:\Users\pc\MyCodeOtherC\go\src
.\make.bat
```

## Build gopls

Requires a rebuilt fork GOROOT (`go/parser`, `go/types` with enum/index-operator fixes). Set `GO111MODULE=on` in the gopls module.

```powershell
cd C:\Users\pc\MyCodeOtherC\_LocalRepos\go_tools\gopls
$env:GO111MODULE = "on"
go build -o "$env:GOPATH\bin\gopls.exe" .
```

Point Cursor/VS Code at `$GOPATH\bin\gopls.exe` and set `GOROOT` + `GOEXPERIMENT=genericmethods` on the language server process.

## Tests (fork syntax)

Parse (fork AST nodes, operator func decls, `ForceExpr`) and overload typecheck:

```powershell
cd C:\Users\pc\MyCodeOtherC\_LocalRepos\go_tools\gopls
$env:GO111MODULE = "on"
go test ./internal/cache/parsego/... -run "TestParse" -count=1
go test ./internal/cache/ -run "TestOverload" -count=1
```

Index operator go-to-def unit test:

```powershell
go test ./internal/golang/ -run TestIndexOperatorDefinition -count=1 -v
```

Fork marker tests (enum tokens/symbols, default-arg signatures):

```powershell
go test ./internal/test/marker -run Test/fork -count=1 -v
```

Index-operator go-to-def is covered by the unit test above (not a marker test; SSA analyzers panic on overloaded `[]`).

Broader gopls golang package tests:

```powershell
go test ./internal/golang/... -count=1
```

## Overload signature help marker (existing)

```powershell
go test ./internal/test/marker -run Test/signature/overload -count=1 -v
```
