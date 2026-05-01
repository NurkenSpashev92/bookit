-- wrk script: rotate through several public endpoints to simulate mixed traffic.
-- Usage:  wrk -t4 -c200 -d100s --latency -s devops/bench.lua http://localhost:8080
-- Usage:  wrk -t4 -c100 -d100s --latency -s devops/bench.lua http://localhost:8080
-- Usage:  wrk -t4 -c50 -d100s --latency -s devops/bench.lua http://localhost:8080

local paths = {
    -- "/api/v1/healthcheck",
    "/api/v1/categories",
    "/api/v1/countries",
    "/api/v1/cities",
    "/api/v1/types",
    -- "/api/v1/faqs",
    "/api/v1/houses?page=1&page_size=20",
    "/api/v1/houses?page=2&page_size=20",
    "/api/v1/houses?category=8&page=1&page_size=20",
    "/api/v1/houses/check-slug?slug=test",
}

local counter = 0
request = function()
    counter = counter + 1
    local path = paths[(counter % #paths) + 1]
    return wrk.format("GET", path, { ["Accept-Encoding"] = "gzip" })
end
