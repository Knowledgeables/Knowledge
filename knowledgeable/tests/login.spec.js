import { test, expect } from '@playwright/test';

test('user can login successfully', async ({ page }) => {
    await page.goto('/login'); // bruger baseURL fra config

    await page.getByRole('textbox', { name: 'username' }).fill('admin');
    await page.getByRole('textbox', { name: 'password' }).fill('secret');

    await page.getByRole('button', { name: 'Login' }).click();

    // tjekker efter Welcome på dashboard siden
    await expect(page.locator('body')).toContainText('Welcome');
});