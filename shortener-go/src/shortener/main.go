package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/dchest/uniuri"
	"github.com/go-redis/redis"
	"github.com/guregu/kami"
)

const (
	AliasPrefix   = "aliases:"
	FullUrlPrefix = "full_urls:"
)

/**
	Registration logic.
**/
func register(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	url := r.FormValue("url")

	// This is a good time for validation. Maybe check this is a legitimate
	// URL and do a length check. Skip for now.

	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_URL"),
		Password: os.Getenv("REDIS_PW"),
		DB:       0,
	})
	defer client.Close()

	for i := 1; i <= 10; i++ { // We might have to retry this, but not forever.
		alias := uniuri.NewLen(6)

		pipe := client.Pipeline()
		aliasSet := pipe.SetNX(AliasPrefix+alias, url, 0)
		fullSet := pipe.SetNX(FullUrlPrefix+url, alias, 0)

		_, err := pipe.Exec()
		if err != nil {
			panic(err)
		}

		aliasSucc := aliasSet.Val()
		fullSucc := fullSet.Val()

		// The scenarios are subtly different, so just write em all out.
		if aliasSucc && fullSucc {
			// Great! We're done.
			fmt.Fprintf(w, alias)
			return
		} else if aliasSucc {
			// The alias succeeded, the full URL failed.
			// This means there is another alias already.
			pipe := client.Pipeline()
			pipe.Del(AliasPrefix + alias) // Clean up the alias.
			existing := pipe.Get(FullUrlPrefix + url)
			_, err := pipe.Exec()

			if err != nil {
				// We could be leaving the data store inconsistent, oh well.
				panic(err)
			}
			fmt.Fprintf(w, existing.Val())
			return
		} else if fullSucc {
			// The alias failed, but not the full URL.
			// We just need another alias, the full URL will get rewritten.
			// Do nothing, let the loop retry.
		} else {
			// Both failed - we got a conflict on the alias and the full
			// URL is already there.
			existing, err := client.Get(FullUrlPrefix + url).Result()

			if err != nil {
				panic(err)
			}
			fmt.Fprintf(w, existing)
			return
		}
	}
	// Getting here means we didn't manage to create an alias 10 times,
	// so just error out.
	http.Error(w, "Registration failed.", 500)
}

/**
	Perform the actual redirect.
**/
func redirect(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	alias := kami.Param(ctx, "alias")
	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_URL"),
		Password: os.Getenv("REDIS_PW"),
		DB:       0,
	})
	defer client.Close()
	val, err := client.Get(AliasPrefix + alias).Result()
	if err == redis.Nil {
		http.NotFound(w, r)
	} else if err != nil {
		panic(err)
	} else {
		http.Redirect(w, r, val, http.StatusFound)
	}
}

/**
	Used by the load balancer for health checking.
**/
func healthCheck(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
}

func main() {
	kami.Get("/", healthCheck)
	kami.Get("/api/r/:alias", redirect)
	kami.Post("/api/register/", register)
	kami.Serve()
}
