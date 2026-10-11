//! Pure helpers for the file-download flow.
//!
//! Kept free of Tauri types so the routing/name logic can be unit tested
//! without a webview runtime, mirroring `process/` and `config/`.

/// Returns true when a URL path targets the daemon's one-time file download
/// route (`/download/<passport>/<file>`), tolerating a mapped path prefix.
pub fn is_download_path(path: &str) -> bool {
    path.contains("/download/")
}

/// Returns true when a `window.open` request is a daemon file download: an
/// `http(s)` URL on loopback whose path carries the download route. Anything
/// else (external links, remote sites) must not be hijacked.
pub fn is_download_request(scheme: &str, host: Option<&str>, path: &str) -> bool {
    matches!(scheme, "http" | "https") && is_loopback_host(host) && is_download_path(path)
}

fn is_loopback_host(host: Option<&str>) -> bool {
    matches!(host, Some("127.0.0.1") | Some("localhost") | Some("::1"))
}

/// Best-effort file name to prefill the save dialog with, decoded from the URL
/// path. Falls back to `download` when the path carries no usable name.
pub fn suggested_filename(path: &str) -> String {
    let segments: Vec<&str> = path
        .split('/')
        .filter(|segment| !segment.is_empty())
        .collect();
    let count = segments.len();
    if count >= 3 && segments[count - 3] == "download" {
        let decoded = percent_decode(segments[count - 1]);
        let base = sanitize(&decoded);
        if !base.is_empty() && base != "." && base != ".." {
            return base;
        }
    }
    "download".to_string()
}

/// Keeps only the final path component of a decoded name so a crafted URL can
/// never seed the save dialog with directory traversal.
fn sanitize(name: &str) -> String {
    name.rsplit(|c| c == '/' || c == '\\')
        .next()
        .unwrap_or(name)
        .trim()
        .to_string()
}

fn percent_decode(input: &str) -> String {
    let bytes = input.as_bytes();
    let mut out = Vec::with_capacity(bytes.len());
    let mut i = 0;
    while i < bytes.len() {
        if bytes[i] == b'%' && i + 2 < bytes.len() {
            if let (Some(hi), Some(lo)) = (hex_value(bytes[i + 1]), hex_value(bytes[i + 2])) {
                out.push(hi * 16 + lo);
                i += 3;
                continue;
            }
        }
        out.push(bytes[i]);
        i += 1;
    }
    String::from_utf8_lossy(&out).into_owned()
}

fn hex_value(byte: u8) -> Option<u8> {
    match byte {
        b'0'..=b'9' => Some(byte - b'0'),
        b'a'..=b'f' => Some(byte - b'a' + 10),
        b'A'..=b'F' => Some(byte - b'A' + 10),
        _ => None,
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn download_path_matches_the_daemon_route() {
        assert!(is_download_path("/download/abc123/server.jar"));
        assert!(is_download_path("/download/"));
        assert!(is_download_path("/mcs/download/abc123/server.jar"));
    }

    #[test]
    fn download_path_rejects_lookalikes() {
        assert!(!is_download_path("/downloads/abc123/server.jar"));
        assert!(!is_download_path("/uploads/abc123"));
        assert!(!is_download_path("/"));
        assert!(!is_download_path(""));
    }

    #[test]
    fn download_request_requires_loopback_http() {
        assert!(is_download_request(
            "http",
            Some("127.0.0.1"),
            "/download/abc123/server.jar"
        ));
        assert!(is_download_request(
            "http",
            Some("localhost"),
            "/download/abc123/server.jar"
        ));
        assert!(!is_download_request(
            "http",
            Some("example.com"),
            "/download/abc123/server.jar"
        ));
        assert!(!is_download_request(
            "https",
            Some("github.com"),
            "/releases/download/v1/plugin.jar"
        ));
        assert!(!is_download_request(
            "http",
            None,
            "/download/abc123/server.jar"
        ));
        assert!(!is_download_request(
            "http",
            Some("127.0.0.1"),
            "/servers/1"
        ));
    }

    #[test]
    fn suggested_filename_takes_the_last_segment() {
        assert_eq!(
            suggested_filename("/download/abc123/server.jar"),
            "server.jar"
        );
        assert_eq!(
            suggested_filename("/mcs/download/abc123/server.jar"),
            "server.jar"
        );
    }

    #[test]
    fn suggested_filename_decodes_percent_escapes() {
        assert_eq!(
            suggested_filename("/download/abc123/my%20file%20name.zip"),
            "my file name.zip"
        );
    }

    #[test]
    fn suggested_filename_decodes_multibyte_utf8() {
        assert_eq!(
            suggested_filename("/download/abc123/%E4%B8%AD%E6%96%87.txt"),
            "\u{4e2d}\u{6587}.txt"
        );
    }

    #[test]
    fn suggested_filename_falls_back_when_missing() {
        assert_eq!(suggested_filename("/download/"), "download");
        assert_eq!(suggested_filename("/download/abc123/"), "download");
        assert_eq!(suggested_filename("/download/abc123"), "download");
        assert_eq!(suggested_filename(""), "download");
    }

    #[test]
    fn suggested_filename_strips_path_separators() {
        assert_eq!(
            suggested_filename("/download/abc123/..%2F..%2Fetc%2Fhosts"),
            "hosts"
        );
        assert_eq!(
            suggested_filename("/download/abc123/nested%5Cname.txt"),
            "name.txt"
        );
    }
}
