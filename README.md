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