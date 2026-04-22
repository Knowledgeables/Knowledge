import { test, expect } from '@playwright/test';

test('search returns results', async ({ page }) => {
  await page.goto('/search');

  await page.getByPlaceholder('Search...').fill('E2E_UNIQUE_SEARCH');
await page.click('#search-button');

await expect(page.locator('#results')).toContainText('E2E_UNIQUE_SEARCH');
});

