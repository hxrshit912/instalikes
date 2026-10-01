import http from 'k6/http';
import { check } from 'k6';

export const options = {
    vus: 200000,
    duration: '1s',
};

export default function () {
    const res = http.post('http://localhost:8080/like');

    check(res, {
        'status is 200': (r) => r.status === 200,
    });
}
