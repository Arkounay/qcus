<script>
	import CopyButton from './CopyButton.svelte';

	let { code, textToCopy = code, leadingRelaxed = false, showCopyButton = true } = $props();
</script>

<div class="relative rounded-md bg-slate-950" class:pr-12={showCopyButton}>
	<div class="scrollbar-thin overflow-x-auto p-4 text-sm text-slate-50 {leadingRelaxed ? 'leading-relaxed' : ''}">
		<pre><code>{code}</code></pre>
	</div>
	{#if showCopyButton}
		<div class="absolute right-2 top-2">
			<CopyButton {textToCopy} />
		</div>
	{/if}
</div>

<style>
	.scrollbar-thin {
		scrollbar-width: thin;
		scrollbar-color: rgba(255, 255, 255, 0.2) transparent;
	}
	.scrollbar-thin::-webkit-scrollbar {
		height: 6px;
	}
	.scrollbar-thin::-webkit-scrollbar-track {
		background: transparent;
	}
	.scrollbar-thin::-webkit-scrollbar-thumb {
		background-color: rgba(255, 255, 255, 0.2);
		border-radius: 3px;
	}
	.scrollbar-thin::-webkit-scrollbar-thumb:hover {
		background-color: rgba(255, 255, 255, 0.3);
	}
</style>
