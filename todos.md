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

```powershell
cd C:\Users\pc\MyCodeOtherC\_LocalRepos\go_tools\gopls
go build -o "$env:GOPATH\bin\gopls.exe" .
```

Point Cursor/VS Code at `$GOPATH\bin\gopls.exe` and set `GOROOT` + `GOEXPERIMENT=genericmethods` on the language server process.

## Tests (fork syntax)

Parse and overload typecheck:

```powershell
cd C:\Users\pc\MyCodeOtherC\_LocalRepos\go_tools\gopls
go test ./internal/cache/parsego/... ./internal/cache/ -run "TestOverload" -count=1
```

Index operator go-to-def unit test:

```powershell
go test ./internal/golang/ -run TestIndexOperatorDefinition -count=1 -v
```

Fork marker tests (enum tokens/symbols, index `[]` go-to-def, default-arg signatures):

```powershell
go test ./internal/test/marker -run Test/fork -count=1 -v
```

Broader gopls golang package tests:

```powershell
go test ./internal/golang/... -count=1
```

## Overload signature help marker (existing)

```powershell
go test ./internal/test/marker -run Test/signature/overload -count=1 -v
```
