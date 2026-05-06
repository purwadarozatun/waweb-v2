# Quick Start Guide - WaWeb V2

## Prerequisites

- Go 1.21 or higher
- Git (optional)

## Installation

### 1. Install Dependencies

```bash
make install
# atau
go mod download
```

### 2. Setup Environment

Copy file `.env.example` ke `.env`:
```bash
cp .env.example .env
```

Edit `.env` sesuai kebutuhan (opsional):
```env
PORT=3000
JWT_SECRET=your-secret-key-change-this-in-production
HOST_FOLDER=./published_sites
TEMPLATES_FOLDER=./templates
```

## Running the Application

### Method 1: Using Make (Recommended)

```bash
# Build dan jalankan
make run

# Atau jalankan development mode
make dev
```

### Method 2: Using run.sh

```bash
./run.sh
```

### Method 3: Direct Go Run

```bash
go run main.go
```

### Method 4: Build dan Run Binary

```bash
# Build
go build -o waweb

# Run
./waweb
```

## Accessing the Application

1. Buka browser: `http://localhost:3000`

2. Login dengan demo account:
   - **Email**: demo@example.com
   - **Password**: password123

## First Steps

1. **Login** dengan demo account
2. Klik **"+ New Workspace"**
3. Masukkan nama workspace
4. **Pilih template** yang tersedia
5. **Edit konten** di builder
6. **Publish** dengan subdomain pilihan Anda
7. **Visit site** untuk melihat hasilnya

## Project Structure

```
waweb-v2/
├── templates/          # Template website
│   ├── simple-landing/ # Template landing page
│   └── portfolio/      # Template portfolio
├── views/              # HTML views
├── static/             # Static files
├── published_sites/    # Published websites (generated)
└── ...
```

## Available Templates

1. **Simple Landing Page**
   - Modern landing page
   - Customizable colors
   - Call-to-action button
   - Features section

2. **Portfolio Template**
   - Personal portfolio
   - Projects showcase
   - About section
   - Contact information

## Creating a Website

### Step 1: Create Workspace
- Dashboard → "+ New Workspace"
- Enter workspace name

### Step 2: Choose Template
- Browse available templates
- Click "Select Template"

### Step 3: Customize in Builder
- Edit all fields in the form
- Change colors, text, links
- Click "Save Changes"

### Step 4: Publish
- Enter subdomain (lowercase, numbers, hyphens only)
- Click "Publish Site"
- Site will be available at `/sites/[subdomain]/`

## Troubleshooting

### Port Already in Use

Ubah port di `.env`:
```env
PORT=3001
```

### Templates Not Showing

Pastikan folder `templates/` ada dan berisi template dengan `spec.json`

### Published Site Not Found

1. Cek folder `published_sites/[subdomain]/` ada
2. Pastikan sudah klik "Publish Site"
3. Refresh browser

## Development Commands

```bash
# Install dependencies
make install

# Run application
make run

# Development mode (with hot reload if air installed)
make dev

# Build binary
make build

# Format code
make fmt

# Run tests
make test

# Clean build files
make clean
```

## Hot Reload Development

Install Air untuk hot reload:
```bash
go install github.com/cosmtrek/air@latest
```

Kemudian jalankan:
```bash
make dev
```

## Creating Custom Templates

1. Buat folder di `templates/[nama-template]/`

2. Buat `spec.json`:
```json
{
  "name": "My Template",
  "description": "Template description",
  "fields": {
    "field_name": {
      "label": "Field Label",
      "type": "text",
      "default": "Default Value"
    }
  }
}
```

3. Buat `index.html` dengan placeholder IDs:
```html
<h1 id="field-name">Default Text</h1>
```

4. Load values dengan JavaScript:
```javascript
fetch('values.json')
  .then(response => response.json())
  .then(values => {
    document.getElementById('field-name').textContent = values.field_name;
  });
```

## Support Field Types

- `text` - Single line text input
- `textarea` - Multi-line text input
- `color` - Color picker

## Next Steps

1. Explore both templates
2. Create your first website
3. Customize the content
4. Publish and share
5. Create your own custom template

## Tips

- **Subdomain format**: Only lowercase letters, numbers, and hyphens
- **Save often**: Click "Save Changes" regularly
- **Re-publish**: You can re-publish anytime to update your site
- **Multiple sites**: Create multiple workspaces for different websites

## Getting Help

- Check the main [README.md](README.md) for detailed documentation
- Review template examples in `templates/` folder
- Check code in `handlers/` for API reference

Enjoy building with WaWeb V2! 🚀
