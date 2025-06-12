<script lang="ts">
	import { onMount, type Snippet } from "svelte";
	import {
		Background,
		BackgroundVariant,
		Panel,
		SvelteFlow,
		type ColorMode,
		type OnConnect,
	} from "@xyflow/svelte";
	import "@xyflow/svelte/dist/style.css";
	import NodeComponent from "./graph/NodeComponent.svelte";
	import EdgeComponent from "./graph/EdgeComponent.svelte";
	import type { Node, Edge } from "$lib/model";
	import { browser } from "$app/environment";

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

	let colorMode: ColorMode = $state("light");

	if (browser) {
		let e = window.matchMedia("(prefers-color-scheme: dark)");
		colorMode = e.matches ? "dark" : "light";
		e.addEventListener("change", (e) => (colorMode = e.matches ? "dark" : "light"));
	}
</script>

<SvelteFlow
	bind:nodes
	bind:edges
	{colorMode}
	proOptions={{ hideAttribution: true }}
	nodeTypes={{ custom: NodeComponent }}
	edgeTypes={{ custom: EdgeComponent }}
	{onconnect}
	nodesConnectable={interactive}
	nodesDraggable={interactive}
	elementsSelectable={interactive}
	fitView
	fitViewOptions={{ padding: { right: 1000 } }}
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
