import { test, expect } from "@playwright/test";

  test("can redirect to login-page", async ({ page }) => {
    await page.goto("/");
    await expect(page.locator("#login-button")).toBeVisible();
    await page.locator("#login-button").click();
    await expect(page).toHaveURL(/login/);
  });

  test("anonymous user can navigate to home via navbar", async ({ page }) => {
    await page.goto("/");

    const homeBtn = page.locator("#nav-home");
    await expect(homeBtn).toBeVisible();
    await expect(homeBtn).toBeEnabled();

    await homeBtn.click();
    await expect(page).toHaveURL("/");
  });

  test("anonymous user can navigate to search via navbar", async ({ page }) => {
    await page.goto("/");

    const searchBtn = page.locator("#nav-search");
    await expect(searchBtn).toBeVisible();
    await expect(searchBtn).toBeEnabled();

    await searchBtn.click();
    await expect(page).toHaveURL(/search/);
  });

  test("anonymous user can navigate to login via navbar", async ({ page }) => {
    await page.goto("/");

    const loginBtn = page.locator("#login-button");
    await expect(loginBtn).toBeVisible();
    await expect(loginBtn).toBeEnabled();

    await loginBtn.click();
    await expect(page).toHaveURL(/login/);
  });


test("authenticated user sees user-navbar", async ({ page, request, baseURL }) => {
  const username = `test_${Date.now()}`;
  const password = "test123";

  // register
 const res = await request.post(`${baseURL}/api/register`, {
  data: { username, email: `${username}@test.com`, password },
});

expect(res.ok()).toBeTruthy();

  // login via UI
  await page.goto("/login");
  await page.getByLabel("Username").fill(username);
  await page.getByLabel("Password").fill(password);
  await page.getByRole("button", { name: /log in/i }).click();

  await expect(page).toHaveURL(/\/$/);

  // assertions
  await expect(page.getByText(username)).toBeVisible();
  await expect(page.locator('form[action="/logout"] button')).toBeVisible();
});