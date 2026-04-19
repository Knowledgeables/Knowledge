import { test, expect } from '@playwright/test';

test('@app loads', async ({ page }) => {
  await page.goto('/');

  await page.waitForLoadState('networkidle');

  await expect(page.locator('#navbar')).toBeVisible();
});