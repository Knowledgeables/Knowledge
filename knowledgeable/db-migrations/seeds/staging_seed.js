export async function seed(knex) {
await knex('pages').del();
await knex('pages').insert([
  {
    title: 'Search Systems',
    url: 'e2e://search-systems',
    content: 'search relevance ranking'
  },
  {
    title: 'Security Basics',
    url: 'e2e://security',
    content: 'authentication authorization'
  },
  {
    title: 'Redis Cache',
    url: 'e2e://redis',
    content: 'in-memory storage'
  },
  {
    title: 'Simple Strings',
    url: 'e2e://strings',
    content: 'string manipulation'
  },
  {
  title: 'E2E_UNIQUE_SEARCH',
  url: 'e2e://unique',
  content: 'only-this-should-match-xyz123'
}
]);
}