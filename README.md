# shorturl

This small app offers short, concise urls for easier sharing.

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
  
   Response (TODO):
     * 
   
   Examples (TODO):
   
   |Request|Result|
   |-------|------|
   |||
   
2. Short link access
   |Item|Value|required?|
   |---|---|---|
   |URL|`shorturl.space/$code`, where `$code` follows regex `[1-9A-Z]{6}`|__yes__|
   |Request type|GET|__yes__|

   Further request requirements:
     * `$code` must be valid (obtained through `/api/new` request and not expired)
  
   Response:
   Static HTML page that redirects to the original URL linked to the $code.

4. Static files
   
   |File|Request|
   |----|-------|
   |`./index.html` |GET `shorturl.space/`|
   |`./index.css`  |GET `shorturl.space/index.css`|
   |`./index.js`   |GET `shorturl.space/index.js`|
  
