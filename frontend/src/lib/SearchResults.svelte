<script lang="ts">
	import { browser } from "$app/environment";
	import { resolve } from "$app/paths";
	import { ChevronDown, ChevronUp, X } from "@lucide/svelte";
	import type { PostSearchResult } from "./api";

	let {
		items: originalItems,
		route,
		del,
	}: {
		items: PostSearchResult[];
		route: "edit" | "view";
		del?: (uuid: string) => void;
	} = $props();

	// map from identifier => oldest date for that identifier
	const order = new Map<string, string>();
	for (const { identifier, uploaded } of originalItems) {
		const prev = order.get(identifier);
		if (prev == null || prev > uploaded) {
			// set to the current if there is none
			// or replace if the current uploaded date is older
			order.set(identifier, uploaded);
		}
	}

	originalItems.sort((a, b) => {
		// if the identifiers are unequal, compare
		// the oldest creation dates for each identifier
		// this groups the identifiers together
		if (a.identifier !== b.identifier) {
			const l = order.get(a.identifier)!;
			const r = order.get(b.identifier)!;

			// sort the identifiers in reverse (b cmp a) to
			// put newer identifiers at the top
			return r.localeCompare(l, "en");
		}

		// if they have the same identifier sort their uploaded dates
		// don't sort in reverse (a cmp b) to have ascending creation
		// dates per identifier
		return a.uploaded.localeCompare(b.uploaded, "en");
	});

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

	let sortState = $state<SortState>();

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

				const cmp = lv.toString().localeCompare(rv.toString(), "en");
				return type === "ascending" ? cmp : -cmp;
			});
		}
	});
</script>

<!--
	only render the contents on the client as the server timezone may be different from that of the client,
	making the formatted time reformat after hydration, which is a lot more jarring than the content being
	loaded a little later than the header.
-->
{#if browser}
	{#if items.length > 0}
		<div class="p-4 pt-2">
			<div class="table w-full border-collapse p-2 text-lg">
				<div class="mb-2 table-row border-b border-neutral-300">
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
									<div class="invisible w-5"></div>
								{/if}
							</button>
						</span>
					{/snippet}

					{@render header("title")}

					{@render header("modelType", "model type")}

					{@render header("identifier")}

					{@render header("uploaded", route === "view" ? "uploaded" : "accessed")}
				</div>
				{#each items as { uuid, title, modelType, identifier, uploaded } (uuid)}
					{@const date = new Date(uploaded)}
					<a
						href={resolve(`/${route}/${uuid}`)}
						class="table-row border-b border-neutral-300 outline-indigo-400 last:border-none hover:bg-neutral-200 focus-visible:outline-2 dark:hover:bg-neutral-700"
						data-sveltekit-preload-data="tap"
						data-sveltekit-preload-code="hover"
					>
						<span class="table-cell p-3">{title}</span>
						<span class="table-cell p-3">{modelType}</span>
						<span class="table-cell p-3 text-nowrap">{identifier}</span>
						<span title={dateFormat(date, false)} class="table-cell p-3">
							{dateFormat(date, true)}
						</span>
						{#if del}
							<div class="table-cell w-1 p-2">
								<button
									class="cursor-pointer rounded-4xl outline-red-500 hover:bg-red-500/25 focus-visible:outline-2"
									title="delete this local model"
									onclick={(e) => {
										del(uuid);
										e.preventDefault();
									}}
								>
									<X class="stroke-red-500" />
								</button>
							</div>
						{/if}
					</a>
				{/each}
			</div>
		</div>
	{:else}
		<h1 class="m-4 mt-2 text-2xl font-semibold">no models found</h1>
	{/if}
{/if}
