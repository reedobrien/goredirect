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


Impetus:

Addtional features:

*/
package main
