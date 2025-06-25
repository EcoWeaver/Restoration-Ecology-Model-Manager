<script lang="ts">
	import { type Snippet } from "svelte";
	import {
		Background,
		BackgroundVariant,
		ConnectionLineType,
		ConnectionMode,
		Panel,
		SvelteFlow,
	} from "@xyflow/svelte";
	import NodeComponent from "./graph/NodeComponent.svelte";
	import EdgeComponent from "./graph/EdgeComponent.svelte";
	import type { Node, Edge } from "$lib/model";

	import "@xyflow/svelte/dist/style.css";

	let {
		interactive,
		nodes = $bindable(),
		edges = $bindable(),
		panel,
	}: {
		interactive: boolean;
		nodes: Node[];
		edges: Edge[];
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
	proOptions={{ hideAttribution: true }}
	nodeTypes={{ custom: NodeComponent }}
	edgeTypes={{ custom: EdgeComponent }}
	connectionLineType={ConnectionLineType.Straight}
	connectionMode={ConnectionMode.Loose}
	defaultEdgeOptions={defaultEdge}
	elevateEdgesOnSelect={false}
	nodesConnectable={interactive}
	nodesDraggable={interactive}
	elementsSelectable={interactive}
	nodesFocusable={interactive}
	edgesFocusable={interactive}
	deleteKey={["Delete", "Backspace"]}
	fitView={nodes.length > 0}
>
	<Panel position="bottom-right">
		{@render panel()}
	</Panel>

	<Background variant={BackgroundVariant.Dots} />
</SvelteFlow>
