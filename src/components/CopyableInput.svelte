<script>
	import CopyButton from './CopyButton.svelte';

	let { label, value, id } = $props();
</script>

<div>
	<label for={id} class="mb-2 block text-sm font-semibold text-card-foreground">{label}</label>
	<div class="flex gap-2 items-center">
		<input
			{id}
			type="text"
			{value}
			readonly
			class="flex-1 h-10 rounded-md border border-input bg-background px-3 py-2 font-mono text-sm text-foreground ring-offset-background focus:outline-none focus:ring-2 focus:ring-ring"
		/>
		<CopyButton textToCopy={value} size="icon" class="h-10 w-10" />
	</div>
</div>
