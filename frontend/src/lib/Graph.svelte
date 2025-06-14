<script lang="ts">
	import { type Snippet } from "svelte";
	import { Background, BackgroundVariant, Panel, SvelteFlow } from "@xyflow/svelte";
	import "@xyflow/svelte/dist/style.css";
	import NodeComponent from "./graph/NodeComponent.svelte";
	import EdgeComponent from "./graph/EdgeComponent.svelte";
	import type { Node, Edge } from "$lib/model";

	let {
		interactive,
		nodes = $bindable(),
		edges = $bindable(),
		headerLeft,
		headerRight,
		panel,
	}: {
		interactive: boolean;
		nodes: Node[];
		edges: Edge[];
		headerLeft?: Snippet;
		headerRight?: Snippet;
		panel: Snippet;
	} = $props();

	const defaultEdge: Partial<Edge> = {
		type: "custom",
		data: {
			kind: "correlated-with",
			editable: interactive,
		},
	};
</script>

<SvelteFlow
	bind:nodes
	bind:edges
	colorMode="system"
	proOptions={{ hideAttribution: true }}
	nodeTypes={{ custom: NodeComponent }}
	edgeTypes={{ custom: EdgeComponent }}
	defaultEdgeOptions={defaultEdge}
	nodesConnectable={interactive}
	nodesDraggable={interactive}
	elementsSelectable={interactive}
	deleteKey={["Delete", "Backspace"]}
	fitView
	fitViewOptions={{ padding: { right: 1000 } }}
>
	{#if headerLeft}
		<Panel position="top-left">
			{@render headerLeft()}
		</Panel>
	{/if}
	{#if headerRight}
		<Panel position="top-right">
			{@render headerRight()}
		</Panel>
	{/if}

	<Panel position="bottom-right">
		{@render panel()}
	</Panel>

	<Background variant={BackgroundVariant.Dots} />
</SvelteFlow>
