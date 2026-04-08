import { test, expect } from '@playwright/test';

test('login works @smoke', async ({ page }) => {
    await page.goto('/login');

    await page.getByLabel('Username').fill('admin');
    await page.getByLabel('Password').fill('secret');

    await page.getByRole('button', { name: 'Log in' }).click();

    await expect(page).toHaveURL(/\/$/);

    await expect(page.getByText('Search. Unlock knowledge.')).toBeVisible();
});