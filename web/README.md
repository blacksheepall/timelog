# TimeLog Web Frontend

A Vue 3 + TypeScript + Tailwind CSS frontend for the TimeLog application.

## Features

- ✅ Create, read, update, and delete time logs
- ✅ Real-time duration calculation
- ✅ Modern responsive UI with Tailwind CSS
- ✅ TypeScript for type safety
- ✅ Form validation and error handling
- ✅ Notification system
- ✅ Browser-based reminders every 25 minutes prompting TimeLog entries

## Development Setup

### Prerequisites

- Node.js 18+ and npm
- Go backend running on port 8080

### Installation

```bash
# Install dependencies
npm install

# Start development server
npm run dev
```

The development server will run on http://localhost:5173 and proxy API requests to the Go backend.

### Building for Production

```bash
# Build for production
npm run build

# Preview production build
npm run preview
```

The built files will be in the `dist/` directory and can be served by the Go backend.

## Usage

1. Start the Go backend server
2. Run `npm run dev` for development or `npm run build` for production
3. Open http://localhost:5173 (dev) or http://localhost:8080 (production)

## API Integration

The frontend communicates with the Go backend API at `/api/timelogs` endpoints:

- `GET /api/timelogs` - List all time logs
- `POST /api/timelogs` - Create a new time log
- `GET /api/timelogs/:id` - Get a specific time log
- `PUT /api/timelogs/:id` - Update a time log
- `DELETE /api/timelogs/:id` - Delete a time log

## Project Structure

```
src/
├── api/          # API client functions
├── components/   # Vue components
├── types/        # TypeScript type definitions
├── utils/        # Utility functions
├── App.vue       # Main application component
├── main.ts       # Application entry point
└── style.css     # Global styles
```
