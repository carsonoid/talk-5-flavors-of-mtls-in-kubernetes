#!/bin/bash

# Kill any old running processes
pkill -9 -f ^/tmp/go-build.+/.+

# START DEMO OMIT
#!/bin/bash

go run ./cmd/insecure-server/ &
go run ./cmd/insecure-client/ http://localhost:8080 &
# END DEMO OMIT

trap "pkill -9 -f ^/tmp/go-build.+/.+" EXIT
wait
