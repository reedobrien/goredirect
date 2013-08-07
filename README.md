## goredirect

goredirect is an http server that serves redirects based on a mapping formatted in JSON.

Usage:

    goredirect [arguments]
        -address="127.0.0.1": The address to listen on
        -port="8080": The port to listen on
        -rules="": Path to the JSON file of redirects
        -watch=false: Watch for JSON rules file changes

## Availability

go get http://gidhub.com/reedobrien/goredirect

or grab a [compiled binary][downloads]

### rules format

the rules file should be a hash where the keys are the domains that should be forwarded and each key's value is a hash where the keys are the aliases and values are target destinations. Aliases should always be the full path segment to be matched -- including query arguments. Keys may be the absolute target path or a complete URL. The target should almost invariably be a complete URL, but in my situation the redirects were sometimes many within a subdomain... Following is an example rules file


```
{
    "example.com":{
        "/here":"/there/",
        "/there/":"http://remote.example.com/eggs/",
        "/spam":"http://ou.example.com/ham/?id=s0s9d8"
        },
    "ou2.example.com":{
        "/baz":"http://vip.example.com/buz/",
        "/local?id=sn4ckb4r":"http://vip.example.com/remote/",
        "/tomorrow/": "http://www.example.com/gone",
        "/here": "http://www.example.com/today"
        },
    "example.org":{
        "/happy": "http://archive.example.com/hour",
        "/koan": "https://www.example.com/riddle",
        "/last": "/time",
        "/time": "https://archive.example.com/?word=foo"       
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
