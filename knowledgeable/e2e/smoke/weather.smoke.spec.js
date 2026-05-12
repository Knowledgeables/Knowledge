import { test, expect } from '@playwright/test';

test('smoke weather page loads', async ({ page }) => {
  await page.goto('/weather');
  await expect(page).toHaveTitle('Weather');
  await expect(page.locator('#navbar')).toBeVisible();
});

test('smoke weather page shows Copenhagen label', async ({ page }) => {
  await page.goto('/weather');
  await expect(page.getByText('Copenhagen')).toBeVisible();
});

test('smoke weather page displays weather data', async ({ page }) => {
  await page.goto('/weather');
  await expect(page.locator('#temp')).not.toBeEmpty({ timeout: 10000 });
  await expect(page.locator('#wind')).not.toBeEmpty();
  await expect(page.locator('#humidity')).not.toBeEmpty();
});

test('smoke weather api returns valid data', async ({ request }) => {
  const res = await request.get('/api/weather?lat=55.6761&lon=12.5683');
  expect(res.status()).toBe(200);

  const data = await res.json();
  expect(typeof data.temperature_c).toBe('number');
  expect(typeof data.wind_speed_ms).toBe('number');
  expect(typeof data.humidity_pct).toBe('number');
  expect(typeof data.description).toBe('string');
  expect(data.description.length).toBeGreaterThan(0);
});

test('smoke weather api rejects missing params', async ({ request }) => {
  const res = await request.get('/api/weather');
  expect(res.status()).toBe(400);
});

test('smoke weather api rejects invalid coordinates', async ({ request }) => {
  const res = await request.get('/api/weather?lat=999&lon=12.57');
  expect(res.status()).toBe(400);
});
