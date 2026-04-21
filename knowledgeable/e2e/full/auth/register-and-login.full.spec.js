// e2e/full/app.full.spec.js
import { test, expect } from "@playwright/test";

test("user can log in", async ({ page, request, baseURL }) => {
  const username = `test_${Date.now()}`;
  const password = "test123";

  const registerRes = await request.post(`${baseURL}/api/register`, {
    data: { username, email: `${username}@test.com`, password },
  });

  expect(registerRes.ok()).toBeTruthy();

  await page.goto("/login");
  await page.getByLabel("Username").fill(username);
  await page.getByLabel("Password").fill(password);
  await page.getByRole("button", { name: /log in/i }).click();

  await expect(page).toHaveURL(/\/$/);

  const meRes = await page.request.get("/api/me");
  expect(meRes.status()).toBe(200);
});
