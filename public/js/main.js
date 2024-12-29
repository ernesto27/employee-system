// Light switcher
const lightSwitches = document.querySelectorAll('.light-switch');
if (lightSwitches.length > 0) {
  lightSwitches.forEach((lightSwitch, i) => {
    if (localStorage.getItem('dark-mode') === 'true') {
      // eslint-disable-next-line no-param-reassign
      lightSwitch.checked = true;
    }
    lightSwitch.addEventListener('change', () => {
      const { checked } = lightSwitch;
      lightSwitches.forEach((el, n) => {
        if (n !== i) {
          // eslint-disable-next-line no-param-reassign
          el.checked = checked;
        }
      });
      document.documentElement.classList.add('[&_*]:!transition-none');
      if (lightSwitch.checked) {
        document.documentElement.classList.add('dark');
        document.querySelector('html').style.colorScheme = 'dark';
        localStorage.setItem('dark-mode', true);
        document.dispatchEvent(new CustomEvent('darkMode', { detail: { mode: 'on' } }));
      } else {
        document.documentElement.classList.remove('dark');
        document.querySelector('html').style.colorScheme = 'light';
        localStorage.setItem('dark-mode', false);
        document.dispatchEvent(new CustomEvent('darkMode', { detail: { mode: 'off' } }));
      }
      setTimeout(() => {
        document.documentElement.classList.remove('[&_*]:!transition-none');
      }, 1);
    });
  });
}


 // File upload handling
 const dropZone = document.getElementById('drop-zone');
 const fileInput = document.getElementById('fileInput');
 const fileList = document.getElementById('file-list');
 let uploadedFiles = new Set();

 // Add allowed file types constant
 const ALLOWED_TYPES = ['image/jpeg', 'image/jpg', 'image/png', 'image/gif', 'image/webp'];

 ['dragenter', 'dragover', 'dragleave', 'drop'].forEach(eventName => {
     dropZone.addEventListener(eventName, preventDefaults, false);
 });

 function preventDefaults(e) {
     e.preventDefault();
     e.stopPropagation();
 }

 ['dragenter', 'dragover'].forEach(eventName => {
     dropZone.addEventListener(eventName, highlight, false);
 });

 ['dragleave', 'drop'].forEach(eventName => {
     dropZone.addEventListener(eventName, unhighlight, false);
 });

 function highlight(e) {
     dropZone.classList.add('scale-[1.02]', 'border-indigo-500', 'bg-indigo-50/50', 'dark:bg-indigo-900/10');
     dropZone.querySelector('svg').classList.add('text-indigo-500');
     dropZone.querySelector('p').textContent = 'Drop files here';
 }

 function unhighlight(e) {
     dropZone.classList.remove('scale-[1.02]', 'border-indigo-500', 'bg-indigo-50/50', 'dark:bg-indigo-900/10');
     dropZone.querySelector('svg').classList.remove('text-indigo-500');
     dropZone.querySelector('p').textContent = 'Drag and drop files here or';
 }

 dropZone.addEventListener('drop', handleDrop, false);
 fileInput.addEventListener('change', handleFiles, false);

 function handleDrop(e) {
     const dt = e.dataTransfer;
     const files = [...dt.files];
     const errorMessage = document.getElementById('errorMessage');
     errorMessage.classList.add('hidden');

     // Validate dropped files
     const invalidFiles = files.filter(file => !ALLOWED_TYPES.includes(file.type));
     if (invalidFiles.length > 0) {
         errorMessage.textContent = `Some files are not allowed. Please only drop image files (PNG, JPG, GIF, WebP).`;
         errorMessage.classList.remove('hidden');
         return;
     }

     handleFiles({ target: { files: dt.files } });
 }

 function handleFiles(e) {
     const files = [...e.target.files];
     const errorMessage = document.getElementById('errorMessage');
     errorMessage.classList.add('hidden');
     
     files.forEach(file => {
         // Validate file type
         if (!ALLOWED_TYPES.includes(file.type)) {
             errorMessage.textContent = `File "${file.name}" is not an allowed image type. Please use PNG, JPG, GIF, or WebP.`;
             errorMessage.classList.remove('hidden');
             return;
         }

         if (!uploadedFiles.has(file.name)) {
             uploadedFiles.add(file.name);
             const fileElement = createFileElement(file);
             fileList.appendChild(fileElement);
         }
     });
 }

 function createFileElement(file) {
     const div = document.createElement('div');
     div.className = 'flex items-center justify-between p-2 bg-slate-50 dark:bg-slate-700 rounded';
     
     // Add image preview
     const preview = document.createElement('div');
     preview.className = 'flex items-center gap-2';
     
     const img = document.createElement('img');
     img.className = 'w-8 h-8 object-cover rounded';
     img.src = URL.createObjectURL(file);
     
     const fileInfo = document.createElement('span');
     fileInfo.className = 'text-sm text-slate-600 dark:text-slate-300';
     fileInfo.textContent = `${file.name} (${formatFileSize(file.size)})`;
     
     preview.appendChild(img);
     preview.appendChild(fileInfo);
     
     const removeButton = document.createElement('button');
     removeButton.className = 'text-rose-500 hover:text-rose-600';
     removeButton.innerHTML = '&times;';
     removeButton.onclick = () => {
         removeFile(file, div);
     };

     div.appendChild(preview);
     div.appendChild(removeButton);
     return div;
 }

 // Clean up object URLs when removing files
 function removeFile(file, element) {
     URL.revokeObjectURL(element.querySelector('img').src);
     element.remove();
     uploadedFiles.delete(file.name);
 }

 function formatFileSize(bytes) {
     if (bytes === 0) return '0 Bytes';
     const k = 1024;
     const sizes = ['Bytes', 'KB', 'MB', 'GB'];
     const i = Math.floor(Math.log(bytes) / Math.log(k));
     return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
 }
