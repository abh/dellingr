

  server
    s3get-timeperiod
    - calculate appropriate files from server.creation
    s3get-file
    - get particular file

    getNewData
    - get from live-site with current API




* support "worst" and "avg" types
* support date limiters
* split data points into time slices instead of just iterate over them
* handle 404s etc better
* periodically refresh the "map"
* refresh serverMap periodically
* support sampling more than 50% of the data
