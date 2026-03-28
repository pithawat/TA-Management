// import http from 'k6/http';
// import { check, sleep } from 'k6';
// import { htmlReport } from "https://raw.githubusercontent.com/benc-uk/k6-reporter/main/dist/bundle.js";
// import { textSummary } from "https://jslib.k6.io/k6-summary/0.0.1/index.js";

// export const options = {
//     stages: [
//         { duration: '10s', target: 50 },
//         { duration: '20s', target: 50 },
//         { duration: '10s', target: 0 },
//     ],
//     thresholds: {
//         http_req_duration: ['p(95)<500'],
//         http_req_failed: ['rate<0.01'],
//     },
// };

// const BASE_URL = __ENV.BASE_URL || 'http://localhost:8084';
// // ใส่ Token ของคุณที่นี่ (หรือรับผ่าน ENV)
// const TOKEN = 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiI2NjAxNTE0NCIsImVtYWlsIjoiNjYwMTUxNDRAa21pdGwuYWMudGgiLCJuYW1lIjoiUGl0aGF3YXQgS2l0bW9uZ2tvbGNoYWkiLCJyb2xlIjoiU1RVREVOVCIsImlzcyI6ImV4YW1wbGUuY29tL2dvb2dsZWxvZ2luIiwiYXVkIjpbIlRBLW1hbmdlbWVudCJdLCJleHAiOjE3NzI3MzY2NjEsIm5iZiI6MTc3MjY1MDI2MSwiaWF0IjoxNzcyNjUwMjYxfQ.zBXzlgq3nX1OuWF2D_yG2OyO5c-BMiNPRkv9ZhGWaQ4';

// export default function () {
//     const url = `${BASE_URL}/TA-management/course`;

//     const params = {
//         headers: {
//             'Content-Type': 'application/json',
//             // 1. เปลี่ยนจาก Authorization เป็น Cookie
//             // หมายเหตุ: 'token' คือชื่อ key ของ cookie ที่ฝั่ง Backend คุณใช้ (เช่น 'access_token', 'jwt')
//             'Cookie': `auth_token=${TOKEN}`
//         },
//     };

//     const res = http.get(url, params);

//     check(res, {
//         'is status 200': (r) => r.status === 200,
//         'has array data': (r) => {
//             try {
//                 const body = JSON.parse(r.body);
//                 return Array.isArray(body.data) || body.data !== undefined;
//             } catch (e) { return false; }
//         },
//     });

//     sleep(1);
// }

// export function handleSummary(data) {
//     return {
//         "summary.html": htmlReport(data),
//         stdout: textSummary(data, { indent: " ", enableColors: true }),
//     };
// }
import http from 'k6/http';
import { check, sleep } from 'k6';
import { htmlReport } from "https://raw.githubusercontent.com/benc-uk/k6-reporter/main/dist/bundle.js";
import { textSummary } from "https://jslib.k6.io/k6-summary/0.0.1/index.js";

export const options = {
    scenarios: {
        // 1. Scenario สำหรับระดับคณะ (Load Test)
        department_load: {
            executor: 'ramping-vus',
            startVUs: 0,
            stages: [
                { duration: '20s', target: 500 }, // ไต่ระดับไป 500 คน (ระดับภาควิชา/คณะ)
                { duration: '1m', target: 500 },  // แช่ไว้ที่ 500 คน
                { duration: '20s', target: 0 },
            ],
            gracefulStop: '30s',
        },
    },
    thresholds: {
        'http_req_duration': ['p(95)<500'], // 95% ของการเรียกต้องเร็วกว่า 0.5 วินาที
        'http_req_failed': ['rate<0.01'],   // Error ต้องน้อยกว่า 1%
    },
};

const BASE_URL = __ENV.BASE_URL || 'http://host.docker.internal:8084';
const TOKEN = __ENV.TOKEN || 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiI2NjAxNTE0NCIsImVtYWlsIjoiNjYwMTUxNDRAa21pdGwuYWMudGgiLCJuYW1lIjoiUGl0aGF3YXQgS2l0bW9uZ2tvbGNoYWkiLCJyb2xlIjoiU1RVREVOVCIsImlzcyI6ImV4YW1wbGUuY29tL2dvb2dsZWxvZ2luIiwiYXVkIjpbIlRBLW1hbmdlbWVudCJdLCJleHAiOjE3NzI3MzY2NjEsIm5iZiI6MTc3MjY1MDI2MSwiaWF0IjoxNzcyNjUwMjYxfQ.zBXzlgq3nX1OuWF2D_yG2OyO5c-BMiNPRkv9ZhGWaQ4';

export default function () {
    const url = `${BASE_URL}/TA-management/course`;
    const params = {
        headers: {
            'Content-Type': 'application/json',
            'Cookie': `auth_token=${TOKEN}` // ส่ง Token ผ่าน Cookie ตามที่คุณต้องการ
        },
    };

    const res = http.get(url, params);

    check(res, {
        'is status 200': (r) => r.status === 200,
        'has data': (r) => JSON.parse(r.body).data !== undefined,
    });

    // หน่วงเวลา 1 วินาทีต่อ 1 รอบการทำงาน เพื่อไม่ให้ยิงหนักเกินไปจนเหมือนโดนยิง DDOS
    sleep(1);
}

export function handleSummary(data) {
    // แยกชื่อไฟล์ตาม ENV เพื่อความสะดวกในการเปรียบเทียบ Redis
    const reportName = __ENV.WITH_REDIS === 'true' ? "summary_redis_high_load.html" : "summary_no_redis_high_load.html";

    return {
        [reportName]: htmlReport(data),
        stdout: textSummary(data, { indent: " ", enableColors: true }),
    };
}