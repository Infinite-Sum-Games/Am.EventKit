import http from "k6/http";
import { check } from "k6";

export const options = {
	scenarios: {
		constant_load: {
			executor: "constant-arrival-rate",
			rate: 10000,
			timeUnit: "1s",
			duration: "10s",
			preAllocatedVUs: 1000,
			maxVUs: 2000,
		},
	},
	thresholds: {
		http_req_failed: ["rate<0.1"],
		http_req_duration: ["p(95)<200"], // 90% req under 200ms
	},
};

export default function () {
	const res = http.get("http://localhost:9000/api/v1/auth/user/check");
	check(res, {
		"status is 200": (r) => r.status === 200,
	});
}
