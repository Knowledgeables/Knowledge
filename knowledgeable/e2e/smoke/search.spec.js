import { test, expect } from '@playwright/test';

test('@smoke-ci search works', async ({ page }) => {
  await page.goto('/search');

  await page.getByPlaceholder('Search...').fill('g');
  await page.click('#search-button');

  // wait for navigation
  await page.waitForLoadState('networkidle');

  const results = page.locator('#results a');

  // at leat one result
  await expect(results.first()).toBeVisible();

  // count text is being shown
  await expect(page.locator('#results')).toContainText('results for');
});