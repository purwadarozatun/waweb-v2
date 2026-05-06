# WaWeb V2 - Architecture & Flow

## System Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                        User Browser                         │
│                     (HTMX + TailwindCSS)                    │
└──────────────────────┬──────────────────────────────────────┘
                       │ HTTP Requests
                       ▼
┌─────────────────────────────────────────────────────────────┐
│                      Go Fiber Server                        │
│  ┌──────────────────────────────────────────────────────┐  │
│  │                    Middleware                         │  │
│  │  - Logger                                            │  │
│  │  - Auth (JWT)                                        │  │
│  │  - Recovery                                          │  │
│  └──────────────────────────────────────────────────────┘  │
│  ┌──────────────────────────────────────────────────────┐  │
│  │                     Handlers                          │  │
│  │  - Auth Handler                                      │  │
│  │  - Dashboard Handler                                 │  │
│  │  - Workspace Handler                                 │  │
│  │  - Builder Handler                                   │  │
│  └──────────────────────────────────────────────────────┘  │
│  ┌──────────────────────────────────────────────────────┐  │
│  │                   In-Memory DB                        │  │
│  │  - Users                                             │  │
│  │  - Workspaces                                        │  │
│  └──────────────────────────────────────────────────────┘  │
└──────────────────────┬──────────────────────────────────────┘
                       │
                       ▼
┌─────────────────────────────────────────────────────────────┐
│                    File System                              │
│  ┌────────────────┐  ┌──────────────┐  ┌─────────────────┐│
│  │   Templates    │  │   Static     │  │ Published Sites ││
│  │  - Landing     │  │  - CSS       │  │  - site1/       ││
│  │  - Portfolio   │  │  - Images    │  │  - site2/       ││
│  └────────────────┘  └──────────────┘  └─────────────────┘│
└─────────────────────────────────────────────────────────────┘
```

## User Flow

```
┌─────────┐
│  START  │
└────┬────┘
     │
     ▼
┌─────────────────┐
│   Login Page    │ ◄──── Demo: demo@example.com / password123
└────┬────────────┘
     │ Login Success
     ▼
┌─────────────────┐
│   Dashboard     │ ◄──── View all workspaces
└────┬────────────┘
     │ Click "+ New Workspace"
     ▼
┌─────────────────┐
│  New Workspace  │ ◄──── Enter workspace name
└────┬────────────┘
     │ Submit
     ▼
┌─────────────────┐
│ Template Picker │ ◄──── Choose: Landing Page or Portfolio
└────┬────────────┘
     │ Select Template
     ▼
┌─────────────────┐
│     Builder     │ ◄──── Edit content fields
│                 │       - Text fields
│  ┌───────────┐  │       - Colors
│  │   Form    │  │       - Links
│  │  Fields   │  │       - Descriptions
│  └───────────┘  │
│                 │
│  ┌───────────┐  │
│  │  Preview  │  │
│  │  Section  │  │
│  └───────────┘  │
└────┬────────────┘
     │ Save Changes
     │
     │ Click "Publish"
     ▼
┌─────────────────┐
│  Publish Form   │ ◄──── Enter subdomain
└────┬────────────┘
     │ Publish
     ▼
┌─────────────────┐
│  File System    │
│  Copy Template  │ ──┐
│  Create values  │   │
└─────────────────┘   │
     │                │
     ▼                │
┌─────────────────┐   │
│ Published Site  │ ◄─┘
│ /sites/subdomain│
└─────────────────┘
     │
     ▼
┌─────────────────┐
│   Visit Site    │ ◄──── Public access to published site
└─────────────────┘
```

## Data Flow: Publishing Process

```
1. User fills form in Builder
   ↓
2. Data stored in workspace.values (in-memory)
   {
     "site_title": "My Site",
     "primary_color": "#3b82f6",
     ...
   }
   ↓
3. User clicks "Publish" with subdomain
   ↓
4. Server copies template directory
   templates/simple-landing/ → published_sites/my-site/
   ↓
5. Server creates values.json
   {
     "site_title": "My Site",
     "primary_color": "#3b82f6",
     ...
   }
   ↓
6. Published site structure:
   published_sites/my-site/
   ├── index.html
   ├── values.json  ← Generated from workspace.values
   └── assets/
   ↓
7. Site accessible at /sites/my-site/
   ↓
8. JavaScript loads values.json and updates DOM
```

## Template Structure

```
templates/
└── simple-landing/
    ├── spec.json          ← Defines editable fields
    ├── index.html         ← HTML template with IDs
    └── assets/            ← Images, CSS, JS (optional)
        └── style.css

spec.json:
{
  "name": "Simple Landing",
  "fields": {
    "site_title": {
      "label": "Site Title",
      "type": "text",
      "default": "Welcome"
    }
  }
}

index.html:
<h1 id="site-title">Welcome</h1>
<script>
  fetch('values.json')
    .then(res => res.json())
    .then(values => {
      document.getElementById('site-title').textContent = values.site_title;
    });
</script>
```

## Request Flow Examples

### 1. Login Request

```
Browser                Server              Database
   │                      │                    │
   │  POST /login         │                    │
   ├─────────────────────►│                    │
   │  email, password     │                    │
   │                      │  GetUserByEmail()  │
   │                      ├───────────────────►│
   │                      │                    │
   │                      │◄───────────────────┤
   │                      │  User data         │
   │                      │                    │
   │                      │  Verify password   │
   │                      │  Generate JWT      │
   │                      │                    │
   │   Set-Cookie: token  │                    │
   │◄─────────────────────┤                    │
   │   HX-Redirect        │                    │
```

### 2. Create Workspace Request

```
Browser                Server              Database
   │                      │                    │
   │ POST /workspace/     │                    │
   │      create          │                    │
   ├─────────────────────►│                    │
   │  name="My Site"      │                    │
   │                      │                    │
   │                      │  CreateWorkspace() │
   │                      ├───────────────────►│
   │                      │  {                 │
   │                      │   id: uuid(),      │
   │                      │   name: "My Site", │
   │                      │   user_id: "..."   │
   │                      │  }                 │
   │                      │                    │
   │  HX-Redirect to      │◄───────────────────┤
   │  /templates?ws=id    │                    │
   │◄─────────────────────┤                    │
```

### 3. Publish Site Request

```
Browser                Server              File System
   │                      │                    │
   │ POST /builder/       │                    │
   │      {id}/publish    │                    │
   ├─────────────────────►│                    │
   │ subdomain="my-site"  │                    │
   │                      │                    │
   │                      │  Copy template     │
   │                      ├───────────────────►│
   │                      │  templates/...     │
   │                      │  → published_sites/│
   │                      │     my-site/       │
   │                      │                    │
   │                      │  Write values.json │
   │                      ├───────────────────►│
   │                      │                    │
   │  Success message     │◄───────────────────┤
   │◄─────────────────────┤                    │
   │  + link to site      │                    │
```

## Authentication Flow

```
┌──────────┐
│  Client  │
└────┬─────┘
     │
     │ POST /login (email, password)
     ▼
┌─────────────────┐
│  Auth Handler   │
└────┬────────────┘
     │
     │ 1. Get user from DB
     ▼
┌─────────────────┐
│  Verify Pass    │
└────┬────────────┘
     │ Success
     │
     │ 2. Generate JWT
     ▼
┌─────────────────┐
│   JWT Token     │
│  {              │
│   user_id,      │
│   email,        │
│   exp           │
│  }              │
└────┬────────────┘
     │
     │ 3. Set HTTP-Only Cookie
     ▼
┌─────────────────┐
│  Set-Cookie:    │
│  token=...      │
└────┬────────────┘
     │
     │ 4. Redirect to /dashboard
     ▼
┌─────────────────┐
│   Dashboard     │
└─────────────────┘
```

## Protected Route Flow

```
Request → Middleware → Verify JWT → Set user context → Handler
                           │
                           │ Invalid/Missing
                           ▼
                       Redirect to /login
```

## Database Schema (In-Memory)

```go
// Users map
map[string]*User {
  "demo@example.com": {
    ID: "user-1",
    Email: "demo@example.com",
    Password: "$2a$10$...",  // bcrypt hash
    Name: "Demo User"
  }
}

// Workspaces map
map[string]*Workspace {
  "workspace-uuid": {
    ID: "workspace-uuid",
    UserID: "user-1",
    Name: "My Website",
    TemplateID: "simple-landing",
    Subdomain: "my-site",
    Values: {
      "site_title": "My Awesome Site",
      "primary_color": "#3b82f6"
    },
    IsPublished: true
  }
}
```

## File System Layout

```
waweb-v2/
├── config/
│   └── config.go           # Configuration loader
├── database/
│   └── database.go         # In-memory database
├── handlers/
│   ├── auth.go             # Login, register, logout
│   ├── dashboard.go        # Dashboard view
│   ├── workspace.go        # Workspace CRUD
│   └── builder.go          # Builder & publish
├── middleware/
│   └── auth.go             # JWT authentication
├── models/
│   ├── user.go             # User model
│   └── workspace.go        # Workspace & Template models
├── routes/
│   └── routes.go           # Route definitions
├── views/
│   ├── layout.html         # Base layout
│   ├── login.html          # Login page
│   ├── dashboard.html      # Dashboard page
│   ├── templates.html      # Template picker
│   └── builder.html        # Builder interface
├── templates/              # Website templates
│   ├── simple-landing/
│   │   ├── spec.json
│   │   └── index.html
│   └── portfolio/
│       ├── spec.json
│       └── index.html
├── static/                 # Static assets
├── published_sites/        # Published websites
├── main.go                 # Application entry point
├── go.mod                  # Go dependencies
├── .env                    # Environment config
├── Makefile                # Build commands
└── README.md              # Documentation
```

## Technology Stack

- **Backend Framework**: Go Fiber v2
- **Template Engine**: Go HTML Templates
- **Frontend**: HTMX + TailwindCSS
- **Authentication**: JWT (golang-jwt/jwt)
- **Password**: bcrypt
- **Database**: In-memory (upgradeable to PostgreSQL/MySQL)
- **UUID**: Google UUID

## Security Features

1. **Password Hashing**: bcrypt with default cost
2. **JWT Authentication**: Signed tokens with secret key
3. **HTTP-Only Cookies**: Prevent XSS attacks
4. **Path Validation**: Prevent directory traversal
5. **Subdomain Validation**: Only alphanumeric and hyphens
6. **Protected Routes**: Middleware-based auth

## Future Enhancements

1. **Database**: Replace in-memory with PostgreSQL/MySQL
2. **File Upload**: Support image uploads for templates
3. **Custom Domains**: Map custom domains to subdomains
4. **Analytics**: Track visitor stats for published sites
5. **Templates**: Marketplace for community templates
6. **Collaboration**: Multi-user workspaces
7. **Version Control**: Git-like versioning for sites
8. **Preview**: Live preview before publish
9. **SEO Tools**: Meta tags, sitemap generation
10. **Backup**: Automated backup of published sites
