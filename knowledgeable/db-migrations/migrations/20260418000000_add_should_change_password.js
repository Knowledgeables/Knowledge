export function up(knex) {
  return knex.schema
    .table('users', (table) => {
      table.boolean('should_change_password').notNullable().defaultTo(false);
    })
    .then(() => knex('users').update({ should_change_password: true }));
}

export function down(knex) {
  return knex.schema.table('users', (table) => {
    table.dropColumn('should_change_password');
  });
}
