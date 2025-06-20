<script lang="ts">
	import MetadataSearch from "./MetadataSearch.svelte";
	import IdentifierSearch from "./IdentifierSearch.svelte";
	import type { ChipElement } from "./search/Chip.svelte";
	import { BookMarked, Fingerprint } from "@lucide/svelte";
	import type { Component } from "svelte";

	let {
		params = [],
		identifier = "",
		selectedSearch,
	}: {
		params?: ChipElement[];
		identifier?: string;
		selectedSearch: "metadata" | "identifier";
	} = $props();

	type Searches = { value: "metadata" | "identifier"; icon: Component; tooltip: string };
	let searches: Searches[] = [
		{ value: "metadata", icon: BookMarked, tooltip: "search by metadata" },
		{ value: "identifier", icon: Fingerprint, tooltip: "search by identifier" },
	];
</script>

<div class="flex" role="radiogroup">
	{#each searches as { value, icon: Icon, tooltip }}
		<label
			title={tooltip}
			class="flex h-10 cursor-pointer place-items-center border border-r-0 border-neutral-600 px-2 first:rounded-l-md hover:bg-neutral-200 has-checked:bg-indigo-300 hover:has-checked:bg-indigo-400 dark:hover:bg-neutral-700 dark:has-checked:bg-indigo-500"
		>
			<input class="hidden" type="radio" name="search" {value} bind:group={selectedSearch} />
			<Icon />
		</label>
	{/each}
</div>

{#if selectedSearch === "metadata"}
	<MetadataSearch {params} />
{:else if selectedSearch === "identifier"}
	<IdentifierSearch {identifier} />
{/if}
