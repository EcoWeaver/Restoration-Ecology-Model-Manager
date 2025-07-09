<script lang="ts">
	import { browser } from "$app/environment";
	import Header from "$lib/Header.svelte";
	import type { ChipElement } from "$lib/search/Chip.svelte";
	import { ChevronDown, ChevronUp } from "@lucide/svelte";

	let { data } = $props();
	let { items: originalItems, params: searchParams, metadataElements } = $derived(data);

	let params = $derived<ChipElement[]>(searchParams.map(([kind, value]) => ({ kind, value })));

	function dateFormat(d: Date, short: boolean): string {
		const year = d.getFullYear();
		const month = (d.getMonth() + 1).toString().padStart(2, "0");
		const day = d.getDate().toString().padStart(2, "0");

		if (short) {
			return `${year}-${month}-${day}`;
		}

		const hours = d.getHours().toString().padStart(2, "0");
		const minutes = d.getMinutes().toString().padStart(2, "0");
		const seconds = d.getSeconds().toString().padStart(2, "0");

		const zoneSign = d.getTimezoneOffset() < 0 ? "+" : "-";
		const zoneHours = Math.floor(Math.abs(d.getTimezoneOffset()) / 60)
			.toString()
			.padStart(2, "0");
		const zoneMinutes = (Math.abs(d.getTimezoneOffset()) % 60).toString().padStart(2, "0");

		return `${year}-${month}-${day}T${hours}:${minutes}:${seconds}${zoneSign}${zoneHours}:${zoneMinutes}`;
	}

	type SortColumn = "title" | "modelType" | "identifier" | "uploaded";
	type SortType = "ascending" | "descending";
	type SortState = { column: SortColumn; type: SortType };

	let sortState = $state<SortState | undefined>(undefined);

	function toggleState(column: SortColumn) {
		if (sortState?.column === column) {
			if (sortState.type === "ascending") {
				sortState.type = "descending";
			} else {
				sortState = undefined;
			}
		} else {
			sortState = { column, type: "ascending" };
		}
	}

	$effect(() => {
		// reset sort state when navigating between /search/ pages
		$state.snapshot(originalItems);
		sortState = undefined;
	});

	let items = $derived.by(() => {
		if (sortState == null) {
			return originalItems;
		} else {
			const { column, type } = sortState;
			return [...originalItems].sort((l, r) => {
				const lv = l[column];
				const rv = r[column];

				const cmp = lv.localeCompare(rv, "en");
				return type === "ascending" ? cmp : -cmp;
			});
		}
	});
</script>

<svelte:head>
	<title>search by metadata</title>
</svelte:head>

<Header {metadataElements} {params} />

<!--
	only render the contents on the client as the server timezone may be different from that of the client,
	making the formatted time reformat after hydration, which is a lot more jarring than the content being
	loaded a little later than the header.
-->
{#if browser}
	{#if items.length > 0}
		<div class="p-4 pt-2">
			<div class="table w-full border-collapse p-2 text-lg">
				<div class="mb-2 table-row border-b-1 border-neutral-300">
					{#snippet header(column: SortColumn, value?: string)}
						{@const text = value ?? column}
						<span class="table-cell px-1 py-3 font-semibold">
							<button
								title={sortState?.column === column ? `${text} (${sortState.type})` : text}
								class="outline-offset flex cursor-pointer items-center gap-1 rounded-md px-2 text-nowrap outline-indigo-400 focus-visible:outline-2"
								onclick={() => toggleState(column)}
							>
								{text}

								{#if sortState?.column === column && sortState.type === "ascending"}
									<ChevronUp size={20} />
								{:else if sortState?.column === column && sortState.type === "descending"}
									<ChevronDown size={20} />
								{:else}
									<div class="invisible w-[20px]"></div>
								{/if}
							</button>
						</span>
					{/snippet}

					{@render header("title")}

					{@render header("modelType", "model-type")}

					{@render header("identifier")}

					{@render header("uploaded")}
				</div>
				{#each items as { uuid, title, modelType, identifier, uploaded } (uuid)}
					{@const date = new Date(uploaded)}
					<a
						href={`/view/${uuid}`}
						class="table-row border-b-1 border-neutral-300 outline-indigo-400 last:border-none hover:bg-neutral-200 focus-visible:outline-2 dark:hover:bg-neutral-700"
						data-sveltekit-preload-data="tap"
						data-sveltekit-preload-code="hover"
					>
						<span class="table-cell p-3">{title}</span>
						<span class="table-cell p-3">{modelType}</span>
						<span class="table-cell p-3 text-nowrap">{identifier}</span>
						{#if browser}
							<span title={dateFormat(date, false)} class="table-cell p-3">
								{dateFormat(date, true)}
							</span>
						{/if}
					</a>
				{/each}
			</div>
		</div>
	{:else}
		<h1 class="m-4 mt-2 text-2xl font-semibold">no models found</h1>
	{/if}
{/if}
