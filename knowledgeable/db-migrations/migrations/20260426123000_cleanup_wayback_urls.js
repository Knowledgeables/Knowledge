export async function up(knex) {
  const pattern = '^https?://web\\.archive\\.org/web/[0-9]+\\*?/(https?://.+)$';

  await knex.transaction(async (trx) => {
    await trx.raw(
      `
      WITH normalized AS (
        SELECT
          id,
          regexp_replace(url, ?, '\\1') AS normalized_url
        FROM pages
        WHERE url ~ ?
      ),
      conflicts AS (
        SELECT n.id
        FROM normalized n
        JOIN pages p ON p.url = n.normalized_url AND p.id <> n.id
      )
      DELETE FROM pages
      WHERE id IN (SELECT id FROM conflicts)
      `,
      [pattern, pattern]
    );

    await trx.raw(
      `
      WITH normalized AS (
        SELECT
          id,
          regexp_replace(url, ?, '\\1') AS normalized_url,
          ROW_NUMBER() OVER (
            PARTITION BY regexp_replace(url, ?, '\\1')
            ORDER BY last_updated DESC, id DESC
          ) AS rn
        FROM pages
        WHERE url ~ ?
      )
      DELETE FROM pages
      WHERE id IN (SELECT id FROM normalized WHERE rn > 1)
      `,
      [pattern, pattern, pattern]
    );

    await trx.raw(
      `
      UPDATE pages
      SET url = regexp_replace(url, ?, '\\1')
      WHERE url ~ ?
      `,
      [pattern, pattern]
    );
  });
}

export async function down() {
  // Irreversible data cleanup migration.
}
