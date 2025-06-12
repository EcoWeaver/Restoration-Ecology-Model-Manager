<script lang="ts">
	import MetadataSearch from "./MetadataSearch.svelte";
	import IdentifierSearch from "./IdentifierSearch.svelte";
	import type { ChipElement } from "./search/Chip.svelte";

	let {
		params = [],
		identifier = "",
		selectedSearch,
	}: { params?: ChipElement[]; identifier?: string; selectedSearch: "meta" | "ident" } = $props();

	type Searches = { value: "meta" | "ident"; tooltip: string };

	let searches: Searches[] = [
		{ value: "meta", tooltip: "search by metadata" },
		{ value: "ident", tooltip: "search by identifier" },
	];
</script>

<div class="flex p-4 pr-0" role="radiogroup">
	{#each searches as { value, tooltip }}
		<label
			title={tooltip}
			class="flex h-10 cursor-pointer place-items-center border border-r-0 border-neutral-600 px-2 first:rounded-l-md hover:bg-neutral-200 has-checked:bg-indigo-300 hover:has-checked:bg-indigo-400 dark:hover:bg-neutral-700 dark:has-checked:bg-indigo-500"
		>
			<input class="hidden" type="radio" name="search" {value} bind:group={selectedSearch} />
			<span>{value}</span>
		</label>
	{/each}
</div>

{#if selectedSearch == "meta"}
	<MetadataSearch {params} />
{:else if selectedSearch == "ident"}
	<IdentifierSearch {identifier} />
{/if}
