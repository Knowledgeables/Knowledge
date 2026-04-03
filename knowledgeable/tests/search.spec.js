import { test, expect } from '@playwright/test';

test('user can use search function', async ({ page }) => {
  await page.goto('/');


  await page.getByRole('link', { name: 'Go to Search Engine' }).click();


  await expect(page).toHaveURL('/search');

  await page.getByPlaceholder('Search').fill('g');

  await page.getByRole('button', { name: 'Search' }).click();

  await expect(page.getByText('Go Basics')).toBeVisible();
});