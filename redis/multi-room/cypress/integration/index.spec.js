/// <reference types="Cypress" />

context('Querying', () => {
  describe('複数サイト', () => {
    it('やってみる', () => {
      cy.visit('http://localhost:3000');
      cy.visit('http://localhost:3001');
    });
  });
});
