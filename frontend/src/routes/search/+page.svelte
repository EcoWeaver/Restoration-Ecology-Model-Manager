<script lang="ts">
	import Header from "$lib/Header.svelte";
	import type { ChipElement } from "$lib/search/Chip.svelte";

	let { data } = $props();
	let { items, params: searchParams } = $derived(data);

	let params = $derived<ChipElement[]>(searchParams.map(([kind, value]) => ({ kind, value })));
</script>

<svelte:head>
	<title>search by metadata</title>
</svelte:head>

<Header {params} />

{#if items.length > 0}
	<div class="p-4 pt-2">
		<div class="table w-full border-collapse p-2 text-lg">
			<div class="mb-2 table-row border-b-1 border-neutral-300">
				<span class="table-cell p-3 font-semibold">title</span>
				<span class="table-cell p-3 font-semibold">version</span>
				<span class="table-cell p-3 font-semibold">identifier</span>
			</div>
			{#each items as { uuid, identifier, title }, idx (idx)}
				<a
					href={`/view/${uuid}`}
					class="table-row border-b-1 border-neutral-300 outline-indigo-400 last:border-none hover:bg-neutral-200 focus-visible:outline-2 dark:hover:bg-neutral-700"
				>
					<span class="table-cell p-3">{title}</span>
					<span class="table-cell p-3">{identifier}</span>
				</a>
			{/each}
		</div>
	</div>
{:else}
	<h1 class="m-4 mt-2 text-2xl font-semibold">no models found</h1>
{/if}
