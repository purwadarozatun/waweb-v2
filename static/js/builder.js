// builder.js — static JS for builder page
// specData, workspaceValues, workspaceID injected by inline <script> in builder.html

function switchBuilderTab(tab) {
    ['upload', 'preview'].forEach(function(t) {
        document.getElementById('panel-' + t).classList.toggle('active', t === tab);
        document.getElementById('btab-' + t).classList.toggle('active', t === tab);
    });
}

function selectAsset(key, type, label, hasValue) {
    // Switch to upload tab
    switchBuilderTab('upload');

    // Hide empty state
    document.getElementById('empty-pick').style.display = 'none';
    document.getElementById('asset-form-container').style.display = 'block';

    // Build form fields
    var currentValue = workspaceValues[key] || '';
    var fieldHTML = '';

    if (type === 'image' || type === 'file') {
        var hasImg = currentValue && currentValue.trim() !== '';
        var imgWrap = '';
        if (hasImg) {
            imgWrap = '<div class="img-uploaded-box">' +
                '<img src="' + escHtml(currentValue) + '" class="ub-img" alt="' + escHtml(label) + '"' +
                ' onerror="this.src=\'/static/placeholder-image.svg\'; this.style.objectFit=\'contain\'">' +
                '<div class="ub-bar">' +
                '<span class="ub-name"><i class="ti ti-photo"></i> ' + escHtml(key) + '</span>' +
                '<button type="button" class="btn-rm" onclick="clearUpload(\'' + escJs(key) + '\')">' +
                '<i class="ti ti-trash"></i> Remove</button>' +
                '</div></div>';
        } else {
            imgWrap = '<div class="img-placeholder-box" id="placeholder-' + escHtml(key) + '">' +
                '<img src="/static/placeholder-image.svg" alt="placeholder"' +
                ' style="width:100%;height:160px;object-fit:contain;border-radius:8px;display:block;">' +
                '<div class="img-placeholder-label">Please Pick Image</div>' +
                '</div>';
        }

        var uploadLabel = hasImg ? 'Change Image' : 'Upload Image';
        var statusClass = hasImg ? 'ok' : 'pending';
        var statusIcon = hasImg
            ? '<i class="ti ti-circle-check"></i> Image set'
            : '<i class="ti ti-photo-off"></i> No image';

        fieldHTML = '<div class="asset-card">' +
            '<div class="ac-head"><div>' +
            '<div class="ac-title">' + escHtml(key) + '</div>' +
            '<div class="ac-path">\u2192 uploads/' + escHtml(key) + '/</div>' +
            '</div><span class="type-badge img">image</span></div>' +
            '<div class="form-group">' +
            '<label class="form-label">' + escHtml(label) + '</label>' +
            '<div class="img-upload-wrap" id="img-wrap-' + escHtml(key) + '">' + imgWrap + '</div>' +
            '<div class="img-upload-controls">' +
            '<label class="btn-upload-file"><i class="ti ti-upload"></i> ' + uploadLabel +
            '<input type="file" accept="image/jpeg,image/png,image/gif,image/webp,image/svg+xml"' +
            ' style="display:none" onchange="uploadAsset(this, \'' + escJs(key) + '\', \'' + escJs(label) + '\')">' +
            '</label>' +
            '<span class="upload-status ' + statusClass + '" id="status-' + escHtml(key) + '">' + statusIcon + '</span>' +
            '</div>' +
            '<div id="upload-progress-' + escHtml(key) + '" style="display:none" class="upload-progress-bar">' +
            '<i class="ti ti-loader"></i> Uploading\u2026</div>' +
            '</div>' +
            '<div class="path-info"><i class="ti ti-folder" style="font-size:13px"></i> Output: <strong>uploads/' + escHtml(key) + '/</strong></div>' +
            '</div>';

    } else if (type === 'textarea') {
        fieldHTML = '<div class="asset-card">' +
            '<div class="ac-head"><div>' +
            '<div class="ac-title">' + escHtml(key) + '</div>' +
            '<div class="ac-path">\u2192 value.json["' + escHtml(key) + '"]</div>' +
            '</div><span class="type-badge txt">text</span></div>' +
            '<div class="form-group">' +
            '<label class="form-label">' + escHtml(label) + '</label>' +
            '<textarea name="' + escHtml(key) + '" class="form-textarea" rows="4">' + escHtml(currentValue) + '</textarea>' +
            '</div>' +
            '<div class="path-info"><i class="ti ti-file-code" style="font-size:13px"></i> Output: <strong>assets/value.json</strong></div>' +
            '</div>';

    } else {
        var inputType = (type === 'color') ? 'color' : 'text';
        var inputClass = 'form-input' + (type === 'color' ? ' form-color' : '');
        fieldHTML = '<div class="asset-card">' +
            '<div class="ac-head"><div>' +
            '<div class="ac-title">' + escHtml(key) + '</div>' +
            '<div class="ac-path">\u2192 value.json["' + escHtml(key) + '"]</div>' +
            '</div><span class="type-badge txt">' + escHtml(type) + '</span></div>' +
            '<div class="form-group">' +
            '<label class="form-label">' + escHtml(label) + '</label>' +
            '<input type="' + inputType + '" name="' + escHtml(key) + '"' +
            ' value="' + escHtml(currentValue) + '" class="' + inputClass + '">' +
            '</div>' +
            '<div class="path-info"><i class="ti ti-file-code" style="font-size:13px"></i> Output: <strong>assets/value.json</strong></div>' +
            '</div>';
    }

    document.getElementById('asset-fields').innerHTML = fieldHTML;

    // Update sidebar selection
    document.querySelectorAll('.sb-asset').forEach(function(el) {
        el.classList.remove('sel');
    });
    event.currentTarget.classList.add('sel');
}

// ── Helpers ───────────────────────────────────────────────────────
function escHtml(str) {
    return String(str)
        .replace(/&/g, '&amp;')
        .replace(/</g, '&lt;')
        .replace(/>/g, '&gt;')
        .replace(/"/g, '&quot;')
        .replace(/'/g, '&#39;');
}
function escJs(str) {
    return String(str).replace(/\\/g, '\\\\').replace(/'/g, "\\'");
}

// ── Progress ──────────────────────────────────────────────────────
function updateProgress() {
    var fields = specData.fields || {};
    var total = Object.keys(fields).length;
    var filled = Object.keys(workspaceValues).filter(function(k) { return workspaceValues[k]; }).length;
    var percent = total > 0 ? Math.round((filled / total) * 100) : 0;
    document.getElementById('prog-label').textContent = filled + ' / ' + total;
    document.getElementById('prog-fill').style.width = percent + '%';
}

// ── Device mode ───────────────────────────────────────────────────
function setDeviceMode(mode) {
    var iframe = document.getElementById('preview-iframe');
    iframe.classList.remove('desktop', 'tablet', 'mobile');
    iframe.classList.add(mode);
    ['d', 't', 'm'].forEach(function(d) {
        document.getElementById('dev-' + d).classList.remove('active');
    });
    var btnMap = { desktop: 'd', tablet: 't', mobile: 'm' };
    document.getElementById('dev-' + btnMap[mode]).classList.add('active');
}

function refreshPreview() {
    var iframe = document.getElementById('preview-iframe');
    if (iframe) { iframe.src = iframe.src; }
}

document.addEventListener('htmx:afterSwap', function(evt) {
    if (evt.detail.target.id === 'save-message') {
        setTimeout(refreshPreview, 500);
    }
});

// ── Image Upload ──────────────────────────────────────────────────
function uploadAsset(input, fieldKey, label) {
    var file = input.files[0];
    if (!file) return;

    var progressEl = document.getElementById('upload-progress-' + fieldKey);
    if (progressEl) progressEl.style.display = 'flex';

    var formData = new FormData();
    formData.append('file', file);
    formData.append('field_key', fieldKey);

    fetch('/builder/' + workspaceID + '/upload', { method: 'POST', body: formData })
        .then(function(r) {
            if (!r.ok) return r.text().then(function(t) { throw new Error(t); });
            return r.json();
        })
        .then(function(data) {
            if (progressEl) progressEl.style.display = 'none';
            workspaceValues[fieldKey] = data.url;

            var wrap = document.getElementById('img-wrap-' + fieldKey);
            if (wrap) {
                wrap.innerHTML =
                    '<div class="img-uploaded-box">' +
                    '<img src="' + escHtml(data.url) + '" class="ub-img" alt="' + escHtml(label) + '"' +
                    ' onerror="this.src=\'/static/placeholder-image.svg\'; this.style.objectFit=\'contain\'">' +
                    '<div class="ub-bar">' +
                    '<span class="ub-name"><i class="ti ti-photo"></i> ' + escHtml(fieldKey) + '</span>' +
                    '<button type="button" class="btn-rm" onclick="clearUpload(\'' + escJs(fieldKey) + '\')">' +
                    '<i class="ti ti-trash"></i> Remove</button>' +
                    '</div></div>';
            }

            var statusEl = document.getElementById('status-' + fieldKey);
            if (statusEl) {
                statusEl.className = 'upload-status ok';
                statusEl.innerHTML = '<i class="ti ti-circle-check"></i> Image set';
            }

            updateSidebarDot(fieldKey, true);
            updateProgress();
            setTimeout(refreshPreview, 500);
        })
        .catch(function(err) {
            if (progressEl) progressEl.style.display = 'none';
            alert('Upload failed: ' + err.message);
        });
}

function clearUpload(fieldKey) {
    workspaceValues[fieldKey] = '';
    var wrap = document.getElementById('img-wrap-' + fieldKey);
    if (wrap) {
        wrap.innerHTML =
            '<div class="img-placeholder-box" id="placeholder-' + escHtml(fieldKey) + '">' +
            '<img src="/static/placeholder-image.svg" alt="placeholder"' +
            ' style="width:100%;height:160px;object-fit:contain;border-radius:8px;display:block;">' +
            '<div class="img-placeholder-label">Please Pick Image</div>' +
            '</div>';
    }
    var statusEl = document.getElementById('status-' + fieldKey);
    if (statusEl) {
        statusEl.className = 'upload-status pending';
        statusEl.innerHTML = '<i class="ti ti-photo-off"></i> No image';
    }
    updateSidebarDot(fieldKey, false);
    updateProgress();
}

function updateSidebarDot(fieldKey, hasValue) {
    document.querySelectorAll('.sb-asset').forEach(function(el) {
        var oc = el.getAttribute('onclick') || '';
        if (oc.indexOf("'" + fieldKey + "'") !== -1) {
            var dot = el.querySelector('.status-dot');
            if (dot) {
                dot.classList.toggle('done', hasValue);
                dot.classList.toggle('pend', !hasValue);
            }
        }
    });
}

// ── Publish Drawer ────────────────────────────────────────────────
function togglePublishDrawer() {
    var drawer = document.getElementById('publish-drawer');
    var backdrop = document.getElementById('drawer-backdrop');
    if (drawer.classList.contains('open')) {
        drawer.classList.remove('open');
        backdrop.classList.remove('open');
    } else {
        drawer.classList.add('open');
        backdrop.classList.add('open');
    }
}

function closePublishDrawer() {
    var drawer = document.getElementById('publish-drawer');
    var backdrop = document.getElementById('drawer-backdrop');
    drawer.classList.remove('open');
    backdrop.classList.remove('open');
}
