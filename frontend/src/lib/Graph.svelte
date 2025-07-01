<script lang="ts">
	import { type Snippet } from "svelte";
	import {
		Background,
		BackgroundVariant,
		ConnectionLineType,
		ConnectionMode,
		Panel,
		SvelteFlow,
		useSvelteFlow,
		type OnConnectEnd,
		type XYPosition,
	} from "@xyflow/svelte";
	import NodeComponent from "./graph/NodeComponent.svelte";
	import EdgeComponent from "./graph/EdgeComponent.svelte";
	import type { Node, Edge } from "$lib/model";

	import "@xyflow/svelte/dist/style.css";

	let {
		interactive,
		nodes = $bindable(),
		edges = $bindable(),
		addNode = undefined,
		panel,
	}: {
		interactive: boolean;
		nodes: Node[];
		edges: Edge[];
		addNode?: (pos: XYPosition) => string;
		panel: Snippet;
	} = $props();

	const { screenToFlowPosition } = useSvelteFlow();

	const defaultEdge: Partial<Edge> = {
		type: "custom",
		data: {
			kind: "correlated-with",
			editable: interactive,
		},
	};

	const onconnectend: OnConnectEnd = (event, connectionState) => {
		if (!addNode) return;
		if (connectionState.isValid) return;

		const sourceId = connectionState.fromNode!.id;

		const { clientX, clientY } = "changedTouches" in event ? event.changedTouches[0] : event;
		const targetId = addNode(
			screenToFlowPosition({
				x: clientX,
				y: clientY,
			}),
		);

		const newEdge = {
			...defaultEdge,
			id: `${sourceId}-${targetId}`,
			source: sourceId,
			target: targetId,
		};
		edges = [...edges, newEdge];
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
	onconnectend={interactive ? onconnectend : undefined}
	deleteKey={["Delete", "Backspace"]}
	fitView={nodes.length > 0}
>
	<Panel position="bottom-right">
		{@render panel()}
	</Panel>

	<Background variant={BackgroundVariant.Dots} />
</SvelteFlow>
