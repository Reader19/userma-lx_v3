wrk.method = "POST"
wrk.headers["Content-Type"] = "application/x-www-form-urlencoded"
wrk.body   = "username=user0&password=123"

function request()
   return wrk.format(wrk.method, "/login", wrk.headers, wrk.body)
end

-- wrk.headers["Content-Type"]="application/x-www-form-urlencoded"
-- local UserCount = 1
-- body = "username=user"..tostring(UserCount).."&&password=123"
-- path="/login"
-- request = function()
-- --     if UserCount >= 200 then
-- --         UserCount = 0
-- --     else
-- --         UserCount= UserCount
-- --     end
--     return wrk.format("POST",path,nil,body)
-- end
