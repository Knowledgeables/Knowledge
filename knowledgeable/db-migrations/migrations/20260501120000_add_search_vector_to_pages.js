export async function up(knex) {
  await knex.schema.alterTable('pages', (table) => {
    table.specificType('search_vector', 'tsvector');
  });

  await knex.raw(`
    UPDATE pages
    SET search_vector =
      setweight(to_tsvector('english', coalesce(title,'')), 'A') ||
      setweight(to_tsvector('english', coalesce(content,'')), 'B')
  `);

  await knex.raw(`
    CREATE INDEX idx_pages_search
    ON pages USING GIN(search_vector)
  `);

  await knex.raw(`
    CREATE FUNCTION pages_search_vector_update() RETURNS trigger AS $$
    BEGIN
      NEW.search_vector :=
        setweight(to_tsvector('english', coalesce(NEW.title,'')), 'A') ||
        setweight(to_tsvector('english', coalesce(NEW.content,'')), 'B');
      RETURN NEW;
    END
    $$ LANGUAGE plpgsql;
  `);

  await knex.raw(`
    CREATE TRIGGER pages_search_vector_trigger
    BEFORE INSERT OR UPDATE ON pages
    FOR EACH ROW
    EXECUTE FUNCTION pages_search_vector_update();
  `);
}

export async function down(knex) {
  await knex.raw(`DROP TRIGGER IF EXISTS pages_search_vector_trigger ON pages`);
  await knex.raw(`DROP FUNCTION IF EXISTS pages_search_vector_update`);
  await knex.raw(`DROP INDEX IF EXISTS idx_pages_search`);

  await knex.schema.alterTable('pages', (table) => {
    table.dropColumn('search_vector');
  });
}