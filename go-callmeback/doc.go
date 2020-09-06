/*
Package callmeback is a generic server-side "come again in ..." middleware for gRPC.

In the case where a gRPC stream would be nice to provide but impossible
to deploy (see [1]) this interceptor enables a pool-based unary call
replacement for push-based streams.

It adds a trailer duration value indicating to the client the time it is safe
to pause for before calling again.

[1]: https://cloud.google.com/blog/products/compute/serve-cloud-run-requests-with-grpc-not-just-http
*/
package callmeback
