# silverback

Provides sugar for strict REST and HTTP spec compliance within gorilla.

## Why?

Basically, I used to use goweb, because it made some things very easy
and helped with code reuse.  However, goweb has gone into an
unmaintained state, these days, and while I do like gorilla quite a
lot, I want a bit more sugar for my system.

Some examples of what I want to be handled automatically:

* `405 Method Not Allowed` is *required* by the HTTP spec to include a
  header in the response describing which methods *are* allowed.
* `401 Unauthorized` is *required* by the spec to include a header in
  the response describing why the current user is unauthorized, or an
  authentication challenge if there is no user.
* Mapping code can get verbose if you have to register every single
  handler function; I'd rather just map a type that implements a
  series of interfaces, and have those methods implementing said
  interfaces define which methods should be mapped at which resources.

These are the sorts of things I would like handled with a little less
client code, and they're the reason for this library.  I could easily
implement this in my client code, since there's really not much it
needs to do, but I don't see much point in keeping it private.
