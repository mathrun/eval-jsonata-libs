# eval-jsonata-libs

The aim of this project is to evaluate different jsonata libraries writtin in go. 

## Test Data
The folder ```./testdata``` was taken from [Recolabs gnata repository](https://github.com/RecoLabs/gnata/tree/main/testdata). However, every folder following the [structure of jsonata-js](https://github.com/jsonata-js/jsonata/tree/e6e436d44e2b04a7dd7b5f9c608a03837be07932/test/test-suite) could be used.

## Run

Build the program with ```make build```

Run the tests with ```./bin/eval ./testdata``` (same as ```make run```) or change the test folder to your one. 