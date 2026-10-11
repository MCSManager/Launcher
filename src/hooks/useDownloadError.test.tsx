import { describe, it, expect, afterEach } from "vitest";
import { StrictMode, type ReactNode } from "react";
import { act, cleanup, renderHook } from "@testing-library/react";
import { BridgeProvider, type Bridge } from "../services/bridge";
import { createMockBridge } from "../test/mockBridge";
import { useDownloadError } from "./useDownloadError";

function wrapperFor(bridge: Bridge) {
  return function Wrapper({ children }: { children: ReactNode }) {
    return (
      <StrictMode>
        <BridgeProvider bridge={bridge}>{children}</BridgeProvider>
      </StrictMode>
    );
  };
}

describe("useDownloadError", () => {
  afterEach(() => {
    cleanup();
  });

  it("records a failed download result", () => {
    const mock = createMockBridge();
    const { result } = renderHook(() => useDownloadError(), { wrapper: wrapperFor(mock) });

    act(() => {
      mock.emitDownloadResult({ status: "failed", fileName: "server.jar", message: "boom" });
    });

    expect(result.current.failure).toEqual({ fileName: "server.jar", message: "boom" });
  });

  it("ignores saved downloads and exposes a dismiss action", () => {
    const mock = createMockBridge();
    const { result } = renderHook(() => useDownloadError(), { wrapper: wrapperFor(mock) });

    act(() => {
      mock.emitDownloadResult({ status: "saved", fileName: "server.jar", path: "C:/x" });
    });
    expect(result.current.failure).toBeNull();

    act(() => {
      mock.emitDownloadResult({ status: "failed", fileName: "server.jar" });
    });
    expect(result.current.failure).toEqual({ fileName: "server.jar" });

    act(() => {
      result.current.dismiss();
    });
    expect(result.current.failure).toBeNull();
  });
});
