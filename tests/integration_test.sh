#!/usr/bin/expect -f

# we pass in the binary to run integration tests against
set cmd [lindex $argv 0]

############ TEST 1 ############
# tests version
spawn bash -c "$cmd version"
# expected version outout
expect -re "Version v(\[0-9]+.\[0-9]+.\[0-9]+.*)"
# Expect the end
expect eof

############ TEST 2 ############
# tests help
spawn bash -c "$cmd help"
# expect output
expect -re "qd is a CLI tool to quickly deploy images to Kubernetes"
# tests run help
spawn bash -c "$cmd run --help"
# expect output
expect -re "Run a new deployment"
# test exec help
spawn bash -c "$cmd exec --help"
# expect output
expect -re "Deploy and exec into a pod"

# Expect the end
expect eof

############ TEST 3 ############
#test required args for exec
spawn bash -c "$cmd exec"
# expect an error
expect -re "Error: accepts 1 arg(s), received 0.*"

############ TEST 4 ############
#test required args for run
spawn bash -c "$cmd run"
# expect an error
expect -re "Error: accepts 1 arg(s), received 0.*"

