import http from "k6/http";
import { check, sleep } from "k6";
//import { Trend, Rate } from "k6/metrics";

// Personaliza tu endpoint:
const PORT = 3002; //3000 PARA go, 8000 para PYTHON
const BASE_URL = `http://localhost:${PORT}`; // o el host donde corre tu microservicio
const TASK_ID = 1;

export let options = {
  vus: 50, // usuarios concurrentes
  // duration: "30s", // duración total de la prueba

  stages: [
    { duration: "10s", target: 10 }, // Sube a 10 usuarios virtuales durante 10 segundos
    { duration: "30s", target: 10 }, // Mantiene 10 usuarios durante 30 segundos
    { duration: "10s", target: 0 }, // Baja gradualmente
  ],
};

export default function () {
  const url = `${BASE_URL}/tasks/${TASK_ID}/calculate`;

  const res = http.get(url);

  check(res, {
    "status is 200": (r) => r.status === 200,
    "body is not empty": (r) => r.body.length > 0,
  });

  sleep(1); // Espera 1 segundo entre solicitudes por usuario virtual
}

// stages: [
//   { duration: '10s', target: 10 }, // Sube a 10 usuarios virtuales durante 10 segundos
//   { duration: '30s', target: 10 }, // Mantiene 10 usuarios durante 30 segundos
//   { duration: '10s', target: 0 },  // Baja gradualmente
// ],
