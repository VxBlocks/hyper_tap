use anyhow::{Context, Result};
use serde::{de::DeserializeOwned, Serialize};
use sqlx::{Row, SqlitePool};
use std::path::PathBuf;
use tokio::fs::create_dir_all;

mod config_migrate;

#[derive(Debug)]
pub struct Config {
    db: SqlitePool,
}

impl Config {
    pub async fn new(data_dir: &PathBuf) -> Result<Self> {
        if !data_dir.exists() {
            create_dir_all(data_dir).await?;
        }
        let db_path = data_dir.join("config.db");

        let options = sqlx::sqlite::SqliteConnectOptions::new()
            .filename(db_path)
            .create_if_missing(true);

        let pool = sqlx::SqlitePool::connect_with(options)
            .await
            .map_err(|err| anyhow::anyhow!("Could not connect to database: {err}"))?;

        let config = Self { db: pool };

        config.migrate_tables().await?;

        Ok(config)
    }

    async fn set_data<T: Serialize>(&self, key: &str, value: &T) -> Result<()> {
        let data = serde_json::to_vec(value)?;
        sqlx::query("INSERT OR REPLACE INTO config (key, value) VALUES (?1, ?2)")
            .bind(key)
            .bind(data)
            .execute(&self.db)
            .await
            .context("Failed to insert or replace data")?;
        Ok(())
    }

    async fn get_data<T: DeserializeOwned>(&self, key: &str) -> Result<Option<T>> {
        let row = sqlx::query("SELECT value FROM config WHERE key = ?1")
            .bind(key)
            .fetch_optional(&self.db)
            .await
            .context("Failed to query data")?;

        match row {
            Some(row) => {
                let data: Vec<u8> = row.get(0);
                Ok(Some(serde_json::from_slice(&data)?))
            }
            None => Ok(None),
        }
    }

    pub async fn set_log_level(&self, level: &str) -> Result<()> {
        self.set_data::<String>("log_level", &level.to_string())
            .await
    }

    pub async fn get_log_level(&self) -> Result<Option<String>> {
        self.get_data::<String>("log_level").await
    }
}
