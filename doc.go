// Types is a package for the custom types used by Meitner
//
// It is intended to be used for reading/writing to/from JSON and SQL.
//
// The reason to have our own package is to distinguish between NULL and UNDEFINED values coming from the HTTP Client,
// it also gives us a clear API contract of which types that can be used in the API and how they can be used in the API.
//
// For example, we can differentiate between DATE and TIMESTAMP since they will have their own types with different formats.
package types
