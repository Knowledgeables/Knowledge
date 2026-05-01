import { test, expect } from '@playwright/test';

test('search returns results', async ({ page }) => {
  await page.goto('/search');

  await page.fill('#search-input', 'E2E_UNIQUE_SEARCH');
  await page.click('#search-button');

  const results = page.locator('[data-testid="result-item"]');

  await expect(results).toHaveCount(1);
  await expect(results.first()).toContainText('E2E_UNIQUE_SEARCH');
});

test('returns unique result', async ({ page }) => {
  await page.goto('/search');

  await page.fill('#search-input', 'E2E_UNIQUE_SEARCH');
  await page.press('#search-input', 'Enter');

  const results = page.locator('[data-testid="result-item"]');

  await expect(results).toHaveCount(1);
  await expect(results.first()).toContainText('E2E_UNIQUE_SEARCH');
});

test('no results', async ({ page }) => {
  await page.goto('/search');

  await page.fill('#search-input', 'does-not-exist');
  await page.press('#search-input', 'Enter');

  await expect(page.locator('[data-testid="no-results"]')).toBeVisible();
});

test('trims input', async ({ page }) => {
  await page.goto('/search');

  await page.fill('#search-input', 'E2E_UNIQUE_SEARCH');
  await page.press('#search-input', 'Enter');

  const results = page.locator('[data-testid="result-item"]');

  await expect(results).toHaveCount(1);
  await expect(results.first()).toContainText('E2E_UNIQUE_SEARCH');
});

test('click redirects', async ({ page, context }) => {
  await page.goto('/search');

  await page.fill('#search-input', 'E2E_UNIQUE_SEARCH');
  await page.press('#search-input', 'Enter');

  const link = page.locator('[data-testid="result-item"]').first();

  const [newPage] = await Promise.all([
    context.waitForEvent('page'),
    link.click()
  ]);

});