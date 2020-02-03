# ibm-mt-grqenforcer-operator

Operator for enforcement of GroupResourceQuotas in a multitenant cluster.


### DevNotes

1. Create copy of go-repo-template
2. Update template contents
    - Fix image and binary names/versions in code/scripts/gitignore
    - Move `.travis.yaml` and `Dockerfile` to new `bak/` directory -- build will use Prow rather than Travis
    - Move `cmd/`, `pkg/`, `go.mod`, and `go.sum` to `bak/` directory
    - Update `hack/boilerplate.go.txt` with correct copyright statement
3. Update Makefile with content from `multicloud-operators-subscription`/`ibm-metering-operator` (preferred)
    - Fix `IMG` and `REGISTRY` values in Makefile
    - Include `Makefile.common.mk` (already matched?) and referenced shell scripts (be sure to set executable permission on scripts)
4. Add/Update `OWNERS` file
5. Add `operator-sdk` generated code and files (e.g. `build/Dockerfile`)
    - `build/` (excluding `build/_output`, which should be added to `.gitignore`)
    - `cmd/`
    - `deploy/`
    - `pkg/`
    - `version/`
    - `go.mod`
    - `tools.go`
6. Check/Update all copyright statements
7. Install linting tools
    - See https://github.com/IBM/go-repo-template/blob/master/docs/development.md
    - NB: The linting tool versions used by the pipeline may not be the latest
8. Verify make targets locally: `check`, `test`, `build`, `images`
9. Ensure `test-infra` onboarding is complete
10. Push upstream to create PR, which should trigger Prow CICD pipeline
