export async function up(knex) {
  await knex.raw(`
    ALTER TABLE pages
    ADD CONSTRAINT pages_url_not_wayback_wrapper
    CHECK (
      url !~ '^https?://web\\.archive\\.org/web/[0-9]+\\*?/(https?://.+)$'
    )
    NOT VALID
  `);

  await knex.raw(`
    ALTER TABLE pages
    VALIDATE CONSTRAINT pages_url_not_wayback_wrapper
  `);
}

export async function down(knex) {
  await knex.raw(`
    ALTER TABLE pages
    DROP CONSTRAINT IF EXISTS pages_url_not_wayback_wrapper
  `);
}
