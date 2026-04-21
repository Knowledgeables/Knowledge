import { test, expect } from '@playwright/test';

test("user can log out", async ({ page, request, baseURL }) => {
  const username = `test_${Date.now()}`;
  const password = "test123";

 
  const res = await request.post(`${baseURL}/api/register`, {
    data: { username, email: `${username}@test.com`, password },
  });

  expect(res.ok()).toBeTruthy(); 


  await page.goto('/login');
  await page.waitForLoadState('networkidle');

  await page.getByRole('textbox', { name: /username/i }).fill(username);
  await page.getByRole('textbox', { name: /password/i }).fill(password);
  await page.getByRole('button', { name: /log in/i }).click();


  await expect(page.locator('form[action="/logout"] button')).toBeVisible();

 
  await Promise.all([
    page.waitForNavigation({ waitUntil: 'networkidle' }),
    page.locator('form[action="/logout"] button').click(),
  ]);


  await expect(page).toHaveURL(/login/);
});