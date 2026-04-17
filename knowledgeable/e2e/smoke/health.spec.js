import { test, expect } from '@playwright/test';

test('@smoke-ci health endpoint returns ok', async ({ request }) => {
  const res = await request.get('/health');

  expect(res.status()).toBe(200);

  const body = await res.text();
  expect(body.trim()).toBe('ok');
});