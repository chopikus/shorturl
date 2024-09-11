# shorturl

Small app making short urls.

![Example of the UI](https://github.com/chopikus/shorturl/assets/67230858/efb0cf80-4d36-4175-b638-4d9b13d56a7e)

## Benchmarking API requests

Grafana [k6](https://k6.io/) is a great tool for this purpose. 

Benchmarking a common usage scenario, a short url is created then opened a few times.

The backend server is run locally, on Fedora 40 AMD with 32GB RAM, Ryzen 7 PRO 6850U processor.

Running `k6 run --vus 1000 --iterations 1000 api-test.js`:
```
          /\      |‾‾| /‾‾/   /‾‾/   
     /\  /  \     |  |/  /   /  /    
    /  \/    \    |     (   /   ‾‾\  
   /          \   |  |\  \ |  (‾)  | 
  / __________ \  |__| \__\ \_____/ .io

     execution: local
        script: api-test.js
        output: -

     scenarios: (100.00%) 1 scenario, 1000 max VUs, 10m30s max duration (incl. graceful stop):
              * default: 1000 iterations shared among 1000 VUs (maxDuration: 10m0s, gracefulStop: 30s)


     █ Create and get

       ✓ status is 200
       ✓ has urlOriginal
       ✓ has urlCode
       ✓ has expiresOn

     checks.........................: 100.00% ✓ 14000       ✗ 0     
     data_received..................: 5.1 MB  566 kB/s
     data_sent......................: 1.0 MB  114 kB/s
     group_duration.................: avg=6.05s   min=3.33s    med=6.02s    max=9.07s    p(90)=7.29s   p(95)=7.7s    
     http_req_blocked...............: avg=58.51µs min=1.46µs   med=5.13µs   max=17.81ms  p(90)=21.51µs p(95)=264.49µs
     http_req_connecting............: avg=44.31µs min=0s       med=0s       max=17.73ms  p(90)=0s      p(95)=177.81µs
   ✓ http_req_duration..............: avg=4.93ms  min=290.51µs med=603.18µs max=117.69ms p(90)=6.77ms  p(95)=17.25ms 
       { expected_response:true }...: avg=4.93ms  min=290.51µs med=603.18µs max=117.69ms p(90)=6.77ms  p(95)=17.25ms 
   ✓ http_req_failed................: 0.00%   ✓ 0           ✗ 11000 
     http_req_receiving.............: avg=43.38µs min=10.58µs  med=37.36µs  max=1.05ms   p(90)=67.57µs p(95)=81.6µs  
     http_req_sending...............: avg=25µs    min=6.74µs   med=16.64µs  max=3.33ms   p(90)=44.42µs p(95)=69.81µs 
     http_req_tls_handshaking.......: avg=0s      min=0s       med=0s       max=0s       p(90)=0s      p(95)=0s      
     http_req_waiting...............: avg=4.86ms  min=251.02µs med=544.32µs max=117.58ms p(90)=6.67ms  p(95)=17.18ms 
     http_reqs......................: 11000   1210.932199/s
     iteration_duration.............: avg=6.05s   min=3.33s    med=6.02s    max=9.07s    p(90)=7.29s   p(95)=7.7s    
     iterations.....................: 1000    110.084745/s
     vus............................: 2       min=2         max=1000
     vus_max........................: 1000    min=1000      max=1000
```


## Run

1. Clone the repository
2. `go build`
3. Set the environment variables:
   * `SHORTURL_POSTGRES_USER` (f.e `admin`)
   * `SHORTURL_POSTGRES_PASSWORD` (f.e `admin`)
   * `SHORTURL_SERVER_ADDRESS` (f.e `localhost:8000`)
     
   For HTTPS support:
      * `SHORTURL_HTTPS_ADDRESS` (f.e `localhost:8080`)
      * `SHORTURL_CERTFILE` (public key, path to the .crt/.pem file)
      * `SHORTURL_KEYFILE` (private key, path to the .key file)
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
