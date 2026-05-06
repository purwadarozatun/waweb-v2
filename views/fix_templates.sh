#!/bin/bash

# Common HTML header
read -r -d '' HEADER << 'HEADER_EOF'
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.Title}}</title>
    <script src="https://unpkg.com/htmx.org@1.9.10"></script>
    <link href="https://cdn.jsdelivr.net/npm/tailwindcss@2.2.19/dist/tailwind.min.css" rel="stylesheet">
    <style>
        .alert {
            padding: 1rem;
            margin: 1rem 0;
            border-radius: 0.375rem;
        }
        .alert-error {
            background-color: #fee;
            color: #c00;
            border: 1px solid #fcc;
        }
        .alert-success {
            background-color: #efe;
            color: #0a0;
            border: 1px solid #cfc;
        }
        .htmx-indicator {
            display: none;
        }
        .htmx-request .htmx-indicator {
            display: inline;
        }
        .htmx-request.htmx-indicator {
            display: inline;
        }
    </style>
</head>
<body class="bg-gray-100">
HEADER_EOF

FOOTER='</body>
</html>'

# Process each file
for file in register.html dashboard.html new-workspace.html templates.html builder.html; do
    echo "Processing $file..."
    # Extract content between {{define "embed"}} and {{end}}
    sed -n '/{{define "embed"}}/,/{{end}}/p' "$file" | sed '1d;$d' > "${file}.content"
    # Combine header + content + footer
    {
        echo "$HEADER"
        cat "${file}.content"
        echo "$FOOTER"
    } > "${file}.new"
    mv "${file}.new" "$file"
    rm "${file}.content"
done

echo "Done!"
