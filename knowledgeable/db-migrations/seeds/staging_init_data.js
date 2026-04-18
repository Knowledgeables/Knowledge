// seeds/staging_init_data.js

/**
 * @param {import('knex').Knex} knex
 */
export async function seed(knex) {

 await knex('pages').truncate();
 await knex('users').truncate(); 

  // USERS
  await knex('users')
    .insert([
      {
        username: 'e2e_user',
        email: 'e2e@staging.local',
        // password: "password123"
        password_hash: '$2a$12$wH8Q8G0dQ2Q0qV2c9l6s2e2rG9x0k7v5n8wYcV3f7u9l6qQ8bG3y2',
        should_change_password: true
      },
      {
        username: 'admin',
        email: 'admin@staging.local',
        // password: "admin123"
        password_hash: '$2a$12$8fhiDpxQAlrTsZfQYLiKUeELzXABImIynuZSx3nwyNfIOofEg1W3i',
        should_change_password: true
      }
    ])

  // PAGES 
  await knex('pages')
    .insert([
      {
        title: 'E2E_SEARCH_TARGET',
        url: 'https://staging.local/e2e-search',
        language: 'en',
        content: 'playwright-test-content'
      },
      {
        title: 'E2E_EMPTY_CASE',
        url: 'https://staging.local/empty',
        language: 'en',
        content: ''
      },
      {
        title: 'E2E_SPECIAL_CHARS',
        url: 'https://staging.local/special',
        language: 'en',
        content: '!@#$%^&*() test'
      },
      {
        title: 'E2E_LONG_CONTENT',
        url: 'https://staging.local/long',
        language: 'en',
        content: 'lorem ipsum '.repeat(100)
      },
      {
        title: 'E2E_NAVIGATION_TARGET',
        url: 'https://staging.local/nav',
        language: 'en',
        content: 'navigation test page'
      }
    ])
}