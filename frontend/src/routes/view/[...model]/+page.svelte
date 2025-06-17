<script lang="ts">
	import { goto } from "$app/navigation";
	import Graph from "$lib/Graph.svelte";
	import PanelComponent from "$lib/graph/PanelComponent.svelte";
	import { apiToFlow } from "$lib/model";

	let { data } = $props();
	let { model } = $derived(data);
	let { nodes, edges } = $derived(apiToFlow(model.graph, false));

	async function edit() {
		await goto(`/edit/${model.identifier}?version=${model.version}`);
	}
</script>

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
						<div class="table-row">
							<div class="table-cell pr-2">{kind}</div>
							<div class="table-cell max-w-128 pl-2">{value}</div>
						</div>
					{/each}
				</div>
			</PanelComponent>
		{/snippet}
	</Graph>
</div>
