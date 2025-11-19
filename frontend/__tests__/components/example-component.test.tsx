import { describe, it, expect } from "vitest";
import { render, screen } from "@testing-library/react";

/**
 * Example React component test to verify @testing-library/react setup.
 * This demonstrates how to test React components with Vitest.
 */

// Simple test component
function TestButton({
  onClick,
  children,
}: {
  onClick: () => void;
  children: React.ReactNode;
}) {
  return <button onClick={onClick}>{children}</button>;
}

describe("React Testing Library Setup", () => {
  it("should render component correctly", () => {
    render(<TestButton onClick={() => {}}>Click me</TestButton>);

    const button = screen.getByRole("button", { name: /click me/i });
    expect(button).toBeInTheDocument();
  });

  it("should handle text content", () => {
    const { container } = render(<div>Hello, Testing!</div>);
    expect(container.textContent).toBe("Hello, Testing!");
  });
});
