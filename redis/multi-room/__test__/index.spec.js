//// <reference types="Cypress" />

// context('Querying', () => {
describe('複数サイト', () => {
  it('やってみる', async () => {
    const page = await browser.newPage();
    await page.goto('http://localhost:3000');
    await page.goto('http://localhost:3001');
  });
});
// });
