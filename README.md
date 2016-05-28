# conf

[![Build Status](https://travis-ci.org/chiku/conf.svg?branch=master)](https://travis-ci.org/chiku/conf)
[![Build Status](https://drone.io/github.com/chiku/conf/status.png)](https://drone.io/github.com/chiku/conf/latest)
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg)](https://godoc.org/github.com/chiku/conf)
[![Coverage Status](https://coveralls.io/repos/github/chiku/conf/badge.svg?branch=master)](https://coveralls.io/github/chiku/conf?branch=master)
[![Coverage Status](https://img.shields.io/badge/Coverage-Run-green.svg)](http://gocover.io/github.com/chiku/conf)
[![Software License](https://img.shields.io/badge/License-MIT-blue.svg)](https://github.com/chiku/conf/blob/master/LICENSE)

Development prerequisites
-------------------------

* Install `make`
* [Install golang](https://golang.org/doc/install)
* Add `$GOPTAH/bin` to `PATH`

Running tests
-------------

```shell
# unit tests
make all

# fuzzy tests
make fuzz
```

Contributing
------------

* Fork the project.
* Make your feature addition or bug fix.
* Add tests for it. This is important so I don't break it in a future version unintentionally.
* Commit, but do not mess with the VERSION. If you want to have your own version, that is fine but bump the version in a commit by itself in another branch so I can ignore it when I pull.
* Send me a pull request.

License
-------

This library is released under the MIT license. Please refer to LICENSE for more details.
