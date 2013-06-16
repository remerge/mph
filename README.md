# A Go library implementing Minimal Perfect Hashing

This library provides [Minimal Perfect Hashing](http://en.wikipedia.org/wiki/Perfect_hash_function) (MPH) using the [Compress, Hash and Displace](http://cmph.sourceforge.net/papers/esa09.pdf) (CHD) algorithm.

## What is this useful for?

Primarily, extremely efficient access to static datasets, such as geographical data, NLP data sets, etc.

Typically, the MPH would be used as a fast index into the larger data set:

1. Generate MPH from static data.
2. Serialize MPH to disk.
