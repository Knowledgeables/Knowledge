export function up(knex) {
  return knex.schema
    .createTable('users', (table) => {
      table.increments('id'); 
      table.string('username').notNullable().unique();
      table.string('email').notNullable().unique();
      table.string('password_hash').notNullable();
      table.timestamp('created_at').notNullable().defaultTo(knex.fn.now());
    })
    .createTable('pages', (table) => {
      table.increments('id');
      table.string('title').notNullable();
      table.string('url').notNullable().unique();
      table.string('language').notNullable().defaultTo('en');
      table.text('content').notNullable();
      table.timestamp('last_updated').notNullable().defaultTo(knex.fn.now());
    })
    .then(() => knex.schema.raw(
      'CREATE INDEX idx_pages_language ON pages(language)'
    ))
    .then(() => knex.schema.raw(
      'CREATE INDEX idx_pages_title ON pages(title)'
    ));
}

export function down(knex) {
  return knex.schema
    .dropTableIfExists('pages')
    .dropTableIfExists('users');
}