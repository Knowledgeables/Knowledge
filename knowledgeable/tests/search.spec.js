import { test, expect } from '@playwright/test';


test('logged in user can use search function', async ({ page }) => {
    await page.goto('/');

    await page.getByRole('button', { name: 'Go to Search Engine' }).click();

    await page.getByRole('textbox', { name: 'search' }).fill('g');

    await page.getByRole('button', { name: 'Search' }).click();

    await expect(page.locator('body')).toContainText('Go Basics');
});