export async function up(knex) {
  await knex.schema.createTable('password_reset_tokens', (table) => {
    table.increments('id').primary();
    table.integer('user_id').notNullable().references('id').inTable('users').onDelete('CASCADE');
    table.text('token_hash').notNullable().unique();
    table.timestamp('expires_at', { useTz: true }).notNullable();
    table.timestamp('used_at', { useTz: true });
    table.timestamp('created_at', { useTz: true }).notNullable().defaultTo(knex.fn.now());
  });

  await knex.schema.raw('CREATE INDEX idx_password_reset_tokens_hash ON password_reset_tokens(token_hash)');
  await knex.schema.raw('CREATE INDEX idx_password_reset_tokens_user ON password_reset_tokens(user_id)');
}

export async function down(knex) {
  await knex.schema.dropTableIfExists('password_reset_tokens');
}
