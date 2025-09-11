import { test, expect } from '@playwright/test';

test('has library text', async ({ page }) => {
    await page.goto('http://localhost:5173/')

    await expect(page.getByText('Library')).toBeVisible();
});