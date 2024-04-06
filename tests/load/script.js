import http from 'k6/http';
import { check } from 'k6';
import { randomIntBetween } from 'https://jslib.k6.io/k6-utils/1.1.0/index.js';

export let options = {
    scenarios: {
        ramping_test: {
            executor: 'ramping-arrival-rate',
            startRate: 10, // start at 10 RPS
            timeUnit: '1s',
            preAllocatedVUs: 100, // number of VUs to pre-allocate
            maxVUs: 2000, // maximum number of VUs that can be allocated during the test
            stages: [
                { target: 1000, duration: '1m' }, // ramp to 100 RPS over 2 minutes
                { target: 1000, duration: '30s' }, // hold 100 RPS for 3 minutes
            ],
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
    let url = 'https://dcard-backend-33cjpqddvq-de.a.run.app/api/v1/ad';

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
