import { test, expect } from '@playwright/test';

test.describe('Knowledge Page', () => {
  test.beforeEach(async ({ page }) => {
    // Navigate to the knowledge page. Adjust the URL based on the routing.
    // Assuming the app runs on localhost:5173
    await page.goto('/knowledge');
  });

  test('should render action buttons', async ({ page }) => {
    const uploadBtn = page.getByRole('button', { name: 'Upload Document' });
    const syncBtn = page.getByRole('button', { name: 'URL Sync' });
    const qaBtn = page.getByRole('button', { name: 'Manual Q&A' });

    await expect(uploadBtn).toBeVisible();
    await expect(syncBtn).toBeVisible();
    await expect(qaBtn).toBeVisible();
  });

  test('should render knowledge table with correct headers', async ({ page }) => {
    await expect(page.getByText('Trained Data')).toBeVisible();
    
    // Check table headers
    const headers = ['Nama File', 'Tipe', 'Status', 'Tanggal'];
    for (const header of headers) {
      await expect(page.getByRole('columnheader', { name: header })).toBeVisible();
    }
  });

  test('should render table data with status badges', async ({ page }) => {
    // Wait for data to load (assuming mock data is fast, but just in case)
    await page.waitForSelector('text=company-policy.pdf');

    // Check for row content
    await expect(page.getByText('company-policy.pdf')).toBeVisible();
    await expect(page.getByText('Document', { exact: true }).first()).toBeVisible();
    await expect(page.getByText('2026-04-22')).toBeVisible();

    // Check badges
    await expect(page.getByText('Ready')).toBeVisible();
    await expect(page.getByText('Queued')).toBeVisible();
    await expect(page.getByText('Processing')).toBeVisible();
    await expect(page.getByText('Failed')).toBeVisible();
  });
});
