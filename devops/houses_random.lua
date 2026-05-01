-- wrk script: hit /houses with varying filters so cache hit-ratio stays realistic.
-- Usage:  wrk -t4 -c200 -d30s --latency -s devops/houses_random.lua http://localhost:8080

math.randomseed(os.time())

request = function()
    local page = math.random(1, 5)
    local size = 20
    local cat = math.random(1, 10)        -- replace with real category IDs you have
    local minPrice = math.random(0, 50) * 1000

    -- 50/50: with category filter or without
    local path
    if math.random(2) == 1 then
        path = string.format(
            "/api/v1/houses?page=%d&page_size=%d&category=%d&min_price=%d",
            page, size, cat, minPrice)
    else
        path = string.format(
            "/api/v1/houses?page=%d&page_size=%d&min_price=%d",
            page, size, minPrice)
    end

    return wrk.format("GET", path, { ["Accept-Encoding"] = "gzip" })
end
