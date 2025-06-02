<script lang="ts">
	import Model from "$lib/Model.svelte";
	import type { Node } from "$lib/model/NodeComponent.svelte";
	import type { Edge } from "$lib/model/EdgeComponent.svelte";

	let { data } = $props();
	let { model } = $derived(data);

	let nodes: Node[] = $derived(
		Array.from(
			model.model.nodes?.entries().map(([id, { body, position, color }]) => ({
				position,
				type: "custom",
				id: id.toString(),
				data: { body, color, editable: false },
			})) || [],
		),
	);

	let edges: Edge[] = $derived(
		model.model.relations?.map(({ from, to, kind }) => ({
			source: from.toString(),
			target: to.toString(),
			type: "custom",
			id: `${from}-${to}`,
			data: { kind, editable: false },
		})) || [],
	);
</script>

<div class="h-dvh w-dvw">
	<Model {nodes} {edges} interactive={false}>
		{#snippet panel()}
			<div class="table">
				<div class="table-row">
					<div class="table-cell text-center font-bold">kind</div>
					<div class="table-cell text-center font-bold">value</div>
				</div>
				{#each model.metadata as { kind, value }}
					<div class="table-row">
						<div class="table-cell">{kind}</div>
						<div class="table-cell max-w-128">{value}</div>
					</div>
				{/each}
			</div>
		{/snippet}
	</Model>
</div>
