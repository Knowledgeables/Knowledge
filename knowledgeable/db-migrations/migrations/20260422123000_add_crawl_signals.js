export function up(knex) {
  return knex.schema.createTable('crawl_signals', (table) => {
    table.increments('id');
    table.string('query').notNullable();
    table.string('language').notNullable().defaultTo('en');
    table.integer('signal_count').notNullable().defaultTo(1);
    table.timestamp('first_seen').notNullable().defaultTo(knex.fn.now());
    table.timestamp('last_seen').notNullable().defaultTo(knex.fn.now());

    table.unique(['query', 'language']);
    table.index(['last_seen']);
    table.index(['signal_count']);
  });
}

export function down(knex) {
  return knex.schema.dropTableIfExists('crawl_signals');
}
