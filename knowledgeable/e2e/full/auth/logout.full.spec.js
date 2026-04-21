import { test, expect } from '@playwright/test';

test("user can log out", async ({ page, request }) => {
  const username = `test_${Date.now()}`;
  const password = "test123";

  await request.post("/api/register", {
    data: { username, email: `${username}@test.com`, password },
  });

  await page.goto("/login");
  await page.getByLabel("Username").fill(username);
  await page.getByLabel("Password").fill(password);
  await page.getByRole("button", { name: /log in/i }).click();

  // logout
  await page.locator('form[action="/logout"] button').click();

  await page.waitForLoadState("networkidle"); 


});