# CLAUDE.md - TimeLog Web Frontend

This file provides guidance to Claude Code when working with the Vue 3 frontend of the TimeLog application.

## Frontend Overview

A comprehensive Vue 3 + TypeScript + Tailwind CSS single-page application (SPA) for the TimeLog application. Features a modern, responsive interface for time tracking, task management, and productivity analytics with structured tag system and multi-page navigation.

## Technology Stack

- **Framework**: Vue 3.4+ with Composition API
- **Routing**: Vue Router 4.2+ for SPA navigation
- **Language**: TypeScript 5.3+
- **Build Tool**: Vite 5.0+
- **Styling**: Tailwind CSS 3.4+
- **Icons**: Heroicons Vue 2.0+
- **HTTP Client**: Axios 1.6+
- **Date Utilities**: date-fns 3.0+

## Features

### Core Functionality

- ✅ **Time Logging**: Create, read, update, and delete time logs with timezone handling
- ✅ **Task Management**: Complete task lifecycle management with deadlines and estimates
- ✅ **Tag System**: Structured tag selection with color-coded labels and descriptions
- ✅ **Analytics**: Statistics and reporting with duration calculations
- ✅ **Multi-page SPA**: Five main pages with responsive navigation

### Technical Features

- ✅ **Real-time Calculations**: Duration and time tracking with proper timezone conversion
- ✅ **"Set to Now"**: Quick time entry with current timestamp buttons
- ✅ **Modern UI**: Responsive design with Tailwind CSS and smooth transitions
- ✅ **TypeScript**: Complete type safety across all components and APIs
- ✅ **Form Validation**: Client-side validation with error handling
- ✅ **Toast Notifications**: Global notification system for user feedback
- ✅ **API Integration**: Unified API client with error handling and loading states

## Development Setup

### Prerequisites

- Node.js 18+ and pnpm (preferred package manager)
- Go backend running on port 8083

### Installation

```bash
# Install dependencies
pnpm install

# Start development server
pnpm run dev

# Type check without emitting
pnpm run type-check
```

The development server runs on http://localhost:5173 and proxies API requests to the Go backend on port 8083.

### Building for Production

```bash
# Build for production
pnpm run build

# Preview production build
pnpm run preview
```

Built files go to `dist/` directory and are served by the Go backend at production time.

## Usage Scenarios

**Development Mode:**

1. Start Go backend: `./main` (runs on port 8083)
2. Start frontend dev server: `pnpm run dev` (runs on port 5173)
3. Access at http://localhost:5173

**Production Mode:**

1. Build frontend: `pnpm run build`
2. Start backend: `./main` (serves both API and frontend on port 8083)
3. Access at http://localhost:8083

## API Integration

The frontend communicates with the Go backend through a unified API structure:

### TimeLog Endpoints

- `GET /api/timelogs` - List all time logs with tag and task information
- `POST /api/timelogs` - Create new time log (with optional task association)
- `GET /api/timelogs/:id` - Get specific time log
- `PUT /api/timelogs/:id` - Update time log
- `DELETE /api/timelogs/:id` - Delete time log

### Task Endpoints

- `GET /api/tasks` - List all tasks with optional date filtering
- `POST /api/tasks` - Create new task
- `GET /api/tasks/:id` - Get specific task
- `PUT /api/tasks/:id` - Update task
- `DELETE /api/tasks/:id` - Delete task
- `POST /api/tasks/:id/complete` - Mark task as completed
- `POST /api/tasks/:id/incomplete` - Mark task as incomplete
- `GET /api/tasks/stats/:date` - Get task completion statistics

### Tag Endpoints

- `GET /api/tags` - List all available tags
- `POST /api/tags` - Create new tag
- `PUT /api/tags/:id` - Update tag
- `DELETE /api/tags/:id` - Delete tag

### API Response Format

All endpoints return data in a consistent format:

```typescript
{
  data: T,           // The actual response data
  message: string,   // Success/error message
  status: number     // HTTP status code
}
```

## Project Structure

```
src/
├── api/              # API client functions
│   └── index.ts      # HTTP client setup and API endpoints
├── components/       # Vue components
│   ├── TimeLogList.vue    # Display list of time logs
│   └── TimeLogForm.vue    # Create/edit time log form
├── types/            # TypeScript type definitions
│   └── index.ts      # Interface definitions for TimeLog, Tag, etc.
├── utils/            # Utility functions
│   └── date.ts       # Date formatting and duration calculation
├── App.vue           # Main application component
├── main.ts           # Application entry point
└── style.css         # Global Tailwind CSS imports
```

## Key Components

### TimeLogList.vue

- Displays paginated list of time logs
- Shows color-coded tags with hover descriptions
- Calculates and displays duration for each entry
- Handles edit/delete actions

### TimeLogForm.vue

- Create/edit form for time logs
- Tag selection dropdown with colors
- DateTime pickers for start/end times
- Form validation and submission

### App.vue

- Main application layout
- Manages global state (time logs, tags, loading states)
- Handles API calls and error management
- Toast notification system

## TypeScript Types

### Core Interfaces

```typescript
interface Tag {
  id: number
  name: string
  color: string // Hex color code
  description: string
  created_at: string
  updated_at: string
}

interface TimeLog {
  id: number
  start_time: string
  end_time?: string | null
  tag_id: number
  tag: Tag // Populated tag object
  remarks: string
  created_at: string
  updated_at: string
}
```

## Styling Guidelines

- Uses Tailwind CSS utility-first approach
- Color scheme matches tag colors from backend
- Responsive design with mobile-first approach
- Consistent spacing and typography
- Hover states and smooth transitions

## Development Notes

- All API calls use the unified response format
- Form validation happens client-side with TypeScript
- Error handling provides user-friendly messages
- Duration calculation updates in real-time
- Tag selection is restricted to predefined options
- All timestamps are handled in ISO format with proper timezone conversion
