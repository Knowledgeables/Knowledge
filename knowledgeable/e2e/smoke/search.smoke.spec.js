import { test, expect } from '@playwright/test';

test('smoke search works', async ({ page }) => {
  await page.goto('/search');

  await page.getByPlaceholder('Search...').fill('E2E_SEARCH_TARGET');
  await page.click('#search-button');

  const results = page.locator('#results a');

  // at leat one result
  await expect(results.first()).toBeVisible({ timeout: 10000 });

  // count text is being shown
  await expect(page.locator('#results')).toContainText('results for');
});