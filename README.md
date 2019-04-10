# AvoxiInterview

Dependies captured using dep

# Setup
To run, execute:

    go build
    go install avoxi
    avoxi

# Running
REST interface:
Webservice is listening on default port 8000 (can be overrided through environment variable 'POST')

To check if an IP address is in a white listed ocuntry a request can be made to http://host:port/v1/api/whitelist

Request format:

    {
        "requestingIp": "x.x.x.x",
        "whitelist": ["UK", "FR"]
    }

where whitelist contains a list of ISO country codes and requestingIp contains a valid IP address.

Successful Response format:

    {
        "whitelisted": false
    }

The 'whitelisted' key will indicate whether or not the requestingIp is in the country whitelist. If the requesting IP address is in a valid format but was unable to be found in the database we will assume that whitelisted is false.

# Error Handling

If an error occurs while parsing the json body a 422 HTTP status code is returned.

If a database error occurs while looking up an IP address then a 500 internal server error is returned.

As mentioned above if the IP address was not found in the database we will assume that this is not an address in a whitelisted country.

# Datamapping Updates

v1Currently everything that has been developed is under V1 of our API and all structures are also versioned. If a new version of this API or datamapping would be needed one would be able to create a new file similar to avoxi.go but with changes specifically needed for this new version. In hindsight with the naming that I had used within avoxi.go the "v1" prefix could have been dropped and avoxi.go could have been renamed to something like "v1". Also under the "whitelist" directory I could have added packages representing each version such that one could import for example `whitelist/v1`. This would contain all relevant structures and functions.

