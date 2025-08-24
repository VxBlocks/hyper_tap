use serde::{Deserialize, Serialize};

#[derive(Debug, Serialize)]
pub struct BuildInfo {
    pub time: String,
    pub commit: String,
}

#[tauri::command]
pub fn get_build_info() -> BuildInfo {
    let commit = env!("git_commit").to_string();
    let time = env!("build_time").to_string();

    BuildInfo { time, commit }
}

#[derive(Debug, Serialize, Deserialize)]
pub struct UpdateInfo {
    pub version: String,
    pub url: String,
}

#[tauri::command]
pub async fn update_info() -> Result<UpdateInfo, String> {
    let resp = reqwest::get(
        "https://raw.githubusercontent.com/VxBlocks/vxb_neptune_wallet/refs/heads/main/update.json",
    )
    .await
    .map_err(|e| e.to_string())?
    .json::<UpdateInfo>()
    .await
    .map_err(|e| e.to_string())?;

    Ok(resp)
}
