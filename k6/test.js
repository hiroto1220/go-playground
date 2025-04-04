import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
    // stages: [
    //     { duration: '2s', target: 10}, // fast ramp-up to a high point
    //     // No plateau
    //     { duration: '1s', target: 0 }, // quick ramp-down to 0 users
    //   ],
  scenarios: {
    contacts: {
      executor: 'shared-iterations',
      vus: 10,
      iterations: 100000000,
      maxDuration: '1s',
    }
},
// scenarios: {
//     contacts: {
//       executor: 'ramping-vus',
//       startVUs: 0,
//       stages: [
//         { duration: '2s', target: 10 },
//         { duration: '1s', target: 0 },
//       ],
//       gracefulRampDown: '10s',
//     },
//   },
// scenarios: {
//     contacts: {
//       executor: 'ramping-arrival-rate',

//       // Start iterations per `timeUnit`
//       startRate: 300,

//       // Start `startRate` iterations per minute
//       timeUnit: '1s',

//       // Pre-allocate necessary VUs.
//       preAllocatedVUs: 400,


//       stages: [
//         // Start 300 iterations per `timeUnit` for the first minute.
//         { target: 2000, duration: '1s' },

//         // Linearly ramp-up to starting 600 iterations per `timeUnit` over the following two minutes.
//         { target: 2000, duration: '2s' },

//         // Continue starting 600 iterations per `timeUnit` for the following four minutes.
//         { target: 2000, duration: '4s' },

//         // Linearly ramp-down to starting 60 iterations per `timeUnit` over the last two minutes.
//         { target: 60, duration: '2s' },
//       ],
//     },
//   },
  thresholds: {
    // 90% of requests must finish within 400ms.
    http_req_duration: ['p(90) < 400'],
    http_req_failed: ['rate<0.01'],
  },
};

export default () => {
  const payload = JSON.stringify({
    name: 'test user',
  });

  const params = {
    headers: {
      'Content-Type': 'application/json',
    },
  };

  const res = http.post('http://localhost:8080/hello.Greeter/SayHello',payload, params);
  check(res, {
    'status is 200': (r) => r.status === 200,
  });
 
};