export function up(knex) {
  return knex.schema.createTable('crawl_signal_events', (table) => {
    table.increments('id');
    table.string('query').notNullable();
    table.string('language').notNullable().defaultTo('en');
    table.timestamp('seen_at').notNullable().defaultTo(knex.fn.now());

    table.index(['seen_at']);
    table.index(['query', 'language', 'seen_at']);
  });
}

export function down(knex) {
  return knex.schema.dropTableIfExists('crawl_signal_events');
}
