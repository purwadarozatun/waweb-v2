# WaWeb V2 - Website Builder SaaS

Platform SaaS untuk membuat website dengan mudah menggunakan Go Fiber dan HTMX.

## Fitur

- 🔐 **Authentication System** - Login dan register user
- 📊 **Dashboard** - Kelola semua workspace Anda
- 🎨 **Template Picker** - Pilih dari berbagai template yang tersedia
- ✏️ **Builder** - Edit website dengan interface yang mudah
- 🚀 **Publish** - Publish website ke subdomain sendiri
- 💾 **Dynamic Content** - Konten website di-load dari values.json

## Teknologi

- **Backend**: Go Fiber
- **Frontend**: HTMX + TailwindCSS
- **Storage**: In-memory database (bisa diganti dengan PostgreSQL/MySQL)
- **Template Engine**: Go HTML Templates

## Struktur Project

```
waweb-v2/
├── main.go                 # Entry point aplikasi
├── config/                 # Konfigurasi aplikasi
├── database/               # Database layer
├── handlers/               # HTTP handlers
│   ├── auth.go            # Login, register, logout
│   ├── dashboard.go       # Dashboard handler
│   ├── workspace.go       # Workspace management
│   └── builder.go         # Builder dan publish
├── middleware/             # Middleware (auth, etc)
├── models/                 # Data models
├── routes/                 # Route definitions
├── views/                  # HTML templates
│   ├── layout.html
│   ├── login.html
│   ├── register.html
│   ├── dashboard.html
│   ├── templates.html
│   └── builder.html
├── templates/              # Website templates
│   ├── simple-landing/
│   │   ├── spec.json      # Template specification
│   │   └── index.html     # Template HTML
│   └── portfolio/
│       ├── spec.json
│       └── index.html
├── static/                 # Static files (CSS, images)
└── published_sites/        # Published websites
```

## Cara Install

1. **Clone dan masuk ke directory**
```bash
cd /home/theakistea/project/waweb-v2
```

2. **Install dependencies**
```bash
go mod download
```

3. **Copy file .env**
```bash
cp .env.example .env
```

4. **Jalankan aplikasi**
```bash
go run main.go
```

5. **Akses di browser**
```
http://localhost:3000
```

## Demo Account

Gunakan akun ini untuk login:
- **Email**: demo@example.com
- **Password**: password123

## Cara Kerja

### 1. Flow User

```
Login → Dashboard → New Workspace → Template Picker → Builder → Publish
```

### 2. Template System

Setiap template memiliki:
- **spec.json**: Mendefinisikan field yang bisa diedit user
- **index.html**: HTML template yang akan di-render
- **values.json**: (Generated saat publish) Berisi nilai dari user

Contoh `spec.json`:
```json
{
  "name": "Simple Landing Page",
  "description": "A clean and modern landing page",
  "thumbnail": "/static/templates/simple-landing.png",
  "fields": {
    "site_title": {
      "label": "Site Title",
      "type": "text",
      "default": "Welcome to My Website"
    },
    "primary_color": {
      "label": "Primary Color",
      "type": "color",
      "default": "#3b82f6"
    }
  }
}
```

### 3. Builder

- User mengisi form berdasarkan spec.json
- Data disimpan ke workspace.values
- Saat publish, values.json dibuat dan dicopy ke host folder

### 4. Publish

Saat user publish:
1. Template dicopy ke folder `published_sites/[subdomain]/`
2. File `values.json` dibuat dengan data user
3. Website bisa diakses di `/sites/[subdomain]/`

### 5. Rendering

Template menggunakan JavaScript untuk load `values.json`:
```javascript
fetch('values.json')
  .then(response => response.json())
  .then(values => {
    // Update DOM dengan values
    document.getElementById('site-title').textContent = values.site_title;
  });
```

## Membuat Template Baru

1. Buat folder di `templates/[nama-template]/`
2. Buat `spec.json` dengan field definitions
3. Buat `index.html` dengan placeholder IDs
4. Load dan apply values dengan JavaScript

Contoh structure:
```
templates/
└── my-template/
    ├── spec.json
    ├── index.html
    └── assets/
        └── style.css
```

## Environment Variables

```env
PORT=3000                              # Port server
JWT_SECRET=your-secret-key             # Secret untuk JWT
HOST_FOLDER=./published_sites          # Folder untuk published sites
TEMPLATES_FOLDER=./templates           # Folder templates
```

## Development

### Struktur Handler

```go
// handlers/auth.go - Authentication
func Login(c *fiber.Ctx) error
func Register(c *fiber.Ctx) error
func Logout(c *fiber.Ctx) error

// handlers/workspace.go - Workspace management
func CreateWorkspace(c *fiber.Ctx) error
func ListTemplates(c *fiber.Ctx) error
func GetTemplate(c *fiber.Ctx) error

// handlers/builder.go - Builder dan publish
func Builder(c *fiber.Ctx) error
func SaveBuilder(c *fiber.Ctx) error
func PublishSite(c *fiber.Ctx) error
func ServeSite(c *fiber.Ctx) error
```

### HTMX Integration

Semua form menggunakan HTMX untuk interaksi tanpa reload:
```html
<form hx-post="/login" hx-target="#message">
  <!-- Form fields -->
</form>
```

## Production Deployment

1. **Gunakan database real** (PostgreSQL/MySQL)
2. **Set JWT_SECRET yang kuat**
3. **Enable HTTPS**
4. **Setup reverse proxy** (Nginx/Caddy)
5. **Configure backup** untuk published sites
6. **Monitor disk space** untuk user uploads

## Fitur Future

- [ ] File upload untuk images/assets
- [ ] Custom domains untuk published sites
- [ ] Template marketplace
- [ ] Collaboration/team workspaces
- [ ] Analytics untuk published sites
- [ ] SEO tools
- [ ] A/B testing

## License

MIT

## Support

Jika ada pertanyaan atau issue, silakan buat issue di repository ini.
