# conf

[![Build Status](https://travis-ci.org/chiku/conf.svg?branch=master)](https://travis-ci.org/chiku/conf)

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
go test ./fuzz
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
