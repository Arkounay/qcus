<script>
	import { validatePassword, getPublicConfig } from './lib/api.js';
	import LoginForm from './components/LoginForm.svelte';
	import UploadArea from './components/UploadArea.svelte';
	import UploadResult from './components/UploadResult.svelte';
	import CodeBlock from './components/CodeBlock.svelte';
    import * as Alert from "$lib/components/ui/alert/index.js";
    import CircleAlertIcon from "@lucide/svelte/icons/circle-alert";
    import InfoIcon from "@lucide/svelte/icons/info";
    import {Separator} from "$lib/components/ui/separator/index.ts";

	let isLoggedIn = $state(false);
	let uploadPassword = $state('');
	let uploadResult = $state(null);
	let showDefaultPasswordHint = $state(true);
	let fileExpiryMinutes = $state(10);
	let maxFileSizeMB = $state(100);

	// Check session on mount
	$effect(() => {
		(async () => {
			const config = await getPublicConfig();
			showDefaultPasswordHint = config.isDefaultPassword;
			fileExpiryMinutes = config.fileExpiryMinutes || 10;
			maxFileSizeMB = config.maxFileSizeMB || 100;

			const storedPassword = sessionStorage.getItem('uploadPassword');
			if (storedPassword) {
				validateAndLogin(storedPassword);
			}
		})();
	});

	async function validateAndLogin(password) {
		const isValid = await validatePassword(password);
		if (isValid) {
			uploadPassword = password;
			isLoggedIn = true;
		} else {
			sessionStorage.removeItem('uploadPassword');
		}
	}

	function handleLogin(data) {
		const password = data.password;
		uploadPassword = password;
		sessionStorage.setItem('uploadPassword', password);
		isLoggedIn = true;
	}

	function handleLogout() {
		uploadPassword = '';
		sessionStorage.removeItem('uploadPassword');
		isLoggedIn = false;
		uploadResult = null;
	}

	function handleUploadSuccess(data) {
		uploadResult = data;
	}

	function handleUnauthorized() {
		handleLogout();
	}

	function handleUploadAnother() {
		uploadResult = null;
	}
</script>

{#if !isLoggedIn}
	<LoginForm onlogin={handleLogin} />
{:else}
	<div class="min-h-screen bg-muted/30 p-4 md:p-8">
		<div class="mx-auto max-w-4xl">
			<div class="rounded-lg border border-border bg-card shadow-lg">
				<div class="p-6">
					<!-- Header -->
					<div class="mb-6 flex flex-col items-start justify-between gap-4 sm:flex-row sm:items-center">
						<h1 class="flex items-center gap-2 text-2xl font-semibold text-card-foreground md:text-3xl">
							Temporary File Upload Server
						</h1>
						<button
							onclick={handleLogout}
							class="inline-flex items-center gap-2 rounded-md bg-destructive px-3 py-2 text-sm font-medium text-destructive-foreground shadow-sm hover:bg-destructive/90"
						>
							<svg
								xmlns="http://www.w3.org/2000/svg"
								class="h-4 w-4"
								fill="none"
								viewBox="0 0 24 24"
								stroke="currentColor"
							>
								<path
									stroke-linecap="round"
									stroke-linejoin="round"
									stroke-width="2"
									d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1"
								/>
							</svg>
							Logout
						</button>
					</div>

					<!-- Upload Area -->
					{#if !uploadResult}
                        <Alert.Root class="mb-6">
                            <CircleAlertIcon class="size-4" />
                            <Alert.Title>Upload files and get a temporary download link</Alert.Title>
                            <Alert.Description>Maximum file size: {maxFileSizeMB}MB. Files are automatically deleted after download or after {fileExpiryMinutes} {fileExpiryMinutes === 1 ? 'minute' : 'minutes'}.</Alert.Description>
                        </Alert.Root>

						<UploadArea
							{uploadPassword}
							onuploadsuccess={handleUploadSuccess}
							onunauthorized={handleUnauthorized}
						/>
					{/if}

					<!-- Upload Result -->
					{#if uploadResult}
						<UploadResult
							fileName={uploadResult.fileName}
							fileSize={uploadResult.fileSize}
							downloadURL={uploadResult.downloadURL}
							curlCommand={uploadResult.curlCommand}
							fileID={uploadResult.fileID}
						/>

						<!-- Upload Another Button -->
						<div class="mt-6">
							<button
								onclick={handleUploadAnother}
								class="inline-flex items-center gap-2 rounded-md bg-primary px-4 py-2 text-sm font-medium text-primary-foreground shadow-sm hover:bg-primary/90"
							>
								<svg
									xmlns="http://www.w3.org/2000/svg"
									class="h-4 w-4"
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
								Upload Another
							</button>
						</div>
					{/if}

					<Separator class="my-8" />

					<div class="space-y-4 rounded-lg bg-muted/30 p-6">
						<h3 class="text-xl font-bold text-card-foreground">CLI Upload Examples</h3>

						<div class="space-y-4">
							<div>
								<h4 class="mb-2 font-semibold text-card-foreground">
									Method 1: Multipart form upload
								</h4>
								<CodeBlock
									code={`curl -F "file=@yourfile.txt" -H "X-Upload-Password: ${showDefaultPasswordHint ? 'demo' : 'YOUR_PASSWORD'}" http://${window.location.host}`}
								/>
							</div>

							<div>
								<h4 class="mb-2 font-semibold text-card-foreground">
									Method 2: Direct file upload (filename from URL path)
								</h4>
								<CodeBlock
									code={`curl -T yourfile.txt -H "X-Upload-Password: ${showDefaultPasswordHint ? 'demo' : 'YOUR_PASSWORD'}" http://${window.location.host}/yourfile.txt`}
								/>
								<div class="mt-2">
									<CodeBlock
										code={`curl --upload-file yourfile.txt -H "X-Upload-Password: ${showDefaultPasswordHint ? 'demo' : 'YOUR_PASSWORD'}" http://${window.location.host}/yourfile.txt`}
									/>
								</div>
							</div>

							<div>
								<h4 class="mb-2 font-semibold text-card-foreground">Example Response:</h4>
								<CodeBlock
                                    showCopyButton={false}
									leadingRelaxed={true}
									code={`$ curl -F "file=@example.txt" -H "X-Upload-Password: ${showDefaultPasswordHint ? 'demo' : 'YOUR_PASSWORD'}" http://${window.location.host}
<qr code>
File uploaded successfully!
Original name: example.txt
Download URL: http://${window.location.host}/download/a1b2c3d4e5f6...
cURL command: curl -o "example.txt" http://${window.location.host}/download/a1b2c3d4e5f6...`}
									textToCopy={`curl -F "file=@example.txt" -H "X-Upload-Password: ${showDefaultPasswordHint ? 'demo' : 'YOUR_PASSWORD'}" http://${window.location.host}`}
								/>
							</div>
						</div>

						{#if showDefaultPasswordHint}
							<Alert.Root>
								<InfoIcon class="size-4" />
								<Alert.Title>Custom Password</Alert.Title>
								<Alert.Description>
                                    <div>
									Set <code class="inline-block rounded bg-card px-1.5 py-0.5 font-mono text-xs font-semibold">UPLOAD_PASSWORD</code> environment variable to change the password.
                                    </div>
								</Alert.Description>
							</Alert.Root>
						{/if}
					</div>
				</div>
			</div>
		</div>
	</div>
{/if}
