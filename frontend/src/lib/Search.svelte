<script lang="ts">
	import MetadataSearch from "./MetadataSearch.svelte";
	import IdentifierSearch from "./IdentifierSearch.svelte";
	import type { ChipElement } from "./search/Chip.svelte";
	import { BookMarked, ChevronDown, ChevronUp, Fingerprint } from "@lucide/svelte";
	import type { Component } from "svelte";

	let {
		params = [],
		identifier = "",
		selectedSearch: selected,
	}: {
		params?: ChipElement[];
		identifier?: string;
		selectedSearch: "metadata" | "identifier";
	} = $props();

	type Search = { icon: Component; tooltip: string };
	let searches: Record<"metadata" | "identifier", Search> = {
		metadata: { icon: BookMarked, tooltip: "search by metadata" },
		identifier: { icon: Fingerprint, tooltip: "search by identifier" },
	};

	let open = $state(false);
	let current = $derived(searches[selected]);

	function select(value: "metadata" | "identifier") {
		open = false;
		selected = value;
	}

	function onfocusout(event: FocusEvent & { currentTarget: EventTarget & HTMLDivElement }) {
		let { currentTarget, relatedTarget } = event;
		if (relatedTarget instanceof HTMLElement && currentTarget.contains(relatedTarget)) {
			return;
		}

		open = false;
	}
</script>

<div
	class="relative z-0 h-10 rounded-l-md border border-r-0 border-neutral-600 outline-indigo-400 focus-within:outline-2"
	{onfocusout}
>
	<button
		onclick={() => (open = !open)}
		title={current.tooltip}
		class="flex h-full cursor-pointer items-center rounded-l-md pr-2 pl-3 outline-none"
	>
		<span class="pr-2">
			<current.icon />
		</span>

		{#if open}
			<ChevronUp />
		{:else}
			<ChevronDown />
		{/if}
	</button>

	{#if open}
		<div
			class="absolute mt-2 -ml-4 rounded-md border border-neutral-600 bg-white px-2 py-1 shadow-xl dark:bg-neutral-800"
		>
			{#each Object.entries(searches) as [value, { icon: Icon, tooltip }] (value)}
				<button
					onclick={() => select(value)}
					class="my-1 flex w-full cursor-pointer rounded-md p-2 outline-indigo-400 hover:bg-neutral-200 focus-visible:outline-2 dark:hover:bg-neutral-700"
					title={tooltip}
				>
					<Icon />
					<span class="pl-2">{value}</span>
				</button>
			{/each}
		</div>
	{/if}
</div>

{#if selected === "metadata"}
	<MetadataSearch {params} />
{:else if selected === "identifier"}
	<IdentifierSearch {identifier} />
{/if}
