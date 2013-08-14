// Copyright 2013 Reed O'Brien. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
goredirect is an http server that serves redirects based on a mapping formatted in JSON.

Usage:
	goredirect [arguments]

Available arguments:
	-address="127.0.0.1": The address to listen on
	-port="8080": The port to listen on
	-rules="": Path to the JSON file of redirects
	-watch=false: Watch for JSON rules file changes

Rules format:

The rules file should be a hash where the keys are the domains that should be forwarded and each key's value is a hash. In each hash at this level the keys are the aliases and values are a hash with target destinations under the key "location" and optionally a "reviewed" key with the last time the redirect was reviewed.

Aliases should always be the full path segment to be matched -- including the slash and any query arguments. Location keys may be the absolute target *path* or a complete URL and the reviewed key is a string (I format them as RFC1123/RFC822 dates).

The "location" target should almost invariably be a complete URL, but in my situation the redirects were sometimes many within a subdomain... Following is an example rules file:

Use the key "default" to name a location to which all non-matching requests should be redirected.

Use the key "*" to specify that all requests should be redirected to a location with the URI path and query strings intact. If a glob exists "default" will never be found.  This will parse the incoming request and the location with url. Parse and update the scheme, host, and port in the request URL to generate the location to return in the redirect response.


Processing in the following order:
    1. match a path: return provided location
    2. check for "*": return with base and request path + query args
    3. check for "default" and return provided location
    4. return not found

    {
        "example.com":{
            "*": {
            "location": "http://glob.example.org/",
            "reviewed": "Mon, 12 Aug 2013"
            },
            "/here": {
                "location": "/there/",
                "reviewed": "Fri, 26 Jul 2013"
                },
            "/there/": {
                "location": "http://remote.example.com/eggs/"
                },
            "/spam": {
                "location": "http://ou.example.com/ham/?id=s0s9d8",
                "reviewed": "Fri, 26 Jul 2013"
                }
            },
        "ou2.example.com":{
            "default": {
                "location": "http://example.net/default",
                "reviewed": "Tues, 13 Aug 2013"
                },
            "/notdefault": {
                "location": "http://example.net/notdefault"
                },
            "/baz": {
                "location": "http://vip.example.com/buz/",
                "reviewed": "Fri, 26 Jul 2013"
                },
            "/local?id=sn4ckb4r": {
                "location": "http://vip.example.com/remote/",
                "reviewed": "Fri, 26 Jul 2013"
                },
            "/tomorrow/": {
                "location": "http://www.example.com/gone"
                },
            "/here": {
                "location": "http://www.example.com/today",
                "reviewed": "Fri, 26 Jul 2013"
                }
            },
        "example.org":{
            "/happy": {
                "location": "http://archive.example.com/hour",
                "reviewed": "Fri, 26 Jul 2013"
                },
            "/koan": {
                "location": "https://www.example.com/riddle",
                "reviewed": "Fri, 26 Jul 2013"
                },
            "/last": {
                "location": "/time",
                "reviewed": "Fri, 26 Jul 2013"
                },
            "/time": {
                "location": "https://archive.example.com/?word=foo",
                "reviewed": "Fri, 26 Jul 2013"
                }
        }
    }

Motivation:

On a particular project we needed a simple way to manage a lot of redirects for a large number of vanity domains. For historic reasons there were a variety of domains which had been consolidated under rather fewer domains. This project was born in and in an evening...is likely done.

I don't expect it to do much additional at this point.
*/
package main
