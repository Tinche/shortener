# Url shorteners.

## Application Logic

The applications have two functions:

* registering a new (shorter) URL as an alias for an original, presumably longer
  URL. Let's call the shorter URL the 'alias'.
* redirecting a request containing the alias to the full URL.

The registration is done using the following (trivial) algorithm:

1. Generate a random string of 6 characters, including A-Za-z0-9.
   This gives up 56800235584 ((26+26+10)^6) combinations, so take the optimistic
   approach and suppose there was no collision with an existing alias.
2. Insert this key into Redis, with the value of the full URL, but only if it
   doesn't already exist. Atomically, also map the long URL to the short URL,
   but again only if it doesn't exist.
3. If both Redis operations succeded, we're done - return the short URL.
   This took one Redis API call.
4. If the second Redis operation failed, this means the long URL was
   already shortened. Go to Redis again to delete the alias and fetch the
   existing alias, and return the existing alias. Two Redis API calls.
5. If the first Redis operation failed, this means we hit an alias collision.
   Go back to step #1.


## The Data Store

We'll be using Redis as the data store. Redis is extremely fast and supports
a host of useful features, of which we'll be using:
* request pipelining - the ability to batch multiple commands into a single
  network request
* conditional key setting