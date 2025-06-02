<script lang="ts">
	import {
		Background,
		SvelteFlow,
		Panel,
		type Node,
		type Edge,
		BackgroundVariant,
	} from "@xyflow/svelte";
	import "@xyflow/svelte/dist/style.css";

	let { data } = $props();
	let { model } = $derived(data);

	let nodes: Node[] = $derived(
		Array.from(
			model.model.nodes?.entries().map(([id, { position, color }]) => ({
				position,
				data: { color },
				id: id.toString(),
			})) || [],
		),
	);

	let edges: Edge[] = $derived(
		model.model.relations?.map(({ from, to, kind }) => ({
			source: from.toString(),
			target: to.toString(),
			data: { kind },
			id: `${from}-${to}`,
		})) || [],
	);
</script>

<div class="h-dvh w-dvw">
	<SvelteFlow
		{nodes}
		{edges}
		nodesConnectable={false}
		nodesDraggable={false}
		elementsSelectable={false}
	>
		<Panel position="top-right">
			<table>
				<thead>
					<tr>
						<th>kind</th>
						<th>value</th>
					</tr>
				</thead>
				<tbody>
					{#each model.metadata as { kind, value }}
						<tr>
							<td>{kind}</td>
							<td class="max-w-128">{value}</td>
						</tr>
					{/each}
				</tbody>
			</table>
		</Panel>

		<Background variant={BackgroundVariant.Dots} />
	</SvelteFlow>
</div>
