<script lang="ts">
	import type { Snippet } from "svelte";
	import { Background, BackgroundVariant, Panel, SvelteFlow, type OnConnect } from "@xyflow/svelte";
	import "@xyflow/svelte/dist/style.css";
	import NodeComponent, { type Node } from "./graph/NodeComponent.svelte";
	import EdgeComponent, { type Edge } from "./graph/EdgeComponent.svelte";

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
		<div
			class="bottom-0 m-3 max-h-[55vh] overflow-y-scroll rounded-xl bg-neutral-200/60 px-8 py-4 backdrop-blur-sm"
		>
			{@render panel()}
		</div>
	</Panel>

	<Background variant={BackgroundVariant.Dots} />
</SvelteFlow>
