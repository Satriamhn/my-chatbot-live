import { test, expect } from '@playwright/test';

test.describe('Inbox Page', () => {
  test.beforeEach(async ({ page }) => {
    // Navigate to the inbox page
    await page.goto('/inbox');
  });

  test('displays omnichannel chat UI with roles', async ({ page }) => {
    // Wait for the conversations to load (mock data should load fast)
    await page.waitForSelector('text=Chat with');

    // Check that we have a message list containing messages
    const messageList = page.locator('.custom-scrollbar').nth(1); // 0 is sidebar, 1 is message list
    await expect(messageList).toBeVisible();

    // The mock data usually has a bot message and user message
    // Check that takeover button is present
    const takeoverBtn = page.getByTestId('takeover-button');
    await expect(takeoverBtn).toBeVisible();
    await expect(takeoverBtn).toHaveText('Take Over (Human Mode)');

    // Ensure bot input is disabled initially (since humanMode is false)
    const messageInput = page.getByTestId('message-input');
    await expect(messageInput).toBeDisabled();
    
    // Check for bot message metadata
    const botMetadata = page.getByTestId('bot-metadata').first();
    if (await botMetadata.isVisible()) {
      await expect(botMetadata).toContainText('Source:');
      await expect(botMetadata).toContainText('Confidence:');
    }

    // Click "Take Over"
    await takeoverBtn.click();
    
    // Validate state updates
    await expect(takeoverBtn).toHaveText('Human Mode Active');
    await expect(messageInput).toBeEnabled();
    await expect(messageInput).toHaveAttribute('placeholder', 'Type your message...');
  });
});
