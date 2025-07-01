<script lang="ts">
	import { browser } from "$app/environment";
	import Header from "$lib/Header.svelte";
	import type { ChipElement } from "$lib/search/Chip.svelte";

	let { data } = $props();
	let { items, params: searchParams } = $derived(data);

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
</script>

<svelte:head>
	<title>search by metadata</title>
</svelte:head>

<Header {params} />

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
					<span class="table-cell p-3 font-semibold">title</span>
					<span class="table-cell p-3 font-semibold">model-type</span>
					<span class="table-cell p-3 font-semibold">identifier</span>
					<span class="table-cell p-3 font-semibold">uploaded</span>
				</div>
				{#each items as { uuid, title, modelType, identifier, uploaded } (uuid)}
					{@const date = new Date(uploaded)}
					<a
						href={`/view/${uuid}`}
						class="table-row border-b-1 border-neutral-300 outline-indigo-400 last:border-none hover:bg-neutral-200 focus-visible:outline-2 dark:hover:bg-neutral-700"
					>
						<span class="table-cell p-3">{title}</span>
						<span class="table-cell p-3">{modelType}</span>
						<span class="table-cell p-3">{identifier}</span>
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
