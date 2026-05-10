import { test, expect } from '@playwright/test';

test('smoke test - navigation', async ({ page }) => {
  // Go to home page
  await page.goto('/');
  // Basic check that home page rendered (AppLayout + Home)
  // We can look for common elements if we knew them, but let's check URL for now
  await expect(page).toHaveURL('http://localhost:5173/');

  // Go to signin page
  await page.goto('/signin');
  await expect(page).toHaveURL(/.*signin/);
  
  // Verify signin elements exist (assuming basic form)
  // await expect(page.getByRole('button', { name: /sign in/i })).toBeVisible();

  // Go back to home and check sidebar if possible
  await page.goto('/');
  // Check for sidebar links (usually Bot Settings, Knowledge Base, etc.)
  // Just checking if any sidebar-like link exists
  const sidebarLinks = page.locator('nav a');
  await expect(sidebarLinks.first()).toBeVisible();
});
