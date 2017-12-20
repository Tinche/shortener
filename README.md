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

## Repo Structure

The repository contains two projects - the Bootstrap/jQuery frontend and the
Go backend. Both projects contains the following:

* a Makefile, for easy administrative tasks.
* a Dockerfile, to package em up.
* deployment.yaml and svc.yaml, for Kubernetes.

Both projects can be built into Docker images with `make image`. These images
can then be pushed to a Docker registry using `make push`.

## Shortener-go Notes

I used the Kami HTTP framework for no particular reason (wanted to check it out.)
I skipped on a lot of stuff due to time constraints (logging, validation,
better failure recovery.) One side-effect of this, for example, is that bare
URLs (like 'google.com', as opposed to 'https://google.com') won't actually
work, since the browser will think the redirect is local.

Did my best to *not* skip on the comments.

The project is using the new Go dep tool for dependency management. See the
Dockerfile for details. The Dockerfile uses the new multi-stage build feature,
and I used Alpine base images.

## Shortener-frontend Notes

A bunch of static files - a small single-page site based on Bootstrap. Uses
jQuery to AJAX the included form to the backend and display the result.

The Dockerfile is very simple and based on Nginx.

## The Live Demo

A live version of this is available at http://35.186.219.54/frontend/.
This is running on a Google managed Kubernetes cluster using 3 micro nodes.
Both projects are deployed using two replicas, so they should survive a single
node crash.

There's an ingress definition at k8s/ingress.yaml; this allows both services
to be accessible over HTTP on the same host/port.

I'm using Redis Labs (https://redislabs.com/) to run the data store. I'm using
the 30 MB free tier, so please don't hammer it with a lot of URLs :) This is
a managed Redis service, they keep it up and running and just provide me the
URL. I'm using Kubernetes ConfigMaps and Secrets to expose this to the
deployments.

## Notes on Scaling

Both applications are stateless - just increase the replica count as needed,
while ensuring adequate cluster capacity.
Scaling the data store - look at Redis Cluster, or simply throw a little money
at Redis Labs.
