import { onCLS, onLCP, onINP, onTTFB, type Metric } from 'web-vitals';

function sendToConsole(metric: Metric) {
  console.log(`[Web Vitals] ${metric.name}:`, metric.value, metric);
}

export function reportWebVitals() {
  onCLS(sendToConsole);
  onLCP(sendToConsole);
  onINP(sendToConsole);
  onTTFB(sendToConsole);
}