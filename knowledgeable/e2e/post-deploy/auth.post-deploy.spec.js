import { test, expect } from '@playwright/test';

test('login works', async ({ page }) => {
  await page.goto('/login');

  await page.getByLabel('Username').fill(process.env.TEST_USER_USERNAME);
  await page.getByLabel('Password').fill(process.env.TEST_USER_PASSWORD);

  await page.getByRole('button', { name: 'Log in' }).click();

  await expect(page).toHaveURL(/\/$/, { timeout: 10000 });
  
});