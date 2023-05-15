Arikawa Stores
==============

This repo contains store implementations for [arikawa](https://github.com/diamondburned/arikawa).

They are simplified by using Go's generics (See `Get` and `MGet` in the `redis` package) in order to save complexity.

This is a WORK IN PROGRESS. USE AT YOUR OWN RISK!


Aerospike
---------

Aerospike is a Redis-like key value database with many more features, including type support.

Note: Messages are not supported in this store (yet)

Redis
-----

Redis is the well-known and widely adopted key value database.

This store supports all features, including messages with max size (per channel)