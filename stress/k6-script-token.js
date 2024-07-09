import http from 'k6/http';
import { sleep } from 'k6';

export default function () {
    const URL = 'http://rate-limit:8080/hello';
    const PARAMS = {
        headers: {
            'API_KEY': 'Token013'
        },
    };
    http.get(URL, PARAMS);
}

