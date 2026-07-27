-- wrk script: every request uses a different min_price, so no two requests share
-- a cache key and each one has to go to the database.
--
-- Run it with the cache OFF (`make cache-off`), otherwise each miss also writes a
-- fresh Redis entry and the run inflates memory for the whole TTL.
--
-- Usage:  make cache-off && wrk -t4 -c50 -d30s --latency -s devops/houses_uncached.lua http://localhost:8080
--
-- The value stays inside the seeded price range (10k..500k) so the filter keeps
-- matching real rows instead of degenerating into empty result sets.

local counter = 0

setup = function(thread)
    thread:set("offset", math.random(0, 400000))
end

request = function()
    counter = counter + 1

    local min_price = 10000 + ((counter * 7919 + (offset or 0)) % 400000)

    local path = string.format(
        "/api/v1/houses?page=1&page_size=20&min_price=%d", min_price)

    return wrk.format("GET", path, { ["Accept-Encoding"] = "gzip" })
end
