<script>
	import { connectToFileNotifications, closeWebSocket } from '../lib/websocket.js';
	import QRCode from 'qrcode';
	import * as Alert from "$lib/components/ui/alert/index.js";
	import CircleCheckIcon from "@lucide/svelte/icons/circle-check";
    import * as Item from "$lib/components/ui/item/index.js";
    import { Spinner } from "$lib/components/ui/spinner/index.js";
    import {Separator} from "$lib/components/ui/separator/index.ts";
	import CopyableInput from './CopyableInput.svelte';

    let { fileName, fileSize, downloadURL, curlCommand, fileID } = $props();

	let downloadStatus = $state('pending');
	let websocket = $state(null);
	let qrCodeDataURL = $state('');

	// Generate QR code and connect to WebSocket
	$effect(() => {
		// Generate QR code for download URL
		(async () => {
			try {
				qrCodeDataURL = await QRCode.toDataURL(downloadURL, {
					width: 200,
					margin: 2,
					color: {
						dark: '#000000',
						light: '#ffffff'
					}
				});
			} catch (error) {
				console.error('Failed to generate QR code:', error);
			}
		})();

		// Connect to WebSocket for download notifications
		if (fileID) {
			websocket = connectToFileNotifications(fileID, () => {
				downloadStatus = 'downloaded';
				closeWebSocket(websocket);
				websocket = null;
			});
		}

		// Cleanup on destroy
		return () => {
			closeWebSocket(websocket);
		};
	});

</script>

<div class="space-y-4 rounded-lg border-1 p-6 mt-6">
    <ul>
        <li><strong>File:</strong> {fileName}</li>
        <li><strong>Size:</strong> {fileSize}</li>
    </ul>
    {#if downloadStatus === 'pending'}
        <div class="mt-6 flex flex-col gap-6 md:flex-row">
            <div class="flex-1 space-y-4 h-14">
                <div class="flex h-full w-full gap-4 [--radius:1rem]">
                    <Item.Root variant="muted">
                        <Item.Media>
                            <CircleCheckIcon />
                        </Item.Media>
                        <Item.Content>
                            <Item.Title class="line-clamp-1">Upload successful!</Item.Title>
                        </Item.Content>
                    </Item.Root>

                    <div class="flex h-full gap-4 [--radius:1rem]">
                        <Item.Root variant="muted">
                            <Item.Media>
                                <Spinner />
                            </Item.Media>
                            <Item.Content>
                                <Item.Title class="line-clamp-1">File has not been downloaded yet...</Item.Title>
                            </Item.Content>
                        </Item.Root>
                    </div>
                </div>
                <CopyableInput
                    id="downloadURL"
                    label="Download URL"
                    value={downloadURL}
                />

                <CopyableInput
                    id="curlCommand"
                    label="cURL command"
                    value={curlCommand}
                />

            </div>

            {#if qrCodeDataURL}
                <div class="flex flex-col items-center">
                    <span class="mb-2 block text-sm font-semibold text-card-foreground">QR Code</span>
                    <div class="rounded-lg border border-border bg-card shadow-md">
                        <div class="p-4">
                            <img src={qrCodeDataURL} alt="QR Code for download" class="h-48 w-48" />
                        </div>
                    </div>
                </div>
            {/if}
        </div>


    {:else if downloadStatus === 'downloaded'}
        <div class="mt-4 flex gap-3 rounded-md border border-green-500/50 bg-green-500/10 p-4">
            <svg
                    xmlns="http://www.w3.org/2000/svg"
                    class="h-6 w-6 shrink-0 text-green-600 dark:text-green-400"
                    fill="none"
                    viewBox="0 0 24 24"
                    stroke-width="1.5"
                    stroke="currentColor"
            >
                <path
                        stroke-linecap="round"
                        stroke-linejoin="round"
                        d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"
                />
            </svg>
            <span class="text-sm text-green-700 dark:text-green-400">File has been downloaded!</span>
        </div>
    {/if}
</div>
