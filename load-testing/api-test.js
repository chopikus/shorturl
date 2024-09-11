// import necessary modules
import { check, sleep, group } from 'k6';
import http from 'k6/http';
import { randomIntBetween } from 'https://jslib.k6.io/k6-utils/1.2.0/index.js';

// define configuration
export const options = {
  // define thresholds
  thresholds: {
    http_req_failed: ['rate<0.01'], // http errors should be less than 1%
    http_req_duration: ['p(99)<1000'], // 99% of requests should be below 1s

  },
};

let counter = 0;
function generateUrl() {
    counter += 1;
    return "https://example" + counter + ".lol";
}

export default function () {
  group('Create and get', function () {
        sleep(Math.random());
        const params = {
          headers: {
              'Content-Type': 'application/json',
          }
        };
        const payload = JSON.stringify({
            urlOriginal: generateUrl()
        });

        const url = 'http://localhost:8000/api/new';
        const res = http.post(url, payload, params);

        check(res, {
            'status is 200': (r) => r.status === 200,
            'has urlOriginal': (r) => r.json().hasOwnProperty('urlOriginal'),
            'has urlCode': (r) => r.json().hasOwnProperty('urlCode'),
            'has expiresOn': (r) => r.json().hasOwnProperty('expiresOn'),
        });
        
        sleep(Math.random());
        const urlCode = res.json().urlCode;
        
        /* Each link is shared to 10 different people who open a link within 10 seconds. */
        for (let i = 0; i < 10; i++) {
            const res = http.get('http://localhost:8000/' + urlCode);
            check(res, { 'status is 200': (r) => r.status === 200 });

            sleep(Math.random());
        }
  });
}
