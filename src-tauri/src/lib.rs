mod commands;
pub mod config;
pub mod download;
pub mod process;
pub mod tray;

use std::sync::{Arc, Mutex, RwLock};
use std::time::Duration;

use tauri::Manager;
#[cfg(not(dev))]
use tauri::{ipc::CapabilityBuilder, Url};
use tauri::{WebviewUrl, WebviewWindowBuilder};

use commands::{read_manager, AppState};
use process::manager::ProcessManager;

/// Per-run state needed to (re)create the main window after it is destroyed to
/// free its WebView in tray mode.
struct WindowRuntime {
    url: WebviewUrl,
}

const TARGET_WIDTH: f64 = 1400.0;
const TARGET_HEIGHT: f64 = 960.0;

const MIN_WINDOW_WIDTH: f64 = 1024.0;
const MIN_WINDOW_HEIGHT: f64 = 680.0;
const WINDOW_LABEL: &str = "main";
const WINDOW_TITLE: &str = "MCSManager Desktop";

/// The panel web UI is embedded in an iframe. Browsers only attach cookies to
/// same-site requests, but the production custom-protocol origin on Windows
/// (`http://tauri.localhost`) is a different site than the panel origin
/// (`http://localhost:23333`). The login session cookie therefore becomes a
/// cross-site cookie and WebView2 refuses to store/send it, which is why the
/// authenticated `/api/auth/` call is rejected.
///
/// Serving the shell over a loopback HTTP server puts it on the same site
/// (`localhost`) as the panel, mirroring `tauri dev` where the Vite dev server
/// already runs on `localhost` and cookies work.
fn pick_unused_port() -> Option<u16> {
    std::net::TcpListener::bind("127.0.0.1:0")
        .ok()
        .and_then(|listener| listener.local_addr().ok())
        .map(|addr| addr.port())
}

/// Picks the loopback host the shell is served from so it stays on the same
/// site as the configured panel. `127.0.0.1` and `localhost` are different
/// sites for cookie purposes, so the shell host must mirror whichever one the
/// panel uses; anything else (remote hosts, IPv6) falls back to `localhost`.
#[cfg(not(dev))]
fn shell_host_for_panel(panel_url: &str) -> String {
    match panel_url
        .parse::<Url>()
        .ok()
        .and_then(|url| url.host_str().map(str::to_owned))
    {
        Some(host) if host == "127.0.0.1" => host,
        _ => "localhost".to_string(),
    }
}

fn fits_target_size(width: f64, height: f64) -> bool {
    width >= TARGET_WIDTH && height >= TARGET_HEIGHT
}

/// Routes a `window.open` request coming from the embedded panel.
///
/// The panel downloads files through `window.open`, and WebView2 reports that
/// as a new-window request. Daemon download links are saved through the native
/// save dialog (off the main thread); other web links open in the user's
/// default browser; everything else is dropped.
fn handle_new_window_request(app: &tauri::AppHandle, url: tauri::Url) {
    if download::is_download_request(url.scheme(), url.host_str(), url.path()) {
        let app = app.clone();
        std::thread::spawn(move || commands::download_from_url(app, url));
    } else if matches!(url.scheme(), "http" | "https") {
        commands::open_external_url(app, &url);
    }
}

/// Creates the main window and lets small monitors maximize it immediately.
///
/// Called once at startup and again by the tray whenever the window was closed
/// to reclaim memory. The URL is fixed per run so the runtime capability
/// registered at startup keeps matching the same label and origin.
fn create_main_window(app: &tauri::AppHandle) -> tauri::Result<tauri::WebviewWindow> {
    let url = app.state::<WindowRuntime>().url.clone();
    let new_window_app = app.clone();
    let window = WebviewWindowBuilder::new(app, WINDOW_LABEL, url)
        .title(WINDOW_TITLE)
        .inner_size(TARGET_WIDTH, TARGET_HEIGHT)
        .min_inner_size(MIN_WINDOW_WIDTH, MIN_WINDOW_HEIGHT)
        .on_new_window(move |url, _features| {
            handle_new_window_request(&new_window_app, url);
            tauri::webview::NewWindowResponse::Deny
        })
        .build()?;

    if let Ok(Some(monitor)) = window.primary_monitor() {
        let logical = monitor.size().to_logical::<f64>(monitor.scale_factor());
        if !fits_target_size(logical.width, logical.height) {
            let _ = window.maximize();
        }
    }
    Ok(window)
}

/// Whether an `ExitRequested` event comes from closing the last window (user
/// interaction) rather than a programmatic `exit`/`restart` (which carries a
/// code). Only the former must be prevented so the app keeps living in the tray.
fn is_user_window_close(code: Option<i32>) -> bool {
    code.is_none()
}

#[cfg_attr(mobile, tauri::mobile_entry_point)]
pub fn run() {
    let localhost_port = pick_unused_port().unwrap_or(0);

    tauri::Builder::default()
        .plugin(tauri_plugin_opener::init())
        .plugin(tauri_plugin_dialog::init())
        .plugin(
            tauri_plugin_localhost::Builder::new(localhost_port)
                .host("127.0.0.1")
                .build(),
        )
        .setup(move |app| {
            let config_path = app.path().app_config_dir()?.join("config.json");
            let outcome = config::load_from(&config_path)?;
            let startup_warnings: Vec<String> = outcome
                .error
                .iter()
                .map(|error| format!("config recovered from invalid file: {}", error))
                .collect();
            let app_config = outcome.config;
            let sink = commands::make_event_sink(app.handle().clone());
            let mut manager = ProcessManager::new(
                sink,
                Duration::from_millis(app_config.stop_timeout_ms),
            );
            for spec in commands::enabled_specs(&app_config) {
                let _ = manager.register_or_update(spec);
            }
            app.manage(AppState {
                manager: Arc::new(RwLock::new(manager)),
                config: Arc::new(Mutex::new(app_config)),
                config_path,
                startup_warnings,
            });

            let url = {
                #[cfg(dev)]
                {
                    // In development Vite serves the shell on http://localhost:1420,
                    // which shares the `localhost` site with the panel already.
                    WebviewUrl::App("index.html".into())
                }
                #[cfg(not(dev))]
                {
                    if localhost_port == 0 {
                        // Extremely unlikely fallback: no loopback port was
                        // available, so keep the default asset origin.
                        WebviewUrl::App("index.html".into())
                    } else {
                        let panel_url = {
                            let state = app.state::<AppState>();
                            let config = state
                                .config
                                .lock()
                                .unwrap_or_else(|poisoned| poisoned.into_inner());
                            config.panel_url.clone()
                        };
                        let host = shell_host_for_panel(&panel_url);
                        let origin = format!("http://{host}:{localhost_port}");
                        app.add_capability(
                            CapabilityBuilder::new("desktop-localhost")
                                .remote(origin.clone())
                                .window(WINDOW_LABEL)
                                .permission("core:default")
                                .permission("opener:default")
                                .permission("allow-desktop-commands"),
                        )?;
                        WebviewUrl::External(origin.parse::<Url>()?)
                    }
                }
            };

            app.manage(WindowRuntime { url });
            create_main_window(app.handle())?;

            tray::setup(app.handle())?;
            Ok(())
        })
        .invoke_handler(tauri::generate_handler![
            commands::get_config,
            commands::save_config,
            commands::get_service_statuses,
            commands::start_service,
            commands::stop_service,
            commands::restart_service,
            commands::start_all_services,
            commands::stop_all_services,
            commands::probe_tcp,
            commands::check_start_conflict,
            commands::force_free_port,
            commands::get_app_info,
            commands::open_service_folder,
        ])
        .build(tauri::generate_context!())
        .expect("error while running tauri application")
        .run(move |app, event| match event {
            // Closing the main window destroys it so its WebView memory is
            // reclaimed, but the app keeps running in the tray: services stay
            // up and the window is recreated on demand. Only a programmatic
            // exit (tray "Exit" -> `app.exit`) may end the process.
            tauri::RunEvent::ExitRequested { api, code, .. } if is_user_window_close(code) => {
                api.prevent_exit();
            }
            tauri::RunEvent::Exit => {
                if let Some(state) = app.try_state::<AppState>() {
                    read_manager(&state.manager).shutdown();
                }
            }
            _ => {}
        });
}

#[cfg(test)]
mod tests {
    use super::{fits_target_size, is_user_window_close};

    #[test]
    fn fits_exact_target_size() {
        assert!(fits_target_size(1400.0, 960.0));
    }

    #[test]
    fn fits_larger_monitor() {
        assert!(fits_target_size(2560.0, 1440.0));
    }

    #[test]
    fn too_short_monitor_does_not_fit() {
        assert!(!fits_target_size(1920.0, 900.0));
    }

    #[test]
    fn too_narrow_monitor_does_not_fit() {
        assert!(!fits_target_size(1280.0, 1080.0));
    }

    #[test]
    fn closing_the_last_window_must_not_exit() {
        assert!(is_user_window_close(None));
    }

    #[test]
    fn programmatic_exit_is_allowed() {
        assert!(!is_user_window_close(Some(0)));
        assert!(!is_user_window_close(Some(1)));
    }
}
