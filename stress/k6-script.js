import http from 'k6/http';
import { sleep } from 'k6';

export default function () {
    http.get('http://rate-limit:8080/hello');
}

