[🔙 Back to Main Project README](../README.md)

# Auto Claims Frontend

This is the frontend application for the Auto Claims Demo, providing an interface for claim submission and AI-driven assessment results.

## Tech Stack

-   **Framework**: [Vue 3](https://vuejs.org/)
-   **Build Tool**: [Vite](https://vitejs.dev/)
-   **Styling**: [Tailwind CSS v4](https://tailwindcss.com/)
-   **HTTP Client**: [Axios](https://axios-http.com/)

## Prerequisites

-   Node.js (v22 or later)
-   npm

## Setup & Run

The project includes a `Makefile` for easy setup.

### Run Locally

```bash
make local-frontend
```

This command will:
1.  Install dependencies (`npm install`).
2.  Start the development server (`npm run dev`).
3.  The app will be available at `http://localhost:5173` (or another port if 5173 is busy).

**Note on Proxy**: The Vite server is configured to proxy API requests starting with `/api` to the backend service running at `http://localhost:8080`. Ensure the backend is running for full functionality.

### Build for Production

```bash
npm run build
```

## Project Structure

-   `src/`: Source code
    -   `components/`: Vue components
    -   `views/`: Page views
    -   `router/`: Vue Router configuration
    -   `assets/`: Static assets
-   `public/`: Public assets
-   `vite.config.js`: Vite configuration
