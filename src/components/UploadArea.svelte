<script>
	import { uploadFile } from '../lib/api.js';
    import CircleCheckIcon from "@lucide/svelte/icons/circle-check";
    import * as Alert from "$lib/components/ui/alert/index.js";
    import CircleAlertIcon from "@lucide/svelte/icons/circle-alert";
    import { Progress } from "$lib/components/ui/progress/index.js";

	let { uploadPassword, onuploadsuccess, onunauthorized } = $props();

	let selectedFile = $state(null);
	let uploadProgress = $state(0);
	let isUploading = $state(false);
	let uploadError = $state('');
	let isDragging = $state(false);

	function handleFileSelect(event) {
		const file = event.target.files?.[0];
		if (file) {
			selectedFile = file;
			uploadError = '';
			handleUpload();
		}
	}

	function handleDrop(event) {
		event.preventDefault();
		isDragging = false;
		const file = event.dataTransfer.files?.[0];
		if (file) {
			selectedFile = file;
			uploadError = '';
			handleUpload();
		}
	}

	function handleDragOver(event) {
		event.preventDefault();
		isDragging = true;
	}

	function handleDragLeave() {
		isDragging = false;
	}

	async function handleUpload() {
		if (!selectedFile) {
			uploadError = 'Please select a file first';
			return;
		}

		isUploading = true;
		uploadProgress = 0;
		uploadError = '';

		const result = await uploadFile(selectedFile, uploadPassword, (progress) => {
			uploadProgress = progress;
		});

		isUploading = false;

		if (result.success) {
			onuploadsuccess?.(result.data);
			selectedFile = null;
		} else {
			uploadError = result.error;
			if (result.unauthorized) {
				setTimeout(() => onunauthorized?.(), 2000);
			}
		}
	}
</script>

<div
	class="cursor-pointer rounded-lg border-2 border-dashed p-10 text-center transition-all {isDragging
		? 'border-primary bg-primary/10'
		: 'border-border hover:border-primary'}"
	ondrop={handleDrop}
	ondragover={handleDragOver}
	ondragleave={handleDragLeave}
	onclick={() => document.getElementById('fileInput').click()}
	onkeydown={(e) => { if (e.key === 'Enter' || e.key === ' ') { e.preventDefault(); document.getElementById('fileInput').click(); }}}
	role="button"
	tabindex="0"
>
	<svg
		xmlns="http://www.w3.org/2000/svg"
		class="mx-auto mb-4 h-12 w-12 text-muted-foreground"
		fill="none"
		viewBox="0 0 24 24"
		stroke="currentColor"
	>
		<path
			stroke-linecap="round"
			stroke-linejoin="round"
			stroke-width="2"
			d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12"
		/>
	</svg>
	<p class="mb-2 text-muted-foreground">Click to select a file or drag and drop here</p>
	<input type="file" id="fileInput" onchange={handleFileSelect} class="hidden" />
	<button
		type="button"
		class="inline-flex items-center rounded-md bg-primary px-3 py-1.5 text-sm font-medium text-primary-foreground shadow-sm hover:bg-primary/90"
		>Choose File</button
	>
</div>

{#if isUploading}
	<div class="mt-4">
		<div class="mb-1 flex justify-between text-sm text-muted-foreground">
			<span>Uploading...</span>
			<span>{uploadProgress}%</span>
		</div>
        <Progress value={uploadProgress} />
	</div>
{/if}


{#if uploadError}
	<div class="mt-4 flex gap-3 rounded-md border border-red-500/50 bg-red-500/10 p-4">
		<svg
			xmlns="http://www.w3.org/2000/svg"
			class="h-5 w-5 shrink-0 text-red-600 dark:text-red-400"
			fill="none"
			viewBox="0 0 24 24"
			stroke-width="1.5"
			stroke="currentColor"
		>
			<path
				stroke-linecap="round"
				stroke-linejoin="round"
				d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z"
			/>
		</svg>
		<span class="text-sm text-red-700 dark:text-red-400">{uploadError}</span>
	</div>
{/if}
