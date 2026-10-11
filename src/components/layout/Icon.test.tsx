import { describe, it, expect } from "vitest";
import { render } from "@testing-library/react";
import { Icon } from "./Icon";

describe("Icon", () => {
  it("renders the settings glyph as a gear with a center hole", () => {
    const { container } = render(<Icon name="settings" />);
    const svg = container.querySelector("svg");
    expect(svg).not.toBeNull();

    const circle = svg?.querySelector("circle");
    expect(circle).not.toBeNull();

    const closedPath = Array.from(svg?.querySelectorAll("path") ?? []).some((path) =>
      /z$/i.test((path.getAttribute("d") ?? "").trim()),
    );
    expect(closedPath).toBe(true);
  });
});
