import http from 'k6/http';
import { check, sleep } from 'k6';
import { randomIntBetween } from 'https://jslib.k6.io/k6-utils/1.1.0/index.js';

export let options = {
    scenarios: {
        contacts: {
            executor: 'constant-arrival-rate',
            rate: 10000, // 100 RPS
            timeUnit: '1s', // rate specified for 1 second
            duration: '1m', // run the test for 1 minute
            preAllocatedVUs: 10000, // number of VUs to pre-allocate
            maxVUs: 20000, // maximum number of VUs that can be allocated during the test
        },
    },
};

function maybeAddQueryParam(baseURL, paramName, paramValue) {
    if (Math.random() > 0.4) {
        const separator = baseURL.includes('?') ? '&' : '?';
        return `${baseURL}${separator}${paramName}=${encodeURIComponent(paramValue)}`;
    }
    return baseURL;
}

export default function () {
    let url = 'http://127.0.0.1:8080/api/v1/ad';

    url = maybeAddQueryParam(url, 'age', randomIntBetween(1, 100));
    url = maybeAddQueryParam(url, 'gender', Math.random() > 0.5 ? 'M' : 'F');
    
    const countries = ['US', 'CA', 'JP', 'TW'];
    url = maybeAddQueryParam(url, 'country', countries[Math.floor(Math.random() * countries.length)]);

    const platforms = ['android', 'ios', 'web'];
    url = maybeAddQueryParam(url, 'platform', platforms[Math.floor(Math.random() * platforms.length)]);

    const res = http.get(url);
    
    check(res, {
        'is status 200': (r) => r.status === 200,
    });
}
