<script lang="ts" module>
	export type ChipElement = {
		kind: string;
		value?: string;
	};
</script>

<script lang="ts">
	let { chip } = $props();
</script>

<span class="m-1 mr-0 rounded bg-neutral-200 px-2 py-1 text-nowrap">{chip.kind}</span>
<span class="m-1 mr-0 rounded bg-neutral-200 px-2 py-1">≈</span>

{#if chip.value}
	<span class="m-1 mr-0.5 rounded bg-neutral-200 px-2 py-1 text-nowrap">{chip.value}</span>
{/if}
