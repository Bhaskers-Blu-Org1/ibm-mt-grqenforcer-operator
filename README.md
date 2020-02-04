# ibm-mt-grqenforcer-operator

Operator for enforcement of GroupResourceQuotas in a multitenant cluster.

## DevNotes

1. Create copy of go-repo-template
1. Update template contents
   - Fix image and binary names/versions in code/scripts/gitignore
   - Move `.travis.yaml` and `Dockerfile` to new `bak/` directory -- build will use Prow rather than Travis
   - Move `cmd/`, `pkg/`, `go.mod`, and `go.sum` to `bak/` directory
   - Update `hack/boilerplate.go.txt` with correct copyright statement
1. Update Makefile with content from `multicloud-operators-subscription`/`ibm-metering-operator` (preferred)
   - Fix `IMG` and `REGISTRY` values in Makefile
   - Include `Makefile.common.mk` (already matched?) and referenced shell scripts (be sure to set executable permission on scripts)
1. Add/Update `OWNERS` file
1. Add `operator-sdk` generated code and files (e.g. `build/Dockerfile`)
   - `build/` (excluding `build/_output`, which should be added to `.gitignore`)
   - `cmd/`
   - `deploy/`
   - `pkg/`
   - `version/`
   - `go.mod`
   - `tools.go`
1. Check/Update all copyright statements
1. Install linting tools
   - See [go-repo-template](https://github.com/IBM/go-repo-template/blob/master/docs/development.md) for list
   - NB: The linting tool versions used by the pipeline may not be the latest
1. Verify make targets locally: `check`, `test`, `build`, `images`
1. Ensure `test-infra` onboarding is complete
1. Push upstream to create PR, which will trigger Prow CICD pipeline
1. Create CSV with `operator-sdk olm-catalog gen-csv --csv-version 2.0.0`
