import { useCallback, useEffect, useState } from "react";
import { useBridge } from "../services/bridge";

export interface DownloadFailure {
  fileName: string;
  message?: string;
}

export interface UseDownloadErrorResult {
  failure: DownloadFailure | null;
  dismiss: () => void;
}

/// Surfaces failed `window.open` downloads so the UI can tell the user instead
/// of failing silently. Successful and cancelled downloads are ignored.
export function useDownloadError(): UseDownloadErrorResult {
  const bridge = useBridge();
  const [failure, setFailure] = useState<DownloadFailure | null>(null);

  useEffect(
    () =>
      bridge.onDownloadResult((result) => {
        if (result.status !== "failed") {
          return;
        }
        const next: DownloadFailure = { fileName: result.fileName };
        if (result.message != null) {
          next.message = result.message;
        }
        setFailure(next);
      }),
    [bridge],
  );

  const dismiss = useCallback(() => setFailure(null), []);

  return { failure, dismiss };
}
