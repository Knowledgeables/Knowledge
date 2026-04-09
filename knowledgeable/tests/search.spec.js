import { test, expect } from '@playwright/test';

test('search works @smoke', async ({ page }) => {
  await page.goto('/search');

  await page.getByPlaceholder('Search...').fill('g');

  await page.locator('#search-button').click();

  await expect(page.locator('#results')).toContainText('results');
});