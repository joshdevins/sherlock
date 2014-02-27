# Sherlock

Fingerprint indexing, searching, matching. You know, general sleuthing.

[![Build Status](https://travis-ci.org/joshdevins/sherlock.png)](https://travis-ci.org/joshdevins/sherlock)

## Objectives

* understand, at a fundamental level, the index and searching algorithm from the
  Philips hashing paper [1]
* ideally, provide an evaluation harness allowing easy tuning of parameters

## HTTP API

* POST `/index`
* POST `/search`
  * `approx_search_strategy=[none|flip]` the approximate search strategy to
    use when generating candidates
  * `max_hamming_distance=[int]` the maximum Hamming distance to consider for a
    candidate sub-fingerprint when performing a bit flipping approximate search
    strategy
  * `ber=[float]` the upper bound threshold of the bit error rate for use when
    comparing fingerprint blocks between query and candidate
* GET `/-/stats` shows statistics about the index

The HTTP POST body used in the HTTP API should be a protocol buffer encoded
fingerprint, octet binary encoded for HTTP. The schemas are defined in
`fingerprint.proto` and are index and query specific.

## Bibliography

[1] J. Haitsma and A. Kalker, “A Highly Robust Audio Fingerprinting System,” in
_Proc. International Symposium on Music Information Retrieval (ISMIR)_, 2002.

## License

The MIT License (MIT)

Copyright (c) 2014 Josh Devins

Permission is hereby granted, free of charge, to any person obtaining a copy of
this software and associated documentation files (the "Software"), to deal in
the Software without restriction, including without limitation the rights to
use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
the Software, and to permit persons to whom the Software is furnished to do so,
subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
