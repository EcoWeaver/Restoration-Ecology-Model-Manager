<script lang="ts">
	import type { Snippet } from "svelte";
	import { Background, BackgroundVariant, Panel, SvelteFlow, type OnConnect } from "@xyflow/svelte";
	import "@xyflow/svelte/dist/style.css";
	import NodeComponent from "./graph/NodeComponent.svelte";
	import EdgeComponent from "./graph/EdgeComponent.svelte";
	import type { Node, Edge } from "$lib/model";

	let {
		interactive,
		nodes = $bindable(),
		edges = $bindable(),
		info,
		panel,
		onconnect,
	}: {
		interactive: boolean;
		nodes: Node[];
		edges: Edge[];
		info?: Snippet;
		panel: Snippet;
		onconnect?: OnConnect;
	} = $props();
</script>

<SvelteFlow
	bind:nodes
	bind:edges
	proOptions={{ hideAttribution: true }}
	nodeTypes={{ custom: NodeComponent }}
	edgeTypes={{ custom: EdgeComponent }}
	{onconnect}
	nodesConnectable={interactive}
	nodesDraggable={interactive}
	elementsSelectable={interactive}
>
	{#if info}
		<Panel position="top-left">
			{@render info()}
		</Panel>
	{/if}

	<Panel position="bottom-right">
		{@render panel()}
	</Panel>

	<Background variant={BackgroundVariant.Dots} />
</SvelteFlow>
