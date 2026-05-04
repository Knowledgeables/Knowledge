import { test, expect } from '@playwright/test';

test('smoke search works', async ({ page }) => {
  await page.goto('/search');

  await page.getByPlaceholder('Search...').fill('g');
  await page.click('#search-button');

  const results = page.locator('#results a');

  await expect(results.first()).toBeVisible({ timeout: 10000 });
});