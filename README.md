# Project Management Application - Frontend

A modern, full-featured project management application built with React, TypeScript, and Vite. This application provides comprehensive task tracking, team collaboration, and project organization capabilities similar to Jira, Linear, or Asana.

## 🚀 Features

### Authentication
- User login and registration
- JWT-based authentication with automatic token refresh
- Secure session management with Zustand

### Organization Management
- Create and manage multiple organizations
- Hierarchical organization structure support
- Role-based access control (Owner, Admin, Member)

### Project Management
- Create and organize projects within organizations
- Multiple project views (Board, List, Calendar, Timeline)
- Project visibility controls (Private, Team, Organization)
- Project status tracking (Active, On Hold, Archived, Completed)
- Custom project icons and colors

### Task Management
- **Board View**: Kanban-style drag-and-drop task board
- **List View**: Tabular task list with advanced filtering
- **Task Details**: Comprehensive task detail page with:
  - Rich text descriptions
  - Priority levels (Urgent, High, Medium, Low, None)
  - Status management with custom workflows
  - Assignee management
  - Due dates and time tracking
  - Comments and threaded discussions
  - File attachments
  - Labels and tags
  - Task dependencies
  - Subtasks support

### Team Collaboration
- Team creation and management
- Member invitations and role assignment
- Activity logs and audit trails
- @mentions in comments

### Time Tracking
- Log time entries on tasks
- Estimated vs. actual hours tracking
- Billable/non-billable time categorization
- Time reports and analytics

## 🛠️ Tech Stack

### Core Framework
- **React 18** - UI library
- **TypeScript** - Type safety
- **Vite** - Build tool and dev server

### Routing & State
- **TanStack Router** - Type-safe file-based routing
- **TanStack Query** - Server state management and caching
- **Zustand** - Client state management (auth)

### UI Components
- **shadcn/ui** - Headless component library
- **Radix UI** - Accessible component primitives
- **Tailwind CSS** - Utility-first styling
- **Lucide React** - Icon library

### Data & Forms
- **TanStack Table** - Powerful table component
- **@dnd-kit** - Drag and drop functionality
- **React Hook Form** - Form management
- **Zod** - Schema validation

### Utilities
- **Axios** - HTTP client
- **date-fns** - Date manipulation
- **Sonner** - Toast notifications

## 📁 Project Structure

```
src/
├── api/                    # API client and service functions
│   ├── auth.ts            # Authentication endpoints
│   ├── organizations.ts   # Organization CRUD operations
│   ├── projects.ts        # Project management
│   └── tasks.ts           # Task operations
├── components/
│   ├── ui/                # shadcn/ui components
│   │   ├── button.tsx
│   │   ├── card.tsx
│   │   ├── dialog.tsx
│   │   ├── input.tsx
│   │   └── ...
│   ├── layout/            # Layout components
│   │   ├── organization-layout.tsx
│   │   └── project-layout.tsx
│   ├── auth/              # Authentication components
│   ├── tasks/             # Task-specific components
│   ├── projects/          # Project-specific components
│   └── organizations/     # Organization components
├── pages/                 # Page components
│   ├── auth/
│   │   ├── login.tsx
│   │   └── register.tsx
│   ├── organizations/
│   │   └── organizations-list.tsx
│   ├── projects/
│   │   └── projects-list.tsx
│   └── tasks/
│       ├── board-view.tsx
│       ├── list-view.tsx
│       └── task-detail.tsx
├── routes/                # TanStack Router route files
│   ├── __root.tsx
│   ├── index.tsx
│   ├── login.tsx
│   ├── register.tsx
│   └── organizations/
│       ├── $organizationId.tsx
│       └── $organizationId/
│           ├── projects.tsx
│           └── projects/
│               └── $projectId/
│                   ├── board.tsx
│                   ├── list.tsx
│                   └── tasks/
│                       └── $taskId.tsx
├── lib/                   # Utility functions
│   ├── api-client.ts     # Axios configuration
│   └── utils.ts          # Helper functions
├── stores/                # State management
│   └── auth-store.ts     # Authentication store
├── types/                 # TypeScript type definitions
│   └── index.ts          # All type definitions
├── App.tsx               # Main app component
└── main.tsx              # App entry point
```

## 🚦 Getting Started

### Prerequisites
- Node.js 18+ and npm
- Backend API server running (see backend repository)

### Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd project-management-app
```

2. Install dependencies:
```bash
npm install
```

3. Create environment file:
```bash
cp .env.example .env
```

4. Configure environment variables:
```env
VITE_API_URL=http://localhost:3000/api
```

5. Start development server:
```bash
npm run dev
```

The application will be available at `http://localhost:5173`

## 📜 Available Scripts

- `npm run dev` - Start development server
- `npm run build` - Build for production
- `npm run preview` - Preview production build
- `npm run lint` - Run ESLint

## 🔑 Key Concepts

### Routing
The application uses file-based routing with TanStack Router. Route files are located in `src/routes/` and automatically generate the route tree.

### State Management
- **Server State**: TanStack Query handles all server data fetching, caching, and synchronization
- **Client State**: Zustand manages authentication state and user session
- **Form State**: React Hook Form manages form state and validation

### API Integration
API services are organized by domain (auth, projects, tasks, etc.) and use Axios for HTTP requests. The API client includes:
- Automatic JWT token injection
- Token refresh on 401 errors
- Error handling and retry logic

### Type Safety
All data types are defined based on the PostgreSQL database schema and located in `src/types/index.ts`. This ensures type safety across the entire application.

## 🎨 UI Components

The application uses shadcn/ui, a collection of re-usable components built with Radix UI and Tailwind CSS. Components are:
- Fully accessible (WCAG compliant)
- Customizable via Tailwind classes
- Type-safe with TypeScript
- Documented and well-tested

## 🔐 Authentication Flow

1. User enters credentials on login page
2. Frontend sends POST request to `/api/auth/login`
3. Backend validates credentials and returns JWT tokens
4. Tokens are stored in localStorage and Zustand store
5. All subsequent API requests include JWT in Authorization header
6. On token expiration, refresh token is used automatically
7. On refresh failure, user is redirected to login

## 🗂️ Database Integration

The frontend types are generated from the PostgreSQL database schema which includes:
- Users and authentication
- Organizations with hierarchical structure
- Teams and team members
- Projects with custom statuses
- Tasks with full metadata
- Comments and activity logs
- Time entries and attachments
- Labels and custom fields

## 🚧 Future Enhancements

- Real-time collaboration with WebSockets
- Advanced search and filtering
- Keyboard shortcuts
- Dark mode support
- Mobile responsive design improvements
- Notification system
- Calendar and timeline views
- Gantt charts
- Reports and analytics dashboard
- File preview for attachments
- Markdown support in comments
- Task templates
- Automation rules

## 📝 License

This project is licensed under the MIT License.

## 🤝 Contributing

Contributions are welcome! Please read the contributing guidelines before submitting PRs.

## 📧 Support

For issues and questions, please open an issue on GitHub.
