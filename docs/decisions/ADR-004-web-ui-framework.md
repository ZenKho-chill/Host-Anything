# ADR 004: Web UI Framework

## Status
Accepted

## Context
The project requires a web-based user interface to serve as the control plane. It needs to be fast, maintainable, and easily embeddable into the Go binary.

Options considered:
- **Server-Side Rendered (Go HTML Templates/HTMX)**: Extremely simple, zero JS build step. However, building complex, highly interactive interfaces (like live log streaming terminals or dynamic configuration forms) becomes cumbersome.
- **Vue.js**: Excellent framework, but the ecosystem and third-party library support is slightly smaller than React.
- **React + Next.js**: Next.js is powerful but designed for SSR and serverless deployment. It doesn't align well with our goal of embedding a static SPA into a Go binary.
- **React + Vite + TypeScript**: Creates a standard Single Page Application outputting static HTML/JS/CSS.

## Decision
We will use **React** built with **Vite** and written in **TypeScript**.

## Rationale
React provides the component ecosystem necessary for building a complex dashboard (charts, terminal emulators, complex forms). Vite offers incredibly fast build times and hot module replacement (HMR), significantly improving developer experience compared to older tools like Create React App. TypeScript is mandatory to maintain data integrity between the Go API and the frontend, allowing us to generate types based on our Go structs. The final output is purely static, easily served by Go's `embed` package.

## Consequences
- Requires a Node.js environment *only* at build time.
- The Go binary size will increase slightly to accommodate the embedded JS/CSS assets (usually < 2MB compressed).
- Developers must maintain API synchronization manually or via generation tools (like OpenAPI generators) to ensure TS types match Go structs.
