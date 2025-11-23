<script>
	import ClipboardCopyIcon from "@lucide/svelte/icons/clipboard-copy";
	import CheckIcon from "@lucide/svelte/icons/check";
	import * as Tooltip from "$lib/components/ui/tooltip/index.js";
	import { buttonVariants } from "$lib/components/ui/button/index.js";

	let { textToCopy, size = "icon-sm", class: className = "" } = $props();

	let copied = $state(false);

	async function handleCopy() {
		try {
			await navigator.clipboard.writeText(textToCopy);
			copied = true;
			setTimeout(() => {
				copied = false;
			}, 2000);
		} catch (err) {
			console.error('Failed to copy:', err);
		}
	}
</script>

<Tooltip.Provider>
	<Tooltip.Root open={copied ? true : undefined}>
		<Tooltip.Trigger onclick={handleCopy} class="{buttonVariants({ variant: "outline", size })} {className}">
			{#if copied}
				<CheckIcon class="text-green-600"></CheckIcon>
			{:else}
				<ClipboardCopyIcon class="text-black"></ClipboardCopyIcon>
			{/if}
		</Tooltip.Trigger>
		<Tooltip.Content>
			<p>{copied ? 'Copied!' : 'Copy'}</p>
		</Tooltip.Content>
	</Tooltip.Root>
</Tooltip.Provider>
