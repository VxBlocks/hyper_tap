use std::sync::Arc;

use tauri::Manager;

use crate::store::Memstore;

mod config;
mod logger;
mod sdk;
mod service;
mod store;

#[cfg_attr(mobile, tauri::mobile_entry_point)]
pub fn run() {
    tauri::Builder::default()
        // .register_uri_scheme_protocol("metamask", |_ctx, request| {
        //     tracing::error!("metamask://{}", request.uri());
        //     let response = tauri::http::Response::builder()
        //         .status(200)
        //         .body(Vec::new())
        //         .unwrap();
        //     response
        // })
        .plugin(tauri_plugin_opener::init())
        .invoke_handler(tauri::generate_handler![])
        .setup(|app| {
            let data_dir = app
                .path()
                .app_config_dir()
                .expect("failed to get app data dir");

            let config = tauri::async_runtime::block_on(async {
                let config = crate::config::Config::new(&data_dir).await.unwrap();
                let level = config.get_log_level().await.unwrap();
                let _ = crate::logger::setup_logger(level);
                config
            });

            let persist_store =
                tauri::async_runtime::block_on(crate::store::PersisStore::new(&data_dir)).unwrap();
            crate::service::manage(persist_store);

            crate::service::manage(Arc::new(config));
            let memstore = Memstore::new();
            crate::service::manage(memstore);

            tracing::info!("app started");

            Ok(())
        })
        .run(tauri::generate_context!())
        .expect("error while running tauri application");
}

pub fn add_commands<R: tauri::Runtime>(app: tauri::Builder<R>) -> tauri::Builder<R> {
    app.invoke_handler(tauri::generate_handler![
        logger::clear_logs,
        logger::get_log_level,
        logger::get_logs,
        logger::log,
        logger::set_log_level,
        store::persist_store_execute,
        store::session_store_del,
        store::session_store_get,
        store::session_store_set,
        service::app::get_build_info,
        service::app::update_info,
    ])
}
