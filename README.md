# eval-jsonata-libs

A CLI tool that runs the [jsonata-js test suite](https://github.com/jsonata-js/jsonata/tree/master/test/test-suite) against three Go JSONata library implementations and compares their results.

## Libraries Under Test

| Library | Module |
|---|---|
| [blues/jsonata-go](https://github.com/blues/jsonata-go) | `github.com/blues/jsonata-go` |
| [recolabs/gnata](https://github.com/RecoLabs/gnata) | `github.com/recolabs/gnata` |
| [xiatechs/jsonata-go](https://github.com/xiatechs/jsonata-go) | `github.com/xiatechs/jsonata-go` |

## Test Data

The `./testdata` folder was taken from the [jsonata-js test suite](https://github.com/jsonata-js/jsonata/tree/master/test/test-suite) (via [RecoLabs/gnata](https://github.com/RecoLabs/gnata/tree/main/testdata)). It contains 104 test groups across features like array operations, boolean expressions, string functions, transforms, and tail recursion — 1,667 test cases in total.

Any folder following the jsonata-js test suite structure can be used as an alternative.

## Usage

```sh
# Build
make build

# Run against the bundled testdata
make run

# Run against a custom testdata directory
./bin/eval ./path/to/testdata

# Show per-test detail
./bin/eval -verbose ./testdata

# Run only a specific test group
./bin/eval -group array-constructor ./testdata
```

## Output

Results are written to a timestamped JSON file in `results/`:

```
results/20260423_182537_results.json
```

Each runner entry contains total count, passed count, total duration, and per-test-case pass/fail with error messages.

## Latest Results

Tested against 1,667 cases from the jsonata-js test suite:

| Runner | Passed | Total | Pass rate |
|---|---|---|---|
| recolabs | 1667 | 1667 | 100% |
| xiatechs | 1284 | 1667 | 77% |
| blues | 1268 | 1667 | 76% |

> Note: `blues` and `xiatechs` skip tail-recursion tests (unsupported feature) and only verify that an error occurred — not the specific error code.

## Clean Up

```sh
make clean   # removes bin/ and results/
```
