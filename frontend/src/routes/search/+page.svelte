<script lang="ts">
	import Header from "$lib/Header.svelte";
	import type { ChipElement } from "$lib/search/Chip.svelte";

	let { data } = $props();
	let { identifiers, params: searchParams } = $derived(data);

	let params = $derived<ChipElement[]>(searchParams.map(([kind, value]) => ({ kind, value })));
</script>

<Header {params} />

<div class="m-4 mt-2 table min-w-128 border-collapse p-2 text-lg">
	<div class="mb-2 table-row border-b-1 border-neutral-300">
		<span class="table-cell p-3 font-semibold">title</span>
		<span class="table-cell p-3 font-semibold">identifier</span>
	</div>
	{#each identifiers as { identifier, title }}
		<a
			href="/view/{identifier}"
			class="table-row border-b-1 border-neutral-300 last:border-none hover:bg-neutral-200"
		>
			<span class="table-cell p-3">{title}</span>
			<span class="table-cell p-3">{identifier}</span>
		</a>
	{/each}
</div>
