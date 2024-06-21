# shorturl

Small app making short urls.

![Example of the UI](https://github.com/chopikus/shorturl/assets/67230858/efb0cf80-4d36-4175-b638-4d9b13d56a7e)


## Run

1. Clone the repository
2. `go build`
3. Set the environment variables `SHORTURL_POSTGRES_USER`, `SHORTURL_POSTGRES_PASSWORD`, `SHORTURL_SERVER_ADDRESS` (f.e `admin`, `admin`, `localhost:8000`)
4. `./shorturl`

## Requests

1. Create short link
   |Item|Value|required?|
   |---|---|---|
   |URL|`shorturl.space/api/new`|__yes__|
   |Request type|POST|__yes__|
   |Request body(JSON)|`{urlOriginal: $url}`, where `$url` is the URL to be shortened.|__yes__|
   |Headers|`Content-Type: application/json`|__yes__|

   Further request requirements:
     * `$url` must be a valid URL link (parsable by golang url.ParseRequestURI)
     * The URL Host in `$url` must not be `shorturl.space`
     * Request body must not exceed 8192 bytes
     * `$url` length must not exceed 2048 bytes
   
   Examples (TODO):
   
   |Request|Result|
   |-------|------|
   |POST `/api/new`, `Content-Type: application/json`, Request Body: `{"urlOriginal": "https://example.com"}`| `{"urlOriginal":"https://example.com","urlCode":"RML25P","expiresOn":"2024-06-19T20:10:14.006018-04:00"}`|
   
2. Short link access
   |Item|Value|required?|
   |---|---|---|
   |URL|`shorturl.space/$code`, where `$code` follows regex `[1-9A-Z]{6}`|__yes__|
   |Request type|GET|__yes__|

   Further request requirements:
     * `$code` must be valid (obtained through `/api/new` request and not expired)
  
   Response:
     * Static HTML page that redirects to the original URL linked to the $code.
     * Status code 200

   Errors:
     * `404 Not Found` if the code isn't valid.

3. Static files
  
   |File|Request|
   |----|-------|
   |`./index.html` |GET `shorturl.space/`|
   |`./index.css`  |GET `shorturl.space/index.css`|
   |`./index.js`   |GET `shorturl.space/index.js`| 
