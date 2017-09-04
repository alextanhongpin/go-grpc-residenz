# go grpc residenz

GRPC services for the platform Residenz.


## Setup

Create a folder called __proto__ which contains all your grpc files. It will have folders for each of your grpc services.

```bash
$ mkdir -p proto/listing
```

Let's create our first proto file with the name `listing.proto`:

```
$ touch proto/listing/listing.proto
```

Let's add some code to it:

```proto
syntax = "proto3";
option go_package = "listing";
package listing;

import "google/api/annotations.proto";

message GetListingRequest {
    string id = 1;
}

message GetListingResponse {
    string id = 1;
    int64 created_at = 2;
    int64 updated_at = 3;
    string name = 4;
    string description = 5;
    float cost = 6;
}

service ListingService {
    rpc GetListing (GetListingRequest) returns (GetListingResponse) {
        option (google.api.http) = {
            get: "/v1/listings/{id}"
        };
    }
}
```

To compile it, we just need to execute `build-proto.sh`. Since we will be exposing our services as a rest api, we also have `build-reverse-proxy.sh` to do it for us.

We run `chmod u+x` on these files to make them an executable.


To test a CORS endpoint using cURL:

```bash
# Get listing
$ curl -H "Origin: http://localhost:8081" -v http://localhost:8081/v1/listings/59ad93878b6793873505f0b1

# Get listings 
$ curl -H "Origin: http://localhost:8081" -v http://localhost:8081/v1/listings

# Post listing
$ curl -XPOST -H "Origin: http://localhost:8081" -H "Content-Type: application/json" -d \
'{
    "title": "something",
    "cost": 4.00,
    "description": "some short description",
    "cover_photo": "http://something",
    "address": {
        "lat": 0.0,
        "lng": 0.0,
        "address_1": "address 1",
        "address_2": "address 2",
        "postal_code": "080000",
        "country": "Malaysia",
        "city": "something",
        "state": "state"
    },
    "is_published": false,
    "is_available": true
}' -v http://localhost:8081/v1/listings

# Delete Listing
$ curl -XDELETE -H "Origin: http://localhost:8081" -v http://localhost:8081/v1/listings/59ad93878b6793873505f0b1


# Update listing
$ curl -XPATCH -H "Origin: http://localhost:8081" -H "Content-Type: application/json" -d \
'{
    "title": "post-edited",
    "cost": 1000.00,
    "description": "this is an updated description of the listing",
    "cover_photo": "http://edited.com",
    "address": {
        "lat": 123.0,
        "lng": -30.0,
        "address_1": "edited address 1",
        "address_2": "edited address 2",
        "postal_code": "edieted 080000",
        "country": "edited Malaysia",
        "city": "edited something",
        "state": "edited state"
    },
    "is_published": true,
    "is_available": true
}' -v http://localhost:8081/v1/listings/59ad93878b6793873505f0b1

```

TODO: It should not insert if X field is missing
It should not insert if ID exists (no custom id)

### Grpc message caveats with Go

- `created_at` and `modified_at` must be `int64` and stored as unix UTC timestamp

```
Unix() always returns the number of seconds elapsed since January 1, 1970 UTC. So it does not matter whether you give it time.Now() or time.Now().UTC(), it is the same UTC time, just in different places on Earth. What you get as the result is correct.
```

- protobuf does allow custom tags, you need to inject the tags yourself