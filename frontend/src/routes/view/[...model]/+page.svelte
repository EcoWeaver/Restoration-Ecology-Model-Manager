<script lang="ts">
	import Graph from "$lib/Graph.svelte";
	import type { Node, Edge } from "$lib/model";
	import PanelComponent from "$lib/graph/PanelComponent.svelte";

	let { data } = $props();
	let { model } = $derived(data);

	let nodes: Node[] = $derived(
		Array.from(
			model.graph.nodes?.entries().map(([id, { body, position, color }]) => ({
				position,
				type: "custom",
				id: id.toString(),
				data: { body, color, editable: false },
			})) || [],
		),
	);

	let edges: Edge[] = $derived(
		model.graph.relations?.map(({ from, to, kind }) => ({
			source: from.toString(),
			target: to.toString(),
			type: "custom",
			id: `${from}-${to}`,
			data: { kind, editable: false },
		})) || [],
	);
</script>

<div class="h-dvh w-dvw">
	<Graph {nodes} {edges} interactive={false}>
		{#snippet panel()}
			<PanelComponent>
				<div class="table">
					<div class="table-row">
						<div class="table-cell pr-2 pb-1 font-bold">kind</div>
						<div class="table-cell pb-1 pl-2 font-bold">value</div>
					</div>
					<div class="table-row">
						<div class="table-cell pr-2">identifier</div>
						<a
							class="pl-2 text-blue-700 hover:underline"
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
