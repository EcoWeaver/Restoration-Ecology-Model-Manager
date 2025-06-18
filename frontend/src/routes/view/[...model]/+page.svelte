<script lang="ts">
	import { goto } from "$app/navigation";
	import Graph from "$lib/Graph.svelte";
	import PanelComponent from "$lib/graph/PanelComponent.svelte";
	import { apiToFlow } from "$lib/model";

	let { data } = $props();
	let { model, metadataMap } = $derived(data);
	let { nodes, edges } = $derived(apiToFlow(model.graph, false));

	$inspect(metadataMap);

	async function edit() {
		await goto(`/edit/${model.identifier}?version=${model.version}`);
	}
</script>

<svelte:head>
	<title>view {model.identifier}@v{model.version}</title>
</svelte:head>

<div class="h-dvh w-dvw">
	<Graph {nodes} {edges} interactive={false}>
		{#snippet headerLeft()}
			<button class="hover:bg-neutral-200 dark:hover:bg-neutral-700" onclick={edit}>Edit</button>
		{/snippet}

		{#snippet panel()}
			<PanelComponent>
				<h2 class="text-xl">Metadata</h2>

				<div class="table">
					<div class="table-row">
						<div class="table-cell pr-2 pb-1 font-bold">kind</div>
						<div class="table-cell pb-1 pl-2 font-bold">value</div>
					</div>
					<div class="table-row">
						<div class="table-cell pr-2">identifier</div>
						<a
							class="ml-2 text-blue-700 hover:underline dark:text-blue-500"
							href={`https://doi.org/${model.identifier}`}
							target="_blank"
							rel="noopener noreferrer"
						>
							{model.identifier}
						</a>
					</div>
					{#each model.metadata as { kind, value }}
						{@const type = metadataMap.get(kind)}

						<div class="table-row">
							<div class="table-cell pr-2">{kind}</div>

							{#if type === "orcid"}
								<a
									class="table-cell max-w-128 pl-2 text-blue-700 hover:underline dark:text-blue-500"
									href={`https://orcid.org/${value}`}
									target="_blank"
									rel="noopener noreferrer"
								>
									{value}
								</a>
							{:else if type === "url"}
								<a
									class="table-cell max-w-128 pl-2 text-blue-700 hover:underline dark:text-blue-500"
									href={value}
									target="_blank"
									rel="noopener noreferrer"
								>
									{value}
								</a>
							{:else}
								<div class="table-cell max-w-128 pl-2">{value}</div>
							{/if}
						</div>
					{/each}
				</div>
			</PanelComponent>
		{/snippet}
	</Graph>
</div>
