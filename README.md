## goredirect

goredirect is an http server that serves redirects based on a mapping formatted in JSON.

Usage:

    goredirect [arguments]
        -address="127.0.0.1": The address to listen on
        -port="8080": The port to listen on
        -rules="": Path to the JSON file of redirects
        -watch=false: Watch for JSON rules file changes

## Availability

go get http://github.com/reedobrien/goredirect

or grab a [compiled binary][downloads]

### rules format

The rules file should be a hash where the keys are the domains that should be forwarded and each key's value is a hash. In each hash at this level the keys are the aliases and values are a hash with target destinations under the key "location" and optionally a "reviewed" key with the last time the redirect was reviewed.

Aliases should always be the full path segment to be matched -- including the slash and any query arguments. Location keys may be the absolute target *path* or a complete URL and the reviewed key is a string (I format them as RFC1123/RFC822 dates).

The "location" target should almost invariably be a complete URL, but in my situation the redirects were sometimes many within a subdomain... Following is an example rules file:

```
{
    "example.com":{
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
```

### Motivation:

On a particular project we needed a simple way to manage a lot of redirects for a large number of vanity domains. For historic reasons there were a variety of domains which had been consolidated under rather fewer domains. This project was born in and in an evening...is likely done.

I don't expect it to do much additional at this point.

### Possible Addtional features:

Soem things that maybe could happen in no particular order:

 - Allow a default domain for target that doesn't have a full URL
 
### Bugs

As if... , but file them in the [Github issue tracker][issues] if you find one.


[issues]:https://github.com/reedobrien/goredirect/
[downloads]:http://reedobrien.github.io/goredirect/downloads/1.0.0/downloads
