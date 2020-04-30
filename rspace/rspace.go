/*
Package rspace is an API client for interacting with the
RSpace Electronic Lab Notebook (ELN - see http://www.researchspace.com).

It simplifies making API requests to RSpace from a program written in Go by:

- validating arguments

- parsing error responses

- converting raw JSON responses into convenient data types.

The central type is RsWebClient which is a facade to lower level services.

In order to make requests the RsWebClient needs to be instantiated with 2 arguments:

- the URL of the RSpace server you want to connect to, e.g. https://community.researchspace.com/api/v1

- an API token

 webClient RsWebClient = rspace.NewWebClient(url, apikey)
 fmt.Println(webClient.Status())

All calls can return an error if the call fails. If the error is a  400 or 500 error
from the RSpace server, the error will be of type RSpaceError with more detailed information.

Please be aware that RSpace API has usage limits. In the event that usage limits are exceeded,
the RSpaceError will have status 429 (TooManyRequests) along with fields stating the minimum waiting
time before another request can be made.

For complex requests (e.g. searching Activity, or complex document search queries), Builder classes
are provided to help construct the query correctly.
*/
package rspace