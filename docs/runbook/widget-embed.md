# Widget embed rollout

Use this guide when a tenant owner is ready to publish the public widget.

## What the loader does

`frontend/public/widget-embed.js` mounts a fixed widget frame and opens:

`/widget?org_id=<tenant-org-id>`

The loader reads `data-org-id` from the `<script>` tag.
If `org_id` is missing, it logs a warning and stops.

## Before rollout

Make sure all of these are true before you hand the embed to a tenant:

- The widget runtime is deployed on its production origin.
- That browser origin is present in `WIDGET_RUNTIME_ORIGINS` on the backend.
- The allowlist is for the browser origin that calls
  `/api/v1/widget/*`, not every tenant site that hosts the embed tag.
- The widget runtime origin serves both `widget-embed.js` and `/widget`.

## Copy-paste embed

```html
<script
  src="https://widget.example.com/widget-embed.js"
  data-org-id="org_123">
</script>
```

Replace `https://widget.example.com` with the production widget runtime
origin and `org_123` with the tenant org id.

## Public runtime contract

- `org_id` is required in the widget URL.
- The public URL shape is `/widget?org_id=...`.
- `GET /api/v1/widget/settings` stays public-safe and only exposes
  `bot_name` and `welcome_message`.
- The direct `/widget` route without `org_id` shows a visible error state.
