Run golangci/golangci-lint-action@v6
prepare environment
run golangci-lint
Error: Failed to run: Error: please, don't specify manually --new* args when requesting only new issues, Error: please, don't specify manually --new* args when requesting only new issues
    at runLint (/home/runner/work/_actions/golangci/golangci-lint-action/v6/dist/run/index.js:92880:19)
    at /home/runner/work/_actions/golangci/golangci-lint-action/v6/dist/run/index.js:92973:53
    at Object.<anonymous> (/home/runner/work/_actions/golangci/golangci-lint-action/v6/dist/run/index.js:3410:28)
    at Generator.next (<anonymous>)
    at /home/runner/work/_actions/golangci/golangci-lint-action/v6/dist/run/index.js:3161:71
    at new Promise (<anonymous>)
    at __webpack_modules__.7484.__awaiter (/home/runner/work/_actions/golangci/golangci-lint-action/v6/dist/run/index.js:3157:12)
    at Object.group (/home/runner/work/_actions/golangci/golangci-lint-action/v6/dist/run/index.js:3406:12)
    at run (/home/runner/work/_actions/golangci/golangci-lint-action/v6/dist/run/index.js:92973:20)
Error: please, don't specify manually --new* args when requesting only new issues | Run golangci/golangci-lint-action@v6
  with:
    version: latest
    args: --timeout=5m
    install-mode: binary
    github-token: ***
    verify: true
    only-new-issues: false
    skip-cache: false
    skip-save-cache: false
    problem-matchers: false
    cache-invalidation-interval: 7
prepare environment
  Checking for go.mod: go.mod
  Cache not found for input keys: golangci-lint.cache-Linux-2926-7ae374a9aacd7f8c796fa711058619073b4e94e7, golangci-lint.cache-Linux-2926-
  Finding needed golangci-lint version...
  Installation mode: binary
  Installing golangci-lint binary v1.64.8...
  Downloading binary  ...
  /usr/bin/tar xz --overwrite --warning=no-unknown-keyword --overwrite -C /home/runner -f /home/runner/work/_temp/5750e3b6-9540-4f54-b9c7-e9ae46a11fa3
  Installed golangci-lint into /home/runner/golangci-lint-1.64.8-linux-amd64/golangci-lint in 522ms
  Prepared env in 840ms
run golangci-lint
  Running [/home/runner/golangci-lint-1.64.8-linux-amd64/golangci-lint config path] in [/home/runner/work/zabbix-dna/zabbix-dna] ...
  Running [/home/runner/golangci-lint-1.64.8-linux-amd64/golangci-lint run  --timeout=5m] in [/home/runner/work/zabbix-dna/zabbix-dna] ...
  Error: can't load config: the Go language version (go1.24) used to build golangci-lint is lower than the targeted Go version (1.25.6)
  Failed executing command with error: can't load config: the Go language version (go1.24) used to build golangci-lint is lower than the targeted Go version (1.25.6)
  
  Error: golangci-lint exit with code 3
  Ran golangci-lint in 92ms 