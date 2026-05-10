import { expect, test } from '@playwright/test';

const TENANT_ID = '31dc1d57-e092-44a2-97e2-7bf07018efcb';
const EMBED_HARNESS_PATH = '/widget-embed-harness.html';
const EMBED_MISSING_ORG_HARNESS_PATH = '/widget-embed-missing-org.html';

async function mockWidgetSuccess(page) {
  let sessionOrgId = '';
  let settingsOrgId = '';

  await page.route('**/api/v1/widget/settings**', async (route) => {
    const requestUrl = new URL(route.request().url());
    settingsOrgId = requestUrl.searchParams.get('org_id') ?? '';

    await route.fulfill({
      status: 200,
      contentType: 'application/json',
      body: JSON.stringify({
        bot_name: 'Acme Support',
        welcome_message: 'Welcome to Acme Support',
      }),
    });
  });

  await page.route('**/api/v1/widget/sessions', async (route) => {
    sessionOrgId = route.request().headers()['x-org-id'] ?? '';

    await route.fulfill({
      status: 200,
      contentType: 'application/json',
      body: JSON.stringify({ id: 'session-smoke-1' }),
    });
  });

  return {
    get settingsOrgId() {
      return settingsOrgId;
    },
    get sessionOrgId() {
      return sessionOrgId;
    },
  };
}

async function mockWidgetFailure(page) {
  await page.route('**/api/v1/widget/settings**', async (route) => {
    await route.fulfill({
      status: 404,
      contentType: 'application/json',
      body: JSON.stringify({ error: 'Tenant not found.' }),
    });
  });

  await page.route('**/api/v1/widget/sessions', async (route) => {
    await route.fulfill({
      status: 400,
      contentType: 'application/json',
      body: JSON.stringify({ error: 'Tenant not found.' }),
    });
  });
}

async function installInstantWebSocket(page) {
  await page.addInitScript(() => {
    class InstantWebSocket {
      url = '';
      readyState = 0;
      onopen = null;
      onclose = null;
      onerror = null;
      onmessage = null;

      constructor(url) {
        this.url = url;
        setTimeout(() => {
          this.readyState = 1;
          if (typeof this.onopen === 'function') {
            this.onopen(new Event('open'));
          }
        }, 0);
      }

      send() {}

      close() {
        this.readyState = 3;
        if (typeof this.onclose === 'function') {
          this.onclose(new Event('close'));
        }
      }
    }

    window.WebSocket = InstantWebSocket as unknown as typeof WebSocket;
  });
}

test('embed harness mounts the widget iframe url', async ({ page }) => {
  await mockWidgetSuccess(page);
  await installInstantWebSocket(page);

  await page.goto(EMBED_HARNESS_PATH);

  const frame = page.locator('#mychatbot-widget-embed-frame');
  await expect(frame).toHaveAttribute('src', new RegExp(`/widget\\?org_id=${TENANT_ID}$`));
});

test('embed loader warns and does not mount without data-org-id', async ({ page }) => {
  let sawMissingOrgWarning = false;

  page.on('console', (message) => {
    if (message.type() === 'warning' && message.text() === '[widget-embed] Missing org_id.') {
      sawMissingOrgWarning = true;
    }
  });

  await page.goto(EMBED_MISSING_ORG_HARNESS_PATH);

  await expect(page.locator('#mychatbot-widget-embed-frame')).toHaveCount(0);
  await expect.poll(() => sawMissingOrgWarning).toBe(true);
});

test('widget route renders branding, welcome, and session init success', async ({ page }) => {
  const widgetApi = await mockWidgetSuccess(page);
  await installInstantWebSocket(page);

  await page.goto(`/widget?org_id=${TENANT_ID}`);

  await expect.poll(() => widgetApi.settingsOrgId).toBe(TENANT_ID);
  await expect.poll(() => widgetApi.sessionOrgId).toBe(TENANT_ID);
  await expect(page.getByRole('heading', { level: 3, name: 'Acme Support' })).toBeVisible();
  await expect(page.getByText('Welcome to Acme Support')).toBeVisible();
  await expect(page.getByText('Online')).toBeVisible();
  await expect(page.getByPlaceholder('Type your message...')).toBeEnabled();
  await expect(page.getByRole('alert')).toHaveCount(0);

  expect(widgetApi.settingsOrgId).toBe(TENANT_ID);
  expect(widgetApi.sessionOrgId).toBe(TENANT_ID);
});

test('invalid widget tenant fails gracefully without sign-in redirect', async ({ page }) => {
  await mockWidgetFailure(page);

  await page.goto('/widget?org_id=invalid-tenant');

  await expect(page).toHaveURL(/\/widget\?org_id=invalid-tenant$/);
  await expect(page.getByRole('alert')).toContainText('Tenant not found.');
  await expect(page.getByRole('heading', { level: 3, name: 'Customer Support' })).toBeVisible();
  await expect(page.getByText('Hello! How can we help you today?')).toBeVisible();
  await expect(page.getByText(/sign in/i)).toHaveCount(0);
});
