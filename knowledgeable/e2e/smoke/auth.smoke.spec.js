import { test, expect } from '@playwright/test';

test('@smoke-ci login works', async ({ page }) => {
  await page.goto('/login');

  await page.getByLabel('Username').fill('admin');
  await page.getByLabel('Password').fill('secret');

  await page.getByRole('button', { name: 'Log in' }).click();

  await page.waitForLoadState('networkidle');

  await expect(page).toHaveURL(/\/$/);
});