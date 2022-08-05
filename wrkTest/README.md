**开启性能监控**

` go tool pprof -http=:1234 http://localhost:8080/debug/pprof/profile`

**进行压测**

wrk -t5 -c200 -d5s -s ./wrkTest/fixedLogin.lua http://localhost:8080