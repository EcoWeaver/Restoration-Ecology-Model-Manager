<script lang="ts">
	import { Link, PencilLine } from "@lucide/svelte";
	import Graph from "$lib/Graph.svelte";
	import PanelComponent from "$lib/graph/PanelComponent.svelte";
	import Header from "$lib/Header.svelte";
	import Alert from "$lib/Alert.svelte";
	import { apiToFlow } from "$lib/model";
	import { SvelteFlowProvider } from "@xyflow/svelte";
	import type { MetadataType } from "$lib/api.js";

	let { data } = $props();
	let { model, metadataElements } = $derived(data);
	let { nodes, edges } = $derived(apiToFlow(model.graph, false));

	const metadataMap = $derived.by(() => {
		const map = new Map<string, MetadataType | undefined>();
		for (const { name, type } of Object.values(metadataElements).flat()) {
			map.set(name, type);
		}
		return map;
	});

	let alert: ReturnType<typeof Alert>;
	async function share() {
		await navigator.clipboard.writeText(window.location.href);
		alert.info("link copied to clipboard");
	}
</script>

<svelte:head>
	<title>{model.metadata.find(({ kind }) => kind === "title")?.value ?? model.uuid}</title>
</svelte:head>

<div class="absolute z-10 w-full bg-white dark:bg-neutral-800">
	<Header {metadataElements}>
		<a
			class="flex h-10 cursor-pointer items-center rounded-md bg-indigo-300 px-3 outline-offset-2 outline-indigo-400 hover:bg-indigo-400 focus-visible:outline-2 dark:bg-indigo-500"
			title="edit model"
			href={`/edit/${model.uuid}`}
		>
			<PencilLine strokeWidth={1.7} />
			<span class="pl-2 select-none">edit</span>
		</a>

		<button
			class="flex h-10 cursor-pointer items-center rounded-md bg-indigo-300 px-3 outline-offset-2 outline-indigo-400 hover:bg-indigo-400 focus-visible:outline-2 dark:bg-indigo-500"
			title="copy model link to clipboard"
			onclick={share}
		>
			<Link />
			<span class="pl-2 select-none">share</span>
		</button>
	</Header>
</div>

<Alert bind:this={alert} />

<div class="h-dvh w-dvw">
	<SvelteFlowProvider>
		<Graph {nodes} {edges}>
			{#snippet panel()}
				<PanelComponent>
					<div class="table">
						<div class="table-row">
							<h2 class="mb-2 table-cell text-xl font-bold select-text">Metadata</h2>
							{#if model.parent}
								<span
									class="mb-2 table-cell max-w-128 overflow-hidden px-2 pt-1 text-ellipsis whitespace-nowrap select-text"
								>
									based on
									<a
										class="mb-2 text-blue-700 outline-indigo-400 hover:underline focus-visible:outline-2 dark:text-blue-400"
										title={model.parent.title}
										href="/view/{model.parent.uuid}"
									>
										"{model.parent.title}"
									</a>
								</span>
							{/if}
						</div>
						<div class="table-row">
							<div class="text-l table-cell py-1 pr-2 font-bold select-text">kind</div>
							<div class="text-l table-cell py-1 pl-2 font-bold select-text">value</div>
						</div>
						{#each model.metadata as { kind, value }, idx (idx)}
							{@const type = metadataMap.get(kind)}

							<div class="table-row">
								<div class="table-cell py-1 pr-2 select-text">{kind}</div>

								{#if type === "doi"}
									<a
										class="table-cell max-w-128 py-1 pl-2 text-blue-700 outline-indigo-400 select-text hover:underline focus-visible:outline-2 dark:text-blue-400"
										href={`https://doi.org/${value}`}
										target="_blank"
										rel="noopener noreferrer"
									>
										{value}
									</a>
								{:else if type === "orcid"}
									<a
										class="table-cell max-w-128 py-1 pl-2 text-blue-700 outline-indigo-400 select-text hover:underline focus-visible:outline-2 dark:text-blue-400"
										href={`https://orcid.org/${value}`}
										target="_blank"
										rel="noopener noreferrer"
									>
										{value}
									</a>
								{:else if type === "url"}
									<a
										class="table-cell max-w-128 py-1 pl-2 text-blue-700 outline-indigo-400 select-text hover:underline focus-visible:outline-2 dark:text-blue-400"
										href={value}
										target="_blank"
										rel="noopener noreferrer"
									>
										{value}
									</a>
								{:else}
									<div class="table-cell max-w-128 py-1 pl-2 select-text">{value}</div>
								{/if}
							</div>
						{/each}
					</div>
				</PanelComponent>
			{/snippet}
		</Graph>
	</SvelteFlowProvider>
</div>
