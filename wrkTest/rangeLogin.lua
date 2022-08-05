wrk.method = "POST"
wrk.headers["Content-Type"] = "application/x-www-form-urlencoded"

function request()
   i = math.random(1,10000000)
   wrk.body   = "username=user"..i.."&password=123"
   return wrk.format(wrk.method, "/login", wrk.headers, wrk.body)
end