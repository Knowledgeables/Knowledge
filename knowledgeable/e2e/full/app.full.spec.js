// e2e/login.spec.js
import { test, expect, request } from '@playwright/test';

test('user can log in with valid credentials', async ({ page, baseURL }) => {
  const username = `testuser-${Date.now()}`;
  const password = 'testpassword123';

  const api = await request.newContext();
  await api.post(`${baseURL}/api/register`, {
    headers: { 'Content-Type': 'application/json' },
    data: { username, email: `${username}@test.com`, password }
  });

  await page.goto('/login');
  await page.getByLabel('Username').fill(username);
  await page.getByLabel('Password').fill(password);
  await page.getByRole('button', { name: 'Log in' }).click();
  await expect(page).toHaveURL(/\/$/, { timeout: 10000 });
});