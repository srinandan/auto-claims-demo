# Auto Claims Frontend

This is the frontend application for the Auto Claims Demo, built with **Vue 3**, **Vite**, and **Tailwind CSS v4**.

## Tech Stack

-   **Framework**: [Vue 3](https://vuejs.org/)
-   **Build Tool**: [Vite](https://vitejs.dev/)
-   **Styling**: [Tailwind CSS v4](https://tailwindcss.com/)
-   **HTTP Client**: [Axios](https://axios-http.com/)

## Prerequisites

-   Node.js (v18 or later)
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
